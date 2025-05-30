package router

import (
	"database/sql"
	"net/http"
	authmodule "online-judge/internal/auth_module"
	"online-judge/internal/handlers"
	custom_middleware "online-judge/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewChiRouter(db *sql.DB, auth *authmodule.JWTAuth,
	handler handlers.Handler,
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
	r.Post("/signup", handler.Signup)
	r.Post("/login", handler.Login)

	// Protected routes
	r.Route("/api", func(r chi.Router) {
		r.Use(custom_middleware.JWTAuthMiddleware(auth))
		// r.Get("/profile", handler.ProfileHandler(db))

		r.Get("/problems", handler.GetProblemList)
		r.Get("/problems/{problemID}", handler.ViewProblem)
		r.Post("/problems", handler.CreateProblem)
		r.Put("/problems", handler.UpdateProblem)

		r.Post("/submit", handler.SubmitCode)
		r.Get("/submissions/{submissionID}", handler.GetSubmissionResultByID)

		r.Route("/contests", func(r chi.Router) {
			r.Get("/", handler.ListContests)
			r.Get("/{contestID}", handler.GetContest)
			r.Post("/", handler.CreateContest)
			r.Post("/{contestID}/problems/{problemID}", handler.AddProblemToContest)
			r.Post("/{contestID}/submit", handler.SubmitContestSolution)
			r.Get("/{contestID}/leaderboard", handler.GetLeaderboard)
		})

	})

	return r
}
