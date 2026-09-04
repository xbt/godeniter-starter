package services

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/xbt/godeniter/utils/str"
	"godeniter-starter/app/models"
)

// ArticleService 文章业务服务
type ArticleService struct {
	mu       sync.RWMutex
	articles map[int]*models.Article
	nextID   int
}

// NewArticleService 实例化文章服务并填充 8 条初始数据（展示每页 5 条的分页场景）
func NewArticleService() *ArticleService {
	s := &ArticleService{
		articles: make(map[int]*models.Article),
		nextID:   1,
	}

	// 预设 8 篇精品文章，覆盖丰富的检索与分页场景
	s.Create("欢迎使用 Godeniter 2.0 框架", "基于纯 Go 标准库与依赖注入容器打造的零依赖 Web 框架，兼具高性能与 CodeIgniter 优雅易用性。", "admin@godeniter.dev", "")
	s.Create("跨平台编译与 Windows 单文件打包指南", "支持通过 embed.FS 将前端静态页面完全打入独立二进制中，交付客户机双击直接运行，免装 Go 环境。", "13800138001", "")
	s.Create("0 依赖轻量依赖注入 (DI) 容器深度剖析", "支持声明式类型映射 (Map) 与接口抽象映射 (MapTo)，中间件与控制器按需声明依赖，框架自动反射装配。", "13900139002", "")
	s.Create("ActiveRecord 链式查询构造器实战", "类似 PHP CodeIgniter 3 体验，开箱支持 SQLite 与 MySQL，提供 Like 模糊查询、批量插入与一键分页。", "dev@godeniter.dev", "")
	s.Create("服务端 Session 会话管理与闪存消息 (Flash Data)", "基于 HMAC-SHA256 安全签名的防篡改 CookieStore，支持一次性提示消息，跨请求读取后自动销毁。", "13700137003", "")
	s.Create("Go 原生 html/template 模板渲染技巧", "采用清晰的经典 MVC 分层，结合模板继承与循环分支，打造高可读性的服务端渲染 (SSR) 系统。", "editor@godeniter.dev", "")
	s.Create("微服务环境下的动态配置与 Sidecar 机制", "通过 config.json 与环境变量实现分层配置覆盖，首次启动就近自生成模板，现场运维零门槛改端口。", "13600136004", "")
	s.Create("企业级数据安全：敏感信息脱敏与 XSS 过滤", "集成手机号、邮箱、身份证自动脱敏工具，配合 HTML 转义防御跨站脚本攻击，护航业务系统稳定上线。", "security@godeniter.dev", "")

	return s
}

// maskAuthor 处理作者联系方式脱敏
func maskAuthor(author string) string {
	if strings.Contains(author, "@") {
		return str.MaskEmail(author)
	} else if len(author) == 11 {
		return str.MaskPhone(author)
	}
	return author
}

// ListPaginate 分页与模糊搜索（按 ID 倒序排列）
func (s *ArticleService) ListPaginate(keyword string, page, pageSize int) ([]*models.Article, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	matched := make([]*models.Article, 0)
	kw := strings.ToLower(strings.TrimSpace(keyword))

	// 先将匹配的文章收集到切片中
	for _, a := range s.articles {
		if kw == "" || strings.Contains(strings.ToLower(a.Title), kw) || strings.Contains(strings.ToLower(a.Content), kw) {
			clone := *a
			clone.AuthorMask = maskAuthor(clone.Author)
			matched = append(matched, &clone)
		}
	}

	// 稳定排序：按 ID 倒序（最新发布的排在最前）
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].ID > matched[j].ID
	})

	total := len(matched)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 5
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

// GetByID 根据 ID 查询文章并自增阅读量
func (s *ArticleService) GetByID(id int) (*models.Article, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.articles[id]
	if ok {
		a.Views++
		clone := *a
		clone.AuthorMask = maskAuthor(clone.Author)
		return &clone, true
	}
	return nil, false
}

// FindByID 根据 ID 查询原始文章（不自增阅读量，供编辑使用）
func (s *ArticleService) FindByID(id int) (*models.Article, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.articles[id]
	if ok {
		clone := *a
		clone.AuthorMask = maskAuthor(clone.Author)
		return &clone, true
	}
	return nil, false
}

// Create 创建新文章
func (s *ArticleService) Create(title, content, author, coverURL string) *models.Article {
	s.mu.Lock()
	defer s.mu.Unlock()

	if author == "" {
		author = "管理员"
	}

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

// Update 更新文章
func (s *ArticleService) Update(id int, title, content, author string) (*models.Article, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	a, ok := s.articles[id]
	if !ok {
		return nil, false
	}

	if title != "" {
		a.Title = title
	}
	if content != "" {
		a.Content = content
	}
	if author != "" {
		a.Author = author
	}

	clone := *a
	clone.AuthorMask = maskAuthor(clone.Author)
	return &clone, true
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
