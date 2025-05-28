package repo

import (
	"context"
	"errors"
	"log"
	"online-judge/internal/models"
	// "sync" // Uncomment if thread safety is needed
)

type SubmissionRepo struct {
	db []models.SubmissionDB
}

func NewSubmissionRepo() *SubmissionRepo {
	return &SubmissionRepo{db: make([]models.SubmissionDB, 0)}
}

// Creates a new submission and appends it to the DB
func (r *SubmissionRepo) NewSubmission(ctx context.Context, problemId, userId int) (int, error) {
	submission := models.SubmissionDB{
		ID:        len(r.db) + 1,
		Status:    "pending",
		UserID:    userId,
		ProblemID: problemId,
		RuntimeMS: 0,
		MemoryKB:  0,
	}
	r.db = append(r.db, submission)
	return submission.ID, nil
}

// Retrieves a submission by ID (only if not pending)
func (r *SubmissionRepo) GetSubmission(ctx context.Context, submissionId int) (*models.SubmissionDB, error) {
	log.Println("\n\n Submissions:", r.db)

	for i := range r.db {
		if r.db[i].ID == submissionId {
			if r.db[i].Status == "pending" {
				log.Println("\nSubmission found. Status:", r.db[i].Status, r.db[i].MemoryKB, r.db[i].RuntimeMS)
				return nil, errors.New("submission is still pending")
			}
			return &r.db[i], nil // FIX: return pointer to actual element
		}
	}
	return nil, errors.New("submission not found")
}

// Updates the runtime, memory, and status of a submission
func (r *SubmissionRepo) UpdateSubmission(ctx context.Context, submissionId, runtime, memory int, status string) error {
	log.Println("\n\nUpdate Submission request received:", submissionId, runtime, memory, status)
	log.Println("Existing submissions: ", r.db)

	for i := range r.db {
		log.Println("Checking: ", r.db[i].ID, submissionId)
		if r.db[i].ID == submissionId {
			log.Println("\nSubmission found. Updating...")
			r.db[i].RuntimeMS = runtime
			r.db[i].MemoryKB = memory
			r.db[i].Status = status
			return nil
		}
	}
	return errors.New("submission not found")
}
