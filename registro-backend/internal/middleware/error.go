package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			// Change this to use structured logger later
			for _, e := range c.Errors {
				log.Printf("Error: %s", e.Error())
			}
			c.JSON(-1, c.Errors)
		}
	}
}
