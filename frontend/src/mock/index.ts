import type {
    Tag,
    Difficulty,
    ProgrammingLanguage,
    ProblemInfo,
    ProblemExample,
    ProblemDetail,
    SubmissionResult,
    TestCase,
    TestCaseResult,
    SubmissionPayload,
} from '../types' // Adjust import path as needed

// Tags
export const tags: Tag[] = [
    { id: 1, name: "Array" },
    { id: 2, name: "String" },
    { id: 3, name: "Dynamic Programming" },
    { id: 4, name: "Graph" },
    { id: 5, name: "Math" },
]

// Difficulties
export const difficulties: Difficulty[] = [
    { id: 1, name: "Easy" },
    { id: 2, name: "Medium" },
    { id: 3, name: "Hard" },
]

// Programming Languages
export const programmingLanguages: ProgrammingLanguage[] = [
    { id: 1, name: "Python" },
    { id: 2, name: "JavaScript" },
    { id: 3, name: "C++" },
]

// Problem Info
export const problemInfos: ProblemInfo[] = [
    {
        id: 101,
        is_solved: true,
        title: "Two Sum",
        difficulty: difficulties[0],
        tags: [tags[0], tags[4]],
        acceptance_rate: 45.3,
    },
    {
        id: 102,
        is_solved: false,
        title: "Longest Substring Without Repeating Characters",
        difficulty: difficulties[1],
        tags: [tags[1]],
        acceptance_rate: 34.2,
    },
    {
        id: 103,
        is_solved: false,
        title: "Longest Palindromic Substring",
        difficulty: difficulties[1],
        tags: [tags[1], tags[2]],
        acceptance_rate: 30.5,
    },
    {
        id: 104,
        is_solved: true,
        title: "Median of Two Sorted Arrays",
        difficulty: difficulties[2],
        tags: [tags[0], tags[4]],
        acceptance_rate: 29.6,
    },
    {
        id: 105,
        is_solved: false,
        title: "Course Schedule",
        difficulty: difficulties[1],
        tags: [tags[3]],
        acceptance_rate: 40.0,
    },
]

// Problem Examples
export const problemExamples: ProblemExample[] = [
    {
        id: 1,
        input: "nums = [2,7,11,15], target = 9",
        output: "[0,1]",
        explanation: "Because nums[0] + nums[1] == 9, we return [0, 1].",
    },
    {
        id: 2,
        input: "s = 'abcabcbb'",
        output: "3",
        explanation: "The answer is 'abc', with the length of 3.",
    },
    {
        id: 3,
        input: "5 2\n2 5 7 0 3",
        output: "3",
        explanation: "The answer is 'abc', with the length of 3.",
    },
]

// Problem Detail
export const problemDetail: ProblemDetail = {
    id: 101,
    is_solved: true,
    title: "Two Sum",
    description:
        "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.",
    examples: problemExamples,
    constraints: "2 <= nums.length <= 10^4, -10^9 <= nums[i] <= 10^9",
}

// Submission Results
export const submissionResults: SubmissionResult[] = [
    {
        id: 201,
        Verdict: "Accepted",
        runtime_ms: 5,
        memory_kb: 1024,
        message: "Success",
    },
    {
        id: 202,
        Verdict: "Wrong Answer",
        runtime_ms: 4,
        memory_kb: 980,
        message: "Expected output '[0,1]', but got '[1,0]'",
    },
    {
        id: 203,
        Verdict: "Time Limit Exceeded",
        runtime_ms: 1001,
        memory_kb: 1100,
        message: "Execution time exceeded the limit",
    },
]

// Test Cases
export const testCases: TestCase[] = [
    { id: 1, input: "nums = [2,7,11,15], target = 9" },
    { id: 2, input: "nums = [3,2,4], target = 6" },
    { id: 3, input: "nums = [3,3], target = 6" },
    { id: 4, input: "nums = [1,2,3], target = 7" },
    { id: 5, input: "nums = [], target = 0" },
]

// Test Case Results
export const testCaseResults: TestCaseResult[] = [
    {
        ...submissionResults[0],
        id: 1,
        input: testCases[0].input,
        std_out: "[0,1]",
        output: "[0,1]",
        expected_output: "[0,1]",
    },
    {
        ...submissionResults[0],
        id: 2,
        input: testCases[1].input,
        std_out: "[1,2]",
        output: "[1,2]",
        expected_output: "[1,2]",
    },
    {
        ...submissionResults[1],
        id: 3,
        input: testCases[2].input,
        std_out: "[0,1]",
        output: "[1,0]",
        expected_output: "[0,1]",
    },
    {
        ...submissionResults[2],
        id: 4,
        input: testCases[3].input,
        std_out: "",
        output: "",
        expected_output: "[]",
    },
    {
        ...submissionResults[1],
        id: 5,
        input: testCases[4].input,
        std_out: "[]",
        output: "[]",
        expected_output: "[]",
        Verdict: "Memory Limit Exceeded",
    },
]

export const submissionPayload: SubmissionPayload = {
    language_id: programmingLanguages[0].id,
    code: `def twoSum(nums, target):\n    for i in range(len(nums)):\n        for j in range(i + 1, len(nums)):\n            if nums[i] + nums[j] == target:\n                return [i, j]`,
    test_cases: testCases,
}
