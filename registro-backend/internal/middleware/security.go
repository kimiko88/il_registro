package middleware

import (
	"github.com/gin-gonic/gin"
)

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent Clickjacking
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		
		// Prevent MIME sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		
		// Enforce HTTPS (HSTS)
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		
		// Referrer Policy
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// X-XSS-Protection (legacy browsers)
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		// Content Security Policy (CSP)
		// Restricts loading resources to trusted domains
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.quasar.dev; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com https://cdn.quasar.dev; font-src 'self' https://fonts.gstatic.com; img-src 'self' data: https://cdn.quasar.dev; connect-src 'self' https://registro-backend-fdu2.onrender.com http://localhost:8080 http://localhost:5173 ws: wss:;")

		c.Next()
	}
}
