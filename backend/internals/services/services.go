package services

import (
	"algo-arena-be/internals/models"
	"algo-arena-be/internals/repo"
	"database/sql"
	"encoding/json"
	"net/http"
)

type ProblemService interface {
	GetProblems(w http.ResponseWriter, r *http.Request)
	ViewProblem(w http.ResponseWriter, r *http.Request)
}

type SubmissionService interface {
	GetSubmissionResult(w http.ResponseWriter, r *http.Request)
	GetRunResult(w http.ResponseWriter, r *http.Request)
	RunCode(w http.ResponseWriter, r *http.Request)
	SubmitCode(w http.ResponseWriter, r *http.Request)
}

type ExecutionService interface {
	Execute(payload models.ExecutionPayload) (models.ExecutionResponse, error)
}

func withUserID(_ *http.Request) (int, bool) {
	id, ok := 1, true
	// id, ok := r.Context().Value(userIDKey).(int)
	return id, ok
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func GetServices(db *sql.DB) (ProblemService, SubmissionService, error) {
	_, pR, sR, err := repo.GetRepos(db)
	if err != nil {
		return nil, nil, err
	}

	execService := executionService{}

	return &problemService{repo: pR}, &submissionService{repo: sR, executeFunc: execService.Execute}, nil
}