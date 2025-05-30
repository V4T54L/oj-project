package repo

import (
	"context"
	"errors"
	"online-judge/internal/models"
)

type ProblemRepo struct {
	db        []models.ProblemDB
	testCases []models.ProblemTestCase
}

func NewProblemRepo() *ProblemRepo {
	problems := []models.ProblemDB{
		{
			ID:                 1,
			Title:              "Add Two Numbers",
			Description:        "Write a program that takes two integers and prints their sum.",
			Constraints:        []string{"-10^9 <= a, b <= 10^9"},
			Slug:               "add-two-numbers",
			Tags:               []models.Tag{{ID: 1, Name: "Math"}},
			Difficulty:         models.Difficulty{ID: 1, Name: "Easy"},
			AcceptanceRate:     95.5,
			SolutionLanguageID: 1, // Python
			SolutionCode: `a = int(input())
b = int(input())
print(a + b)`,
			Explanation:    "Read two integers from stdin and print their sum.",
			Status:         "Active",
			RuntimeLimitMS: 1000,
			MemoryLimitKB:  65536,
			Examples: []models.ProblemExamples{
				{
					ID:             1,
					Input:          "2\n3",
					ExpectedOutput: "5",
					Explanation:    "2 + 3 = 5",
				},
			},
		},
	}

	// Example test cases
	testCases := []models.ProblemTestCase{
		{ID: 1, Input: "2\n3", ExpectedOutput: "5"},
		{ID: 2, Input: "-1\n1", ExpectedOutput: "0"},
		{ID: 3, Input: "1000000\n2345678", ExpectedOutput: "3345678"},
	}

	// Initialize the ProblemRepo
	problemRepo := ProblemRepo{
		db:        problems,
		testCases: testCases,
	}

	return &problemRepo
}

// GetProblems returns all problems with status "Active" as []ProblemInfo.
func (r *ProblemRepo) GetProblems(ctx context.Context, isAdmin bool) ([]models.ProblemInfo, error) {
	var result []models.ProblemInfo
	for _, p := range r.db {
		if p.Status == "Active" || isAdmin {
			result = append(result, models.ProblemInfo{
				ID:             p.ID,
				Title:          p.Title,
				Slug:           p.Slug,
				Tags:           p.Tags,
				Difficulty:     p.Difficulty,
				AcceptanceRate: p.AcceptanceRate,
				IsSolved:       p.IsSolved,
			})
		}
	}
	return result, nil
}

// GetProblemByID returns a single ProblemDetail with status "Active".
func (r *ProblemRepo) GetProblemByID(ctx context.Context, problemId int) (*models.ProblemDetail, error) {
	for _, p := range r.db {
		if p.ID == problemId && p.Status == "Active" {
			return &models.ProblemDetail{
				ID:             p.ID,
				Title:          p.Title,
				Description:    p.Description,
				Constraints:    p.Constraints,
				Slug:           p.Slug,
				Tags:           p.Tags,
				Difficulty:     p.Difficulty,
				AcceptanceRate: p.AcceptanceRate,
				IsSolved:       p.IsSolved,
				Examples:       p.Examples,
			}, nil
		}
	}
	return nil, errors.New("active problem not found")
}

// CreateProblem adds a new problem with status "In Review".
func (r *ProblemRepo) CreateProblem(ctx context.Context, problem *models.ProblemDB) (int, error) {
	problem.ID = len(r.db) + 1
	problem.Status = "In Review"
	r.db = append(r.db, *problem)
	return problem.ID, nil
}

// UpdateProblemByID updates a problem by ID, sets status to "In Review".
func (r *ProblemRepo) UpdateProblemByID(ctx context.Context, updated *models.ProblemDB) error {
	for i := range r.db {
		if r.db[i].ID == updated.ID {
			updated.Status = "In Review"
			r.db[i] = *updated
			return nil
		}
	}
	return errors.New("problem not found")
}

// UpdateProblemStatusByID updates only the status of the problem by ID.
func (r *ProblemRepo) UpdateProblemStatusByID(ctx context.Context, problemID int, status string) error {
	for i := range r.db {
		if r.db[i].ID == problemID {
			r.db[i].Status = status
			return nil
		}
	}
	return errors.New("problem not found")
}

func (r *ProblemRepo) GetProblemTestCases(ctx context.Context, problemId int) ([]models.ProblemTestCase, error) {
	if len(r.testCases) == 0 {
		return nil, errors.New("test cases not found")
	}

	return r.testCases, nil
}

func (r *ProblemRepo) GetProblemMetadata(ctx context.Context, problemId int) (*models.ProblemDB, error) {
	for _, p := range r.db {
		if p.ID == problemId && p.Status == "Active" {
			return &p, nil
		}
	}
	return nil, errors.New("active problem not found")
}
