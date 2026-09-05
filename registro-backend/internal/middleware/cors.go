package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func buildOriginsMap() map[string]bool {
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsStr != "" {
		allowedOrigins = strings.Split(allowedOriginsStr, ",")
	} else {
		allowedOrigins = []string{
			"http://localhost:5173",
			"http://localhost:3000",
			"http://localhost:8080",
			"http://localhost:9000",
		}

	}
	originsMap := make(map[string]bool)
	for _, o := range allowedOrigins {
		trimmed := strings.TrimSpace(o)
		if trimmed != "" && trimmed != "*" && (strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://")) {
			originsMap[trimmed] = true
		}
	}
	return originsMap
}

func CORSMiddleware() gin.HandlerFunc {
	originsMap := buildOriginsMap()
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Add("Vary", "Origin")

		isAllowed := originsMap[origin]
		if isAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			if origin == "" {
				c.AbortWithStatusJSON(400, gin.H{"error": "Origin header required for preflight OPTIONS request"})
				return
			}
			if !isAllowed {
				c.AbortWithStatus(403)
				return
			}
			reqMethod := c.Request.Header.Get("Access-Control-Request-Method")
			if reqMethod != "" {
				allowedMethods := map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true, "PATCH": true, "OPTIONS": true}
				if !allowedMethods[strings.ToUpper(reqMethod)] {
					c.AbortWithStatus(405)
					return
				}
			}
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
