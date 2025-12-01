package routes

import (
	"student-portal/internal/config"
	"student-portal/internal/enums"
	"student-portal/internal/handler"

	appMiddleware "student-portal/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// SetupRouter configures the Chi router with middlewares and routes.
func SetupRouter(cfg *config.Config, authHandler *handler.AuthHandler, userHandler *handler.UserHandler) *chi.Mux {
	r := chi.NewRouter()

	// Global Middleware
	r.Use(
		appMiddleware.RequestLogger, // Custom structured request logging
		middleware.Recoverer,        // Recover from panics
		cors.New(cors.Options{ // CORS setup
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}).Handler,
	)

	// Public Routes (No authentication required)
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// Protected Routes (Authentication required)
	r.Route("/api", func(r chi.Router) {
		r.Route("/profile", func(r chi.Router) {
			r.Use(appMiddleware.AuthMiddleware(cfg))
			r.Get("/", userHandler.GetOwnProfile)
			r.Put("/", userHandler.UpdateOwnProfile)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(appMiddleware.AuthMiddleware(cfg), appMiddleware.RoleMiddleware(string(enums.RoleAdmin)))
			r.Get("/", userHandler.ListUsers)
			r.Get("/{id}", userHandler.GetUserByID)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})
	})

	return r
}
