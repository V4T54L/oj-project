import { useEffect, useState, useRef, useCallback } from "react";
import { useParams } from "react-router-dom";
import MonacoEditor from "@monaco-editor/react";
import { runCode, getRunResult } from "../api/submissions";
import { getProblemDetail } from "../api/problems";
import { programmingLanguages } from "../mock";
import type { TestCase, TestCaseResult, ProgrammingLanguage, SubmissionResult, ProblemDetail } from "../types";
import { getDiffText } from "../utils";

export default function ViewProblemPage() {
    const params = useParams();
    const [selectedLang, setSelectedLang] = useState<ProgrammingLanguage>(programmingLanguages[0]);
    const [code, setCode] = useState("// Write your code here...");
    const [results, setResults] = useState<TestCaseResult[]>([]);
    const [summary, setSummary] = useState<SubmissionResult | null>(null);
    const [isRunning, setIsRunning] = useState(false);
    const [problem, setProblem] = useState<ProblemDetail | null>(null);

    const [customTestCases, setCustomTestCases] = useState<TestCase[]>([]);
    const [selectedCaseIdx, setSelectedCaseIdx] = useState(0);
    const currentRunId = useRef(0); // Track the current runId for polling

    // Handle test case input changes
    const handleTestCaseChange = useCallback((field: "input" | "expected_output", value: string) => {
        const updated = [...customTestCases];
        updated[selectedCaseIdx] = {
            ...updated[selectedCaseIdx],
            [field]: value,
        };
        setCustomTestCases(updated);
    }, [customTestCases, selectedCaseIdx]);

    // Handle running the code
    const handleRunCode = useCallback(async () => {
        if (isRunning || !code || !customTestCases.length) return;

        setIsRunning(true);
        const problemIdStr = params["id"];
        if (!problemIdStr) return;

        const problemId = parseInt(problemIdStr);
        if (!problemId) return;

        try {
            const { run_id } = await runCode(problemId, {
                code,
                language_id: selectedLang.id,
                test_cases: customTestCases,
            });
            currentRunId.current = run_id; // Store runId to track polling
            pollForResults(run_id);
        } catch (error) {
            console.error("Failed to run code:", error);
            setIsRunning(false); // Stop loading if the code submission fails
        }
    }, [code, customTestCases, isRunning, params, selectedLang]);

    // Polling for results
    const pollForResults = useCallback(async (runId: number) => {
        const intervalId = setInterval(async () => {
            if (currentRunId.current <= 0) {
                clearInterval(intervalId); // Stop polling if no active run
                return;
            }

            try {
                const fetchedResults = await getRunResult(runId);
                if (fetchedResults) {
                    setResults(fetchedResults);
                    // setSummary();
                    currentRunId.current = 0; // Stop polling once results are obtained
                    clearInterval(intervalId); // Clear the polling interval
                    setIsRunning(false); // Stop the running state
                }
            } catch (error) {
                console.error("Error fetching run results:", error);
            }
        }, 2000); // Poll every 2 seconds
    }, []);

    // Fetch the problem details
    useEffect(() => {
        const fetchProblemDetail = async () => {
            const problemIdStr = params["id"];
            if (!problemIdStr) return;

            const problemId = parseInt(problemIdStr);
            if (!problemId) return;

            try {
                const fetchedProblem = await getProblemDetail(problemId);
                setProblem(fetchedProblem);
            } catch (error) {
                console.error("Failed to fetch problem:", error);
            }
        };

        fetchProblemDetail();
    }, [params]);

    // Update custom test cases when the problem data is loaded
    useEffect(() => {
        if (!problem) return;

        const convertedTestCases = problem.examples.map(({ id, input, output }) => ({
            id,
            input,
            expected_output: output,
        }));
        setCustomTestCases(convertedTestCases);
    }, [problem]);

    if (!problem) {
        return <div>Loading problem details...</div>;
    }

    return (
        <div className="bg-gray-900 text-white min-h-screen p-4 md:p-8 space-y-8">
            <ProblemHeader title={problem.title} />
            <ProblemDetails
                description={problem.description}
                constraints={problem.constraints}
                examples={problem.examples}
            />
            <CodeEditor
                code={code}
                onChangeCode={setCode}
                selectedLang={selectedLang}
                onChangeLang={(id) => {
                    const lang = programmingLanguages.find((lang) => lang.id === Number(id));
                    if (lang) setSelectedLang(lang);
                }}
                isRunning={isRunning}
                onRun={handleRunCode}
            />
            {/* Test Case Tabs */}
            <TestCaseTabs
                customTestCases={customTestCases}
                selectedCaseIdx={selectedCaseIdx}
                setSelectedCaseIdx={setSelectedCaseIdx}
                handleTestCaseChange={handleTestCaseChange}
            />
            {summary && <SubmissionSummary summary={summary} />}
            <TestCasesPanel results={results} />
        </div>
    );
}

function ProblemHeader({ title }: { title: string }) {
    return (
        <h1 className="text-4xl font-bold text-center bg-gradient-to-r from-indigo-500 to-purple-600 bg-clip-text text-transparent">
            {title}
        </h1>
    );
}

function ProblemDetails({ description, constraints, examples }: { description: string, constraints: string, examples: any[] }) {
    return (
        <div className="space-y-6 text-gray-200">
            <p className="text-lg">{description}</p>
            <Section title="Constraints">
                <pre className="bg-gray-800 p-4 rounded text-sm text-gray-100 overflow-auto">{constraints}</pre>
            </Section>
            <Section title="Examples">
                {examples.map((example) => (
                    <div key={example.id} className="bg-gray-800 p-4 mb-4 rounded text-sm text-gray-100">
                        <strong>Input:</strong>
                        <pre className="border border-gray-500 rounded-md my-2 p-2">{example.input}</pre>
                        <strong>Output:</strong>
                        <pre className="border border-gray-500 rounded-md my-2 p-2">{example.output}</pre>
                        <p><strong>Explanation:</strong> {example.explanation}</p>
                    </div>
                ))}
            </Section>
        </div>
    );
}

function CodeEditor({
    code,
    onChangeCode,
    selectedLang,
    onChangeLang,
    isRunning,
    onRun,
}: {
    code: string;
    onChangeCode: (val: string) => void;
    selectedLang: ProgrammingLanguage;
    onChangeLang: (id: number) => void;
    isRunning: boolean;
    onRun: () => void;
}) {
    return (
        <div className="space-y-4">
            <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                <DarkSelect
                    options={programmingLanguages.map((lang) => ({ label: lang.name, value: lang.id }))}
                    value={selectedLang.id}
                    onChange={onChangeLang}
                />
                <div>
                    {isRunning ? (
                        <Spinner />
                    ) : (
                        <>
                            <button
                                onClick={onRun}
                                disabled={isRunning}
                                className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded disabled:opacity-50 transition mr-4"
                            >
                                Run
                            </button>
                            <button
                                onClick={() => { }}
                                disabled={isRunning}
                                className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded disabled:opacity-50 transition"
                            >
                                Submit
                            </button>
                        </>
                    )}
                </div>
            </div>
            <div className="h-80 border border-gray-700 rounded overflow-hidden">
                <MonacoEditor
                    height="100%"
                    language={mapLanguage(selectedLang.name)}
                    value={code}
                    theme="vs-dark"
                    onChange={(val) => onChangeCode(val || "")}
                    options={{ fontSize: 14, minimap: { enabled: false } }}
                />
            </div>
        </div>
    );
}

function TestCaseTabs({
    customTestCases,
    selectedCaseIdx,
    setSelectedCaseIdx,
    handleTestCaseChange,
}: {
    customTestCases: TestCase[];
    selectedCaseIdx: number;
    setSelectedCaseIdx: (idx: number) => void;
    handleTestCaseChange: (field: "input" | "expected_output", value: string) => void;
}) {
    return (
        <div className="space-y-4">
            <div className="flex flex-wrap gap-4">
                {customTestCases.map((tc, idx) => (
                    <button
                        key={tc.id}
                        className={`bg-gray-800 px-4 py-2 rounded-md text-sm ${selectedCaseIdx === idx ? "bg-indigo-600" : "hover:bg-gray-700"
                            }`}
                        onClick={() => setSelectedCaseIdx(idx)}
                    >
                        Test Case {tc.id}
                    </button>
                ))}
            </div>
            {customTestCases[selectedCaseIdx] && (
                <div>
                    <textarea
                        value={customTestCases[selectedCaseIdx].input}
                        onChange={(e) => handleTestCaseChange("input", e.target.value)}
                        placeholder="Input"
                        className="w-full h-32 p-2 bg-gray-700 text-gray-100 rounded-md"
                    />
                    <textarea
                        value={customTestCases[selectedCaseIdx].expected_output}
                        onChange={(e) => handleTestCaseChange("expected_output", e.target.value)}
                        placeholder="Expected Output"
                        className="w-full h-32 p-2 mt-4 bg-gray-700 text-gray-100 rounded-md"
                    />
                </div>
            )}
        </div>
    );
}

function Spinner() {
    return (
        <div className="spinner-border animate-spin text-indigo-500" role="status">
            <span className="visually-hidden">Loading...</span>
        </div>
    );
}

function mapLanguage(languageName: string) {
    switch (languageName.toLowerCase()) {
        case "python":
            return "python";
        case "javascript":
            return "javascript";
        case "java":
            return "java";
        // Add more languages as needed
        default:
            return "plaintext";
    }
}

function DarkSelect({
    options,
    value,
    onChange,
}: {
    options: { label: string; value: number }[];
    value: number;
    onChange: (val: number) => void;
}) {
    return (
        <select
            value={value}
            onChange={(e) => onChange(Number(e.target.value))}
            className="bg-gray-800 border border-gray-600 text-white rounded px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
        >
            {options.map((opt) => (
                <option key={opt.value} value={opt.value}>
                    {opt.label}
                </option>
            ))}
        </select>
    );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
    const [open, setOpen] = useState(true);
    return (
        <div>
            <button
                onClick={() => setOpen(!open)}
                className="text-lg font-semibold text-indigo-400 mb-2 focus:outline-none"
            >
                {title} {open ? "▲" : "▼"}
            </button>
            {open && <div className="pl-2">{children}</div>}
        </div>
    );
}


function TestCasesPanel({ results }: { results: TestCaseResult[] }) {
    const [activeIndex, setActiveIndex] = useState(0);

    if (results.length === 0) return null;

    const active = results[activeIndex];

    return (
        <div className="space-y-4 mt-8">
            <h3 className="text-2xl font-bold text-indigo-300">Test Case Results</h3>

            {/* Tabs */}
            <div className="flex flex-wrap gap-2">
                {results.map((res, idx) => {
                    const isActive = idx === activeIndex;
                    const verdict = res.verdict;

                    const dotColor =
                        verdict === "Accepted" ? "bg-green-400" : "bg-red-400";

                    return (
                        <button
                            key={res.id}
                            onClick={() => setActiveIndex(idx)}
                            className={`flex items-center gap-2 px-4 py-1 rounded-full text-sm font-medium transition ${isActive ? "bg-indigo-600 text-white" : "bg-gray-800 text-gray-300 hover:bg-gray-700"
                                }`}
                        >
                            <span className={`w - 2 h-2 rounded-full ${dotColor}`}></span>
                            Test {idx + 1}
                        </button>
                    );
                })}
            </div >


            {/* Active Test Case */}
            < div className="bg-gray-900 border border-gray-700 rounded-md shadow p-4" >
                <div className="gap-4 text-sm text-gray-200 mb-4">
                    <div>
                        <p className="font-semibold mb-1 text-indigo-200">Input:</p>
                        <pre className="bg-gray-800 p-2 rounded whitespace-pre-wrap break-words">{active.input}</pre>
                    </div>
                    {active.std_out && (
                        <div className="md:col-span-2">
                            <p className="font-semibold mb-1 text-indigo-200 mt-4">Std Out:</p>
                            <pre className="bg-gray-800 p-2 rounded whitespace-pre-wrap break-words">{active.std_out}</pre>
                        </div>
                    )}
                    <div>
                        <p className="font-semibold mb-1 text-indigo-200 mt-4">Output</p>
                        <pre className="bg-gray-800 p-2 rounded whitespace-pre-wrap break-words">{active.output + " "}</pre>
                    </div>
                    <div>
                        <p className="font-semibold mb-1 text-indigo-200 mt-4">Expected Output:</p>
                        <pre className="bg-gray-800 p-2 rounded whitespace-pre-wrap break-words">{active.expected_output + " "}</pre>
                    </div>
                    <div className="md:col-span-2">
                        <p className="font-semibold mb-1 text-indigo-200 mt-4">Diff</p>
                        <pre className="bg-gray-800 p-2 rounded whitespace-pre-wrap break-words">
                            {getDiffText(active.expected_output, active.output)}
                        </pre>
                    </div>
                </div>

                <div className="border-t border-gray-700 pt-2 text-sm text-gray-300">
                    <p>
                        <strong>Verdict:</strong>{" "}
                        <span
                            className={`font - semibold ${active.verdict === "Accepted" ? "text-green-400" : "text-red-400"
                                }`}
                        >
                            {active.verdict}
                        </span>
                    </p>
                    <p><strong>Runtime:</strong> {active.runtime_ms} ms</p>
                    <p><strong>Memory:</strong> {active.memory_kb} KB</p>
                </div >
            </div >
        </div >
    );
}

function SubmissionSummary({ summary }: { summary: SubmissionResult }) {
    return (
        <div className="bg-gray-800 rounded p-4 mt-4 space-y-2 border border-gray-700">
            <h3 className="text-xl font-semibold text-indigo-300">Submission Summary</h3>
            <p>
                <strong>Verdict:</strong>{" "}
                <span className={summary.verdict === "Accepted" ? "text-green-400" : "text-red-400"}>
                    {summary.verdict}
                </span>
            </p>
            <p><strong>Runtime:</strong> {summary.runtime_ms} ms</p>
            <p><strong>Memory:</strong> {summary.memory_kb} KB</p>
            <p><strong>Message:</strong> {summary.message}</p>
        </div>
    );
}
