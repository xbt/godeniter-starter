package controllers

import (
	"net/http"
	"strconv"

	"github.com/xbt/godeniter"
	"github.com/xbt/godeniter/session"
	"godeniter-starter/app/models"
	"godeniter-starter/app/services"
)

// AdminController 后台管理控制器 (文章完整 CRUD 标配)
type AdminController struct{}

// checkAuth 统一登录态鉴权，未登录则拦截跳转
func (ctrl *AdminController) checkAuth(c *godeniter.Context, sess session.Session) string {
	if sess == nil {
		c.Redirect(http.StatusFound, "/login")
		return ""
	}
	username := sess.GetString("username")
	if username == "" {
		sess.SetFlash("notice", "请先登录管理员账号。")
		c.Redirect(http.StatusFound, "/login")
		return ""
	}
	return username
}

// List 文章管理列表页
func (ctrl *AdminController) List(c *godeniter.Context, svc *services.ArticleService, sess session.Session) {
	username := ctrl.checkAuth(c, sess)
	if username == "" {
		return
	}

	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize := 10

	items, total := svc.ListPaginate(keyword, page, pageSize)
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	c.HTML(http.StatusOK, "admin.html", godeniter.H{
		"Title":       "文章管理中心",
		"CurrentUser": username,
		"FlashNotice": sess.GetFlashString("notice"),
		"Articles":    items,
		"TotalCount":  total,
		"Keyword":     keyword,
		"CurrentPage": page,
		"TotalPages":  totalPages,
	})
}

// CreateForm 显示新增文章表单页面
func (ctrl *AdminController) CreateForm(c *godeniter.Context, sess session.Session) {
	username := ctrl.checkAuth(c, sess)
	if username == "" {
		return
	}

	c.HTML(http.StatusOK, "article_form.html", godeniter.H{
		"Title":       "发布新文章",
		"CurrentUser": username,
		"IsEdit":      false,
		"ActionURL":   "/admin/articles/create",
		"Article":     &models.Article{Author: username},
	})
}

// CreateSubmit 处理新增文章提交
func (ctrl *AdminController) CreateSubmit(c *godeniter.Context, svc *services.ArticleService, sess session.Session) {
	username := ctrl.checkAuth(c, sess)
	if username == "" {
		return
	}

	title := c.PostForm("title")
	author := c.PostForm("author")
	content := c.PostForm("content")

	if title == "" || content == "" {
		c.HTML(http.StatusOK, "article_form.html", godeniter.H{
			"Title":       "发布新文章",
			"CurrentUser": username,
			"IsEdit":      false,
			"ActionURL":   "/admin/articles/create",
			"Article":     &models.Article{Title: title, Author: author, Content: content},
			"Error":       "文章标题和内容均不能为空！",
		})
		return
	}

	svc.Create(title, content, author, "")
	sess.SetFlash("notice", "🎉 新文章发布成功！")
	c.Redirect(http.StatusFound, "/admin/articles")
}

// EditForm 显示编辑文章表单页面
func (ctrl *AdminController) EditForm(c *godeniter.Context, svc *services.ArticleService, sess session.Session) {
	username := ctrl.checkAuth(c, sess)
	if username == "" {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/articles")
		return
	}

	article, ok := svc.FindByID(id)
	if !ok {
		sess.SetFlash("notice", "文章不存在或已被删除。")
		c.Redirect(http.StatusFound, "/admin/articles")
		return
	}

	c.HTML(http.StatusOK, "article_form.html", godeniter.H{
		"Title":       "编辑文章 #" + strconv.Itoa(id),
		"CurrentUser": username,
		"IsEdit":      true,
		"ActionURL":   "/admin/articles/edit/" + strconv.Itoa(id),
		"Article":     article,
	})
}

// EditSubmit 处理编辑文章提交
func (ctrl *AdminController) EditSubmit(c *godeniter.Context, svc *services.ArticleService, sess session.Session) {
	username := ctrl.checkAuth(c, sess)
	if username == "" {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/articles")
		return
	}

	title := c.PostForm("title")
	author := c.PostForm("author")
	content := c.PostForm("content")

	if title == "" || content == "" {
		c.HTML(http.StatusOK, "article_form.html", godeniter.H{
			"Title":       "编辑文章 #" + strconv.Itoa(id),
			"CurrentUser": username,
			"IsEdit":      true,
			"ActionURL":   "/admin/articles/edit/" + strconv.Itoa(id),
			"Article":     &models.Article{ID: id, Title: title, Author: author, Content: content},
			"Error":       "文章标题和内容均不能为空！",
		})
		return
	}

	svc.Update(id, title, content, author)
	sess.SetFlash("notice", "✏️ 文章修改已保存！")
	c.Redirect(http.StatusFound, "/admin/articles")
}

// Delete 处理删除文章
func (ctrl *AdminController) Delete(c *godeniter.Context, svc *services.ArticleService, sess session.Session) {
	username := ctrl.checkAuth(c, sess)
	if username == "" {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err == nil {
		svc.Delete(id)
		sess.SetFlash("notice", "🗑️ 文章已成功删除！")
	}

	c.Redirect(http.StatusFound, "/admin/articles")
}
