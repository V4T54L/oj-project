import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { getProblemBySlug, runCode, submitSolution, getSubmission } from '../api/endpoints';
import type { ProblemDetail, Language, RunCodePayload, SubmissionPayload, TestCase, Submission } from '../types';
import Editor from '@monaco-editor/react';
import DiscussionModal from '../components/DiscussionsModal';
import { Loader } from 'lucide-react';

const languageToMonaco = (lang: Language) => lang ?? 'plaintext';
const languageOptions: Language[] = ['python', 'cpp', 'java', 'javascript'];

const ProblemDetailPage: React.FC = () => {
  const { slug } = useParams<{ slug: string }>();
  const [problem, setProblem] = useState<ProblemDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [language, setLanguage] = useState<Language>('python');
  const [code, setCode] = useState<string>('');
  const [testCases, setTestCases] = useState<TestCase[]>([]);
  const [activeTab, setActiveTab] = useState(0);
  const [output, setOutput] = useState<string>('');
  const [running, setRunning] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [submissionDetails, setSubmissionDetails] = useState<Submission>();
  const [isDiscussionOpen, setIsDiscussionOpen] = useState(false);

  useEffect(() => {
    if (!slug) return;
    getProblemBySlug(slug)
      .then((res) => {
        setProblem(res.data);
        setCode(res.data.SolutionCode || '');
        setLanguage(res.data.SolutionLanguage || 'python');
        setTestCases(res.data.TestCases);
      })
      .catch(() => {
        setProblem(null);
      })
      .finally(() => setLoading(false));
  }, [slug]);

  const pollSubmission = async (id: number, maxAttempts = 10, delay = 1500): Promise<Submission | null> => {
    let attempts = 0;
    while (attempts < maxAttempts) {
      const { data } = await getSubmission(id);
      if (data?.Results?.length) return data;
      await new Promise((res) => setTimeout(res, delay));
      attempts++;
    }
    return null;
  };

  const handleRun = async () => {
    if (!problem) return;
    setRunning(true);
    setOutput('');
    const payload: RunCodePayload = {
      ProblemID: problem.ID,
      Language: language,
      Code: code,
      Cases: testCases,
    };
    try {
      const { data } = await runCode(payload);
      const result = await pollSubmission(data.run_id);
      if (result) {
        setSubmissionDetails(result);
      } else {
        setOutput('Run timed out. Please try again.');
      }
    } catch (err: unknown) {
      setOutput('Run failed. Please try again. ' + err);
    } finally {
      setRunning(false);
    }
  };

  const handleSubmit = async () => {
    if (!problem) return;
    setSubmitting(true);
    const payload: SubmissionPayload = {
      ProblemID: problem.ID,
      Language: language,
      Code: code,
    };
    try {
      const { data } = await submitSolution(payload);
      const result = await pollSubmission(data.submission_id);
      if (result) {
        setSubmissionDetails(result);
      } else {
        setOutput('Submission timed out. Please try again.');
      }
    } catch {
      setOutput('Submission failed.');
    } finally {
      setSubmitting(false);
    }
  };


  if (loading) return <p className="p-6">Loading problem...</p>;
  if (!problem) return <p className="p-6 text-red-500">Problem not found.</p>;

  return (
    <div className="max-w-5xl mx-auto px-4 py-6 space-y-6">
      <div>
        <h1 className="text-3xl font-bold">{problem.Title}</h1>
        <div className="flex gap-2 text-sm text-gray-600 mt-2">
          <span className="bg-gray-200 px-2 py-0.5 rounded">
            Difficulty: {problem.Difficulty}
          </span>
          {problem?.Tags?.map((tag) => (
            <span key={tag} className="bg-blue-100 text-blue-800 px-2 py-0.5 rounded">
              {tag}
            </span>
          ))}
        </div>

        <div className="mt-4">
          <button
            onClick={() => setIsDiscussionOpen(true)}
            className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded text-sm"
          >
            💬 View Discussions
          </button>
        </div>


        <p className="text-gray-800 whitespace-pre-wrap mt-4">{problem.Description}</p>
        <h2 className="font-semibold mt-4">Constraints:</h2>
        <ul className="list-disc list-inside text-sm text-gray-700">
          {problem?.Constraints?.map((c, idx) => <li key={idx}>{c}</li>)}
        </ul>
      </div>

      <div className="space-y-4">
        <div className="flex justify-between items-center">
          <label htmlFor="language" className="font-medium">Language:</label>
          <select
            id="language"
            value={language}
            onChange={(e) => setLanguage(e.target.value as Language)}
            className="border border-gray-300 rounded px-2 py-1 text-sm"
          >
            {languageOptions?.map((lang) => (
              <option key={lang} value={lang}>{lang}</option>
            ))}
          </select>
        </div>

        <Editor
          height="400px"
          theme="vs-dark"
          defaultLanguage={languageToMonaco(language)}
          language={languageToMonaco(language)}
          value={code}
          onChange={(value) => setCode(value || '')}
          options={{ fontSize: 14, minimap: { enabled: false }, wordWrap: 'on' }}
        />

        <div>
          <h2 className="font-semibold mb-2">Test Cases</h2>
          <div className="flex gap-2 mb-2">
            {testCases.map((_, idx) => (
              <button
                key={idx}
                onClick={() => setActiveTab(idx)}
                className={`px-3 py-1 rounded border ${activeTab === idx ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 border-gray-300'}`}
              >
                Case {idx + 1}
              </button>
            ))}
          </div>
          <div className="space-y-2">
            <textarea
              className="w-full p-2 border rounded font-mono text-sm"
              rows={3}
              value={testCases[activeTab]?.Input || ''}
              onChange={(e) => {
                const updated = [...testCases];
                updated[activeTab].Input = e.target.value;
                setTestCases(updated);
              }}
              placeholder="Input"
            />
            <textarea
              className="w-full p-2 border rounded font-mono text-sm"
              rows={3}
              value={testCases[activeTab]?.ExpectedOutput || ''}
              onChange={(e) => {
                const updated = [...testCases];
                updated[activeTab].ExpectedOutput = e.target.value;
                setTestCases(updated);
              }}
              placeholder="Expected Output"
            />
          </div>
        </div>

        <div className="flex gap-4 mt-4">
          <button
            onClick={handleRun}
            disabled={running || submitting}
            className="bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-2 rounded flex items-center gap-2"
          >
            {running && <Loader className="animate-spin w-4 h-4" />}
            {running ? 'Running...' : 'Run Code'}
          </button>

          <button
            onClick={handleSubmit}
            disabled={submitting || running}
            className="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded flex items-center gap-2"
          >
            {submitting && <Loader className="animate-spin w-4 h-4" />}
            {submitting ? 'Submitting...' : 'Submit'}
          </button>

        </div>

        {output && (
          <div className="mt-6 bg-gray-100 p-4 rounded">
            <h3 className="font-medium mb-2">Output</h3>
            <pre className="whitespace-pre-wrap text-sm">{output}</pre>
          </div>
        )}

        {
          submissionDetails?.Results && (
            <div className="mt-6 bg-white border border-gray-300 p-4 rounded">
              <h3 className="font-semibold text-lg mb-4">Submission Result - Case #{activeTab + 1}</h3>

              {
                submissionDetails.Results[activeTab] ? (
                  <div className="space-y-3 text-sm text-gray-800">
                    <p><strong>Status:</strong> {submissionDetails.Results[activeTab].Status.toUpperCase()}</p>
                    <p><strong>Runtime:</strong> {submissionDetails.Results[activeTab].RuntimeMS} ms</p>
                    <p><strong>Memory:</strong> {submissionDetails.Results[activeTab].MemoryKB} KB</p>

                    <div>
                      <p className="font-semibold">Input</p>
                      <pre className="bg-gray-100 p-2 rounded whitespace-pre-wrap font-mono">{testCases[activeTab]?.Input}</pre>
                    </div>

                    <div>
                      <p className="font-semibold">Expected Output</p>
                      <pre className="bg-gray-100 p-2 rounded whitespace-pre-wrap font-mono">{testCases[activeTab]?.ExpectedOutput}</pre>
                    </div>

                    <div>
                      <p className="font-semibold">StdOut</p>
                      <pre className="bg-green-100 p-2 rounded whitespace-pre-wrap font-mono">{submissionDetails.Results[activeTab].StdOut || 'N/A'}</pre>
                    </div>

                    {submissionDetails.Results[activeTab].StdErr && (
                      <div>
                        <p className="font-semibold text-red-600">StdErr</p>
                        <pre className="bg-red-100 p-2 rounded whitespace-pre-wrap font-mono text-red-800">{submissionDetails.Results[activeTab].StdErr}</pre>
                      </div>
                    )}
                  </div>
                ) : (
                  <p className="text-sm text-gray-500">No result found for this case.</p>
                )}
            </div>
          )}

      </div>

      <DiscussionModal
        problemID={problem.ID}
        isOpen={isDiscussionOpen}
        onClose={() => setIsDiscussionOpen(false)}
      />
    </div>
  );
};

export default ProblemDetailPage;
