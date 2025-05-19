package server

import (
	"algo-arena-be/internals/services"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func NewChiRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByIP(100, time.Minute))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	return r
}

func RegisterRoutes(r *chi.Mux, ps services.ProblemService, sS services.SubmissionService) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	})

	r.Route("/", func(r chi.Router) {
		// r.Use(middlewares.AuthMiddleware(validateTokenFunc))

		r.Get("/problems", ps.GetProblems)
		r.Get("/problems/{problemId}", ps.ViewProblem)
		r.Post("/problems/{problemId}/run", sS.RunCode)
		r.Post("/problems/{problemId}/submit", sS.SubmitCode)

		r.Get("/submissions/{submissionId}", sS.GetSubmissionResult)
		r.Get("/runs/{runId}", sS.GetRunResult)
	})
}
