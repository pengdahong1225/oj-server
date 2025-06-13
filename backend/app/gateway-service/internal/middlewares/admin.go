package middlewares

import (
	"net/http"
)

func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value("claims").(*UserClaims)
		if claims == nil {
			http.Error(w, "无权限", http.StatusForbidden)
			return
		}
		if claims.Authority == 0 {
			http.Error(w, "无权限", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
