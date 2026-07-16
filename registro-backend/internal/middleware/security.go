package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

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

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		nonce := generateNonce()

		// Store the nonce so templates / other middleware can reference it
		c.Set("csp_nonce", nonce)
		// Expose it to the frontend via header so the SPA can pass it to
		// dynamically-injected scripts (Quasar reads this on boot if configured)
		c.Writer.Header().Set("X-CSP-Nonce", nonce)

		// Prevent Clickjacking
		c.Writer.Header().Set("X-Frame-Options", "DENY")

		// Prevent MIME sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		// Enforce HTTPS (HSTS) — 2 years, include subdomains, preload-ready
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		// Referrer Policy
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// X-XSS-Protection (legacy IE/Edge compatibility)
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		// Permissions-Policy: restrict powerful browser APIs not needed by the app
		c.Writer.Header().Set("Permissions-Policy",
			"camera=(), microphone=(), geolocation=(), payment=(), usb=(), magnetometer=(), gyroscope=()")

		// Content Security Policy — nonce-based, no unsafe-inline / unsafe-eval
		// Quasar bundles its scripts as static assets; any runtime eval requirement
		// should be addressed by using 'wasm-unsafe-eval' only if WASM is used.
		csp := fmt.Sprintf(
			"default-src 'self'; "+
				"script-src 'self' 'nonce-%s'; "+
				"style-src 'self' 'nonce-%s' https://fonts.googleapis.com; "+
				"font-src 'self' https://fonts.gstatic.com; "+
				"img-src 'self' data: blob: https://cdn.quasar.dev; "+
				"connect-src 'self' https://registro-backend-fdu2.onrender.com http://localhost:8080 http://localhost:5173 ws: wss:; "+
				"object-src 'none'; "+
				"base-uri 'self'; "+
				"form-action 'self'; "+
				"frame-ancestors 'none';",
			nonce, nonce,
		)
		c.Writer.Header().Set("Content-Security-Policy", csp)

		c.Next()
	}
}
