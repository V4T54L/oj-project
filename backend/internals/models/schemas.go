package models

type UserInfo struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}

type basic struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Tag basic
type Difficulty basic
type ProgrammingLanguage basic

type ProblemInfo struct {
	ID             int        `json:"id"`
	Title          string     `json:"title"`
	IsSolved       bool       `json:"is_solved"`
	Difficulty     Difficulty `json:"difficulty"`
	Tags           []Tag      `json:"tags"`
	AcceptanceRate int        `json:"acceptance_rate"`
}

type ProblemExample struct {
	ID          int    `json:"id"`
	Input       string `json:"input"`
	Output      string `json:"output"`
	Explanation string `json:"explanation"`
}

type ProblemDetail struct {
	ProblemInfo
	Description string           `json:"description"`
	Examples    []ProblemExample `json:"examples"`
	Constraints string           `json:"constraints"`
}

type SubmissionResult struct {
	ID        int    `json:"id"`
	Verdict   string `json:"verdict"`
	RuntimeMS int    `json:"runtime_ms"`
	MemoryKB  int    `json:"memory_kb"`
	Message   string `json:"message"`
}

type RunResult struct {
	SubmissionResult
	Input          string `json:"input"`
	Output         string `json:"output"`
	ExpectedOutput string `json:"expected_output"`
}

type SubmissionPayload struct {
	LanguageID int      `json:"language_id"`
	Code       string   `json:"code"`
	TestCases  []string `json:"test_cases"`
}
