package router

import (
	"context"
	"database/sql"
	"net/http"
	authmodule "online-judge/internal/auth_module"
	"online-judge/internal/handlers"
	"online-judge/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewChiRouter(db *sql.DB, auth *authmodule.JWTAuth,
	executeCodefunc func(ctx context.Context, language string, payload models.ExecuteCodePayload) error,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	})

	// Public routes
	r.Post("/signup", handlers.Signup)
	r.Post("/login", handlers.Login)

	// Protected routes
	r.Route("/api", func(r chi.Router) {
		// r.Use(custom_middleware.JWTAuthMiddleware(auth))
		// r.Get("/profile", handlers.ProfileHandler(db))

		r.Get("/problems", handlers.GetProblemList)
		r.Get("/problems/{problemID}", handlers.ViewProblem)
		r.Post("/problems", handlers.CreateProblem)
		r.Put("/problems", handlers.UpdateProblem)

		r.Post("/submit", handlers.SubmitCode(executeCodefunc))
		r.Get("/submissions/{submissionID}", handlers.GetSubmissionResultByID)
	})

	return r
}
