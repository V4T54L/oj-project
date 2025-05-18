import { useState } from "react";
import { problemDetail, programmingLanguages, submissionResults, problemExamples, testCaseResults } from "../mock";
import type { ProgrammingLanguage, TestCaseResult, SubmissionResult, ProblemExample } from "../types";
import MonacoEditor from "@monaco-editor/react";
import { getDiffText } from "../utils";


export default function ViewProblemPage() {
    const [selectedLang, setSelectedLang] = useState<ProgrammingLanguage>(programmingLanguages[0]);
    const [code, setCode] = useState("// Write your code here...");
    const [results, setResults] = useState<TestCaseResult[]>([]);
    const [summary, setSummary] = useState<SubmissionResult | null>(null);
    const [isRunning, setIsRunning] = useState(false);

    const [customTestCases, setCustomTestCases] = useState<ProblemExample[]>(problemExamples);
    const [selectedCaseIdx, setSelectedCaseIdx] = useState(0);

    const handleTestCaseChange = (field: "input" | "expected_output", value: string) => {
        const updated = [...customTestCases];
        updated[selectedCaseIdx] = {
            ...updated[selectedCaseIdx],
            [field]: value,
        };
        setCustomTestCases(updated);
    };

    const handleRunCode = () => {
        setIsRunning(true);
        setTimeout(() => {
            setResults(testCaseResults)
            setSummary(null);
            setIsRunning(false);
        }, 1000);
    };

    const handleSubmitCode = () => {
        setIsRunning(true);
        setTimeout(() => {
            setResults([])
            setSummary(submissionResults[0]);
            setIsRunning(false);
        }, 1000);
    };


    return (
        <div className="bg-gray-900 text-white min-h-screen p-4 md:p-8 space-y-8">
            <ProblemHeader title={problemDetail.title} />

            <ProblemDetails
                description={problemDetail.description}
                constraints={problemDetail.constraints}
                examples={problemDetail.examples}
            />

            <CodeEditor
                code={code}
                onChangeCode={setCode}
                selectedLang={selectedLang}
                onChangeLang={(id) => {
                    const lang = programmingLanguages.find(lang => lang.id === Number(id));
                    if (lang) setSelectedLang(lang);
                }}
                isRunning={isRunning}
                onRun={handleRunCode}
                onSubmit={handleSubmitCode}
            />

            {/* Test Case Tabs */}
            <div className="space-y-4">
                <h3 className="text-xl font-bold text-indigo-300">Test Cases</h3>

                {/* Badge Tabs */}
                <div className="flex gap-2 flex-wrap mb-4">
                    {customTestCases.map((tc, idx) => (
                        <button
                            key={tc.id}
                            onClick={() => setSelectedCaseIdx(idx)}
                            className={`px-3 py-1 rounded-full border text-sm ${idx === selectedCaseIdx
                                ? "bg-indigo-600 text-white"
                                : "bg-gray-800 text-gray-300 hover:bg-gray-700"
                                }`}
                        >
                            Test {idx + 1}
                        </button>
                    ))}
                </div>

                {/* Editable Fields */}
                <div className="bg-gray-800 border border-gray-700 rounded p-4 space-y-2 text-sm">
                    <div>
                        <label className="block mb-1 font-medium text-indigo-200">Input:</label>
                        <textarea
                            value={customTestCases[selectedCaseIdx].input}
                            onChange={(e) => handleTestCaseChange("input", e.target.value)}
                            className="w-full bg-gray-900 text-white p-2 rounded border border-gray-600"
                            rows={2}
                        />
                    </div>

                    <div>
                        <label className="block mb-1 font-medium text-indigo-200">Expected Output:</label>
                        <textarea
                            value={customTestCases[selectedCaseIdx].output}
                            onChange={(e) => handleTestCaseChange("expected_output", e.target.value)}
                            className="w-full bg-gray-900 text-white p-2 rounded border border-gray-600"
                            rows={2}
                        />
                    </div>
                </div>
            </div>


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

function ProblemDetails({ description, constraints, examples }: { description: string, constraints: string, examples: ProblemExample[] }) {
    return (
        <div className="space-y-6 text-gray-200">
            <p className="text-lg">{description}</p>

            <Section title="Constraints">
                <pre className="bg-gray-800 p-4 rounded text-sm text-gray-100 overflow-auto">
                    {constraints}
                </pre>
            </Section>

            <Section title="Examples">
                {examples.map((example: ProblemExample) => (
                    <div key={example.id} className="bg-gray-800 p-4 mb-4 rounded text-sm text-gray-100">
                        <strong>Input:</strong>
                        <pre className="border border-gray-500 rounded-md my-2 p-2">
                            {example.input}
                        </pre>
                        <strong>Output:</strong>
                        <pre className="border border-gray-500 rounded-md my-2 p-2">
                            {example.output}
                        </pre>
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
    onSubmit,
}: {
    code: string;
    onChangeCode: (val: string) => void;
    selectedLang: ProgrammingLanguage;
    onChangeLang: (id: number) => void;
    isRunning: boolean;
    onRun: () => void;
    onSubmit: () => void;
}) {
    return (
        <div className="space-y-4">
            <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                <DarkSelect
                    options={programmingLanguages.map(lang => ({
                        label: lang.name,
                        value: lang.id,
                    }))}
                    value={selectedLang.id}
                    onChange={onChangeLang}
                />

                <div className="">

                    {isRunning ? <Spinner /> : (
                        <>
                            <button
                                onClick={onRun}
                                disabled={isRunning}
                                className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded disabled:opacity-50 transition mr-4"
                            >
                                Run
                            </button>

                            <button
                                onClick={onSubmit}
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
                    options={{
                        fontSize: 14,
                        minimap: { enabled: true },
                    }}
                />
            </div>
        </div>
    );
}

function SubmissionSummary({ summary }: { summary: SubmissionResult }) {
    return (
        <div className="bg-gray-800 rounded p-4 mt-4 space-y-2 border border-gray-700">
            <h3 className="text-xl font-semibold text-indigo-300">Submission Summary</h3>
            <p>
                <strong>Verdict:</strong>{" "}
                <span className={summary.Verdict === "Accepted" ? "text-green-400" : "text-red-400"}>
                    {summary.Verdict}
                </span>
            </p>
            <p><strong>Runtime:</strong> {summary.runtime_ms} ms</p>
            <p><strong>Memory:</strong> {summary.memory_kb} KB</p>
            <p><strong>Message:</strong> {summary.message}</p>
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
                    const verdict = res.Verdict;

                    const dotColor =
                        verdict === "Accepted" ? "bg-green-400" : "bg-red-400";

                    return (
                        <button
                            key={res.id}
                            onClick={() => setActiveIndex(idx)}
                            className={`flex items-center gap-2 px-4 py-1 rounded-full text-sm font-medium transition ${isActive ? "bg-indigo-600 text-white" : "bg-gray-800 text-gray-300 hover:bg-gray-700"
                                }`}
                        >
                            <span className={`w-2 h-2 rounded-full ${dotColor}`}></span>
                            Test {idx + 1}
                        </button>
                    );
                })}
            </div>


            {/* Active Test Case */}
            <div className="bg-gray-900 border border-gray-700 rounded-md shadow p-4">
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
                            className={`font-semibold ${active.Verdict === "Accepted" ? "text-green-400" : "text-red-400"
                                }`}
                        >
                            {active.Verdict}
                        </span>
                    </p>
                    <p><strong>Runtime:</strong> {active.runtime_ms} ms</p>
                    <p><strong>Memory:</strong> {active.memory_kb} KB</p>
                </div>
            </div>
        </div>
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

function Spinner() {
    return (
        <span className="inline-block animate-spin rounded-full h-5 w-5 border-t-2 border-b-2 border-white" />
    );
}

function mapLanguage(langName: string): string {
    const map: Record<string, string> = {
        "Python": "python",
        "JavaScript": "javascript",
        "TypeScript": "typescript",
        "C++": "cpp",
        "Java": "java",
    };
    return map[langName] || "plaintext";
}
