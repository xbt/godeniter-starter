package main

import (
	"embed"
	"fmt"
	"github.com/xbt/godeniter"
	"github.com/xbt/godeniter/middleware"
	"github.com/xbt/godeniter/session"
	"godeniter-starter/app/controllers"
	"godeniter-starter/app/services"
	"godeniter-starter/config"
	"html/template"
	"io/fs"
	"os"
)

//go:embed views/*
var viewsFS embed.FS

// setupApp 初始化 Godeniter 引擎并挂载所有中间件、模板和路由 (便于测试与复用)
func setupApp(cfg *config.Config) *godeniter.Engine {
	// 确保本地上传与数据目录存在
	_ = os.MkdirAll(cfg.Upload.Dir+"/images", 0755)
	_ = os.MkdirAll("./data", 0755)

	// 初始化 Godeniter 经典引擎 (内置 Logger、Recovery 与优雅停机)
	app := godeniter.Classic()

	// 挂载中间件
	app.Use(middleware.CORS())
	store := session.NewCookieStore(cfg.App.SessionKey)
	app.Use(godeniter.Session(store, "starter_sess"))

	// 静态文件目录映射
	app.Static("/uploads", cfg.Upload.Dir)

	// 加载内嵌 HTML 服务端模板
	subViews, err := fs.Sub(viewsFS, "views")
	if err == nil {
		app.SetHTMLTemplate(template.Must(template.ParseFS(subViews, "*.html")))
	}

	// 注册业务服务与控制器依赖注入 (DI)
	articleSvc := services.NewArticleService()
	app.Map(articleSvc)
	app.Map(cfg)

	homeCtrl := &controllers.HomeController{}
	authCtrl := &controllers.AuthController{}
	adminCtrl := &controllers.AdminController{}
	articleAPICtrl := &controllers.ArticleAPIController{}

	// 注册 Web 页面路由 (SSR)
	app.Get("/", homeCtrl.Index)
	app.Get("/article/:id", homeCtrl.Detail)
	app.Get("/login", authCtrl.LoginForm)
	app.Post("/login", authCtrl.LoginSubmit)
	app.Get("/logout", authCtrl.Logout)

	// 后台文章管理 CRUD 路由
	app.Get("/admin", adminCtrl.List)
	app.Get("/admin/articles", adminCtrl.List)
	app.Get("/admin/articles/create", adminCtrl.CreateForm)
	app.Post("/admin/articles/create", adminCtrl.CreateSubmit)
	app.Get("/admin/articles/edit/:id", adminCtrl.EditForm)
	app.Post("/admin/articles/edit/:id", adminCtrl.EditSubmit)
	app.Get("/admin/articles/delete/:id", adminCtrl.Delete)

	// 注册 RESTful API 路由分组 (/api/v1)
	api := app.Group("/api/v1")
	{
		api.Get("/articles", articleAPICtrl.List)
		api.Get("/articles/:id", articleAPICtrl.Detail)
		api.Post("/articles", articleAPICtrl.Create)
		api.Delete("/articles/:id", articleAPICtrl.Delete)
		api.Post("/upload", articleAPICtrl.Upload)
	}

	// 自定义 404 未命中页面
	app.NotFound = func(c *godeniter.Context) {
		c.Fail(40400, fmt.Sprintf("接口或页面不存在: [%s %s]", c.Method, c.Path))
	}

	return app
}

func main() {
	// 1. 动态加载应用配置 (优先读取本地 config.json，不存在则自动生成；支持环境变量覆盖)
	cfg := config.LoadConfig()

	// 2. 初始化应用引擎
	app := setupApp(cfg)

	// 3. 启动 HTTP 服务 (端口由 config.json 中的 app.port 动态决定，并支持 Ctrl+C 平滑退出)
	_ = app.Run(cfg.App.Port)
}

