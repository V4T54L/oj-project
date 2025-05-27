package mockdata

import "online-judge/internal/models"

var Tags = []models.Tag{
	{ID: 1, Name: "Array"},
	{ID: 2, Name: "Dynamic Programming"},
	{ID: 3, Name: "Binary Search"},
}

var Difficulties = []models.Difficulty{
	{ID: 1, Name: "Easy"},
	{ID: 2, Name: "Medium"},
	{ID: 3, Name: "Hard"},
}

var ProgrammingLanguages = []models.ProgrammingLanguage{
	{ID: 1, Name: "Go"},
	{ID: 2, Name: "Python"},
	{ID: 3, Name: "JavaScript"},
}

var ExampleCases = []models.ProblemExamples{
	{ID: 1, Input: "[2,7,11,15], target=9", ExpectedOutput: "[0,1]", Explaination: "nums[0] + nums[1] == 9"},
	{ID: 2, Input: "[3,2,4], target=6", ExpectedOutput: "[1,2]", Explaination: "nums[1] + nums[2] == 6"},
}

var TestCases = []models.ProblemTestCase{
	{ID: 1, Input: "[2,7,11,15], target=9", ExpectedOutput: "[0,1]"},
	{ID: 2, Input: "[3,2,4], target=6", ExpectedOutput: "[1,2]"},
}

var ProblemsDB = []models.ProblemDB{
	{
		ID:                 101,
		Title:              "Two Sum",
		Description:        "Given an array of integers, return indices of the two numbers such that they add up to a specific target.",
		Slug:               "two-sum",
		Tags:               []models.Tag{Tags[0]},
		Difficulty:         Difficulties[0],
		AcceptanceRate:     45.3,
		IsSolved:           nil,
		Examples:           ExampleCases,
		SolutionLanguageID: 1,
		SolutionCode:       `func twoSum(nums []int, target int) []int { /*...*/ }`,
		Explaination:       "Use a hash map to store indices.",
		Status:             "Active",
		RuntimeLimitMS:     1000,
		MemoryLimitKB:      256000,
	},
}

func ToProblemInfo(db models.ProblemDB) models.ProblemInfo {
	return models.ProblemInfo{
		ID:             db.ID,
		Title:          db.Title,
		Slug:           db.Slug,
		Tags:           db.Tags,
		Difficulty:     db.Difficulty,
		AcceptanceRate: db.AcceptanceRate,
		IsSolved:       db.IsSolved,
	}
}

func ToProblemDetail(db models.ProblemDB) models.ProblemDetail {
	return models.ProblemDetail{
		ID:             db.ID,
		Title:          db.Title,
		Description:    db.Description,
		Slug:           db.Slug,
		Tags:           db.Tags,
		Difficulty:     db.Difficulty,
		AcceptanceRate: db.AcceptanceRate,
		IsSolved:       db.IsSolved,
		Examples:       db.Examples,
	}
}

var Submissions = []models.SubmissionDB{}
