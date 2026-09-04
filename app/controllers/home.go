package controllers

import (
	"net/http"
	"strconv"

	"github.com/xbt/godeniter"
	"github.com/xbt/godeniter/session"
	"github.com/xbt/godeniter/utils/str"
	"godeniter-starter/app/services"
)

// HomeController Web 页面控制器 (服务端渲染)
type HomeController struct{}

// Index 渲染首页 (支持关键词搜索、5条/页分页)
func (ctrl *HomeController) Index(c *godeniter.Context, svc *services.ArticleService, sess session.Session) {
	var username string
	var flashNotice string
	if sess != nil {
		username = sess.GetString("username")
		flashNotice = sess.GetFlashString("notice")
	}

	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize := 5

	items, total := svc.ListPaginate(keyword, page, pageSize)

	// 计算分页元数据
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	pages := make([]int, 0, totalPages)
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	c.HTML(http.StatusOK, "index.html", godeniter.H{
		"Title":       "Godeniter Starter 博客系统",
		"CurrentUser": username,
		"FlashNotice": flashNotice,
		"TotalCount":  total,
		"Articles":    items,
		"Keyword":     keyword,
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"Pages":       pages,
		"HasPrev":     page > 1,
		"HasNext":     page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
		"RandomUUID":  str.UUID(),
	})
}

// Detail 渲染文章详情页 (动态路由参数 :id，自动自增阅读量)
func (ctrl *HomeController) Detail(c *godeniter.Context, svc *services.ArticleService, sess session.Session) {
	var username string
	if sess != nil {
		username = sess.GetString("username")
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	article, ok := svc.GetByID(id)
	if !ok {
		c.HTML(http.StatusNotFound, "detail.html", godeniter.H{
			"Title":       "文章不存在",
			"CurrentUser": username,
			"Error":       "抱歉，该文章不存在或已被作者删除。",
		})
		return
	}

	c.HTML(http.StatusOK, "detail.html", godeniter.H{
		"Title":       article.Title,
		"CurrentUser": username,
		"Article":     article,
	})
}

// Features 渲染框架特性体验中心页面 (/features)
func (ctrl *HomeController) Features(c *godeniter.Context, sess session.Session) {
	var username string
	if sess != nil {
		username = sess.GetString("username")
	}

	c.HTML(http.StatusOK, "features.html", godeniter.H{
		"Title":       "Godeniter 2.0 框架特性全景体验中心",
		"CurrentUser": username,
		"ActiveNav":   "features",
	})
}

// PanicDemo 演示 Recovery 中间件捕获 Panic 不崩服 (/demo/panic)
func (ctrl *HomeController) PanicDemo(c *godeniter.Context) {
	panic("【模拟业务故障】这是一次由 /demo/panic 故意触发的 Runtime Panic！请观察控制台彩色错误堆栈输出与服务端返回的 500 状态码。服务未宕机，依然平稳运行！")
}

