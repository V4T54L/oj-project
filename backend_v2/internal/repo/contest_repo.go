package repo

import (
	"context"
	"errors"
	"online-judge/internal/models"
	"sort"
	"sync"
)

type InMemoryContestRepo struct {
	mu                 sync.RWMutex
	contests           map[int]models.ContestDetail
	contestProblems    map[int][]models.ContestProblem
	contestSubmissions map[int][]models.ContestSubmission
	userScores         map[int]map[int]models.ContestRank // contestID -> userID -> ContestRank
	idCounter          int
}

func NewInMemoryContestRepo() *InMemoryContestRepo {
	return &InMemoryContestRepo{
		contests:           make(map[int]models.ContestDetail),
		contestProblems:    make(map[int][]models.ContestProblem),
		contestSubmissions: make(map[int][]models.ContestSubmission),
		userScores:         make(map[int]map[int]models.ContestRank),
		idCounter:          1,
	}
}

// CreateContest creates a new contest.
func (r *InMemoryContestRepo) CreateContest(ctx context.Context, contest *models.ContestDetail) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	contest.ID = r.idCounter
	r.idCounter++

	r.contests[contest.ID] = *contest
	return contest.ID, nil
}

// AddProblemToContest associates a problem with a contest.
func (r *InMemoryContestRepo) AddProblemToContest(ctx context.Context, contestID, problemID, points int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.contests[contestID]; !exists {
		return errors.New("contest not found")
	}

	r.contestProblems[contestID] = append(r.contestProblems[contestID], models.ContestProblem{
		ID: problemID,
		ProblemInfo: models.ProblemInfo{
			ID: problemID, // Assuming ProblemInfo includes an ID
		},
		Points: points,
	})

	return nil
}

// GetContestByID retrieves full contest details including problems.
func (r *InMemoryContestRepo) GetContestByID(ctx context.Context, contestID int) (*models.ContestDetail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	contest, exists := r.contests[contestID]
	if !exists {
		return nil, errors.New("contest not found")
	}

	return &contest, nil
}

// ListContests returns all contests.
func (r *InMemoryContestRepo) ListContests(ctx context.Context) ([]models.ContestDetail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []models.ContestDetail
	for _, c := range r.contests {
		result = append(result, c)
	}
	return result, nil
}

// AddContestSubmission logs a submission and updates the leaderboard.
func (r *InMemoryContestRepo) AddContestSubmission(ctx context.Context, submission models.ContestSubmission) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Append submission
	r.contestSubmissions[submission.ContestID] = append(r.contestSubmissions[submission.ContestID], submission)

	// Initialize user rank info if not present
	if _, ok := r.userScores[submission.ContestID]; !ok {
		r.userScores[submission.ContestID] = make(map[int]models.ContestRank)
	}

	// userRank := r.userScores[submission.ContestID][submission.UserID]

	// TODO: Implement this
	// // Update the problem status
	// for i, problem := range r.contestProblems[submission.ContestID] {
	// 	if problem.ID == submission.ProblemID {
	// 		status := "not_attempted"
	// 		if submission.Score > 0 {
	// 			status = "solved"
	// 		} else {
	// 			status = "attempted"
	// 		}

	// 		// Update attempt count and time taken (if available)
	// 		problemStatus := models.ContestRank{
	// 			UserInfo: models.UserInfo{
	// 				ID: userRank.UserInfo.ID,
	// 			},
	// 			ProblemStatuses: []struct {
	// 				ProblemID    int
	// 				Status       string
	// 				AttemptCount int
	// 				TimeTaken    *int
	// 				Score        int
	// 			}{
	// 				{
	// 					ProblemID:    problem.ID,
	// 					Status:       status,
	// 					AttemptCount: 1,   // You can modify this logic to handle multiple attempts
	// 					TimeTaken:    nil, // You can modify this to record time if available
	// 					Score:        submission.Score,
	// 				},
	// 			},
	// 		}

	// 		// Store the rank and score info
	// 		r.userScores[submission.ContestID][submission.UserID] = problemStatus

	// 		break
	// 	}
	// }

	return nil
}

// GetLeaderboard returns users sorted by score for a contest.
func (r *InMemoryContestRepo) GetLeaderboard(ctx context.Context, contestID int) ([]models.ContestRank, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	scores, ok := r.userScores[contestID]
	if !ok {
		return nil, nil
	}

	var leaderboard []models.ContestRank
	for _, rank := range scores {
		leaderboard = append(leaderboard, rank)
	}

	// Sorting leaderboard by total score (descending order)
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Score > leaderboard[j].Score
	})

	return leaderboard, nil
}
