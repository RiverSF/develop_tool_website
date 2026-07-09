package middleware

import (
	"develop_tools/pkg/logger"
	"github.com/gin-gonic/gin"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求前
		t := time.Now()

		c.Next()

		// 请求后
		latency := time.Since(t).Microseconds()

		// 获取发送的 status
		status := c.Writer.Status()

		// 记录本次请求日志
		logger.Info("get a request, api=%s, userAgent=%s, status:%d, latency=%d", c.Request.URL, c.Request.UserAgent(), status, latency)
	}
}
