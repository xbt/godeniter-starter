package middleware

import (
	"github.com/xbt/godeniter"
)

// SecurityHeaders 基础 Web 安全防护头中间件
// 自动注入 X-Content-Type-Options、X-Frame-Options、X-XSS-Protection 等行业标准安全标头
func SecurityHeaders() godeniter.HandlerFunc {
	return func(c *godeniter.Context) {
		// 防御 MIME 类型嗅探攻击
		c.Header("X-Content-Type-Options", "nosniff")
		// 防御点击劫持 (Clickjacking)，仅允许同源嵌入 iframe
		c.Header("X-Frame-Options", "SAMEORIGIN")
		// 开启浏览器内置 XSS 过滤防护
		c.Header("X-XSS-Protection", "1; mode=block")
		// 控制 Referer 泄漏
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}
