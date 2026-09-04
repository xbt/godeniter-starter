package services

import (
	"github.com/xbt/godeniter/utils/str"
	"godeniter-starter/app/models"
	"strings"
	"sync"
	"time"
)

// ArticleService 文章业务服务
type ArticleService struct {
	mu       sync.RWMutex
	articles map[int]*models.Article
	nextID   int
}

// NewArticleService 实例化文章服务
func NewArticleService() *ArticleService {
	s := &ArticleService{
		articles: make(map[int]*models.Article),
		nextID:   1,
	}
	s.Create("欢迎使用 Godeniter Starter", "这是一个基于 Godeniter 2.0 框架初始化的官方脚手架项目，0 外部依赖，开箱即用。", "13800138000", "")
	s.Create("微服务与单文件打包指南", "支持直接打包为单个 Windows .exe，客户机双击即可运行。", "admin@godeniter.dev", "")
	return s
}

// ListPaginate 分页与模糊搜索
func (s *ArticleService) ListPaginate(keyword string, page, pageSize int) ([]*models.Article, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	matched := make([]*models.Article, 0)
	kw := strings.ToLower(strings.TrimSpace(keyword))

	for _, a := range s.articles {
		if kw == "" || strings.Contains(strings.ToLower(a.Title), kw) || strings.Contains(strings.ToLower(a.Content), kw) {
			clone := *a
			if strings.Contains(clone.Author, "@") {
				clone.AuthorMask = str.MaskEmail(clone.Author)
			} else if len(clone.Author) == 11 {
				clone.AuthorMask = str.MaskPhone(clone.Author)
			} else {
				clone.AuthorMask = clone.Author
			}
			matched = append(matched, &clone)
		}
	}

	total := len(matched)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	start := (page - 1) * pageSize
	if start >= total {
		return []*models.Article{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	return matched[start:end], total
}

// GetByID 根据 ID 查询并自增阅读量
func (s *ArticleService) GetByID(id int) (*models.Article, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.articles[id]
	if ok {
		a.Views++
	}
	return a, ok
}

// Create 创建文章
func (s *ArticleService) Create(title, content, author, coverURL string) *models.Article {
	s.mu.Lock()
	defer s.mu.Unlock()
	a := &models.Article{
		ID:        s.nextID,
		Title:     title,
		Content:   content,
		Author:    author,
		CoverURL:  coverURL,
		Views:     0,
		CreatedAt: time.Now(),
	}
	s.articles[s.nextID] = a
	s.nextID++
	return a
}

// Delete 删除文章
func (s *ArticleService) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.articles[id]; ok {
		delete(s.articles, id)
		return true
	}
	return false
}
