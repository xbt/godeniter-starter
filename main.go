package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/xbt/godeniter"
	"github.com/xbt/godeniter/daemon"
	"github.com/xbt/godeniter/middleware"
	"github.com/xbt/godeniter/session"
	"github.com/xbt/godeniter/tray"
	"godeniter-starter/app/controllers"
	appMiddleware "godeniter-starter/app/middleware"
	"godeniter-starter/app/services"
	"godeniter-starter/config"
)

//go:embed views/*
var viewsFS embed.FS

//go:embed app.ico
var appIcoBytes []byte

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

	// 加载内嵌 HTML 服务端模板 (原生内置 <!--{{ ... }}--> 无侵入注释模板支持与单文件打包)
	subViews, err := fs.Sub(viewsFS, "views")
	if err == nil {
		app.LoadHTMLFS(subViews, "*.html")
	}

	// 浏览器 Favicon 图标路由 (内嵌单文件打包，返回 0 依赖自定义 ICO)
	app.Get("/favicon.ico", func(c *godeniter.Context) {
		c.Data(http.StatusOK, "image/x-icon", appIcoBytes)
	})


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

	// 3. 判断是否启用桌面系统托盘模式 (支持命令行参数 `tray` 或配置文件 `"tray": true`)
	isTrayMode := cfg.App.Tray
	if len(os.Args) > 1 && strings.ToLower(os.Args[1]) == "tray" {
		isTrayMode = true
	}

	if isTrayMode {
		webURL := "http://127.0.0.1" + cfg.App.Port
		fmt.Printf(">> [TRAY] 正在以桌面系统托盘模式启动 [%s]...\n", cfg.App.Name)
		fmt.Printf(">> [TRAY] 本地后台访问网址: %s\n", webURL)
		fmt.Println(">> [TRAY] 提示: 点击或右键系统托盘/状态栏图标可进行管理")

		// 异步协程启动 Web 服务
		go func() {
			if err := app.Run(cfg.App.Port); err != nil && err != http.ErrServerClosed {
				fmt.Printf(">> [ERROR] Web 服务运行异常: %v\n", err)
			}
		}()

		// 主线程运行跨平台桌面托盘与状态栏菜单 (阻塞至用户退出)
		_ = tray.Run(tray.Options{
			Title:     cfg.App.Name,
			Tooltip:   fmt.Sprintf("%s (%s)", cfg.App.Name, cfg.App.Port),
			IconBytes: appIcoBytes,
			URL:       webURL,
			AppDir:    tray.GetExecutableDir(),
			Version:   "v1.0.0",
			Port:      cfg.App.Port,
			OnExit: func() {
				fmt.Println(">> [TRAY] 托盘已退出，正在安全关闭服务...")
			},
		})
		return
	}

	// 4. 由守护进程管理器统一接管服务启动与生命周期指令 (支持 start/stop/restart/status 与后台静默运行)
	_ = daemon.Run(app, cfg.App.Port, daemon.Config{
		Daemon:  cfg.App.Daemon,
		PIDFile: cfg.App.PIDFile,
		LogFile: cfg.App.LogFile,
	})

}
