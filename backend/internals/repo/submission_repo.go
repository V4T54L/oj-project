package repo

import (
	"algo-arena-be/internals/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type submissionRepo struct {
	db *sql.DB
}

// const statusAccepted = "Accepted"
const statusPending = "Pending"

// const statusFailed = "Failed"

// GetSubmissionResult retrieves metadata for a user’s submission.
func (r *submissionRepo) GetSubmissionResult(ctx context.Context, userId, submissionId int) (*models.SubmissionResult, error) {
	ctx, cancel := context.WithTimeout(ctx, maxQuerySeconds*time.Second)
	defer cancel()

	query := `
		SELECT cs.id, s.name, cs.runtime_ms, cs.memory_kb, cs.message
		FROM code_submissions cs
		JOIN status_ids s ON s.id = cs.status_id
		WHERE cs.id = $1 AND cs.user_id = $2
	`

	var result models.SubmissionResult
	err := r.db.QueryRowContext(ctx, query, submissionId, userId).Scan(
		&result.ID, &result.Verdict, &result.RuntimeMS, &result.MemoryKB, &result.Message,
	)
	// if err == sql.ErrNoRows {
	// 	return nil, nil
	// }
	if err != nil {
		return nil, fmt.Errorf("GetSubmissionResult: %w", err)
	}

	return &result, nil
}

// GetRunResult retrieves a simulated run result, joined with test input/output.
func (r *submissionRepo) GetRunResult(ctx context.Context, userId, runId int) (*models.RunResult, error) {
	ctx, cancel := context.WithTimeout(ctx, maxQuerySeconds*time.Second)
	defer cancel()

	query := `
		SELECT cs.id, s.name, cs.runtime_ms, cs.memory_kb, cs.message,
		       ht.input, cs.code AS output, ht.expected_output
		FROM code_submissions cs
		JOIN status_ids s ON cs.status_id = s.id
		JOIN hidden_test_cases ht ON ht.problem_id = cs.problem_id
		WHERE cs.id = $1 AND cs.user_id = $2
		LIMIT 1
	`

	var result models.RunResult
	err := r.db.QueryRowContext(ctx, query, runId, userId).Scan(
		&result.ID, &result.Verdict, &result.RuntimeMS, &result.MemoryKB, &result.Message,
		&result.Input, &result.Output, &result.ExpectedOutput,
	)
	// if err == sql.ErrNoRows {
	// 	return nil, nil
	// }
	if err != nil {
		return nil, fmt.Errorf("GetRunResult: %w", err)
	}

	return &result, nil
}

// RunCode inserts a temporary submission record for simulated runs (no judge).
func (r *submissionRepo) RunCode(ctx context.Context, userId, problemId, languageId int, code string, inputs []string) (int, error) {
	// In a real system, delegate this to an async execution engine
	statusID, err := r.getStatusID(ctx, statusPending)
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(ctx, maxQuerySeconds*time.Second)
	defer cancel()

	query := `
		INSERT INTO code_submissions (user_id, problem_id, language_id, code, status_id, message)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var runID int
	err = r.db.QueryRowContext(ctx, query, userId, problemId, languageId, code, statusID, "Code run queued").Scan(&runID)
	if err != nil {
		return 0, fmt.Errorf("RunCode: %w", err)
	}

	return runID, nil
}

// SubmitCode stores a submission record to be judged later.
func (r *submissionRepo) SubmitCode(ctx context.Context, userId, problemId, languageId int, code string) (int, error) {
	// Stub for real submission queueing
	statusID, err := r.getStatusID(ctx, statusPending)
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(ctx, maxQuerySeconds*time.Second)
	defer cancel()

	query := `
		INSERT INTO code_submissions (user_id, problem_id, language_id, code, status_id, message)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var submissionID int
	err = r.db.QueryRowContext(ctx, query, userId, problemId, languageId, code, statusID, "Submission queued").Scan(&submissionID)
	if err != nil {
		return 0, fmt.Errorf("SubmitCode: %w", err)
	}

	return submissionID, nil
}

// AddRunResult updates an existing code run submission with results.
func (r *submissionRepo) AddRunResult(ctx context.Context, result *models.RunResult) error {
	return r.updateSubmission(ctx, result.ID, result.Verdict, result.RuntimeMS, result.MemoryKB, result.Message)
}

// AddSubmissionResult updates a real submission with verdict + results.
func (r *submissionRepo) AddSubmissionResult(ctx context.Context, result *models.RunResult) error {
	return r.updateSubmission(ctx, result.ID, result.Verdict, result.RuntimeMS, result.MemoryKB, result.Message)
}

// Helper to update submission records with final results
func (r *submissionRepo) updateSubmission(ctx context.Context, submissionID int, verdict string, runtimeMS, memoryKB int, message string) error {
	statusID, err := r.getStatusID(ctx, verdict)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, maxQuerySeconds*time.Second)
	defer cancel()

	query := `
		UPDATE code_submissions
		SET status_id = $1, runtime_ms = $2, memory_kb = $3, message = $4
		WHERE id = $5
	`

	_, err = r.db.ExecContext(ctx, query, statusID, runtimeMS, memoryKB, message, submissionID)
	if err != nil {
		return fmt.Errorf("updateSubmission: %w", err)
	}
	return nil
}

// Helper to resolve status name to ID
func (r *submissionRepo) getStatusID(ctx context.Context, name string) (int, error) {
	query := `SELECT id FROM status_ids WHERE name = $1`
	var id int
	err := r.db.QueryRowContext(ctx, query, name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("getStatusID: %w", err)
	}
	return id, nil
}
