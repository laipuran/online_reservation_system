package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"ors-be/internal/api/http/handler"
	"ors-be/internal/api/http/middleware"
	"ors-be/internal/auth"
)

func NewServer(
	authH *handler.AuthHandler,
	tokenGen auth.TokenGenerator,
	allowedOrigins string,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{allowedOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", authH.Register())
		r.Post("/auth/login", authH.Login())

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(tokenGen))
			// 后续需要认证的路由在此添加
		})
	})

	return r
}
