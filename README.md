# Godeniter Starter 脚手架项目

**Godeniter Starter** 是官方推荐的标准业务工程骨架。基于 **Godeniter 2.0 框架** 构建，拥有 **0 外部依赖**、**内置秒级/分级 Cron 调度**、**纯 Go 多存储驱动**、**依赖注入**、**跨平台桌面托盘与守护进程**、**主动安全防御**、**单文件独立打包** 与 **极速启动** 特性。

---

## 🚀 快速上手 (Quick Start)

### 1. 本地前台开发模式 (随时 Ctrl+C 停止)
```bash
# 启动 Web 服务 (默认前台运行，终端实时滚动日志与彩色 ASCII Banner)
go run main.go
```
启动后终端将自动打印本机与局域网访问地址，在浏览器中打开：
* 🌐 **前台博客首页**：`http://127.0.0.1:8080` (支持分页、搜索、作者信息脱敏)
* ⚡ **框架特性全景体验中心**：`http://127.0.0.1:8080/features` (14 项核心特性在线交互探针)
* 🔐 **后台管理中心**：`http://127.0.0.1:8080/admin` (默认账号：`admin` / 密码：`123456`)

### 2. Linux / 服务器后台守护模式 (关闭终端/断开 SSH 持续运行)
无需额外编译，源码与打包二进制均原生支持标准服务生命周期管理：
```bash
# 1. 后台静默启动 (自动脱离终端，记录 PID 至 app.pid，重定向日志至 app.log，立即返回命令行)
go run main.go start      # 源码方式
./dist/app start          # 二进制方式

# 2. 查看运行状态 (检测存活性与 PID)
go run main.go status     # 或 ./dist/app status

# 3. 动态查看后台日志 (按 Ctrl+C 仅退出查看，不影响程序运行)
tail -f app.log

# 4. 平滑安全停止服务 (优雅停机并清理 PID 文件)
go run main.go stop       # 或 ./dist/app stop

# 5. 一键平滑重启服务
go run main.go restart    # 或 ./dist/app restart
```

### 3. 跨平台桌面系统托盘 / 状态栏客户端模式 (macOS / Windows 常驻)
若您希望将本 Web 服务作为独立的本地客户端或桌面托盘常驻运行：
* **进入托盘模式**：
  在终端中传入 `tray` 参数启动：
  ```bash
  go run main.go tray      # 源码方式
  ./dist/app tray          # 二进制方式 (Windows 自动隐藏控制台黑框)
  ```
  * **Windows**：Win32 原生自动隐藏控制台黑框，右下角任务栏托盘图标优雅常驻，右键弹出管理菜单，双击图标直接打开浏览器后台；
  * **macOS**：顶部状态栏常驻图标，点击展开管理菜单。
* **常驻原生菜单**：
  * 🌐 **打开管理后台**：调起系统默认浏览器访问 Web 页面；
  * 📁 **打开应用目录**：调起系统文件管理器 (Explorer / Finder) 定位应用目录，方便查找 `config.json` 与数据文件；
  * ℹ️ **关于系统**：弹窗展示应用名称、版本、运行端口与进程 PID 等信息；
  * ⏹️ **退出程序**：平滑优雅关闭 Web 服务并安全退出，托盘图标无感清理。

---

## ⏰ 内置秒级与 Linux Crontab 计划任务 (Cron Scheduler)

脚手架完整展示了 `godeniter` 底层纯 Go 自研的 **`cron` 计划任务调度引擎**：

* **100% 纯 Go 标准库打造**：零外部依赖，对标 Linux Crontab 规范，内置 1 秒高精度时钟轮询；
* **秒级与经典分级全自适应**：
  * **6 段秒级语法**：`*/15 * * * * *`（每 15 秒执行一次）；
  * **5 段经典分级语法**：`0 3 * * *`（每天凌晨 3:00 执行，自动补齐秒为 0）；
  * **常用语义宏**：`@daily`, `@hourly`, `@weekly`, `@every 30s`；
* **Panic 容错隔离**：任何任务执行异常自动捕获并记录彩色错误堆栈，绝对不拖垮 Web 主进程；
* **手动即时触发 (`Trigger`)**：支持在独立协程中随时手动触发一次执行，不影响后续自然计划时间；
* **生命周期与优雅停机**：应用退出时自动执行 `Drain`，等待正在运行的任务安全收尾。

### Starter 中的使用示例 (`main.go`)：
```go
// 1. 注册每 15 秒一次的系统心跳任务
app.Schedule("heartbeat", "系统运行时心跳巡检", "*/15 * * * * *", func() error {
    return nil
})

// 2. 注册每小时一次的临时缓存清理任务
app.Schedule("cleanup", "临时缓存模拟清理", "@hourly", func() error {
    return nil
})
```

### 配套 RESTful API 与在线交互：
* `GET /api/v1/cron/jobs`：返回所有活跃任务的快照列表（运行次数、最后耗时、运行状态、下次时间）；
* `POST /api/v1/cron/trigger/:id`（或 `GET`）：立即触发指定任务，即时执行并刷新状态。

---

## 📴 离线与受限网络开发说明 (Air-gapped & Offline Ready)

本 Starter 工程与 `godeniter` 核心框架均采用 **100% 纯 Go 标准库（0 外部第三方依赖）** 设计：

* **源码包解压即用**：将 `godeniter` 和 `godeniter-starter` 解压到同级目录：
  ```text
  my_workspace/
  ├── godeniter/           <-- 框架源码
  └── godeniter-starter/   <-- 脚手架工程
  ```
* **本地依赖自动生效**：`go.mod` 中已内置 `replace github.com/xbt/godeniter => ../godeniter`，编译直接读取本地源码，**全程 0 网络请求**。
* **终极离线交付**：直接运行 `./build.sh` 生成 `dist/app.exe` 交付给客户，客户机无需安装 Go 环境，双击即可运行。

---

## ⚙️ 动态配置管理 (`config.json`)

项目基于 **纯 Go 标准库（0 外部依赖）** 实现三层动态配置装配：

* **开箱即用与 SQLite 专属随机库名**：首次运行若未检测到配置文件，程序将**自动在当前目录生成一份格式化的 `config.json`**；若使用 SQLite 驱动，系统会自动生成带 8 位随机专属后缀的库名（如 `./data/app_k8x9m2p1.db`，或识别 `{random}` 占位符）并就地写回固化，**既杜绝默认固定库名被外部恶意探测猜测，又确保服务重启后数据持久安全**！
* **修改端口、守护模式与数据库**：直接使用任意文本编辑器（如记事本）打开 `config.json` 即可调整：
  ```json
  {
    "app": {
      "name": "Godeniter Starter Application",
      "port": ":8080",
      "env": "development",
      "session_key": "godeniter-starter-secret-salt-2026",
      "daemon": false,
      "pid_file": "./app.pid",
      "log_file": "./app.log"
    },
    "database": {
      "driver": "sqlite",
      "dsn": "./data/app_{random}.db",
      "max_open_conns": 1,
      "max_idle_conns": 1,
      "conn_max_lifetime": 300
    },
    "upload": {
      "dir": "./uploads",
      "max_size_mb": 5,
      "allowed_exts": [".jpg", ".png", ".jpeg", ".webp"]
    }
  }
  ```
* **连接 MySQL 生产数据库**：若需切换至 MySQL，仅需在 `config.json` 中配置：
  ```json
  "database": {
    "driver": "mysql",
    "dsn": "root:password@tcp(127.0.0.1:3306)/your_db?charset=utf8mb4&parseTime=True&loc=Local",
    "max_open_conns": 50,
    "max_idle_conns": 10,
    "conn_max_lifetime": 3600
  }
  ```
  *(详见 [《数据库与 ActiveRecord 开发手册 (docs/database.md)》](../godeniter/docs/database.md) 中的 MySQL 实战完整示例)*
* **云原生 / 容器化覆盖**：支持通过系统环境变量（如 `PORT=:9000`、`DATABASE_DSN="..."`）动态覆盖对应字段。

---

## 🌟 核心演示功能 (Full Features Demos)

脚手架内置了一套完整的现代化轻量博客/内容管理系统，采用清新典雅的**浅色调 UI（Light Theme）**，全面覆盖 Godeniter 框架在实际企业级业务中用到的核心能力：

1. **框架特性全景体验中心 (`/features`)**：
   - 顶部导航设有“⚡ 框架特性”专区，全景呈现 **14 大全能杀手级特性**；
   - 包含分类筛选切换（全部、计划任务、统一存储、桌面与守护、安全防御、核心架构）；
   - 配套**在线实时交互探针**：可现场探测响应标头、拉取 Cron 任务快照、手动触发心跳任务、探测已挂载存储引擎、模拟黑客恶意探测敏感文件并验证拦截。
2. **纯 Go 内置秒级 / 分级 Cron 计划任务调度**：
   - 支持 6 段秒级（`*/15 * * * * *`）与经典分级；
   - 任务 Panic 容灾隔离；提供 REST API 端点查询快照与手动触发。
3. **纯 Go 统一多存储驱动抽象 (`storage`)**：
   - 抽象 `storage.Driver` 接口，默认本地磁盘存储驱动；
   - 代码内预留纯标准库手写 AWS SigV4 认证的 S3 / MinIO / Cloudflare R2 / 阿里云 OSS 及 WebDAV 云盘驱动接入示例，**0 引入臃肿的三方 SDK**。
4. **文件上传全链路演示 (Upload & Static Serv)**：
   - **封面图片上传**：发布/编辑文章表单支持本地选择封面图片上传，前端即时图片预览；
   - **服务端安全校验**：使用框架内置 `c.SaveUploadedFileWithOptions`，实施 5MB 上限限制、`.jpg/.png/.jpeg/.webp` 格式白名单校验，并自动重命名保存至 `./uploads/images/`；
   - **多端封面展示**：首页列表图文自适应卡片、详情页高清头图以及后台表格封面缩略图；同时提供独立 RESTful 上传接口 `/api/v1/upload`。
5. **自定义业务中间件流水线 (Custom Middleware)**：
   - **路由守卫中间件 (`AuthRequired`)**：保护 `/admin` 路由组，未登录拦截并携带 Flash 提示重定向；
   - **响应耗时中间件 (`ResponseTimer`)**：自动记录每个请求的处理耗时，向响应头注入 `X-Response-Time` 与 W3C `Server-Timing`；
   - **基础安全头中间件 (`SecurityHeaders`)**：自动注入 `X-Content-Type-Options: nosniff`、`X-Frame-Options: SAMEORIGIN`、`X-XSS-Protection` 等安全标头；
   - **敏感文件拦截中间件 (`BlockSensitiveFiles`)**：自动阻断针对 `.db`、`config.json`、`.env` 的恶意扫描，直接返回 403 Forbidden。
6. **Martini 风格强类型反射依赖注入 (Dependency Injection)**：
   - 服务（如 `ArticleService`、`*config.Config`）注册至全局容器 `app.Map(...)`；
   - 控制器 Handler 签名自由，按需直接声明依赖参数，框架在运行时自动按类型反射注入，控制器与底层彻底解耦。
7. **零依赖参数绑定与结构体 Tag 校验 (Binding & Validation)**：
   - 基于纯 Go 标准库与反射实现的轻量验证器，无需第三方库；
   - 表单提交声明 `binding:"required,min=3,max=80"`，校验失败自动在 Web 页面友好回显红字提示。
8. **企业级数据安全与脱敏实战 (Utils & Security)**：
   - **正文防 XSS 过滤**：发布与编辑正文自动执行 `str.XSSFilter` 过滤恶意脚本；
   - **敏感联系方式脱敏**：作者手机号/邮箱智能掩码脱敏（`str.MaskPhone` / `str.MaskEmail`）；
   - **内容摘要截断**：卡片摘要使用 `str.Truncate` 安全截断。
9. **前台 8 篇精品文章预置与 5 条/页分页**：
   - 倒序呈现 8 篇实战文章，直观展示 `第 1 / 2 页` 分页场景，支持上一页/下一页无缝翻页；
   - 顶部搜索框支持标题与内容模糊检索，并与分页深度联动。
10. **Panic 优雅容灾与不宕机 (`Recovery`)**：
    - 提供 `/demo/panic` 测试端点，点击可验证 `Recovery()` 中间件优雅捕获运行时异常并输出彩色堆栈，服务端返回 500，**进程永不宕机**。
11. **无侵入 HTML 注释模板语法实战 (`<!--{{ ... }}-->`)**：
    - 视图使用框架原生 `app.LoadHTMLFS(subViews, "*.html")` 加载；
    - 在 `views/detail.html` 中实战应用 `<!--{{ if .CurrentUser }}-->` 与 `<!--{{ .Article.ID }}-->`。本地直接用浏览器双击打开静态 HTML 原型时不乱码、不破坏按钮排版，Go 服务端渲染时无缝编译为高效 AST。
12. **开箱即用的自动化端到端测试 (`main_test.go`)**：
    - 包含 10 大系统测试用例，覆盖中间件安全头、Server-Timing、Favicon 输出、Panic 恢复、特性全景页、分页检索、详情渲染（含无侵入模板断言）、文件上传、后台完整 CRUD、敏感文件 403 阻断探测、存储驱动以及 Cron 调度与手动触发接口，`go test -v .` 秒级全绿（~0.3s）。

---

## 📂 全工程文件与目录职责全景清单 (File Directory Index)

本脚手架严格遵守清晰分层，以下为工程内全部源码文件与目录的职责清单：

### 1. 核心启动与配置层
| 文件路径 | 职责说明 |
| :--- | :--- |
| [`main.go`](./main.go) | **应用启动总入口**：加载配置、初始化 `godeniter` 引擎、挂载全局中间件、注册 HTML 模板与 Favicon、依赖注入装配、注册 Web/API 路由、挂载 Cron 调度与生命周期管理。 |
| [`main_test.go`](./main_test.go) | **自动化端到端测试**：覆盖中间件安全头、Server-Timing、Favicon 输出、Panic 恢复、文章分页检索、API 文件上传、后台登录与 CRUD、敏感文件 403 阻断探测、存储驱动及 Cron 调度即时触发。 |
| [`config.json`](./config.json) | **外部动态配置文件**：声明端口、Session 密钥、SQLite 驱动与 DSN 路径、文件上传限制。支持冷重启动态调整。 |
| [`config/app.go`](./config/app.go) | **配置装配与随机数据库固化**：纯标准库解析 JSON 配置；识别 `{random}` 占位符自动生成专属 SQLite 随机库名并自动写回固化持久化；支持环境变量覆盖。 |
| [`config/app_test.go`](./config/app_test.go) | **配置单元测试**：验证默认随机 DSN 生成机制、`{random}` 占位符自动固化写回及服务重启数据防丢失。 |

### 2. 控制器层 (`app/controllers/`)
| 控制器文件 | 职责说明 |
| :--- | :--- |
| [`app/controllers/home.go`](./app/controllers/home.go) | **前台页面控制器**：处理网站首页展示 (`/`)、框架特性全景体验中心 (`/features`)、文章详情阅读页 (`/article/:id`) 及故意触发 Panic 验证服务自愈的测试路由 (`/demo/panic`)。 |
| [`app/controllers/admin.go`](./app/controllers/admin.go) | **后台管理控制器**：受 `AuthRequired` 路由守卫保护。处理文章管理列表、带封面图片上传的文章新增 (`create`)、编辑更新 (`edit`) 与删除 (`delete`)。 |
| [`app/controllers/auth.go`](./app/controllers/auth.go) | **会话鉴权控制器**：提供管理员登录页面渲染 (`GET /login`)、账号密码校验与 Session 建立 (`POST /login`) 及安全退出 (`GET /logout`)。 |
| [`app/controllers/api_article.go`](./app/controllers/api_article.go) | **RESTful API 控制器**：提供标准 JSON 接口：文章列表查询分页 (`GET /api/v1/articles`)、文章详情 (`GET /api/v1/articles/:id`)、创建文章 (`POST`)、删除文章 (`DELETE`) 与图片上传 (`POST /api/v1/upload`)。 |

### 3. 中间件层 (`app/middleware/`)
| 中间件文件 | 职责说明 |
| :--- | :--- |
| [`app/middleware/auth.go`](./app/middleware/auth.go) | **路由守卫拦截器**：检查 Session 是否存在登录态，未登录时通过 Flash 提示并 302 重定向至登录页。 |
| [`app/middleware/security.go`](./app/middleware/security.go) | **安全标头与文件保护**：接入 `godeniter/middleware.Security()` 注入行业基准安全头；接入 `BlockSensitive()` 拦截对 `.db`、`config.json`、`.env` 的探测 (403)。 |
| [`app/middleware/timer.go`](./app/middleware/timer.go) | **耗时监控中间件**：接入 `godeniter/middleware.ServerTiming()`，为响应注入 `X-Response-Time` 与 W3C 标准 `Server-Timing`。 |
| [`app/middleware/keyauth.go`](./app/middleware/keyauth.go) | **API Key 鉴权中间件**：接入 `godeniter/middleware.KeyAuth()`，展示如何从 Bearer Token、X-API-Key 或 URL 参数提取密钥并执行鉴权。 |

### 4. 业务服务与数据模型层 (`app/services/` & `app/models/`)
| 文件路径 | 职责说明 |
| :--- | :--- |
| [`app/services/article.go`](./app/services/article.go) | **文章业务服务**：内存线程安全并发管理；预装 8 篇打样文章；提供分页、模糊检索、自增阅读量、XSS 安全过滤及 CRUD 操作。 |
| [`app/services/storage.go`](./app/services/storage.go) | **多存储驱动演示工厂**：展示如何使用 `godeniter/storage.Driver`；默认使用本地磁盘存储驱动，并在代码中演示如何一行切换到 WebDAV 或 S3 / Cloudflare R2 / 阿里云 OSS 云对象存储。 |
| [`app/models/article.go`](./app/models/article.go) | **数据实体与校验 DTO**：定义 `Article` 核心模型；定义 API 请求绑定 `CreateArticleRequest`、分页查询 `ArticleQueryRequest` 及带结构体 Tag 自动校验的表单 `FormArticleRequest`。 |

### 5. 视图模板与静态资源 (`views/` & `uploads/`)
| 目录 / 资源 | 职责说明 |
| :--- | :--- |
| [`views/index.html`](./views/index.html) | **前台首页模板**：浅色优雅卡片设计，包含顶部导航、文章模糊搜索框、封面展示与分页条。 |
| [`views/detail.html`](./views/detail.html) | **文章详情模板**：展示大图封面、发布时间、脱敏作者信息、阅读量计数与正文排版。 |
| [`views/features.html`](./views/features.html) | **特性全景体验中心**：全景呈现 14 项框架核心能力，包含在线任务探针、响应头探测、敏感扫描防御校验与存储驱动探测。 |
| [`views/login.html`](./views/login.html) | **后台登录面板**：支持回车提交、Session Flash 错误提示与优雅卡片居中布局。 |
| [`views/admin.html`](./views/admin.html) | **后台管理列表页**：展示文章封面微缩图、标题、作者、阅读数与编辑/删除操作栏。 |
| [`views/article_form.html`](./views/article_form.html) | **文章发布/编辑表单**：支持选择本地封面即时预览，配合结构体 Tag 字段错误高亮提示。 |
| [`app.ico`](./app.ico) | **应用专属图标**：纯标准库 Windows 资源编译器与浏览器 Favicon 统一图标源文件。 |
| [`uploads/images/`](./uploads/images/) | **运行时上传目录**：内置 `sample_cover.svg` 缺省封面矢量图，作为上传文件的分发目录。 |
| [`build.sh` / `build.bat`](./build.sh) | **一键单文件打包脚本**：自动化探测图标编译 `.syso` 并产出 Windows 64位与 macOS/Linux 独立全能二进制。 |

---

## 🌟 底层能力活字典与特性速查 (Showcase & Best Practices)

作为官方标准起步工程，本 Starter 是探索 `godeniter` 核心引擎所有特性的最佳窗口：

| 特性分类 | 演示模块 / 源码位置 | 特性说明 |
| :--- | :--- | :--- |
| **纯 Go 内置 Cron 计划任务** | `main.go` / `/api/v1/cron/jobs` | 纯 Go 自研秒级/分级 Crontab 调度内核，支持 6 段式与经典 5 段式，支持 Panic 隔离与手动即时触发 (`Trigger`)。 |
| **纯 Go 多存储驱动** | `app/services/storage.go` | 基于纯标准库（0 外部 SDK）抽象的 `storage.Driver`。开箱支持本地磁盘，附带详尽示例演示如何切换至 **WebDAV**（坚果云/Nextcloud）或 **AWS SigV4 规范的 S3 / Cloudflare R2 / 阿里云 OSS / MinIO**。 |
| **特性全景体验中心** | `views/features.html` | 现代化多卡片全景呈现，支持分类筛选，配套任务拉取、即时触发、存储驱动探测、黑客探测拦截等在线交互探针。 |
| **Web 安全防护标头** | `app/middleware/security.go` | 原生接入 `middleware.Security()`，全自动注入 `nosniff`、`SAMEORIGIN`、`XSS-Protection`、`Referrer-Policy` 等行业基线安全头。 |
| **敏感资源探测防御** | `app/middleware/security.go` | 原生接入 `middleware.BlockSensitive()`，自动探测拦截黑客针对 `.db`、`.sqlite`、`config.json`、`.env`、`.git` 的未授权嗅探，直接 403 阻断并记入审计日志。 |
| **Server-Timing 追踪** | `app/middleware/timer.go` | 原生接入 `middleware.ServerTiming()`，为所有 HTTP 响应注入 `X-Response-Time` 和 W3C 标准 `Server-Timing` 标头，方便前端分析后端耗时。 |
| **API Key 鉴权中间件** | `app/middleware/keyauth.go` | 原生接入 `middleware.KeyAuth()`，多通道自动提取 `Bearer Token`、`X-API-Key` 或 Query 参数并执行高效鉴权。 |
| **依赖注入 (DI 容器)** | `main.go` / 控制器 | Martini 风格纯 Go 反射注入容器，Handler 参数自由声明，框架运行时自动按类型反射注入，实现控制器与业务彻底解耦。 |
| **无侵入 HTML 注释模板** | `views/detail.html` | 独创 `<!--{{ ... }}-->` 语法，本地双击 HTML 原型不乱码，线上 Go 原生 AST 极速渲染。 |
| **SQLite 随机库名持久化** | `config/app.go` | 自动识别 `{random}` 占位符或默认库名，生成专属不可预测随机名并自动写回 `config.json` 固化，防止针对固定库名的恶意猜测。 |
| **单文件 Favicon 内嵌** | `main.go` | 纯标准库 0 外部工具内嵌 `app.ico` 并挂载 `/favicon.ico` 路由，消灭浏览器控制台 404 图标告警。 |
| **跨平台控制台自动隐身** | `main.go` | Windows 下桌面双击无黑框闪烁，自动最小化至托盘并提供右键控制菜单；CLI 命令行带参仍可正常输出彩色日志。 |

---

## 🚀 进阶与实战延伸案例：godetype

如果您想了解如何基于本 Starter 脚手架延伸构建高复杂度的真实生产级业务系统，请参考官方旗舰案例：
* **[godetype](https://github.com/xbt/godetype)**：基于 `godeniter-starter` 骨架孵化衍生的 Typecho 极客博客系统（100% 纯 Go 0-CGO 驱动、内置 Crontab 可视化管理看板与每日物理镜像备份、深度兼容 Typecho 官方数据表、Fifty shades of Tux 经典调色板）。

---

## 📦 一键编译单文件交付 (统一全能单二进制)

```bash
# 生成 Windows 64位统一全能单文件可执行程序 (dist/app.exe) 及 macOS/Linux 本地二进制
./build.sh     # macOS / Linux
build.bat      # Windows
```

生成的单文件无需安装任何环境，直接拷贝给客户，**统一且功能完备**：
* `dist/app_tray.exe`：**Windows 纯静默桌面托盘客户端**
  - 基于 Windows GUI 子系统构建，双击直接常驻屏幕右下角任务栏托盘，**100% 彻底无黑框、无闪烁**，右键唤出完整管理菜单；
* `dist/app.exe`：**Windows 统一全能二进制**
  - 在 CMD / PowerShell 中：支持 `run/console` 调试与 `start/status/stop/restart` 守护进程管理；
  - 桌面双击直接运行：自动隐藏控制台黑框进入右下角托盘；
* `dist/app`：**macOS / Linux 统一全能二进制**（支持 CLI 运维与顶部状态栏托盘）。

---

## 🎨 自定义应用与网页图标 (纯标准库 0 依赖 app.ico)

工程原生内置了 **纯 Go 标准库 Windows 资源段编译器** 与 **浏览器 Favicon** 双端图标一体化支持：

* **更换专属图标**：只需将您的定制图标命名为 `app.ico` 放置在项目根目录下（工程已预置精美默认图标）。
* **Windows .exe 桌面图标动态缝合**：执行 `./build.sh` 或 `build.bat` 时，脚本自动动态检测 `app.ico`，通过框架内置纯标准库将其自动转译为 `resource_windows_amd64.syso` 并缝合进 `dist/app.exe`。在 Windows 桌面和资源管理器中呈现专属应用图标，**全程 0 外部依赖、0 外部工具链、断网无网直接可用**！
* **浏览器 Favicon 自动内嵌**：程序内嵌该图标并注册 `/favicon.ico` 路由，浏览器访问时标签页左上角自动展示该图标。

---

## 🔏 (可选) UPX 极速瘦身与 Windows 数字签名 (防拦截)

若需要进一步缩减可执行文件体积，并消除 Windows SmartScreen 蓝底拦截弹窗：

1. **UPX 极速压缩**（可选）：
   ```bash
   upx --best dist/app.exe
   ```
2. **本地生成专属自签名代码证书**（纯 Go 标准库 0 依赖）：
   ```bash
   # 支持参数：-name (发布者名称), -org (机构名称), -years (有效年限), -pass (密码), -out (输出目录)
   go run github.com/xbt/godeniter/cmd/cert -name "我的软件工作室" -org "我的公司" -years 10 -out ./certs
   ```
   *(执行后在本地 `./certs/` 目录下生成专属私钥和公钥 `app_codesign.cer`)*
3. **执行代码签名**：
   * 在 Mac/Linux 上：`osslsigncode sign -pkcs12 certs/app_codesign.pfx -pass 123456 -in dist/app.exe -out dist/app_signed.exe`
   * 在 Windows 上：`signtool sign /f certs\app_codesign.pfx /p 123456 dist\app.exe`
4. **客户机一键信任**：把生成的公钥 `app_codesign.cer` 给客户电脑，以管理员权限执行 `certutil -addstore -f "ROOT" app_codesign.cer`，软件从此双击秒开，永不弹未知发布者拦截！
   *(详见 [《Windows 数字签名与代码证书实战手册》](../godeniter/docs/code_signing.md))*

---

## 📄 开源许可证 (License)

Godeniter Starter 脚手架工程基于宽松友好的 **[MIT License](./LICENSE)** 协议开源，允许任何个人与企业自由用于商业业务系统或闭源软件的研发与分发。
