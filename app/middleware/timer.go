package middleware

import (
	"fmt"
	"time"

	"github.com/xbt/godeniter"
)

// ResponseTimer 响应计时中间件 (性能监控与 Server-Timing 标头注入)
func ResponseTimer() godeniter.HandlerFunc {
	return func(c *godeniter.Context) {
		start := time.Now()

		// 执行后续处理链
		c.Next()

		duration := time.Since(start)
		durMs := fmt.Sprintf("%.2fms", float64(duration.Microseconds())/1000.0)

		// 注入响应头，便于前端或接口调用方监测后端性能
		c.Header("X-Response-Time", durMs)
		c.Header("Server-Timing", fmt.Sprintf("app;dur=%.2f", float64(duration.Microseconds())/1000.0))
	}
}
