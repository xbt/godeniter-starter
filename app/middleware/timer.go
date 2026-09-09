package middleware

import (
	"github.com/xbt/godeniter"
	coreMiddleware "github.com/xbt/godeniter/middleware"
)

// ResponseTimer 响应计时中间件
// 底层复用 godeniter/middleware.ServerTiming()，为响应头注入 X-Response-Time 与 W3C 标准 Server-Timing 标头
func ResponseTimer() godeniter.HandlerFunc {
	return godeniter.WrapMiddleware(coreMiddleware.ServerTiming())
}
