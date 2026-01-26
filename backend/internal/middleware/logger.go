package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns a gin middleware for logging
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Printf("[%s] %s %s %d %v",
			method,
			clientIP,
			path,
			statusCode,
			latency,
		)
	}
}
