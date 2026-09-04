package config

import "os"

// AppConfig 应用全局配置结构
type AppConfig struct {
	AppName    string // 应用名称
	Port       string // 监听端口 (如 :8080)
	Env        string // 运行环境 (development / production)
	UploadDir  string // 本地文件上传存储目录
	SessionKey string // 服务端 Session HMAC 签名密钥
}

// LoadConfig 加载配置（支持通过系统环境变量动态覆盖，适配 Docker / 云原生 / 客户机部署）
func LoadConfig() *AppConfig {
	return &AppConfig{
		AppName:    getEnv("APP_NAME", "Godeniter Starter Application"),
		Port:       getEnv("PORT", ":8080"),
		Env:        getEnv("APP_ENV", "development"),
		UploadDir:  getEnv("UPLOAD_DIR", "./uploads"),
		SessionKey: getEnv("SESSION_KEY", "godeniter-starter-secret-salt-2026"),
	}
}

// getEnv 读取指定环境变量；若未配置或为空则返回默认值 fallback
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
