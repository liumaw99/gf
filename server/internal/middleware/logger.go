package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Logger 结构化请求日志中间件
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s | %3d | %13v | %15s | %-7s %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
		)
	})
}
