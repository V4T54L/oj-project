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
    expected_output: string
    explanation: string
}

export type ProblemDetail = {
    id: number
    title: string
    description: string
    slug: string
    tags: Tag[]
    difficulty: Difficulty
    AcceptanceRate: number
    is_solved?: boolean
    examples: ProblemExample[]
    constraints: string
}

export type SubmissionResult = {
    id: number
    status: string
    runtime_ms: number
    memory_kb: number
    // message: string
}

export type TestCaseResult = SubmissionResult & {
    input: string
    // std_out: string
    output: string
    expected_output: string
    runtime_ms: number
    memory_kb: number
    status: string
}

export type TestCase = {
    id: number
    input: string
    expected_output: string
}

export type SubmissionPayload = {
    problem_id: number
    language_id: number
    code: string
    // test_cases: TestCase[]
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

export type UserInfo = {
    username: string
    email: string
}