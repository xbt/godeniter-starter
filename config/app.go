package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Config 全局配置定义
type Config struct {
	App      AppConfig      `json:"app"`
	Database DatabaseConfig `json:"database"`
	Upload   UploadConfig   `json:"upload"`
}

// AppConfig 应用服务基础配置
type AppConfig struct {
	Name       string `json:"name"`        // 应用名称
	Port       string `json:"port"`        // 服务监听端口 (例如 ":8080")
	Env        string `json:"env"`         // 运行环境: development / production
	SessionKey string `json:"session_key"` // Session 加密签名密钥
}

// DatabaseConfig 数据库连接配置 (支持 sqlite / mysql 等标准 SQL 驱动)
type DatabaseConfig struct {
	Driver          string `json:"driver"`            // 驱动类型: sqlite 或 mysql
	DSN             string `json:"dsn"`               // 连接数据源串 (如 "./data/app.db" 或 "user:pwd@tcp(127.0.0.1:3306)/dbname")
	MaxOpenConns    int    `json:"max_open_conns"`    // 最大连接数 (SQLite 建议 1，MySQL 建议 20~100)
	MaxIdleConns    int    `json:"max_idle_conns"`    // 最大空闲连接数
	ConnMaxLifetime int    `json:"conn_max_lifetime"` // 连接最长存活秒数
}

// UploadConfig 文件上传安全配置
type UploadConfig struct {
	Dir         string   `json:"dir"`          // 存储物理目录
	MaxSizeMB   int64    `json:"max_size_mb"`  // 单文件最大体积上限 (MB)
	AllowedExts []string `json:"allowed_exts"` // 扩展名白名单列表
}

// DefaultConfig 返回开箱即用的默认配置
func DefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:       "Godeniter Starter Application",
			Port:       ":8080",
			Env:        "development",
			SessionKey: "godeniter-starter-secret-salt-2026",
		},
		Database: DatabaseConfig{
			Driver:          "sqlite",
			DSN:             "./data/app.db",
			MaxOpenConns:    1,
			MaxIdleConns:    1,
			ConnMaxLifetime: 300,
		},
		Upload: UploadConfig{
			Dir:         "./uploads",
			MaxSizeMB:   5,
			AllowedExts: []string{".jpg", ".png", ".jpeg", ".webp"},
		},
	}
}

// LoadConfig 加载配置（纯 Go 标准库实现，0 外部依赖）：
// 1. 加载内置默认配置；
// 2. 检测本地 config.json（若不存在则自动生成一份带缩进的样例配置文件）；
// 3. 读取环境变量覆盖对应字段（方便 Docker / 云原生部署）。
func LoadConfig(configPaths ...string) *Config {
	cfg := DefaultConfig()

	filePath := "config.json"
	if len(configPaths) > 0 && configPaths[0] != "" {
		filePath = configPaths[0]
	}

	// 1. 检测外部文件
	if fileBytes, err := os.ReadFile(filePath); err == nil {
		if err := json.Unmarshal(fileBytes, cfg); err != nil {
			fmt.Printf(">> [WARN] 解析配置文件 [%s] 失败: %v，将回退使用默认配置\n", filePath, err)
		} else {
			fmt.Printf(">> [CONFIG] 成功从本地加载配置文件: %s\n", filePath)
		}
	} else if os.IsNotExist(err) {
		// 若配置文件不存在，自动在当前目录生成一份初始模板供开发者或客户参考修改
		if dumpBytes, dumpErr := json.MarshalIndent(cfg, "", "  "); dumpErr == nil {
			_ = os.WriteFile(filePath, dumpBytes, 0644)
			fmt.Printf(">> [CONFIG] 未检测到配置文件，已自动在当前目录创建默认模板: %s\n", filePath)
		}
	}

	// 2. 环境变量动态覆盖 (最高优先级)
	if p := os.Getenv("PORT"); p != "" {
		if p[0] != ':' {
			p = ":" + p
		}
		cfg.App.Port = p
	}
	if env := os.Getenv("APP_ENV"); env != "" {
		cfg.App.Env = env
	}
	if key := os.Getenv("SESSION_KEY"); key != "" {
		cfg.App.SessionKey = key
	}
	if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
		cfg.Database.DSN = dsn
	}
	if driver := os.Getenv("DATABASE_DRIVER"); driver != "" {
		cfg.Database.Driver = driver
	}
	if uploadDir := os.Getenv("UPLOAD_DIR"); uploadDir != "" {
		cfg.Upload.Dir = uploadDir
	}
	if maxMB := os.Getenv("UPLOAD_MAX_MB"); maxMB != "" {
		if n, err := strconv.ParseInt(maxMB, 10, 64); err == nil {
			cfg.Upload.MaxSizeMB = n
		}
	}

	return cfg
}
