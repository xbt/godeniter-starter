package controllers

import (
	"github.com/xbt/godeniter"
	"github.com/xbt/godeniter/session"
	"github.com/xbt/godeniter/utils/str"
	"godeniter-starter/app/services"
	"net/http"
)

// HomeController Web 页面控制器 (服务端渲染)
type HomeController struct{}

// Index 渲染首页
func (ctrl *HomeController) Index(c *godeniter.Context, svc *services.ArticleService, sess session.Session) {
	var username string
	var flashNotice string
	if sess != nil {
		username = sess.GetString("username")
		flashNotice = sess.GetFlashString("notice")
	}

	items, total := svc.ListPaginate("", 1, 10)

	c.HTML(http.StatusOK, "index.html", godeniter.H{
		"Title":       "Godeniter Starter Dashboard",
		"CurrentUser": username,
		"FlashNotice": flashNotice,
		"TotalCount":  total,
		"Articles":    items,
		"RandomUUID":  str.UUID(),
	})
}
