package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

const (
	maxQueryTime = time.Second * 5
)

type serviceImpl struct {
	db                *sql.DB
	submissions       []Submission
	contestSubmission []ContestSolvedProblems
	redis             *RedisService
}

func NewService(db *sql.DB, redis *RedisService) *serviceImpl {
	return &serviceImpl{db: db, redis: redis, submissions: make([]Submission, 0)}
}

func (s *serviceImpl) Register(ctx context.Context, username, email, password string) (int, string, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return 0, "", fmt.Errorf("failed to hash the password: %w", err)
	}

	query := `INSERT INTO users (username, email, hashed_password) VALUES ($1, $2, $3) returning id, role;`

	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	var userID int
	var role string
	err = s.db.QueryRowContext(ctx, query, username, email, hashedPassword).Scan(&userID, &role)
	if err != nil {
		return 0, "", fmt.Errorf("failed to register user: %w", err)
	}

	return userID, role, nil
}

func (s *serviceImpl) Login(ctx context.Context, username, password string) (int, string, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return 0, "", fmt.Errorf("failed to hash the password: %w", err)
	}

	query := `Select id, role FROM users WHERE username=$1 and hashed_password=$2;`

	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	var userID int
	var role string

	log.Println("\n\n Query: ", query)
	log.Println("\n Params: ", username, hashedPassword)
	err = s.db.QueryRowContext(ctx, query, username, hashedPassword).Scan(&userID, &role)
	if err != nil {
		return 0, "", fmt.Errorf("failed to authenticate user: %w", err)
	}

	return userID, role, nil
}

func (s *serviceImpl) GetUserByID(ctx context.Context, userID int) (*User, error) {
	const query = `
		SELECT id, username, email, role, rating
		FROM users WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	var user User
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role, &user.Rating,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Fetch solved problems
	const solvedQuery = `
		SELECT p.id, p.title, p.difficulty
		FROM solved_problems sp
		JOIN problems p ON sp.problem_id = p.id
		WHERE sp.user_id = $1;
	`
	rows, err := s.db.QueryContext(ctx, solvedQuery, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch solved problems: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pi ProblemInfo
		err := rows.Scan(&pi.ID, &pi.Title, &pi.Difficulty)
		if err != nil {
			return nil, fmt.Errorf("failed to scan solved problem: %w", err)
		}
		user.SolvedProblems = append(user.SolvedProblems, pi)
	}

	return &user, nil

}

func (s *serviceImpl) GetUserProfile(ctx context.Context, username string) (*User, error) {
	const query = `
		SELECT id, username, email, role, rating
		FROM users WHERE username = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	var user User
	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role, &user.Rating,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Fetch solved problems
	const solvedQuery = `
		SELECT p.id, p.title, p.difficulty
		FROM solved_problems sp
		JOIN problems p ON sp.problem_id = p.id
		WHERE sp.user_id = $1;
	`
	rows, err := s.db.QueryContext(ctx, solvedQuery, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch solved problems: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pi ProblemInfo
		err := rows.Scan(&pi.ID, &pi.Title, &pi.Difficulty)
		if err != nil {
			return nil, fmt.Errorf("failed to scan solved problem: %w", err)
		}
		user.SolvedProblems = append(user.SolvedProblems, pi)
	}

	return &user, nil

}

func (s *serviceImpl) GetProblems(ctx context.Context) ([]ProblemInfo, error) {
	const query = `
		SELECT id, title, difficulty, slug
		FROM problems
		WHERE status = 'active';
	`

	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch problems: %w", err)
	}
	defer rows.Close()

	var problems []ProblemInfo
	for rows.Next() {
		var p ProblemInfo
		if err := rows.Scan(&p.ID, &p.Title, &p.Difficulty, &p.Slug); err != nil {
			return nil, err
		}

		// Fetch tags
		tagRows, err := s.db.QueryContext(ctx, `SELECT tag FROM problem_tags WHERE problem_id = $1`, p.ID)
		if err == nil {
			for tagRows.Next() {
				var tag string
				tagRows.Scan(&tag)
				p.Tags = append(p.Tags, tag)
			}
			tagRows.Close()
		}

		problems = append(problems, p)
	}
	return problems, nil

}

func (s *serviceImpl) AddProblem(ctx context.Context, problem *ProblemDetail) (int, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert main problem entry
	const insertProblem = `
		INSERT INTO problems (title, description, constraints, difficulty, author_id, status, solution_language, solution_code, failure_reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id;
	`

	constraints := strings.Join(problem.Constraints, "\n")

	var problemID int
	err = tx.QueryRowContext(ctx, insertProblem,
		problem.Title, problem.Description, constraints, problem.Difficulty,
		problem.AuthorID, problem.Status, problem.SolutionLanguage, problem.SolutionCode,
		problem.FailureReason,
	).Scan(&problemID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert problem: %w", err)
	}

	// Insert tags
	for _, tag := range problem.Tags {
		_, err := tx.ExecContext(ctx, `INSERT INTO problem_tags (problem_id, tag) VALUES ($1, $2)`, problemID, tag)
		if err != nil {
			return 0, fmt.Errorf("failed to insert tag '%s': %w", tag, err)
		}
	}

	// Insert examples
	for _, ex := range problem.TestCases {
		_, err := tx.ExecContext(ctx, `INSERT INTO test_cases (problem_id, input, expected_output) VALUES ($1, $2, $3)`,
			problemID, ex.Input, ex.ExpectedOutput)
		if err != nil {
			return 0, fmt.Errorf("failed to insert test case: %w", err)
		}
	}

	// Insert limits
	for _, limit := range problem.Limits {
		_, err := tx.ExecContext(ctx, `INSERT INTO limits (problem_id, language, time_limit_ms, memory_limit_kb) VALUES ($1, $2, $3, $4)`,
			problemID, limit.Language, limit.TimeLimitMS, limit.MemoryLimitKB)
		if err != nil {
			return 0, fmt.Errorf("failed to insert limit: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return problemID, nil
}

func (s *serviceImpl) UpdateProblemByID(ctx context.Context, id int, problem *ProblemDetail) error {
	const query = `
		UPDATE problems
		SET title = $1,
			description = $2,
			constraints = $3,
			status = $4,
			solution_language = $5,
			solution_code = $6,
			failure_reason = $7
		WHERE id = $8;
	`

	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	constraints := strings.Join(problem.Constraints, "\n")

	_, err := s.db.ExecContext(ctx, query,
		problem.Title, problem.Description, constraints,
		problem.Status, problem.SolutionLanguage, problem.SolutionCode,
		problem.FailureReason, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update problem: %w", err)
	}
	return nil
}

func (s *serviceImpl) GetProblemBySlug(ctx context.Context, slug string) (*ProblemDetail, error) {
	const problemQuery = `
		SELECT id, title, description, constraints, difficulty, author_id, status, failure_reason
		FROM problems WHERE slug = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	var pd ProblemDetail
	var constraints *string
	err := s.db.QueryRowContext(ctx, problemQuery, slug).Scan(
		&pd.ID, &pd.Title, &pd.Description, &constraints, &pd.Difficulty,
		&pd.AuthorID, &pd.Status, &pd.FailureReason,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get problem by slug: %w", err)
	}
	if constraints != nil {
		pd.Constraints = strings.Split(*constraints, "\n")
	}

	// Tags
	tagRows, err := s.db.QueryContext(ctx, `SELECT tag FROM problem_tags WHERE problem_id = $1`, pd.ID)
	if err == nil {
		defer tagRows.Close()
		for tagRows.Next() {
			var tag string
			tagRows.Scan(&tag)
			pd.Tags = append(pd.Tags, tag)
		}
	}

	// Test Cases
	tcRows, err := s.db.QueryContext(ctx, `SELECT id, input, expected_output FROM test_cases WHERE problem_id = $1`, pd.ID)
	if err == nil {
		defer tcRows.Close()
		for tcRows.Next() {
			var tc TestCase
			tcRows.Scan(&tc.ID, &tc.Input, &tc.ExpectedOutput)
			pd.TestCases = append(pd.TestCases, tc)
		}
	}

	// Limits
	limitRows, err := s.db.QueryContext(ctx, `SELECT language, time_limit_ms, memory_limit_kb FROM limits WHERE problem_id = $1`, pd.ID)
	if err == nil {
		defer limitRows.Close()
		for limitRows.Next() {
			var l Limits
			l.ProblemID = pd.ID
			limitRows.Scan(&l.Language, &l.TimeLimitMS, &l.MemoryLimitKB)
			pd.Limits = append(pd.Limits, l)
		}
	}

	return &pd, nil
}

func (s *serviceImpl) RunCode(ctx context.Context, userID, problemID int, language Language, code string, testCases []TestCase) (int, error) {
	const limitsQuery = `SELECT time_limit_ms, memory_limit_kb FROM limits WHERE problem_id=$1 and language=$2;`
	// const insertSubmission = `
	// 	INSERT INTO submissions (user_id, problem_id, language, code, status, message)
	// 	VALUES ($1, $2, $3, $4, 'pending', '')
	// 	RETURNING id;
	// `

	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	submissionID := 0
	// err := s.db.QueryRowContext(ctx, insertSubmission, userID, problemID, language, code).Scan(&submissionID)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to insert submission: %w", err)
	// }

	// TODO: Add to the database
	submissionID = len(s.submissions) + 1
	s.submissions = append(s.submissions, Submission{
		ID:        submissionID,
		ProblemID: &problemID,
		UserID:    userID,
		ContestID: nil,
		Language:  language,
		Code:      code,
		Status:    "pending",
		Message:   "",
		Results:   nil,
	})

	// Fetch time/memory limits from DB
	var timeLimitMS, memoryLimitKB int
	err := s.db.QueryRowContext(ctx, limitsQuery, problemID, language).Scan(&timeLimitMS, &memoryLimitKB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			timeLimitMS = defaultTimeLimitMS
			memoryLimitKB = defaultMemoryLimitKB
		} else {
			return 0, fmt.Errorf("failed to fetch limits: %w", err)
		}
	}

	// Cap the values to maximum allowed
	if timeLimitMS > maxTimeLimitMS {
		timeLimitMS = maxTimeLimitMS
	}
	if memoryLimitKB > maxMemoryLimitKB {
		memoryLimitKB = maxMemoryLimitKB
	}

	payload := ExecutionPayload{
		ID:            submissionID,
		Language:      language,
		Code:          code,
		TestCases:     testCases,
		TimeLimitMS:   timeLimitMS,
		MemoryLimitKB: memoryLimitKB,
		ExecutionType: EXECUTION_RUN,
		ContestID:     0,
		ProblemID:     0,
	}

	// Send to executor/queue here if needed
	s.redis.ExecuteCode(ctx, payload)

	return submissionID, nil
}

func (s *serviceImpl) GetRunResult(ctx context.Context, runID int) (Submission, error) {
	// const submissionQuery = `
	// 	SELECT id, user_id, problem_id, contest_id, language, code, status, message
	// 	FROM submissions WHERE id = $1;
	// `

	// var sub Submission
	// ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	// defer cancel()

	// err := s.db.QueryRowContext(ctx, submissionQuery, runID).Scan(
	// 	&sub.ID, &sub.UserID, &sub.ProblemID, &sub.ContestID,
	// 	&sub.Language, &sub.Code, &sub.Status, &sub.Message,
	// )
	// if err != nil {
	// 	return Submission{}, fmt.Errorf("submission not found: %w", err)
	// }

	// // Get test results
	// const resultsQuery = `
	// 	SELECT id, status, stdout, stderr, runtime_ms, memory_kb
	// 	FROM test_results WHERE submission_id = $1;
	// `

	// rows, err := s.db.QueryContext(ctx, resultsQuery, runID)
	// if err != nil {
	// 	return Submission{}, fmt.Errorf("failed to fetch test results: %w", err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var res TestResult
	// 	err := rows.Scan(&res.ID, &res.Status, &res.Output, &res.ExpectedOutput, &res.RuntimeMS, &res.MemoryKB)
	// 	if err != nil {
	// 		return Submission{}, fmt.Errorf("error scanning test result: %w", err)
	// 	}
	// 	sub.Results = append(sub.Results, res)
	// }

	// TODO: Fetch from the database
	sub := s.submissions[runID-1]

	return sub, nil
}

func (s *serviceImpl) SubmitCode(ctx context.Context, userID, problemID, contestID int, language Language, code string) (int, error) {
	// const insertSubmission = `
	// 	INSERT INTO submissions (user_id, problem_id, language, code, status, message)
	// 	VALUES ($1, $2, $3, $4, 'pending', '')
	// 	RETURNING id;
	// `
	const getTestCases = `SELECT id, input, expected_output FROM test_cases WHERE problem_id = $1;`
	const getLimits = `SELECT time_limit_ms, memory_limit_kb FROM limits WHERE problem_id = $1 AND language = $2;`

	const (
		defaultTimeLimitMS   = 2000
		defaultMemoryLimitKB = 65536
		maxTimeLimitMS       = 5000
		maxMemoryLimitKB     = 131072
	)

	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	// Insert the submission
	submissionID := 0
	// err := s.db.QueryRowContext(ctx, insertSubmission, userID, problemID, language, code).Scan(&submissionID)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to insert submission: %w", err)
	// }

	// TODO: Add to the database
	submissionID = len(s.submissions) + 1
	s.submissions = append(s.submissions, Submission{
		ID:        submissionID,
		ProblemID: &problemID,
		UserID:    userID,
		ContestID: &contestID,
		Language:  language,
		Code:      code,
		Status:    "pending",
		Message:   "",
		Results:   nil,
	})

	// Fetch test cases
	rows, err := s.db.QueryContext(ctx, getTestCases, problemID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch test cases: %w", err)
	}
	defer rows.Close()

	var testCases []TestCase
	for rows.Next() {
		var tc TestCase
		if err := rows.Scan(&tc.ID, &tc.Input, &tc.ExpectedOutput); err != nil {
			return 0, fmt.Errorf("failed to scan test case: %w", err)
		}
		testCases = append(testCases, tc)
	}

	// Fetch limits
	var timeLimitMS, memoryLimitKB int
	err = s.db.QueryRowContext(ctx, getLimits, problemID, language).Scan(&timeLimitMS, &memoryLimitKB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			timeLimitMS = defaultTimeLimitMS
			memoryLimitKB = defaultMemoryLimitKB
		} else {
			return 0, fmt.Errorf("failed to fetch resource limits: %w", err)
		}
	}

	// Enforce max caps
	if timeLimitMS > maxTimeLimitMS {
		timeLimitMS = maxTimeLimitMS
	}
	if memoryLimitKB > maxMemoryLimitKB {
		memoryLimitKB = maxMemoryLimitKB
	}

	// Prepare payload
	payload := ExecutionPayload{
		ID:            submissionID,
		Language:      language,
		Code:          code,
		TestCases:     testCases,
		TimeLimitMS:   timeLimitMS,
		MemoryLimitKB: memoryLimitKB,
		ExecutionType: EXECUTION_SUBMIT,
		ContestID:     contestID,
		ProblemID:     problemID,
	}

	// Send for execution
	s.redis.ExecuteCode(ctx, payload)

	return submissionID, nil
}

func (s *serviceImpl) GetSubmissionResult(ctx context.Context, runID int) (Submission, error) {
	return s.GetRunResult(ctx, runID)
}

func (s *serviceImpl) GetUserSubmissions(ctx context.Context, userID, problemID int) ([]Submission, error) {
	const query = `
		SELECT id, user_id, problem_id, contest_id, language, code, status, message
		FROM submissions
		WHERE user_id = $1 AND problem_id = $2
		ORDER BY id DESC;
	`

	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch submissions: %w", err)
	}
	defer rows.Close()

	var subs []Submission
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.ID, &s.UserID, &s.ProblemID, &s.ContestID, &s.Language, &s.Code, &s.Status, &s.Message)
		if err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}

func (s *serviceImpl) CreateContest(ctx context.Context, contest *Contest) (int, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	const insertContest = `
		INSERT INTO contests (name, status, start_time, end_time)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var contestID int
	err = tx.QueryRowContext(ctx, insertContest, contest.Name, contest.Status, contest.StartTime, contest.EndTime).Scan(&contestID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert contest: %w", err)
	}

	const insertProblem = `
		INSERT INTO contest_problems (contest_id, problem_id, max_points)
		VALUES ($1, $2, $3);
	`

	for _, cp := range contest.Problems {
		_, err := tx.ExecContext(ctx, insertProblem, contestID, cp.ID, cp.MaxPoints)
		if err != nil {
			return 0, fmt.Errorf("failed to insert contest problem: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit contest creation: %w", err)
	}

	return contestID, nil
}

func (s *serviceImpl) UpdateContest(ctx context.Context, id int, contest *Contest) error {
	const updateQuery = `
		UPDATE contests
		SET name = $1, status = $2, start_time = $3, end_time = $4
		WHERE id = $5;
	`

	_, err := s.db.ExecContext(ctx, updateQuery, contest.Name, contest.Status, contest.StartTime, contest.EndTime, id)
	if err != nil {
		return fmt.Errorf("failed to update contest: %w", err)
	}
	return nil
}

func (s *serviceImpl) GetAllContests(ctx context.Context) ([]Contest, error) {
	const query = `
		SELECT id, name, status, start_time, end_time
		FROM contests ORDER BY start_time DESC;
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contests: %w", err)
	}
	defer rows.Close()

	var contests []Contest
	for rows.Next() {
		var c Contest
		err := rows.Scan(&c.ID, &c.Name, &c.Status, &c.StartTime, &c.EndTime)
		if err != nil {
			return nil, err
		}
		contests = append(contests, c)
	}
	return contests, nil
}

func (s *serviceImpl) GetContestByID(ctx context.Context, contestID int) (*Contest, error) {
	const baseQuery = `
		SELECT id, name, status, start_time, end_time
		FROM contests WHERE id = $1;
	`

	var c Contest
	err := s.db.QueryRowContext(ctx, baseQuery, contestID).Scan(&c.ID, &c.Name, &c.Status, &c.StartTime, &c.EndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get contest: %w", err)
	}

	// Load problems
	const problemQuery = `
		SELECT p.id, p.title, p.difficulty, cp.max_points, p.slug
		FROM contest_problems cp
		JOIN problems p ON cp.problem_id = p.id
		WHERE cp.contest_id = $1;
	`

	rows, err := s.db.QueryContext(ctx, problemQuery, contestID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contest problems: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cp ContestProblem
		var pi ProblemInfo
		if err := rows.Scan(&pi.ID, &pi.Title, &pi.Difficulty, &cp.MaxPoints, &pi.Slug); err != nil {
			return nil, err
		}
		cp.ProblemInfo = &pi
		c.Problems = append(c.Problems, cp)
	}

	return &c, nil
}

func (s *serviceImpl) JoinContestByID(ctx context.Context, userID, contestID int) error {
	const query = `
		INSERT INTO contest_participants (contest_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`

	_, err := s.db.ExecContext(ctx, query, contestID, userID)
	if err != nil {
		return fmt.Errorf("failed to join contest: %w", err)
	}
	return nil
}

// func (s *serviceImpl) StartContest(ctx context.Context, contestID int) error { return nil } // Add contest details to cache

// func (s *serviceImpl) EndContest(ctx context.Context, contestID int) error { return nil } // Remove contest from thhe cache

func (s *serviceImpl) GetLeaderboard(ctx context.Context, contestID int) ([]ContestParticipant, error) {
	// const query = `
	// 	SELECT u.id, u.username, cp.score, cp.rating_change
	// 	FROM contest_participants cp
	// 	JOIN users u ON cp.user_id = u.id
	// 	WHERE cp.contest_id = $1
	// 	ORDER BY cp.score DESC;
	// `

	// rows, err := s.db.QueryContext(ctx, query, contestID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	// }
	// defer rows.Close()

	// var leaderboard []ContestParticipant
	// for rows.Next() {
	// 	var p ContestParticipant
	// 	if err := rows.Scan(&p.UserID, &p.Username, &p.Score, &p.RatingChange); err != nil {
	// 		return nil, err
	// 	}
	// 	leaderboard = append(leaderboard, p)
	// }

	// TODO: Use Database

	// Step 1: Initialize a map to store participants and their contest data
	participantMap := make(map[int]*ContestParticipant)

	// Step 2: Loop through all contest submissions and populate participant data
	for _, contestSolved := range s.contestSubmission {
		if contestSolved.ContestID == contestID {
			// Check if the participant already exists in the map
			participant, exists := participantMap[contestSolved.UserID]
			if !exists {
				// If not, create a new participant entry
				participant = &ContestParticipant{
					UserID:         contestSolved.UserID,
					Score:          0,
					ProblemsSolved: []ContestProblem{},
					RatingChange:   0, // Assuming rating change can be tracked separately
				}
				participantMap[contestSolved.UserID] = participant
			}

			// Step 3: Update the participant's score and problems solved
			// Assuming each ContestSolvedProblem has a ScoreDelta that represents the points for the problem
			participant.Score += contestSolved.ScoreDelta

			// Create a ContestProblem object that contains problem details (assuming you have such a struct)
			contestProblem := ContestProblem{
				ProblemInfo: &ProblemInfo{ID: contestSolved.ProblemID, Title: "<Placeholder>",
					Tags: []string{"Placeholder"}, Difficulty: DIFFICULTY_EASY, Slug: fmt.Sprintf("placehodler %d", contestSolved.ProblemID)},
			}

			// Add the problem to the list of solved problems
			participant.ProblemsSolved = append(participant.ProblemsSolved, contestProblem)
		}
	}

	// Step 4: Convert the map to a slice for sorting
	leaderboard := make([]ContestParticipant, 0, len(participantMap))
	for _, participant := range participantMap {
		leaderboard = append(leaderboard, *participant)
	}

	// Step 5: Sort the leaderboard by score in descending order
	sort.SliceStable(leaderboard, func(i, j int) bool {
		return leaderboard[i].Score > leaderboard[j].Score
	})

	return leaderboard, nil
}

func (s *serviceImpl) CreateDiscussion(ctx context.Context, discussion *Discussion) (int, error) {
	const query = `
		INSERT INTO discussions (problem_id, title, content, author_id, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	var discussionID int
	err := s.db.QueryRowContext(ctx, query,
		discussion.ID, discussion.Title, discussion.Content,
		discussion.AuthorID, discussion.IsActive,
	).Scan(&discussionID)
	if err != nil {
		return 0, fmt.Errorf("failed to create discussion: %w", err)
	}

	return discussionID, nil

}

func (s *serviceImpl) UpdateDiscussion(ctx context.Context, discussion *Discussion) error {
	const query = `
		UPDATE discussions
		SET title = $1, content = $2, is_active = $3
		WHERE id = $4;
	`

	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query,
		discussion.Title, discussion.Content, discussion.IsActive, discussion.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update discussion: %w", err)
	}
	return nil
}

func (s *serviceImpl) GetDiscussionByID(ctx context.Context, discussionID int) (*Discussion, error) {
	const baseQuery = `
		SELECT id, problem_id, title, content, author_id, is_active
		FROM discussions WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	var d Discussion
	err := s.db.QueryRowContext(ctx, baseQuery, discussionID).Scan(
		&d.ID, new(int), &d.Title, &d.Content, &d.AuthorID, &d.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get discussion: %w", err)
	}

	// Count votes
	err = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(vote), 0) FROM discussion_votes WHERE discussion_id = $1`, discussionID).Scan(&d.Votes)
	if err != nil {
		return nil, fmt.Errorf("failed to get vote count: %w", err)
	}

	// Get comments
	rows, err := s.db.QueryContext(ctx, `
		SELECT content, author_id FROM discussion_comments
		WHERE discussion_id = $1 ORDER BY created_at ASC;
	`, discussionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment DiscussionComment
		if err := rows.Scan(&comment.Content, &comment.AuthorID); err != nil {
			return nil, err
		}
		d.Comments = append(d.Comments, comment)
	}

	return &d, nil
}

func (s *serviceImpl) AddVoteToDiscussion(ctx context.Context, userID, discussionID int, vote Vote) error {
	const query = `
	INSERT INTO discussion_votes (user_id, discussion_id, vote)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id, discussion_id) DO UPDATE SET vote = EXCLUDED.vote;
	`

	_, err := s.db.ExecContext(ctx, query, userID, discussionID, int(vote))
	if err != nil {
		return fmt.Errorf("failed to vote on discussion: %w", err)
	}
	return nil
}

func (s *serviceImpl) UpdateSubmission(ctx context.Context, submission *Submission) error {
	// // Define the SQL queries for updating the submission and inserting test results
	// const submissionQuery = `
	// 	UPDATE submissions
	// 	SET status = $1, message = $2
	// 	WHERE id = $3
	// 	RETURNING id;
	// `

	// // Start a transaction to ensure atomicity
	// tx, err := s.db.BeginTx(ctx, nil)
	// if err != nil {
	// 	return fmt.Errorf("failed to begin transaction: %w", err)
	// }
	// defer tx.Rollback()

	// // Update the submission status and message
	// var updatedID int
	// err = tx.QueryRowContext(ctx, submissionQuery, submission.Status, submission.Message, submission.ID).Scan(&updatedID)
	// if err != nil {
	// 	return fmt.Errorf("failed to update submission: %w", err)
	// }

	// // Ensure that the updated submission ID matches the provided ID
	// if updatedID != submission.ID {
	// 	return fmt.Errorf("submission ID mismatch during update")
	// }

	// // Insert or update the test results for the submission
	// const testResultQuery = `
	// 	INSERT INTO test_results (submission_id, status, stdout, stderr, runtime_ms, memory_kb)
	// 	VALUES ($1, $2, $3, $4, $5, $6)
	// 	ON CONFLICT (submission_id, id)
	// 	DO UPDATE SET status = $2, stdout = $3, stderr = $4, runtime_ms = $5, memory_kb = $6;
	// `

	// // Iterate through the test results and insert/update them
	// for _, result := range submission.Results {
	// 	_, err := tx.ExecContext(ctx, testResultQuery, submission.ID, result.Status, result.Output, result.ExpectedOutput, result.RuntimeMS, result.MemoryKB)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to update test result: %w", err)
	// 	}
	// }

	// // Commit the transaction
	// if err := tx.Commit(); err != nil {
	// 	return fmt.Errorf("failed to commit transaction: %w", err)
	// }

	// TODO: Update into the database
	// Step 1: Validate the submission ID
	if submission.ID > len(s.submissions) || submission.ID <= 0 {
		return errors.New("invalid submission")
	}

	// Step 2: Update the submission details
	sub := &s.submissions[submission.ID-1]
	sub.Message = submission.Message
	sub.Results = submission.Results
	sub.Status = submission.Status

	// Step 3: Handle contest-specific logic
	if submission.ContestID != nil && *submission.ContestID > 0 {
		// Initialize points structure to hold cache data
		points := CachePoints{Points: 0}

		// Step 4: Attempt to fetch the contest problem data from Redis
		err := s.redis.Get(ctx, GetContestProblemKey(*submission.ContestID, *submission.ProblemID), &points)
		if err != nil {
			// Cache lookup failed (either contest is over or invalid submission)
			// We do nothing here as per the provided logic
			return nil
		}

		// Step 5: Check if the contest problem has already been solved by this user
		// Check if the contest and problem combination exists for this user
		for _, contestSolved := range s.contestSubmission {
			if contestSolved.UserID == submission.UserID && contestSolved.ContestID == *submission.ContestID && contestSolved.ProblemID == *submission.ProblemID {
				// This contest problem has already been solved by the user, no need to add it again
				return nil
			}
		}

		// Step 6: If not already solved, add a new ContestSolvedProblems entry
		newContestSolved := ContestSolvedProblems{
			ContestID:  *submission.ContestID,
			UserID:     submission.UserID,
			ProblemID:  *submission.ProblemID,
			SolvedAt:   int(time.Now().Unix()), // Record the timestamp of when the problem was solved
			ScoreDelta: points.Points,          // Using the points fetched from the cache
		}

		// Add to the contestSubmission list
		s.contestSubmission = append(s.contestSubmission, newContestSolved)
	}

	return nil
}
