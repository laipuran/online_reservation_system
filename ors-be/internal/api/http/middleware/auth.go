package middleware

import (
	"context"
	"net/http"
	"strings"

	"ors-be/internal/api/http/response"
	"ors-be/internal/auth"
)

type contextKey string

const UserCtxKey contextKey = "user"

func Auth(tokenGen auth.TokenGenerator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			if tokenStr == "" {
				response.JSON(w, http.StatusUnauthorized, response.Unauthorized("Token 格式错误"))
				return
			}

			claims, err := tokenGen.Validate(tokenStr)
			if err != nil {
				response.JSON(w, http.StatusUnauthorized, response.Unauthorized("Token 无效或已过期"))
				return
			}

			ctx := context.WithValue(r.Context(), UserCtxKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClaims(ctx context.Context) *auth.Claims {
	claims, _ := ctx.Value(UserCtxKey).(*auth.Claims)
	return claims
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r.Context())
			if claims == nil {
				response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
				return
			}
			if claims.Role != role {
				response.JSON(w, http.StatusForbidden, response.Forbidden("权限不足"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
