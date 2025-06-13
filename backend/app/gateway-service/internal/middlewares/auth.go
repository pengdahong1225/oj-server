package middlewares

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func AuthLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 排除登录等无需认证的路径
		if skipAuth(r) {
			next.ServeHTTP(w, r)
			return
		}

		// jwt鉴权头部信息 token
		token := r.Header.Get("x-token")
		if token == "" {
			http.Error(w, "未登录", http.StatusUnauthorized)
			return
		}
		j := NewJWT()
		// 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, TokenExpired) {
				http.Error(w, "授权已过期", http.StatusUnauthorized)
				return
			} else {
				http.Error(w, "token验证失败", http.StatusUnauthorized)
				return
			}
		}

		// 剩余时间
		remain := claims.ExpiresAt - time.Now().Unix()
		logrus.Debugf("剩余：%v, [%v天]\n", remain, remain/86400)

		// token通过
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
