package models

type basic struct {
	ID   int
	Name string
}

type Tag basic
type Difficulty basic
type ProgrammingLanguage basic

type ProblemExamples struct {
	ID             int
	Input          string
	ExpectedOutput string
	Explaination   string
}

type ProblemTestCase struct {
	ID             int
	Input          string
	ExpectedOutput string
}

type ProblemInfo struct {
	ID             int
	Title          string
	Slug           string
	Tags           []Tag
	Difficulty     Difficulty
	AcceptanceRate float32
	IsSolved       *bool
}

type ProblemDetail struct {
	ID             int
	Title          string
	Description    string
	Slug           string
	Tags           []Tag
	Difficulty     Difficulty
	AcceptanceRate float32
	IsSolved       *bool
	Examples       []ProblemExamples
}

type ProblemDB struct {
	ID                 int
	Title              string
	Description        string
	Slug               string
	Tags               []Tag
	Difficulty         Difficulty
	AcceptanceRate     float32
	IsSolved           *bool
	Examples           []ProblemExamples
	SolutionLanguageID int
	SolutionCode       string
	Explaination       string
	Status             string // Active, In Review, Rejected, Inactive
	RuntimeLimitMS     int
	MemoryLimitKB      int
}

type SubmitCodePayload struct {
	ProblemID  int
	LanguageID int
	Code       string
}

type SubmitCodeResponse struct {
	SubmissionID int
}

type TestCaseResult struct {
	ID             int
	Input          string
	Output         string
	ExpectedOutput string
	RuntimeMS      int
	MemoryKB       int
	Status         string // TLE, MLE, Acccepted, Wrong Answer, $Error.message
}

type ExecuteCodePayload struct {
	ID             int
	LanguageID     int
	Code           string
	TestCases      []ProblemTestCase
	RuntimeLimitMS int
	MemoryLimitKB  int
	ExecutionType  string // Run, Submit, Validation
}

type ExecuteCodeResponse struct {
	ID              int
	Status          string // TLE, MLE, Acccepted, Wrong Answer, $Error.message
	TestCaseResults []TestCaseResult
	ExecutionType   string // Run, Submit, Validation
}

type SubmissionDB struct {
	ID        int
	Status    string
	UserID    int
	ProblemID int
	RuntimeMS int
	MemoryKB  int
}

type UserDB struct {
	ID       int
	Username string
	Email    string
	Password string
	IsAdmin  bool
	IsActive bool
}

type UserInfo struct {
	Username string
	Email    string
}

type SignupPayload struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

type LoginPayload struct {
	Username string
	Password string
}

type AuthResponse UserInfo