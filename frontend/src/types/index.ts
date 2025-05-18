type basic = {
    id: number
    name: string
}

export type Tag = basic

export type Difficulty = basic

export type ProgrammingLanguage = basic

export type ProblemInfo = {
    id: number
    is_solved: boolean
    title: string
    difficulty: Difficulty
    tags: Tag[]
    acceptance_rate: number
}

export type ProblemExample = {
    id: number
    input: string
    output: string
    explanation: string
}

export type ProblemDetail = {
    id: number
    is_solved: boolean
    title: string
    description: string
    examples: ProblemExample[]
    constraints: string
}

export type SubmissionResult = {
    id: number
    Verdict: string
    runtime_ms: number
    memory_kb: number
    message: string
}

export type TestCaseResult = SubmissionResult & {
    input: string
    std_out: string
    output: string
    expected_output: string
}

export type TestCase = {
    id: number
    input: string
}

export type SubmissionPayload = {
    language_id: number
    code: string
    test_cases: TestCase[]
}

export type SignupPayload = {
    avatar_url: string
    name: string
    username: string
    email: string
    password: string
    confirm_password: string
}

export type LoginPayload = {
    username: string
    password: string
}
