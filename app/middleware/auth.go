package middleware

import (
	"github.com/xbt/godeniter"
	"github.com/xbt/godeniter/session"
	"net/http"
)

// AuthRequired 登录认证拦截中间件 (路由守卫)
func AuthRequired() godeniter.HandlerFunc {
	return func(c *godeniter.Context) {
		sessVal, ok := c.Session()
		if !ok || sessVal == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		sess, ok := sessVal.(session.Session)
		if !ok || sess.GetString("username") == "" {
			sess.SetFlash("notice", "请先登录管理员账号后再进行后台操作。")
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}

