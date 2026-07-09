package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"develop_tools/pkg/logger"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/assets") {
			c.Next()
			return
		}

		t := time.Now()
		c.Next()

		latency := time.Since(t).Microseconds()
		status := c.Writer.Status()
		logger.Info("get a request, api=%s, userAgent=%s, status:%d, latency=%d", c.Request.URL, c.Request.UserAgent(), status, latency)
	}
}
