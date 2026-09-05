# Godeniter Starter 脚手架项目

**Godeniter Starter** 是官方推荐的标准业务工程骨架。基于 **Godeniter 2.0 框架** 构建，拥有 **0 外部依赖**、**依赖注入**、**单文件打包** 与 **极速启动** 特性。

---

## 🚀 快速上手 (Quick Start)

### 1. 本地前台开发模式 (随时 Ctrl+C 停止)
```bash
# 启动 Web 服务 (默认前台运行，终端实时滚动日志与彩色 Banner)
go run main.go
```
启动后终端将自动打印本机与局域网访问地址，在浏览器中打开：`http://127.0.0.1:8080`
* **默认管理员账号**：`admin`
* **默认管理员密码**：`123456`

### 2. Linux / 服务器后台守护模式 (关闭终端/断开 SSH 持续运行)
无需编译，源码与打包二进制均原生支持标准服务生命周期管理：
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
若您希望将本 Web 服务作为独立的本地客户端或桌面托盘分发给用户使用：
* **默认开箱即用（无参数直接运行 / 双击）**：
  在 Windows 和 macOS 桌面环境下，直接执行或双击程序，**默认自动以系统托盘模式启动**，无需任何命令行参数：
  * **Windows**：Win32 原生自动隐藏控制台黑框，右下角任务栏托盘图标优雅常驻，右键弹出管理菜单，双击图标直接打开浏览器后台；
  * **macOS**：顶部状态栏常驻图标，点击展开管理菜单。
* **常驻原生菜单**：
  * 🌐 **打开管理后台**：调起系统默认浏览器访问 Web 页面；
  * 📁 **打开应用目录**：调起系统文件管理器 (Explorer / Finder) 定位应用目录，方便查找 `config.json` 与数据文件；
  * ℹ️ **关于系统**：弹窗展示应用名称、版本、运行端口与进程 PID 等信息；
  * ⏹️ **退出程序**：平滑优雅关闭 Web 服务并安全退出，托盘图标无感清理。
* **开发者控制台调试模式**：若需要在终端中查看彩色 ASCII Banner 与实时请求日志，只需传入 `console` 或 `run` 参数：
  ```bash
  go run main.go console      # 或 ./dist/app console
  ```

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

* **开箱即用**：首次运行若未检测到配置文件，程序将**自动在当前目录生成一份格式化的 `config.json`**。
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
      "dsn": "./data/app.db",
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

脚手架内置了一套完整的现代化轻量博客/内容管理系统，采用清新典雅的**浅色调 UI（Light Theme）**，全面覆盖 Godeniter 框架在实际企业级业务中可能用到的核心功能点：

1. **文件上传全链路演示 (Upload & Static Serv)**：
   - **封面图片上传**：发布/编辑文章表单支持本地选择封面图片上传，前端即时图片预览；
   - **服务端安全校验**：使用框架内置 `c.SaveUploadedFileWithOptions`，实施 5MB 上限限制、`.jpg/.png/.jpeg/.webp` 格式白名单校验，并自动重命名保存至 `./uploads/images/`；
   - **多端封面展示**：首页列表图文自适应卡片、详情页高清头图以及后台表格封面缩略图；同时提供独立 RESTful 上传接口 `/api/v1/upload`。
2. **自定义业务中间件流水线 (Custom Middleware)**：
   - **路由守卫中间件 (`AuthRequired`)**：保护 `/admin` 路由组，未登录拦截并携带 Flash 提示重定向；
   - **响应耗时中间件 (`ResponseTimer`)**：自动记录每个请求的处理耗时，向响应头注入 `X-Response-Time` 与 `Server-Timing`；
   - **基础安全头中间件 (`SecurityHeaders`)**：自动注入 `X-Content-Type-Options: nosniff`、`X-Frame-Options: SAMEORIGIN`、`X-XSS-Protection` 等安全标头。
3. **零依赖参数绑定与结构体 Tag 校验 (Binding & Validation)**：
   - 基于纯 Go 标准库与反射实现的轻量验证器，无需第三方库；
   - 表单提交声明 `binding:"required,min=3,max=80"`，校验失败自动在 Web 页面友好回显红字提示。
4. **企业级数据安全与脱敏实战 (Utils & Security)**：
   - **正文防 XSS 过滤**：发布与编辑正文自动执行 `str.XSSFilter` 过滤恶意脚本；
   - **敏感联系方式脱敏**：作者手机号/邮箱智能掩码脱敏（`str.MaskPhone` / `str.MaskEmail`）；
   - **内容摘要截断**：卡片摘要使用 `str.Truncate` 安全截断。
5. **前台 8 篇精品文章预置与 5 条/页分页**：
   - 倒序呈现 8 篇实战文章，直观展示 `第 1 / 2 页` 分页场景，支持上一页/下一页无缝翻页；
   - 顶部搜索框支持标题与内容模糊检索，并与分页深度联动。
6. **框架特性体验中心 (`/features`) 与 Panic 优雅容灾**：
   - 顶部导航设有“⚡ 框架特性”专区，可视化展示各功能点；
   - 提供 `/demo/panic` 测试端点，点击可验证 `Recovery()` 中间件优雅捕获运行时异常并输出彩色堆栈，**进程永不宕机**。
7. **无侵入 HTML 注释模板语法实战 (`<!--{{ ... }}-->`)**：
   - 视图使用框架原生 `app.LoadHTMLFS(subViews, "*.html")` 加载；
   - 在 `views/detail.html` 中实战应用 `<!--{{ if .CurrentUser }}-->` 与 `<!--{{ .Article.ID }}-->`。本地直接用浏览器双击打开静态 HTML 原型时不乱码、不破坏按钮排版，Go 服务端渲染时无缝编译为高效 AST。
8. **开箱即用的自动化端到端测试 (`main_test.go`)**：
   - 包含 7 大系统测试用例，覆盖中间件注入、Panic 恢复、特性页、分页检索、详情渲染（含无侵入模板断言）、文件上传及后台完整 CRUD，`go test -v .` 秒级全绿（~0.3s）。

---

## 📁 规范的项目目录结构

```text
godeniter-starter/
├── main.go                 # 应用启动入口 (中间件挂载、路由组注册、依赖注入装配)
├── main_test.go            # 自动化端到端测试 (中间件头、Panic恢复、文件上传、分页、CRUD)
├── config.json             # 外部动态配置文件 (端口、数据库、上传配置)
├── config/                 # 配置装配层
│   └── app.go
├── data/                   # 本地数据库存储目录 (自动创建)
├── app/
│   ├── controllers/        # 控制器层 (API 控制器与 Web 页面控制器)
│   │   ├── home.go         # 前台首页、文章详情页、特性体验中心与 Panic 测试
│   │   ├── admin.go        # 后台管理控制器 (带封面上传的列表、新建、编辑、删除 CRUD)
│   │   ├── auth.go         # Session 登录与注销控制器
│   │   └── api_article.go  # RESTful API 控制器 (含文件上传、分页检索与参数校验)
│   ├── models/             # 数据实体与请求校验 DTO (含 binding 规则)
│   │   └── article.go
│   ├── services/           # 业务逻辑层 (Service)
│   │   └── article.go      # 预设 8 篇文章，提供分页、搜索、自增阅读量、XSS过滤与 CRUD
│   └── middleware/         # 自定义业务中间件
│       ├── auth.go         # AuthRequired 登录认证拦截 (路由守卫)
│       ├── timer.go        # ResponseTimer 响应计时与性能监控
│       └── security.go     # SecurityHeaders 安全防护响应头
├── views/                  # 内嵌浅色调 HTML 模板 (单文件打包)
│   ├── index.html          # 浅色调前台首页 (含搜索框、封面卡片、分页导航)
│   ├── detail.html         # 浅色调文章详情页 (含高清封面图、面包屑与阅读量)
│   ├── features.html       # 浅色调框架特性全景体验中心 (可视化探针与交互测试)
│   ├── login.html          # 浅色调登录面板 (管理员登录)
│   ├── admin.html          # 浅色调后台管理中心 (含封面缩略图列表、操作入口)
│   └── article_form.html   # 浅色调文章发布与编辑表单 (含封面图片选择与本地预览)
├── uploads/                # 运行时文件上传存储目录
│   └── images/             # 上传图片存储目录 (内置 sample_cover.svg)
├── build.sh / build.bat    # 跨平台一键打包单文件脚本
└── go.mod                  # 模块声明 (引入 godeniter)
```


---

## 📦 一键编译单文件交付 (统一全能单二进制)

```bash
# 生成 Windows 64位统一全能单文件可执行程序 (dist/app.exe) 及 macOS/Linux 本地二进制
./build.sh     # macOS / Linux
build.bat      # Windows
```

生成的单文件无需安装任何环境，直接拷贝给客户，**统一且功能完备**：
* `dist/app.exe`：**Windows 统一全能二进制**
  - 在 CMD / PowerShell 中：支持 `start/status/stop/restart` 守护命令与前台日志；
  - 桌面双击或以托盘模式运行：Windows 原生自动隐藏控制台黑框，常驻屏幕右下角托盘，提供完整右键菜单！
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


