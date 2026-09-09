package middleware

import (
	"net/http"

	"github.com/xbt/godeniter"
	coreMiddleware "github.com/xbt/godeniter/middleware"
)

// APIKeyAuth 演示 API Key / Bearer Token 鉴权中间件
// 底层复用 godeniter/middleware.KeyAuth()，支持从 Authorization: Bearer <key>、X-API-Key 头或 ?api_key 参数提取
func APIKeyAuth(expectedKey string) godeniter.HandlerFunc {
	mw := coreMiddleware.KeyAuth(func(key string) bool {
		return key == expectedKey
	}, coreMiddleware.KeyAuthOptions{
		ErrorHandler: func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json; charset=utf-8")
			res.WriteHeader(http.StatusUnauthorized)
			_, _ = res.Write([]byte(`{"code":401,"message":"401 Unauthorized: 无效或缺失 API Key (演示密钥: starter_secret_key_8888)"}`))
		},
	})

	return godeniter.WrapMiddleware(mw)
}
