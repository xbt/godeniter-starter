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

## 🌟 核心演示功能 (Demo Features)

脚手架内置了一套完整的现代化轻量博客/内容管理系统，采用清新典雅的**浅色调 UI（Light Theme）**，涵盖日常 Web 开发标配能力：

1. **首页分页与搜索联动**：
   - 预置 **8 篇**精品架构与实战文章，默认每页展示 **5 篇**（直观呈现 `第 1 页 / 第 2 页` 分页场景）；
   - 支持前台顶部搜索框，输入关键词即可对标题和正文进行模糊检索，并与分页参数深度联动。
2. **文章详情页 (`/article/:id`)**：
   - 点击文章卡片或标题进入专属详情页，展示完整的正文排版、作者脱敏信息、发布时间与阅读计数；
   - 每次访问详情页自动累加阅读量 (`views`)，展示路由动态参数解析能力。
3. **后台文章管理与完整增删改查 (`/admin/articles`)**：
   - 标配登录权限拦截，未登录自动跳转到 `/login`；
   - **列表管理**：表格化直观管理所有文章，支持搜索、查看阅读数；
   - **发布新文章**：支持富文本正文录入与字段校验；
   - **编辑更新**：动态回显旧数据并保存修改；
   - **快速删除**：提供一键确认删除。
4. **全套现代化浅色调设计**：
   - 采用高级浅灰背景 (`#f8fafc`)、白底圆角卡片、柔和阴影与品牌蓝按钮，阅读体验清爽舒适。
5. **开箱即用的自动化测试**：
   - 包含 `main_test.go`，执行 `go test -v .` 即可秒级完成所有页面渲染、搜索、分页、登录鉴权及后台 CRUD 的端到端测试。

---

## 📁 规范的项目目录结构

```text
godeniter-starter/
├── main.go                 # 应用启动入口 (路由注册、依赖注入装配)
├── main_test.go            # 自动化端到端测试 (首页分页、搜索、详情页、后台 CRUD)
├── config.json             # 外部动态配置文件 (端口、数据库、上传配置)
├── config/                 # 配置装配层
│   └── app.go
├── data/                   # 本地数据库存储目录 (自动创建)
├── app/
│   ├── controllers/        # 控制器层 (API 控制器与 Web 页面控制器)
│   │   ├── home.go         # 前台首页与文章详情页控制器 (搜索、分页、详情)
│   │   ├── admin.go        # 后台管理控制器 (文章列表、新建、编辑、删除 CRUD)
│   │   ├── auth.go         # Session 登录与注销控制器
│   │   └── api_article.go  # RESTful API 控制器 (含文件上传、分页检索与参数校验)
│   ├── models/             # 数据实体与请求校验 DTO
│   │   └── article.go
│   ├── services/           # 业务逻辑层 (Service)
│   │   └── article.go      # 预设 8 篇文章，提供分页、搜索、自增阅读量与 CRUD
│   └── middleware/         # 自定义中间件 (如 AuthRequired 登录拦截)
│       └── auth.go
├── views/                  # 内嵌浅色调 HTML 模板 (单文件打包)
│   ├── index.html          # 浅色调前台首页 (含搜索框、文章卡片、分页导航)
│   ├── detail.html         # 浅色调文章详情页 (含面包屑、正文排版与阅读量)
│   ├── login.html          # 浅色调登录面板 (管理员登录)
│   ├── admin.html          # 浅色调后台管理中心 (文章列表、操作入口)
│   └── article_form.html   # 浅色调文章发布与编辑通用表单
├── uploads/                # 运行时文件上传存储目录
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

