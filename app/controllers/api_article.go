package controllers

import (
	"github.com/xbt/godeniter"
	"github.com/xbt/godeniter/router"
	"github.com/xbt/godeniter/utils/upload"
	"godeniter-starter/app/models"
	"godeniter-starter/app/services"
	"math"
	"net/http"
	"strconv"
)

// ArticleAPIController 文章 RESTful API 控制器
type ArticleAPIController struct{}

// List 获取文章列表 (GET /api/v1/articles)
func (ctrl *ArticleAPIController) List(c *godeniter.Context, svc *services.ArticleService) {
	var query models.ArticleQueryRequest
	_ = c.BindQuery(&query)

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}

	items, total := svc.ListPaginate(query.Keyword, query.Page, query.PageSize)
	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	c.Success(godeniter.H{
		"items": items,
		"pagination": godeniter.H{
			"total":       total,
			"page":        query.Page,
			"page_size":   query.PageSize,
			"total_pages": totalPages,
			"has_next":    query.Page < totalPages,
			"has_prev":    query.Page > 1,
		},
	})
}

// Detail 获取文章详情 (GET /api/v1/articles/:id)
func (ctrl *ArticleAPIController) Detail(params router.Params, svc *services.ArticleService) (int, godeniter.H) {
	id, err := strconv.Atoi(params.Get("id"))
	if err != nil {
		return http.StatusBadRequest, godeniter.H{"code": 40001, "message": "非法的文章 ID"}
	}

	article, found := svc.GetByID(id)
	if !found {
		return http.StatusNotFound, godeniter.H{"code": 40401, "message": "文章不存在"}
	}

	return http.StatusOK, godeniter.H{
		"code":    0,
		"message": "ok",
		"data":    article,
	}
}

// Create 创建文章 (POST /api/v1/articles)
func (ctrl *ArticleAPIController) Create(c *godeniter.Context, svc *services.ArticleService) {
	var req models.CreateArticleRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Fail(40002, "参数校验失败: "+err.Error())
		return
	}

	article := svc.Create(req.Title, req.Content, req.Author, req.CoverURL)
	c.Success(article)
}

// Delete 删除文章 (DELETE /api/v1/articles/:id)
func (ctrl *ArticleAPIController) Delete(c *godeniter.Context, svc *services.ArticleService) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Fail(40001, "非法文章 ID")
		return
	}

	if !svc.Delete(id) {
		c.Fail(40401, "文章不存在或已被删除")
		return
	}

	c.Success(godeniter.H{"deleted_id": id})
}

// Upload 上传文件 (POST /api/v1/upload)
func (ctrl *ArticleAPIController) Upload(c *godeniter.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Fail(40010, "请选择上传文件: "+err.Error())
		return
	}

	opts := upload.Options{
		SaveDir:     "./uploads/images",
		MaxBytes:    5 * 1024 * 1024,
		AllowedExts: []string{".jpg", ".png", ".jpeg", ".webp"},
		AutoRename:  true,
	}

	savedPath, err := c.SaveUploadedFileWithOptions(file, opts)
	if err != nil {
		c.Fail(40011, "文件上传失败: "+err.Error())
		return
	}

	c.Success(godeniter.H{
		"filename":   file.Filename,
		"saved_path": savedPath,
		"url":        "/" + savedPath,
	})
}
