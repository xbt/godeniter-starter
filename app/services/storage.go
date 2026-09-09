package services

import (
	"github.com/xbt/godeniter/storage"
)

// StorageDriver 统一多存储驱动接口 (复用底层纯 Go 实现的 godeniter/storage)
// 演示特性：纯标准库 0 外部臃肿 SDK，支持本地文件、WebDAV 以及 AWS SigV4 规范的 S3 / Cloudflare R2 / 阿里云 OSS / MinIO
type StorageDriver = storage.Driver

// NewDefaultStorageDriver 演示创建默认存储驱动实例
// 开发者可根据生产环境变量或配置文件动态切换驱动实现
func NewDefaultStorageDriver(uploadDir string) StorageDriver {
	// 1. 默认演示：本地磁盘文件驱动
	return storage.NewLocalStorageDriver(uploadDir, "/uploads/images")

	// 2. 生产云存储演示 (S3 / Cloudflare R2 / 阿里云 OSS / MinIO，0-CGO 纯 Go 原生签名):
	// return storage.NewS3StorageDriver(storage.S3Config{
	// 	Endpoint:  "s3.us-east-1.amazonaws.com",
	// 	Region:    "us-east-1",
	// 	Bucket:    "my-app-bucket",
	// 	AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
	// 	SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	// 	Domain:    "https://cdn.mycompany.com", // 可选 CDN 加速域名
	// })

	// 3. WebDAV 协议驱动演示 (坚果云 / Nextcloud / InfiniCLOUD):
	// return storage.NewWebDAVStorageDriver("https://dav.jianguoyun.com/dav/", "user@mail.com", "app-password")
}
