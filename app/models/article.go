package models

import "time"

// Article 文章数据实体
type Article struct {
	ID         int       `json:"id" db:"id"`
	Title      string    `json:"title" db:"title"`
	Content    string    `json:"content" db:"content"`
	Author     string    `json:"author" db:"author"`
	AuthorMask string    `json:"author_mask"` // 脱敏后的联系方式
	CoverURL   string    `json:"cover_url" db:"cover_url"`
	Views      int       `json:"views" db:"views"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// CreateArticleRequest 创建文章请求 DTO (用于 RESTful API，支持 JSON 绑定与校验)
type CreateArticleRequest struct {
	Title    string `json:"title" binding:"required,min=2,max=60"`
	Content  string `json:"content" binding:"required,min=5"`
	Author   string `json:"author" binding:"required"`
	CoverURL string `json:"cover_url"`
}

// ArticleQueryRequest 文章列表查询分页 DTO
type ArticleQueryRequest struct {
	Keyword  string `form:"keyword"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

// FormArticleRequest Web 表单提交 DTO (用于服务端渲染表单，支持结构体 Tag 自动校验)
type FormArticleRequest struct {
	Title   string `form:"title" binding:"required,min=3,max=80"`
	Author  string `form:"author" binding:"required,min=2,max=30"`
	Content string `form:"content" binding:"required,min=10"`
}
