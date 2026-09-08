package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfig_RandomSQLiteDSN(t *testing.T) {
	cfg1 := DefaultConfig()
	cfg2 := DefaultConfig()

	if !strings.HasPrefix(cfg1.Database.DSN, "./data/app_") || !strings.HasSuffix(cfg1.Database.DSN, ".db") {
		t.Fatalf("DefaultConfig DSN 预期带随机后缀，实际为: %s", cfg1.Database.DSN)
	}

	if cfg1.Database.DSN == cfg2.Database.DSN {
		t.Fatalf("两次生成的 DefaultConfig SQLite 随机 DSN 不应相同: %s vs %s", cfg1.Database.DSN, cfg2.Database.DSN)
	}
}

func TestLoadConfig_PlaceholderAndPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	testCfgPath := filepath.Join(tmpDir, "config.json")

	// 1. 模拟用户配置了带 {random} 占位符的配置
	initJSON := `{
  "app": { "port": ":9090" },
  "database": {
    "driver": "sqlite",
    "dsn": "./data/custom_{random}.db"
  }
}`
	if err := os.WriteFile(testCfgPath, []byte(initJSON), 0644); err != nil {
		t.Fatal(err)
	}

	// 2. 首次 LoadConfig，预期占位符被自动替换并写回
	cfg := LoadConfig(testCfgPath)
	if strings.Contains(cfg.Database.DSN, "{random}") {
		t.Fatalf("预期 {random} 被替换，实际依然存在: %s", cfg.Database.DSN)
	}
	if !strings.HasPrefix(cfg.Database.DSN, "./data/custom_") || !strings.HasSuffix(cfg.Database.DSN, ".db") {
		t.Fatalf("DSN 格式不符合预期: %s", cfg.Database.DSN)
	}
	persistedDSN := cfg.Database.DSN

	// 3. 再次读取文件内容，验证是否真实固化持久化
	persistedBytes, err := os.ReadFile(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(persistedBytes), persistedDSN) {
		t.Fatalf("配置文件未成功固化已生成的 DSN，文件内容: %s", string(persistedBytes))
	}

	// 4. 第二次 LoadConfig，模拟服务重启，预期 DSN 绝对保持不变，保证数据不丢
	cfgReload := LoadConfig(testCfgPath)
	if cfgReload.Database.DSN != persistedDSN {
		t.Fatalf("服务重启再次加载预期复用同一 DSN，实际发生了改变: %s vs %s", cfgReload.Database.DSN, persistedDSN)
	}
}
