package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

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
	app.Use(appMiddleware.BlockSensitiveFiles())
	app.Use(appMiddleware.SecurityHeaders())
	app.Use(middleware.CORS())
	store := session.NewCookieStore(cfg.App.SessionKey)
	app.Use(godeniter.Session(store, "starter_sess"))

	// 静态文件目录映射 (支持本地物理文件与上传文件分发)
	app.Static("/uploads", cfg.Upload.Dir)

	// 注册全局视图模板辅助函数 (通过 app.SetFuncMap 链式调用)
	app.SetFuncMap(template.FuncMap{
		"formatDate": func(t time.Time, layout string) string {
			if layout == "" {
				layout = "2006-01-02 15:04"
			}
			return t.Format(layout)
		},
	})

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

	// 7. 注册计划任务调度示例 (展示 godeniter 原生内置秒级与分级 Cron 引擎能力)
	_, _ = app.Schedule("heartbeat", "系统运行时心跳巡检", "*/15 * * * * *", func() error {
		return nil
	})
	_, _ = app.Schedule("cleanup", "临时缓存模拟清理", "@hourly", func() error {
		return nil
	})

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
		api.Get("/cron/jobs", func(c *godeniter.Context) {
			c.Success(app.Cron.Jobs())
		})
		triggerHandler := func(c *godeniter.Context) {
			jobID := c.Param("id")
			err := app.Cron.Trigger(jobID)
			if err != nil {
				c.Fail(400, err.Error())
				return
			}
			c.Success(godeniter.H{"message": "任务触发成功: " + jobID, "job_id": jobID})
		}
		api.Post("/cron/trigger/:id", triggerHandler)
		api.Get("/cron/trigger/:id", triggerHandler)
		api.Get("/storage/status", func(c *godeniter.Context) {
			c.Success(godeniter.H{
				"current_driver": "local",
				"upload_dir":     cfg.Upload.Dir,
				"supported_drivers": []godeniter.H{
					{"driver": "local", "name": "Local 本地磁盘驱动", "status": "active", "desc": "零外部依赖，极速本地读写"},
					{"driver": "webdav", "name": "WebDAV 云盘/NAS 驱动", "status": "supported", "desc": "原生 HTTP WebDAV 协议，支持群晖 NAS、Nextcloud、坚果云"},
					{"driver": "s3", "name": "Amazon S3 / MinIO 驱动", "status": "supported", "desc": "手写纯标准库 SigV4 认证，0-SDK 依赖，兼容 Cloudflare R2 / 阿里云 OSS"},
				},
			})
		})
	}

	// 11. 自定义 404 未命中页面
	app.NotFound = func(c *godeniter.Context) {
		c.Fail(40400, fmt.Sprintf("接口或页面不存在: [%s %s]", c.Method, c.Path))
	}

	return app
}

func main() {
	// 0. 若当前为 Windows GUI 模式且用户从现有终端 (CMD/PowerShell) 启动，自动挂载父级控制台输入输出
	tray.AttachConsole()

	// 1. 动态加载应用配置 (优先读取本地 config.json，不存在则自动生成；支持环境变量覆盖)
	cfg := config.LoadConfig()

	// 2. 初始化应用引擎
	app := setupApp(cfg)

	// 3. 命令行参数与运行模式判定:
	// - 显式子命令:
	//     tray: 强制桌面系统托盘模式 (Win32 原生自动隐藏黑框，macOS 顶部状态栏常驻)
	//     console / run: 显式前台控制台调试模式 (输出彩色 ASCII Banner 与实时请求日志)
	//     start / stop / restart / status: 后台守护进程管理器接管
	// - 无参数直接运行 (如 Windows / macOS 桌面双击):
	//     默认以系统托盘模式启动，Win32 原生自动隐藏黑框！
	cmd := ""
	if len(os.Args) > 1 {
		cmd = strings.ToLower(os.Args[1])
	}

	isTrayMode := false
	switch cmd {
	case "tray":
		isTrayMode = true
	case "console", "run":
		isTrayMode = false
	case "start", "stop", "restart", "status":
		isTrayMode = false
	default:
		// 如果文件名包含 tray (如 app_tray.exe)，或在具备图形界面的操作系统 (Windows / macOS) 下双击无参数运行，默认进入系统托盘模式
		if strings.Contains(strings.ToLower(os.Args[0]), "tray") || runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
			isTrayMode = true
		} else {
			isTrayMode = false
		}
	}

	if isTrayMode {
		webURL := "http://127.0.0.1" + cfg.App.Port
		fmt.Printf(">> [TRAY] 正在以桌面系统托盘模式启动 [%s]...\n", cfg.App.Name)
		fmt.Printf(">> [TRAY] 本地后台访问网址: %s\n", webURL)
		fmt.Println(">> [TRAY] 提示: 顶部菜单栏/系统托盘已常驻图标与管理菜单，随时按 Ctrl+C 或点击菜单项安全退出")

		// 1. 同步预检并监听网络端口，若遇端口冲突立即友好提示杀进程方案并退出
		ln, err := net.Listen("tcp", cfg.App.Port)
		if err != nil {
			portNum := strings.TrimPrefix(cfg.App.Port, ":")
			fmt.Println("\n" + strings.Repeat("=", 78))
			fmt.Printf("❌ 【端口被占用】服务启动失败: 端口 %s 已被占用！\n", cfg.App.Port)
			fmt.Printf(">> 详细底层错误: %v\n", err)
			fmt.Println(">> 💡 解决方案 (一键释放端口并杀掉冲突进程):")
			if runtime.GOOS == "windows" {
				fmt.Printf("   👉 Windows CMD 终端执行:        for /f \"tokens=5\" %%a in ('netstat -aon ^| findstr :%s') do taskkill /f /pid %%a\n", portNum)
				fmt.Printf("   👉 Windows PowerShell 终端执行:  Stop-Process -Id (Get-NetTCPConnection -LocalPort %s).OwningProcess -Force\n", portNum)
			} else {
				fmt.Printf("   👉 macOS/Linux 终端执行:         kill -9 $(lsof -ti :%s -sTCP:LISTEN)\n", portNum)
			}
			fmt.Printf("   👉 或者修改 config.json 中的 \"port\": \":8081\" 更换为其它未占用端口\n")
			fmt.Println(strings.Repeat("=", 78) + "\n")

			// 若在桌面双击运行 (无终端黑框模式)，弹出系统原生警告框
			tray.ShowAlert("端口被占用 - Godeniter", fmt.Sprintf(
				"端口 %s 已被占用，启动失败！\n\n如需一键释放端口，请在终端执行：\nkill -9 $(lsof -ti :%s -sTCP:LISTEN)\n\n或修改 config.json 更改端口。",
				cfg.App.Port, portNum,
			))
			os.Exit(1)
		}

		// 2. 将自身 PID 写入文件，方便与 daemon stop / start 统一管理
		pid := os.Getpid()
		_ = os.WriteFile(cfg.App.PIDFile, []byte(strconv.Itoa(pid)), 0644)
		defer func() {
			_ = os.Remove(cfg.App.PIDFile)
		}()

		// 3. 异步启动 Web 服务
		srv := &http.Server{
			Handler: app,
		}
		go func() {
			if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
				fmt.Printf(">> [ERROR] Web 服务运行异常: %v\n", err)
			}
		}()

		// 4. 自动唤起系统默认浏览器打开后台网址 (提供“双击有反应”的即时正向反馈)
		go func() {
			time.Sleep(150 * time.Millisecond)
			_ = tray.OpenURL(webURL)
		}()

		// 5. 主线程运行跨平台桌面托盘与状态栏菜单 (阻塞至用户退出)
		_ = tray.Run(tray.Options{
			Title:       "Godeniter",
			Tooltip:     fmt.Sprintf("%s (%s)", cfg.App.Name, cfg.App.Port),
			IconBytes:   appIcoBytes,
			URL:         webURL,
			AppDir:      tray.GetExecutableDir(),
			Version:     "v1.0.1",
			Port:        cfg.App.Port,
			HideConsole: true,
			OnExit: func() {
				fmt.Println("\n>> [TRAY] 收到退出指令，正在安全平滑关闭 Web 服务...")
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				_ = srv.Shutdown(ctx)
				_ = os.Remove(cfg.App.PIDFile)
				fmt.Println(">> [TRAY] 服务已成功安全停止。")
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
