package repo

import (
	"algo-arena-be/internals/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

const acceptedStatus = "Accepted"

type problemRepo struct {
	db *sql.DB
}

// GetProblems returns a list of problems with user-specific IsSolved flags.
func (r *problemRepo) GetProblems(ctx context.Context, userId int) ([]models.ProblemInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, maxQuerySeconds*time.Second)
	defer cancel()

	query := `
		SELECT
			p.id,
			p.title,
			p.acceptance_rate,
			d.id,
			d.name,
			EXISTS (
				SELECT 1
				FROM code_submissions cs
				JOIN status_ids s ON s.id = cs.status_id
				WHERE cs.problem_id = p.id AND cs.user_id = $1 AND s.name = $2
			) AS is_solved
		FROM problems p
		LEFT JOIN difficulties d ON p.difficulty_id = d.id
	`

	rows, err := r.db.QueryContext(ctx, query, userId, acceptedStatus)
	if err != nil {
		return nil, fmt.Errorf("querying problems: %w", err)
	}
	defer rows.Close()

	var problems []models.ProblemInfo
	for rows.Next() {
		var p models.ProblemInfo
		var difficulty models.Difficulty
		accRate := 0.0
		err := rows.Scan(
			&p.ID, &p.Title, &accRate,
			&difficulty.ID, &difficulty.Name, &p.IsSolved,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning problem row: %w", err)
		}

		p.AcceptanceRate = int(accRate * 100)
		p.Difficulty = difficulty
		p.Tags = []models.Tag{} // not loading tags here for performance
		problems = append(problems, p)
	}

	return problems, nil
}

// ViewProblem returns detailed info about a specific problem for the user.
func (r *problemRepo) ViewProblem(ctx context.Context, userId, problemId int) (*models.ProblemDetail, error) {
	ctx, cancel := context.WithTimeout(ctx, maxQuerySeconds*time.Second)
	defer cancel()

	// Fetch main problem data
	query := `
		SELECT 
			p.id, p.title, p.description, p.acceptance_rate, p.constraints,
			d.id, d.name,
			EXISTS (
				SELECT 1
				FROM code_submissions cs
				JOIN status_ids s ON s.id = cs.status_id
				WHERE cs.problem_id = p.id AND cs.user_id = $1 AND s.name = $2
			) AS is_solved
		FROM problems p
		LEFT JOIN difficulties d ON p.difficulty_id = d.id
		WHERE p.id = $3
	`

	var detail models.ProblemDetail
	var difficulty models.Difficulty

	accRate := 0.0

	err := r.db.QueryRowContext(ctx, query, userId, acceptedStatus, problemId).Scan(
		&detail.ID, &detail.Title, &detail.Description,
		&accRate, &detail.Constraints,
		&difficulty.ID, &difficulty.Name, &detail.IsSolved,
	)
	// if err == sql.ErrNoRows {
	// 	return nil, nil
	// }
	if err != nil {
		return nil, fmt.Errorf("fetching problem detail: %w", err)
	}

	detail.AcceptanceRate = int(accRate * 100)

	detail.Difficulty = difficulty

	// Fetch tags
	tagsQuery := `
		SELECT t.id, t.name
		FROM tags t
		JOIN problem_tags pt ON pt.tag_id = t.id
		WHERE pt.problem_id = $1
	`
	tagRows, err := r.db.QueryContext(ctx, tagsQuery, problemId)
	if err != nil {
		return nil, fmt.Errorf("fetching tags: %w", err)
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var tag models.Tag
		if err := tagRows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, fmt.Errorf("scanning tag: %w", err)
		}
		detail.Tags = append(detail.Tags, tag)
	}

	// Fetch examples
	exampleQuery := `
		SELECT id, input, expected_output, explanation
		FROM problem_examples
		WHERE problem_id = $1
	`
	exampleRows, err := r.db.QueryContext(ctx, exampleQuery, problemId)
	if err != nil {
		return nil, fmt.Errorf("fetching examples: %w", err)
	}
	defer exampleRows.Close()

	for exampleRows.Next() {
		var ex models.ProblemExample
		if err := exampleRows.Scan(&ex.ID, &ex.Input, &ex.Output, &ex.Explanation); err != nil {
			return nil, fmt.Errorf("scanning example: %w", err)
		}
		detail.Examples = append(detail.Examples, ex)
	}

	return &detail, nil
}
