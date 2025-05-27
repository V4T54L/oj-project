package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	mockdata "online-judge/internal/mock_data"
	"online-judge/internal/models"

	"github.com/go-chi/chi/v5"
)

const projectUserReviewStatus = "in review"
const projectActiveStatus = "active"

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
	}

	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func validateProblem(problemID int) {
	go func(probID int) {
		time.Sleep(10 * time.Second) // simulate execution delay

		for i, prob := range mockdata.ProblemsDB {
			if prob.ID == probID {
				mockdata.ProblemsDB[i].Status = projectActiveStatus
				break
			}
		}
	}(problemID)
}

func Signup(w http.ResponseWriter, r *http.Request) {
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

func Login(w http.ResponseWriter, r *http.Request) {
	var payload models.SignupPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// find user from database
	user := models.UserDB{Username: "some", Email: "some@one.in", ID: 1, Password: "", IsAdmin: true, IsActive: true}

	handleAuth(user, w)
	w.WriteHeader(http.StatusCreated)
}

func GetProblemList(w http.ResponseWriter, r *http.Request) {

	// handle filters

	problems := []models.ProblemInfo{}
	for _, problemDB := range mockdata.ProblemsDB {
		problems = append(problems, mockdata.ToProblemInfo(problemDB))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(problems)
}

func ViewProblem(w http.ResponseWriter, r *http.Request) {
	problemID := 0 // from query
	problem := mockdata.ToProblemDetail(mockdata.ProblemsDB[problemID])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(problem)
}

func CreateProblem(w http.ResponseWriter, r *http.Request) {
	var payload models.ProblemDB
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payload.ID = len(mockdata.ProblemsDB) + 1
	payload.Status = projectUserReviewStatus
	mockdata.ProblemsDB = append(mockdata.ProblemsDB, payload)

	validateProblem(payload.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": payload.ID})
}

func UpdateProblem(w http.ResponseWriter, r *http.Request) {
	var payload models.ProblemDB
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// simple linear update
	for i, prob := range mockdata.ProblemsDB {
		if prob.ID == payload.ID {
			mockdata.ProblemsDB[i] = payload
			break
		}
	}

	validateProblem(payload.ID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func SubmitCode(executeCode func(ctx context.Context, language string, payload models.ExecuteCodePayload) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload models.SubmitCodePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		submission := models.SubmissionDB{
			ID:        len(mockdata.Submissions) + 1,
			Status:    "pending",
			UserID:    1, // assuming a dummy user
			ProblemID: payload.ProblemID,
			RuntimeMS: 0,
			MemoryKB:  0,
		}

		mockdata.Submissions = append(mockdata.Submissions, submission)

		// TODO: obtain language from ID
		language := "python"
		if err := executeCode(r.Context(), language, models.ExecuteCodePayload{
			ID: submission.ID,
		}); err != nil {
			http.Error(w, "error submitting the code: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.SubmitCodeResponse{SubmissionID: submission.ID})
	}
}

func GetSubmissionResultByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "submissionID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid submission ID", http.StatusBadRequest)
		return
	}

	for _, sub := range mockdata.Submissions {
		if sub.ID == id && sub.Status != "pending" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(sub)
			return
		}
	}

	http.Error(w, "Submission not found or still pending", http.StatusNotFound)
}
