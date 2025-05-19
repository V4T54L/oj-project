package services

import (
	"algo-arena-be/internals/models"
	"algo-arena-be/internals/repo"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type submissionService struct {
	repo repo.SubmissionRepo
}

func (s *submissionService) GetSubmissionResult(w http.ResponseWriter, r *http.Request) {
	userID, ok := withUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	submissionID, err := strconv.Atoi(chi.URLParam(r, "submissionId"))
	if err != nil {
		http.Error(w, "Invalid submission ID", http.StatusBadRequest)
		return
	}

	result, err := s.repo.GetSubmissionResult(r.Context(), userID, submissionID)
	if err != nil {
		http.Error(w, "Failed to fetch submission result", http.StatusInternalServerError)
		return
	}
	if result == nil {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *submissionService) GetRunResult(w http.ResponseWriter, r *http.Request) {
	userID, ok := withUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	runID, err := strconv.Atoi(chi.URLParam(r, "runId"))
	if err != nil {
		http.Error(w, "Invalid run ID", http.StatusBadRequest)
		return
	}

	result, err := s.repo.GetRunResult(r.Context(), userID, runID)
	if err != nil {
		http.Error(w, "Failed to fetch run result", http.StatusInternalServerError)
		return
	}
	if result == nil {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *submissionService) RunCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := withUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	problemID, err := strconv.Atoi(chi.URLParam(r, "problemId"))
	if err != nil {
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	var payload models.SubmissionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	runID, err := s.repo.RunCode(r.Context(), userID, problemID, payload.LanguageID, payload.Code, payload.TestCases)
	if err != nil {
		http.Error(w, "Failed to run code", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"run_id": runID})
}

func (s *submissionService) SubmitCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := withUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	problemID, err := strconv.Atoi(chi.URLParam(r, "problemId"))
	if err != nil {
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	var payload models.SubmissionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	submissionID, err := s.repo.SubmitCode(r.Context(), userID, problemID, payload.LanguageID, payload.Code)
	if err != nil {
		http.Error(w, "Failed to submit code", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"submission_id": submissionID})
}
