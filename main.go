package main

import (
	"embed"
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
	// 1. 加载应用配置
	cfg := config.LoadConfig()

	// 2. 确保本地上传目录存在
	_ = os.MkdirAll(cfg.UploadDir+"/images", 0755)

	// 3. 初始化 Godeniter 经典引擎
	app := godeniter.Classic()

	// 4. 挂载中间件
	app.Use(middleware.CORS())
	store := session.NewCookieStore(cfg.SessionKey)
	app.Use(godeniter.Session(store, "starter_sess"))

	// 5. 静态文件目录映射
	app.Static("/uploads", cfg.UploadDir)

	// 6. 加载内嵌 HTML 服务端模板
	subViews, err := fs.Sub(viewsFS, "views")
	if err == nil {
		app.SetHTMLTemplate(template.Must(template.ParseFS(subViews, "*.html")))
	}

	// 7. 注册业务服务与控制器依赖
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

	// 10. 启动 HTTP 服务
	_ = app.Run(cfg.Port)
}
