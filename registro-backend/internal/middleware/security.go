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
func generateNonce() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback: should never happen, but prevents nil panic
		return "fallback-nonce-replace-me"
	}
	return base64.StdEncoding.EncodeToString(b)
}

// backendURL returns the backend origin used in connect-src.
// It is read from the BACKEND_URL environment variable so that the production
// URL is never hardcoded in source code. Falls back to localhost for local dev.
func backendURL() string {
	if url := strings.TrimSpace(os.Getenv("BACKEND_URL")); url != "" {
		return url
	}
	return "http://localhost:8080"
}

// backendWSURL returns the WebSocket origin derived from BACKEND_WS_URL or BACKEND_URL.
func backendWSURL() string {
	if ws := strings.TrimSpace(os.Getenv("BACKEND_WS_URL")); ws != "" {
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
		nonce := generateNonce()

		// Store the nonce so templates / other middleware can reference it
		c.Set("csp_nonce", nonce)

		// Prevent Clickjacking
		c.Writer.Header().Set("X-Frame-Options", "DENY")

		// Prevent MIME sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		// Enforce HTTPS (HSTS) — 2 years, include subdomains, preload-ready
		// Only set HSTS header on HTTPS connections or when running in release mode
		isHTTPS := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" || gin.Mode() == gin.ReleaseMode || os.Getenv("SERVER_MODE") == "release"
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

		// Content Security Policy — nonce-based, no unsafe-inline / unsafe-eval
		// Specific backend URL and WebSocket origins are allowed rather than un-scoped wildcard ws:/wss:
		csp := fmt.Sprintf(
			"default-src 'self'; "+
				"script-src 'self' 'nonce-%s'; "+
				"style-src 'self' 'nonce-%s' https://fonts.googleapis.com; "+
				"font-src 'self' https://fonts.gstatic.com; "+
				"img-src 'self' data: blob: https://cdn.quasar.dev; "+
				"connect-src 'self' %s %s; "+
				"object-src 'none'; "+
				"base-uri 'self'; "+
				"form-action 'self'; "+
				"frame-ancestors 'none';",
			nonce, nonce, backendURL(), backendWSURL(),
		)
		c.Writer.Header().Set("Content-Security-Policy", csp)

		c.Next()
	}
}
