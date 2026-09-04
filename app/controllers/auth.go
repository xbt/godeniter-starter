package controllers

import (
	"github.com/xbt/godeniter"
	"github.com/xbt/godeniter/session"
	"net/http"
)

// AuthController 用户认证与 Session 控制器
type AuthController struct{}

// LoginForm 显示登录页面
func (ctrl *AuthController) LoginForm(c *godeniter.Context) {
	c.HTML(http.StatusOK, "login.html", godeniter.H{
		"Title": "系统登录",
	})
}

// LoginSubmit 处理登录表单提交
func (ctrl *AuthController) LoginSubmit(c *godeniter.Context, sess session.Session) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 简单演示校验
	if username == "admin" && password == "123456" {
		if sess != nil {
			sess.Set("username", username)
		}
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "login.html", godeniter.H{
		"Title":    "系统登录",
		"Username": username,
		"Error":    "账号或密码错误 (提示: admin / 123456)",
	})
}

// Logout 注销登录
func (ctrl *AuthController) Logout(c *godeniter.Context, sess session.Session) {
	if sess != nil {
		sess.Delete("username")
	}
	c.Redirect(http.StatusFound, "/")
}
