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
