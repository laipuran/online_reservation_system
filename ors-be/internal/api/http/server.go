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
	userH *handler.UserHandler,
	providerH *handler.ServiceProviderHandler,
	serviceH *handler.ServiceHandler,
	tagH *handler.TagHandler,
	categoryH *handler.CategoryHandler,
	interestH *handler.UserInterestHandler,
	reservationH *handler.ReservationHandler,
	reviewH *handler.ReviewHandler,
	notificationH *handler.NotificationHandler,
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
		r.Get("/categories", categoryH.List())
		r.Get("/providers/{id}", providerH.GetByID())
		r.Get("/providers/{id}/services", serviceH.ListByProvider())
		r.Get("/providers/{id}/reviews", reviewH.ListByProvider())
		r.Get("/services", serviceH.List())
		r.Get("/services/{id}", serviceH.GetByID())
		r.Get("/services/{id}/reviews", reviewH.ListByService())
		r.Get("/services/{id}/tags", serviceH.ListTags())
		r.Get("/tags", tagH.List())
		r.Get("/tags/{id}", tagH.GetByID())

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(tokenGen))
			r.Get("/users/me", userH.GetMine())
			r.Put("/users/me", userH.UpdateMine())
			r.Put("/users/me/password", userH.UpdatePassword())
			r.Get("/users/me/interests", interestH.ListMine())
			r.Put("/users/me/interests", interestH.ReplaceMine())
			r.Get("/notifications", notificationH.ListMine())
			r.Get("/notifications/unread-count", notificationH.CountUnread())
			r.Put("/notifications/{id}/read", notificationH.MarkRead())
			r.Put("/notifications/read-all", notificationH.MarkAllRead())

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("customer"))
				r.Post("/reservations", reservationH.Create())
				r.Get("/reservations", reservationH.ListMine())
				r.Get("/reservations/{id}", reservationH.GetMine())
				r.Put("/reservations/{id}/cancel", reservationH.CancelMine())
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("provider"))
				r.Post("/tags", tagH.Create())
				r.Post("/providers/me", providerH.CreateMine())
				r.Get("/providers/me", providerH.GetMine())
				r.Put("/providers/me", providerH.UpdateMine())
				r.Post("/services", serviceH.Create())
				r.Put("/services/{id}", serviceH.Update())
				r.Patch("/services/{id}/status", serviceH.UpdateStatus())
				r.Put("/services/{id}/tags", serviceH.ReplaceTags())
				r.Get("/provider/reservations", reservationH.ListForProvider())
				r.Put("/provider/reservations/{id}/confirm", reservationH.ConfirmForProvider())
				r.Put("/provider/reservations/{id}/reject", reservationH.RejectForProvider())
			})
		})
	})

	return r
}
