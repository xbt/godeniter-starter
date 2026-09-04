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

## 📁 规范的项目目录结构

```text
godeniter-starter/
├── main.go                 # 应用启动入口 (路由注册、依赖注入装配)
├── config/                 # 应用配置层
│   └── app.go
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
