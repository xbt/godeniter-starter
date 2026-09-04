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

func main() {
	// 1. 动态加载应用配置 (优先读取本地 config.json，不存在则自动生成；支持环境变量覆盖)
	cfg := config.LoadConfig()

	// 2. 确保本地上传与数据目录存在
	_ = os.MkdirAll(cfg.Upload.Dir+"/images", 0755)
	_ = os.MkdirAll("./data", 0755)

	// 3. 初始化 Godeniter 经典引擎 (内置 Logger、Recovery 与优雅停机)
	app := godeniter.Classic()

	// 4. 挂载中间件
	app.Use(middleware.CORS())
	store := session.NewCookieStore(cfg.App.SessionKey)
	app.Use(godeniter.Session(store, "starter_sess"))

	// 5. 静态文件目录映射
	app.Static("/uploads", cfg.Upload.Dir)

	// 6. 加载内嵌 HTML 服务端模板
	subViews, err := fs.Sub(viewsFS, "views")
	if err == nil {
		app.SetHTMLTemplate(template.Must(template.ParseFS(subViews, "*.html")))
	}

	// 7. 注册业务服务与控制器依赖注入 (DI)
	articleSvc := services.NewArticleService()
	app.Map(articleSvc)
	app.Map(cfg)

	homeCtrl := &controllers.HomeController{}
	authCtrl := &controllers.AuthController{}
	articleAPICtrl := &controllers.ArticleAPIController{}

	// 8. 注册 Web 页面路由 (SSR)
	app.Get("/", homeCtrl.Index)
	app.Get("/login", authCtrl.LoginForm)
	app.Post("/login", authCtrl.LoginSubmit)
	app.Get("/logout", authCtrl.Logout)

	// 9. 注册 RESTful API 路由分组 (/api/v1)
	api := app.Group("/api/v1")
	{
		api.Get("/articles", articleAPICtrl.List)
		api.Get("/articles/:id", articleAPICtrl.Detail)
		api.Post("/articles", articleAPICtrl.Create)
		api.Delete("/articles/:id", articleAPICtrl.Delete)
		api.Post("/upload", articleAPICtrl.Upload)
	}

	// 10. 自定义 404 未命中页面
	app.NotFound = func(c *godeniter.Context) {
		c.Fail(40400, fmt.Sprintf("接口或页面不存在: [%s %s]", c.Method, c.Path))
	}

	// 11. 启动 HTTP 服务 (端口由 config.json 中的 app.port 动态决定，并支持 Ctrl+C 平滑退出)
	_ = app.Run(cfg.App.Port)
}
