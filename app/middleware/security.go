package middleware

import (
	"github.com/xbt/godeniter"
	coreMiddleware "github.com/xbt/godeniter/middleware"
)

// SecurityHeaders 基础 Web 安全防护头中间件
// 底层复用 godeniter/middleware.Security()，自动注入 nosniff、SAMEORIGIN、XSS 过滤等行业标准标头
func SecurityHeaders() godeniter.HandlerFunc {
	return godeniter.WrapMiddleware(coreMiddleware.Security())
}

// BlockSensitiveFiles 阻断敏感系统文件与数据库探测中间件
// 底层复用 godeniter/middleware.BlockSensitive()，自动拦截针对 .db、.sqlite、config.json、.env 等的探测，直接返回 403 Forbidden 并中断执行
func BlockSensitiveFiles() godeniter.HandlerFunc {
	return godeniter.WrapMiddleware(coreMiddleware.BlockSensitive())
}
