package models

import "time"

type basic struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Tag basic
type Difficulty basic
type ProgrammingLanguage basic

type ProblemExamples struct {
	ID             int    `json:"id"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	Explanation    string `json:"explanation"`
}

type ProblemTestCase struct {
	ID             int    `json:"id"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
}

type ProblemInfo struct {
	ID             int        `json:"id"`
	Title          string     `json:"title"`
	Slug           string     `json:"slug"`
	Tags           []Tag      `json:"tags"`
	Difficulty     Difficulty `json:"difficulty"`
	AcceptanceRate float32    `json:"acceptance_rate"`
	IsSolved       *bool      `json:"is_solved"`
}

type ProblemDetail struct {
	ID             int               `json:"id"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Constraints    []string          `json:"constraints"`
	Slug           string            `json:"slug"`
	Tags           []Tag             `json:"tags"`
	Difficulty     Difficulty        `json:"difficulty"`
	AcceptanceRate float32           `json:"acceptance_rate"`
	IsSolved       *bool             `json:"is_solved"`
	Examples       []ProblemExamples `json:"examples"`
}

type ProblemDB struct {
	ID                 int               `json:"id"`
	Title              string            `json:"title"`
	Description        string            `json:"description"`
	Constraints        []string          `json:"constraints"`
	Slug               string            `json:"slug"`
	Tags               []Tag             `json:"tags"`
	Difficulty         Difficulty        `json:"difficulty"`
	AcceptanceRate     float32           `json:"acceptance_rate"`
	IsSolved           *bool             `json:"is_solved,omitempty"`
	Examples           []ProblemExamples `json:"examples"`
	SolutionLanguageID int               `json:"solution_language_id"`
	SolutionCode       string            `json:"solution_code"`
	Explanation        string            `json:"explanation"`
	Status             string            `json:"status"` // Active, In Review, Rejected, Inactive
	RuntimeLimitMS     int               `json:"runtime_limit_ms"`
	MemoryLimitKB      int               `json:"memory_limit_kb"`
}

type SubmitCodePayload struct {
	ProblemID  int    `json:"problem_id"`
	LanguageID int    `json:"language_id"`
	Code       string `json:"code"`
}

type SubmitCodeResponse struct {
	SubmissionID int `json:"submission_id"`
}

type TestCaseResult struct {
	ID             int    `json:"id"`
	Input          string `json:"input"`
	Output         string `json:"output"`
	ExpectedOutput string `json:"expected_output"`
	RuntimeMS      int    `json:"runtime_ms"`
	MemoryKB       int    `json:"memory_kb"`
	Status         string `json:"status"` // TLE, MLE, Acccepted, Wrong Answer, $Error.message
}

type ExecuteCodePayload struct {
	ID int `json:"id"`
	// LanguageID     int               `json:"language_id"`
	Code           string            `json:"code"`
	TestCases      []ProblemTestCase `json:"test_cases"`
	RuntimeLimitMS int               `json:"runtime_limit_ms"`
	MemoryLimitKB  int               `json:"memory_limit_kb"`
	ExecutionType  string            `json:"execution_type"` // Run, Submit, Validation
}

type ExecuteCodeResponse struct {
	ID              int              `json:"id"`
	Status          string           `json:"status"` // TLE, MLE, Acccepted, Wrong Answer, $Error.message
	TestCaseResults []TestCaseResult `json:"test_case_results"`
	ExecutionType   string           `json:"execution_type"` // Run, Submit, Validation
}

type SubmissionDB struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	UserID    int    `json:"user_id"`
	ProblemID int    `json:"problem_id"`
	RuntimeMS int    `json:"runtime_ms"`
	MemoryKB  int    `json:"memory_kb"`
}

type UserDB struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
	IsActive bool   `json:"is_active"`
}

type UserInfo struct {
	ID       int
	Username string `json:"username"`
	Email    string `json:"email"`
}

type SignupPayload struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse UserInfo

// type ContestDB struct {
// 	ID        int
// 	Title     string
// 	StartTime time.Time
// 	EndTime   time.Time
// 	IsActive  bool
// }

// type ContestProblemInfo struct {
// 	ProblemID int
// 	MaxScore  int
// }

// type ContestDetail struct {
// 	ID        int                  `json:"id"`
// 	Title     string               `json:"title"`
// 	StartTime time.Time            `json:"start_time"`
// 	EndTime   time.Time            `json:"end_time"`
// 	Problems  []ContestProblemInfo `json:"problems"`
// }

// type ContestInfo struct {
// 	ID        int
// 	Title     string
// 	StartTime time.Time
// 	EndTime   time.Time
// 	Problems  []ContestProblemInfo
// }

type ContestSubmission struct {
	ID           int       `json:"id"`
	SubmissionID int       `json:"submission_id"`
	UserID       int       `json:"user_id"`
	ContestID    int       `json:"contest_id"`
	ProblemID    int       `json:"problem_id"`
	Score        int       `json:"score"`
	SubmittedAt  time.Time `json:"submitted_at"`
}

// type UserScore struct {
// 	UserID     int    `json:"user_id"`
// 	Username   string `json:"username"`
// 	TotalScore int    `json:"total_score"`
// }

type ContestDetail struct {
	ID           int              `json:"id"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	StartTime    time.Time        `json:"start_time"`
	EndTime      time.Time        `json:"end_time"`
	Status       string           `json:"status"` //'upcoming' | 'active' | 'completed'
	CreatedBy    UserInfo         `json:"created_by"`
	Participants []UserInfo       `json:"participants"`
	Problems     []ContestProblem `json:"problems"`
	BannerURI    *string          `json:"banner_uri"`
	Difficulty   string           `json:"difficulty"` //'easy' | 'medium' | 'hard'
	Tags         []Tag            `json:"tags"`
	Prizes       []Prize          `json:"prizes"`
}

type ContestProblem struct {
	ID          int `json:"id"`
	ProblemInfo `json:"problem_info"`
	Points      int `json:"points"`
}

type Prize struct {
	Rank        int     `json:"rank"`
	Description string  `json:"description"`
	Value       *string `json:"value"`
}

type ContestRank struct {
	UserInfo        `json:"user_info"`
	Rank            int `json:"rank"`
	Score           int `json:"score"`
	ProblemStatuses []struct {
		ProblemID    int    `json:"problem_id"`
		Status       string `json:"status"` //'solved' | 'attempted' | 'not_attempted'
		AttemptCount int    `json:"attempt_count"`
		TimeTaken    *int   `json:"time_taken"` // in seconds
		Score        int    `json:"score"`
	} `json:"problem_statuses"`
}

type RunDB struct {
	ID        int              `json:"id"`
	UserID    int              `json:"user_id"`
	ProblemID int              `json:"problem_id"`
	Results   []TestCaseResult `json:"results"`
	Status    string           `json:"status"`
}

type RunPayload struct {
	ID         int               `json:"id"`
	UserID     int               `json:"user_id"`
	ProblemID  int               `json:"problem_id"`
	LanguageID int               `json:"language_id"`
	Code       string            `json:"code"`
	TestCases  []ProblemTestCase `json:"test_cases"`
}

type DiscussionDB struct {
	ID                  int
	ProblemID           int
	CreatedBy           UserInfo
	Title               string
	Content             string
	Tags                []Tag
	Likes               int
	Dislikes            int
	CurrentUserReaction int // -1 = dislike; 1 = like
	Comments            []CommentDB
}

type CommentDB struct {
	ID                  int
	DiscussionID        int
	CreatedBy           UserInfo
	Content             string
	Likes               int
	Dislikes            int
	CurrentUserReaction int // -1 = dislike; 1 = like
}

type NewDiscussionPayload struct {
	ProblemID int
	Title     string
	Content   string
	TagIDs    []int
}

type NewCommentPayload struct {
	DiscussionID int
	Content      string
}

type ReactionPayload struct {
	ID       int
	Reaction int
}
