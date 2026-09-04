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

## 📁 规范的项目目录结构

```text
godeniter-starter/
├── main.go                 # 应用启动入口 (路由注册、依赖注入装配)
├── config.json             # 外部动态配置文件 (端口、数据库、上传配置)
├── config/                 # 配置装配层
│   └── app.go
├── data/                   # 本地数据库存储目录 (自动创建)
├── app/
│   ├── controllers/        # 控制器层 (API 控制器与 Web 页面控制器)
│   │   ├── home.go         # 首页/仪表盘控制器
│   │   ├── api_article.go  # RESTful API 控制器 (含文件上传、分页检索与参数校验)
│   │   └── auth.go         # Session 登录与注销控制器
│   ├── models/             # 数据实体与请求校验 DTO
│   │   └── article.go
│   ├── services/           # 业务逻辑层 (Service)
│   │   └── article.go
│   └── middleware/         # 自定义中间件 (如 AuthRequired 登录拦截)
│       └── auth.go
├── views/                  # 内嵌 HTML 模板 (单文件打包)
│   ├── index.html
│   └── login.html
├── uploads/                # 运行时文件上传存储目录
├── build.sh / build.bat    # 跨平台一键打包单文件脚本
└── go.mod                  # 模块声明 (引入 godeniter)
```

---

## 📦 一键编译单文件交付 (Windows .exe)

```bash
# 生成 Windows 64位单文件可执行程序 (dist/app.exe)
./build.sh
```

生成的 `dist/app.exe` 无需安装任何环境，直接拷贝给客户，**双击即可直接运行**！
