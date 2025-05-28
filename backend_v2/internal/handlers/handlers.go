package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"online-judge/internal/middleware"
	"online-judge/internal/models"
	"online-judge/internal/repo"
	"online-judge/internal/services"

	"github.com/go-chi/chi/v5"
)

func handleAuth(user models.UserDB, w http.ResponseWriter) {
	info := models.UserInfo{
		Username: user.Username,
		Email:    user.Email,
	}

	token := "generate from info"

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour), // Set expiration time
		HttpOnly: true,                           // Make cookie accessible only by the server
		Secure:   false,                          // Ensure cookie is sent only over HTTPS in production
		SameSite: http.SameSiteStrictMode,        // recommended, prevent CSRF
		Path:     "/",
	}

	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

type Handler struct {
	submissionRepo *repo.SubmissionRepo
	problemRepo    *repo.ProblemRepo
	redisService   *services.RedisService
}

func NewHandler(submissionRepo *repo.SubmissionRepo,
	problemRepo *repo.ProblemRepo,
	redisService *services.RedisService) (*Handler, error) {
	return &Handler{
		submissionRepo: submissionRepo,
		problemRepo:    problemRepo,
		redisService:   redisService,
	}, nil
}

func (h *Handler) validateProblem(ctx context.Context, problemID int) error {
	problem, err := h.problemRepo.GetProblemMetadata(ctx, problemID)
	if err != nil {
		return err
	}

	testCases, err := h.problemRepo.GetProblemTestCases(ctx, problemID)
	if err != nil {
		return err
	}

	// TODO: obtain language from ID
	language := "python"
	if err := h.redisService.ExecuteCode(ctx, language, models.ExecuteCodePayload{
		ID:             problemID,
		Code:           problem.SolutionCode,
		TestCases:      testCases,
		RuntimeLimitMS: problem.MemoryLimitKB,
		MemoryLimitKB:  problem.MemoryLimitKB,
		ExecutionType:  "validation",
	}); err != nil {
		return err
	}

	return nil
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var payload models.SignupPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// add user to the database
	user := models.UserDB{Username: "some", Email: "some@one.in", ID: 1, Password: "", IsAdmin: true, IsActive: true}

	handleAuth(user, w)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var payload models.LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// find user from database
	user := models.UserDB{Username: "some", Email: "some@one.in", ID: 1, Password: "", IsAdmin: true, IsActive: true}

	handleAuth(user, w)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetProblemList(w http.ResponseWriter, r *http.Request) {

	// handle filters

	problems, err := h.problemRepo.GetProblems(r.Context(), true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(problems)
}

func (h *Handler) ViewProblem(w http.ResponseWriter, r *http.Request) {
	problemID := 1 // from query
	problem, err := h.problemRepo.GetProblemByID(r.Context(), problemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(problem)
}

func (h *Handler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	var payload models.ProblemDB
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.problemRepo.CreateProblem(r.Context(), &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = h.validateProblem(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *Handler) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	var payload models.ProblemDB
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.problemRepo.UpdateProblemByID(r.Context(), &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = h.validateProblem(r.Context(), payload.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (h *Handler) SubmitCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		log.Println("Raw : ", r.Context().Value(middleware.UserIDKey))
		http.Error(w, "invalid token (userID)", http.StatusBadRequest)
		return
	}

	var payload models.SubmitCodePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	problem, err := h.problemRepo.GetProblemMetadata(r.Context(), payload.ProblemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	testCases, err := h.problemRepo.GetProblemTestCases(r.Context(), payload.ProblemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	submissionId, err := h.submissionRepo.NewSubmission(r.Context(), payload.ProblemID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: obtain language from ID
	language := "python"
	if err := h.redisService.ExecuteCode(r.Context(), language, models.ExecuteCodePayload{
		ID:             submissionId,
		Code:           payload.Code,
		TestCases:      testCases,
		RuntimeLimitMS: problem.MemoryLimitKB,
		MemoryLimitKB:  problem.MemoryLimitKB,
		ExecutionType:  "submission",
	}); err != nil {
		http.Error(w, "error submitting the code: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SubmitCodeResponse{SubmissionID: submissionId})
}

func (h *Handler) GetSubmissionResultByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "submissionID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid submission ID", http.StatusBadRequest)
		return
	}

	result, err := h.submissionRepo.GetSubmission(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
