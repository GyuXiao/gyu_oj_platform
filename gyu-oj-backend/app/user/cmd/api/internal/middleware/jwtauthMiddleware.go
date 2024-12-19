package middleware

import (
	"github.com/zeromicro/go-zero/rest/handler"
	"gyu-oj-backend/common/constant"
	"net/http"
)

type JwtAuthMiddleware struct {
	secret string
}

func NewJwtAuthMiddleware(secret string) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		secret: secret,
	}
}

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 存在 Jwt token
		if len(r.Header.Get(constant.AuthorizationHeader)) > 0 {
			// 解析 token 并将数据写入 context
			authHandler := handler.Authorize(m.secret)
			authHandler(next).ServeHTTP(w, r)
			return
		}
		// 否则继续执行
		next(w, r)
	}
}
