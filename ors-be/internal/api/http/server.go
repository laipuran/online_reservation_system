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
	providerH *handler.ServiceProviderHandler,
	serviceH *handler.ServiceHandler,
	tagH *handler.TagHandler,
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
		r.Get("/providers/{id}", providerH.GetByID())
		r.Get("/providers/{id}/services", serviceH.ListByProvider())
		r.Get("/services", serviceH.List())
		r.Get("/services/{id}", serviceH.GetByID())
		r.Get("/tags", tagH.List())
		r.Get("/tags/{id}", tagH.GetByID())

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(tokenGen))
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("provider"))
				r.Post("/tags", tagH.Create())
				r.Post("/providers/me", providerH.CreateMine())
				r.Get("/providers/me", providerH.GetMine())
				r.Put("/providers/me", providerH.UpdateMine())
				r.Post("/services", serviceH.Create())
				r.Put("/services/{id}", serviceH.Update())
				r.Patch("/services/{id}/status", serviceH.UpdateStatus())
			})
		})
	})

	return r
}
