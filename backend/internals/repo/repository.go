package repo

import (
	"algo-arena-be/internals/models"
	"context"
	"database/sql"
	"errors"
)

const maxQuerySeconds = 5

type UserRepo interface {
	GetUserByUserID(ctx context.Context, userId int) (*models.UserInfo, error)
	CreateUser(ctx context.Context, username, email, avatarUrl, hashedPassword string) (int, error)
	GetUserByCreds(ctx context.Context, username, hashedPassword string) (*models.UserInfo, error)
}

type ProblemRepo interface {
	// TODO: Create, Update and Delete problems
	GetProblems(ctx context.Context, userId int) ([]models.ProblemInfo, error)
	ViewProblem(ctx context.Context, userId, problemId int) (*models.ProblemDetail, error)
}

type SubmissionRepo interface {
	GetSubmissionResult(ctx context.Context, userId, submissionId int) (*models.SubmissionResult, error)
	GetRunResult(ctx context.Context, userId, runId int) (*models.RunResult, error)
	RunCode(ctx context.Context, userId, problemId, languageId int, code string, testCases []models.TestCase) (int, error)
	SubmitCode(ctx context.Context, userId, problemId, languageId int, code string) (int, error)
	AddRunResult(ctx context.Context, result *models.RunResult) error
	AddSubmissionResult(ctx context.Context, result *models.RunResult) error
}

func GetRepos(db *sql.DB) (
	UserRepo, ProblemRepo, SubmissionRepo, error,
) {
	if db == nil {
		return nil, nil, nil, errors.New("nil value provided for db connection")
	}

	mockUser := &models.UserInfo{
		ID:        1,
		Username:  "itsDitto",
		AvatarUrl: "https://dash-tech.netlify.app/assets/pfp_m_2-BOA2oL-X.jpg",
	}

	return &userRepo{mockUser: mockUser}, &problemRepo{db: db}, &submissionRepo{db: db, runs: make(map[int]*models.SubmissionResult)}, nil
}
