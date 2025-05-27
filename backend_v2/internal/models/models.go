package models

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
	Explaination       string            `json:"explaination"`
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
	ID             int               `json:"id"`
	LanguageID     int               `json:"language_id"`
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
