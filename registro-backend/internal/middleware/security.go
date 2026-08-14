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
		nonce, err := generateNonce()
		if err != nil {
			c.JSON(500, gin.H{"error": "internal security error"})
			c.Abort()
			return
		}

		// Store the nonce so templates / other middleware can reference it
		c.Set("csp_nonce", nonce)

		// Prevent Clickjacking
		c.Writer.Header().Set("X-Frame-Options", "DENY")

		// Prevent MIME sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		// Enforce HTTPS (HSTS) — 2 years, include subdomains, preload-ready
		// Only set HSTS header on HTTPS connections or when running behind a TLS proxy
		isHTTPS := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
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
		csp := fmt.Sprintf(
			"default-src 'self'; "+
				"script-src 'self' 'nonce-%s'; "+
				"style-src 'self' 'nonce-%s' https://fonts.googleapis.com; "+
				"font-src 'self' https://fonts.gstatic.com; "+
				"img-src 'self' data: https://cdn.quasar.dev; "+
				"connect-src 'self' %s %s; "+
				"object-src 'none'; "+
				"frame-src 'none'; "+
				"base-uri 'self'; "+
				"form-action 'self'; "+
				"frame-ancestors 'none'; "+
				"upgrade-insecure-requests;",
			nonce, nonce, backendURL(), backendWSURL(),
		)
		c.Writer.Header().Set("Content-Security-Policy", csp)

		c.Next()
	}
}
