package services

import (
	"algo-arena-be/internals/repo"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type problemService struct {
	repo repo.ProblemRepo
}

func (s *problemService) GetProblems(w http.ResponseWriter, r *http.Request) {
	userID, ok := withUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	problems, err := s.repo.GetProblems(r.Context(), userID)
	if err != nil {
		log.Println("Failed to fetch problems: ", err)
		http.Error(w, "Failed to fetch problems", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, problems)
}

func (s *problemService) ViewProblem(w http.ResponseWriter, r *http.Request) {
	userID, ok := withUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	problemIDStr := chi.URLParam(r, "problemId")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	problem, err := s.repo.ViewProblem(r.Context(), userID, problemID)
	if err != nil {
		log.Println("Failed to fetch problem: ", err)
		http.Error(w, "Failed to fetch problem", http.StatusInternalServerError)
		return
	}
	if problem == nil {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, http.StatusOK, problem)
}
