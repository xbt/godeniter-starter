package controllers

import (
	"net/http"
	"strconv"

	"github.com/xbt/godeniter"
	"github.com/xbt/godeniter/session"
	"github.com/xbt/godeniter/utils/upload"
	"godeniter-starter/app/models"
	"godeniter-starter/app/services"
)

// AdminController 后台管理控制器 (已通过 AuthRequired 中间件统一路由鉴权守卫)
type AdminController struct{}

// getUsername 安全提取当前登录管理员账号名
func (ctrl *AdminController) getUsername(sess session.Session) string {
	if sess == nil {
		return "admin"
	}
	u := sess.GetString("username")
	if u == "" {
		return "admin"
	}
	return u
}

// List 文章管理列表页 (GET /admin 或 /admin/articles)
func (ctrl *AdminController) List(c *godeniter.Context, svc *services.ArticleService, sess session.Session) {
	username := ctrl.getUsername(sess)

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

// CreateForm 显示新增文章表单页面 (GET /admin/articles/create)
func (ctrl *AdminController) CreateForm(c *godeniter.Context, sess session.Session) {
	username := ctrl.getUsername(sess)

	c.HTML(http.StatusOK, "article_form.html", godeniter.H{
		"Title":       "发布新文章",
		"CurrentUser": username,
		"IsEdit":      false,
		"ActionURL":   "/admin/articles/create",
		"Article":     &models.Article{Author: username},
	})
}

// CreateSubmit 处理新增文章提交 (POST /admin/articles/create)
// 演示特性：1. 结构体 Tag 自动验证 (BindAndValidate)；2. 封面图片安全上传与限制 (SaveUploadedFileWithOptions)
func (ctrl *AdminController) CreateSubmit(c *godeniter.Context, svc *services.ArticleService, sess session.Session) {
	username := ctrl.getUsername(sess)

	var req models.FormArticleRequest
	// 特性演示：使用框架纯标准库实现的结构体 Tag 校验器
	if err := c.BindAndValidate(&req); err != nil {
		c.HTML(http.StatusOK, "article_form.html", godeniter.H{
			"Title":       "发布新文章",
			"CurrentUser": username,
			"IsEdit":      false,
			"ActionURL":   "/admin/articles/create",
			"Article":     &models.Article{Title: req.Title, Author: req.Author, Content: req.Content},
			"Error":       "表单参数校验失败: " + err.Error(),
		})
		return
	}

	// 特性演示：处理本地封面图片上传与限制
	coverURL := ""
	file, fileErr := c.FormFile("cover")
	if fileErr == nil && file != nil {
		opts := upload.Options{
			SaveDir:     "./uploads/images",
			MaxBytes:    5 * 1024 * 1024, // 5MB 上限
			AllowedExts: []string{".jpg", ".png", ".jpeg", ".webp"},
			AutoRename:  true,
		}
		savedPath, saveErr := c.SaveUploadedFileWithOptions(file, opts)
		if saveErr != nil {
			c.HTML(http.StatusOK, "article_form.html", godeniter.H{
				"Title":       "发布新文章",
				"CurrentUser": username,
				"IsEdit":      false,
				"ActionURL":   "/admin/articles/create",
				"Article":     &models.Article{Title: req.Title, Author: req.Author, Content: req.Content},
				"Error":       "封面图片上传失败: " + saveErr.Error(),
			})
			return
		}
		coverURL = "/" + savedPath
	}

	svc.Create(req.Title, req.Content, req.Author, coverURL)
	sess.SetFlash("notice", "🎉 新文章发布成功！")
	c.Redirect(http.StatusFound, "/admin/articles")
}

// EditForm 显示编辑文章表单页面 (GET /admin/articles/edit/:id)
func (ctrl *AdminController) EditForm(c *godeniter.Context, svc *services.ArticleService, sess session.Session) {
	username := ctrl.getUsername(sess)

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

// EditSubmit 处理编辑文章提交 (POST /admin/articles/edit/:id)
func (ctrl *AdminController) EditSubmit(c *godeniter.Context, svc *services.ArticleService, sess session.Session) {
	username := ctrl.getUsername(sess)

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

	var req models.FormArticleRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.HTML(http.StatusOK, "article_form.html", godeniter.H{
			"Title":       "编辑文章 #" + strconv.Itoa(id),
			"CurrentUser": username,
			"IsEdit":      true,
			"ActionURL":   "/admin/articles/edit/" + strconv.Itoa(id),
			"Article":     &models.Article{ID: id, Title: req.Title, Author: req.Author, Content: req.Content, CoverURL: article.CoverURL},
			"Error":       "表单参数校验失败: " + err.Error(),
		})
		return
	}

	coverURL := article.CoverURL
	// 若用户上传了新封面，则替换
	file, fileErr := c.FormFile("cover")
	if fileErr == nil && file != nil {
		opts := upload.Options{
			SaveDir:     "./uploads/images",
			MaxBytes:    5 * 1024 * 1024,
			AllowedExts: []string{".jpg", ".png", ".jpeg", ".webp"},
			AutoRename:  true,
		}
		savedPath, saveErr := c.SaveUploadedFileWithOptions(file, opts)
		if saveErr != nil {
			c.HTML(http.StatusOK, "article_form.html", godeniter.H{
				"Title":       "编辑文章 #" + strconv.Itoa(id),
				"CurrentUser": username,
				"IsEdit":      true,
				"ActionURL":   "/admin/articles/edit/" + strconv.Itoa(id),
				"Article":     &models.Article{ID: id, Title: req.Title, Author: req.Author, Content: req.Content, CoverURL: article.CoverURL},
				"Error":       "新封面上传失败: " + saveErr.Error(),
			})
			return
		}
		coverURL = "/" + savedPath
	}

	svc.Update(id, req.Title, req.Content, req.Author, coverURL)
	sess.SetFlash("notice", "✏️ 文章修改已保存！")
	c.Redirect(http.StatusFound, "/admin/articles")
}

// Delete 处理删除文章 (GET /admin/articles/delete/:id)
func (ctrl *AdminController) Delete(c *godeniter.Context, svc *services.ArticleService, sess session.Session) {
	id, err := strconv.Atoi(c.Param("id"))
	if err == nil {
		svc.Delete(id)
		sess.SetFlash("notice", "🗑️ 文章已成功删除！")
	}

	c.Redirect(http.StatusFound, "/admin/articles")
}
