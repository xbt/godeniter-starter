package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"

	"github.com/xbt/godeniter"
	"github.com/xbt/godeniter/middleware"
	"github.com/xbt/godeniter/session"
	"godeniter-starter/app/controllers"
	appMiddleware "godeniter-starter/app/middleware"
	"godeniter-starter/app/services"
	"godeniter-starter/config"
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

	// 挂载全局中间件流水线 (展示洋葱圈模型、耗时追踪与安全防护标头)
	app.Use(appMiddleware.ResponseTimer())
	app.Use(appMiddleware.SecurityHeaders())
	app.Use(middleware.CORS())
	store := session.NewCookieStore(cfg.App.SessionKey)
	app.Use(godeniter.Session(store, "starter_sess"))

	// 静态文件目录映射 (支持本地物理文件与上传文件分发)
	app.Static("/uploads", cfg.Upload.Dir)

	// 加载内嵌 HTML 服务端模板 (支持单文件二进制一键打包)
	subViews, err := fs.Sub(viewsFS, "views")
	if err == nil {
		app.SetHTMLTemplate(template.Must(template.ParseFS(subViews, "*.html")))
	}

	// 注册业务服务与控制器依赖注入 (DI 容器)
	articleSvc := services.NewArticleService()
	app.Map(articleSvc)
	app.Map(cfg)

	homeCtrl := &controllers.HomeController{}
	authCtrl := &controllers.AuthController{}
	adminCtrl := &controllers.AdminController{}
	articleAPICtrl := &controllers.ArticleAPIController{}

	// 8. 注册 Web 页面路由 (服务端渲染 SSR)
	app.Get("/", homeCtrl.Index)
	app.Get("/features", homeCtrl.Features)
	app.Get("/article/:id", homeCtrl.Detail)
	app.Get("/demo/panic", homeCtrl.PanicDemo)
	app.Get("/login", authCtrl.LoginForm)
	app.Post("/login", authCtrl.LoginSubmit)
	app.Get("/logout", authCtrl.Logout)

	// 9. 注册后台文章管理 CRUD 路由分组 (使用中间件路由守卫 AuthRequired 统一权限校验)
	admin := app.Group("/admin", appMiddleware.AuthRequired())
	{
		admin.Get("", adminCtrl.List)
		admin.Get("/articles", adminCtrl.List)
		admin.Get("/articles/create", adminCtrl.CreateForm)
		admin.Post("/articles/create", adminCtrl.CreateSubmit)
		admin.Get("/articles/edit/:id", adminCtrl.EditForm)
		admin.Post("/articles/edit/:id", adminCtrl.EditSubmit)
		admin.Get("/articles/delete/:id", adminCtrl.Delete)
	}

	// 10. 注册 RESTful API 路由分组 (/api/v1)
	api := app.Group("/api/v1")
	{
		api.Get("/articles", articleAPICtrl.List)
		api.Get("/articles/:id", articleAPICtrl.Detail)
		api.Post("/articles", articleAPICtrl.Create)
		api.Delete("/articles/:id", articleAPICtrl.Delete)
		api.Post("/upload", articleAPICtrl.Upload)
	}

	// 11. 自定义 404 未命中页面
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
