package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// generateNonce creates a cryptographically secure random nonce (16 bytes → 24-char base64).
func generateNonce() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// sanitizeCSPOrigin strips any control characters, newlines, semicolons, or quotes
// to prevent header/CSP injection via environment variables.
func sanitizeCSPOrigin(raw string) string {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.ReplaceAll(cleaned, "\n", "")
	cleaned = strings.ReplaceAll(cleaned, "\r", "")
	cleaned = strings.ReplaceAll(cleaned, ";", "")
	cleaned = strings.ReplaceAll(cleaned, "'", "")
	cleaned = strings.ReplaceAll(cleaned, "\"", "")
	return cleaned
}

// backendURL returns the backend origin used in connect-src.
func backendURL() string {
	if url := sanitizeCSPOrigin(os.Getenv("BACKEND_URL")); url != "" {
		return url
	}
	if os.Getenv("APP_ENV") == "production" || os.Getenv("GIN_MODE") == "release" {
		return ""
	}
	return "http://localhost:8080"
}

// backendWSURL returns the WebSocket origin derived from BACKEND_WS_URL or BACKEND_URL.
func backendWSURL() string {
	if ws := sanitizeCSPOrigin(os.Getenv("BACKEND_WS_URL")); ws != "" {
		return ws
	}
	bURL := backendURL()
	if strings.HasPrefix(bURL, "https://") {
		return strings.Replace(bURL, "https://", "wss://", 1)
	}
	return strings.Replace(bURL, "http://", "ws://", 1)
}

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent Clickjacking
		c.Writer.Header().Set("X-Frame-Options", "DENY")

		// Prevent MIME sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		// Enforce HTTPS (HSTS) — 2 years, include subdomains, preload-ready
		isHTTPS := c.Request.TLS != nil
		if !isHTTPS && os.Getenv("TRUST_PROXY_HEADERS") == "true" {
			isHTTPS = c.GetHeader("X-Forwarded-Proto") == "https" || c.GetHeader("X-Forwarded-Ssl") == "on"
		}
		if isHTTPS {
			c.Writer.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		// Referrer Policy
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// X-XSS-Protection (legacy IE/Edge compatibility)
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		// Permissions-Policy: restrict powerful browser APIs not needed by the app
		c.Writer.Header().Set("Permissions-Policy",
			"camera=(), microphone=(), geolocation=(), payment=(), usb=(), magnetometer=(), gyroscope=()")

		// Content Security Policy — strict nonce-based CSP without 'unsafe-inline'
		nonce, err := generateNonce()
		if err == nil {
			c.Set("csp_nonce", nonce)
			upgradeInsecure := ""
			if isHTTPS || os.Getenv("APP_ENV") == "production" || os.Getenv("GIN_MODE") == "release" {
				upgradeInsecure = " upgrade-insecure-requests;"
			}

			connectSources := []string{"'self'"}
			if bURL := backendURL(); bURL != "" {
				connectSources = append(connectSources, bURL)
			}
			if wsURL := backendWSURL(); wsURL != "" {
				connectSources = append(connectSources, wsURL)
			}
			connectSources = append(connectSources, "https://*.supabase.co")

			var uniqueSources []string
			seen := make(map[string]bool)
			for _, src := range connectSources {
				src = strings.TrimSpace(src)
				if src != "" && !seen[src] {
					seen[src] = true
					uniqueSources = append(uniqueSources, src)
				}
			}
			connectSrcStr := strings.Join(uniqueSources, " ")

			csp := fmt.Sprintf(
				"default-src 'self'; "+
					"script-src 'self' 'nonce-%s'; "+
					"style-src 'self' 'nonce-%s' https://fonts.googleapis.com; "+
					"font-src 'self' https://fonts.gstatic.com; "+
					"img-src 'self' data: https://cdn.quasar.dev; "+
					"connect-src %s; "+
					"object-src 'none'; "+
					"frame-src 'none'; "+
					"base-uri 'self'; "+
					"form-action 'self'; "+
					"frame-ancestors 'none';%s",
				nonce, nonce, connectSrcStr, upgradeInsecure,
			)
			c.Writer.Header().Set("Content-Security-Policy", csp)
		}

		c.Next()
	}
}
