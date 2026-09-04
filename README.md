# Godeniter Starter 脚手架项目

**Godeniter Starter** 是官方推荐的标准业务工程骨架。基于 **Godeniter 2.0 框架** 构建，拥有 **0 外部依赖**、**依赖注入**、**单文件打包** 与 **极速启动** 特性。

---

## 🚀 快速上手 (Quick Start)

### 1. 本地启动运行
```bash
# 启动 Web 服务
go run main.go
```
启动后终端将自动打印本机与局域网访问地址，在浏览器中打开：`http://127.0.0.1:8080`

* **默认管理员账号**：`admin`
* **默认管理员密码**：`123456`

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
* **修改端口与数据库**：直接使用任意文本编辑器（如记事本）打开 `config.json` 即可调整：
  ```json
  {
    "app": {
      "name": "Godeniter Starter Application",
      "port": ":8080",
      "env": "development",
      "session_key": "godeniter-starter-secret-salt-2026"
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
7. **开箱即用的自动化端到端测试 (`main_test.go`)**：
   - 包含 7 大系统测试用例，覆盖中间件注入、Panic 恢复、特性页、分页检索、详情渲染、文件上传及后台完整 CRUD，`go test -v .` 秒级全绿（~0.3s）。

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

## 📦 一键编译单文件交付 (Windows .exe)

```bash
# 运行单元测试
go test -v .

# 生成 Windows 64位单文件可执行程序 (dist/app.exe) 及本地二进制
./build.sh
```

生成的 `dist/app.exe` 无需安装任何环境，直接拷贝给客户，**双击即可直接运行**！

