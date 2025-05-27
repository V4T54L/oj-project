package router

import (
	"database/sql"
	"net/http"
	authmodule "online-judge/internal/auth_module"
	"online-judge/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewChiRouter(db *sql.DB, auth *authmodule.JWTAuth) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
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

		r.Post("/submit", handlers.SubmitCode)
		r.Get("/submission/{submissionID}", handlers.GetSubmissionResultByID)
	})

	return r
}
