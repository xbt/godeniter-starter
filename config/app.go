package config

// AppConfig 应用全局配置结构
type AppConfig struct {
	AppName    string
	Port       string
	Env        string
	UploadDir  string
	SessionKey string
}

// LoadConfig 加载配置
func LoadConfig() *AppConfig {
	return &AppConfig{
		AppName:    "Godeniter Starter Application",
		Port:       ":8080",
		Env:        "development",
		UploadDir:  "./uploads",
		SessionKey: "godeniter-starter-secret-salt-2026",
	}
}
