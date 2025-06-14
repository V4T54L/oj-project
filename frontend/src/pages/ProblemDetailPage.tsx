import React, { useEffect, useState, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import { getProblemBySlug, runCode, submitSolution, getSubmission, getDiscussions, voteDiscussion, createDiscussion } from '../api/endpoints';
import type { ProblemDetail, Language, RunCodePayload, SubmissionPayload, TestCase, Submission, Discussion, DiscussionComment, AddVotePayload } from '../types';
import Editor from '@monaco-editor/react';
import DiscussionsTab from '../components/DiscussionsTab';
import { Loader } from 'lucide-react';
import ProblemDetails from '../components/ProblemDetails';
import clsx from 'clsx';

const ProblemDetailPage: React.FC = () => {
  const { slug } = useParams<{ slug: string }>();
  const [problem, setProblem] = useState<ProblemDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [language, setLanguage] = useState<Language>('python');
  const [code, setCode] = useState<string>('');
  const [testCases, setTestCases] = useState<TestCase[]>([]);
  const [activeTab, setActiveTab] = useState(0); // 0 for Problem, 1 for Discussions
  const [output, setOutput] = useState<string>('');
  const [running, setRunning] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [submissionDetails, setSubmissionDetails] = useState<Submission>();
  const [discussions, setDiscussions] = useState<Discussion[]>([]);
  const [isDiscussionsLoading, setIsDiscussionsLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchDiscussions = async () => {
      if (!problem) return;
      try {
        const response = await getDiscussions(problem.ID);
        setDiscussions(response.data);
      } catch (error) {
        console.error("Error fetching discussions", error);
      } finally {
        setIsDiscussionsLoading(false);
      }
    };

    fetchDiscussions();
  }, [problem]);

  useEffect(() => {
    if (!slug) return;

    // Fetch Problem Details
    getProblemBySlug(slug)
      .then((res) => {
        const p = res.data;
        p.Tags = p.Tags || [];
        p.Limits = p.Limits || [];
        p.Constraints = p.Constraints || [];
        p.Examples = p.Examples || [];
        setProblem(p);
        setCode('');
        setLanguage('python');
        setTestCases(res.data.Examples);
      })
      .catch(() => {
        setProblem(null);
      })
      .finally(() => setLoading(false));

    // Fetch Discussions when the tab is switched to Discussions
    if (activeTab === 1) {
      getDiscussions(Number(slug))
        .then((res) => {
          setDiscussions(res.data);
        })
        .catch(() => {
          setDiscussions([]);
        });
    }
  }, [slug, activeTab]);

  const handleRun = async () => {
    try {
      const payload: RunCodePayload = {
        ProblemID: problem?.id,
        Language: language,
        Code: code,
        Cases: testCases,
      }
      const response = await runCode(payload);
      const runID = response.run_id;
      pollForResults(runID);
    } catch (error) {
      console.error("Error initiating the run", error);
    }
  };

  const handleSubmit = async () => {
    // Your handleSubmit function logic
    try {
      const payload: SubmissionPayload = {
        ProblemID: problem?.id,
        Language: language,
        Code: code,
      }
      const response = await runCode(payload);
      const submissionID = response.submission_id;
      pollForResults(submissionID);
    } catch (error) {
      console.error("Error submitting the code", error);
    }
  };

  const pollForResults = useCallback(async (runId: number) => {
    const intervalId = setInterval(async () => {
      if (currentRunId.current <= 0) {
        clearInterval(intervalId); // Stop polling if no active run
        return;
      }

      try {
        const fetchedResults = await getSubmissionResult(runId);
        if (fetchedResults && fetchedResults.Status != "pending") {
          // setResults(fetchedResults);
          setSummary(fetchedResults);
          currentRunId.current = 0; // Stop polling once results are obtained
          clearInterval(intervalId); // Clear the polling interval
          setIsRunning(false); // Stop the running state
        }
      } catch (error) {
        console.error("Error fetching run results:", error);
      }
    }, 2000); // Poll every 2 seconds
  }, []);


  const handleCreateDiscussion = async (newDiscussion: Discussion) => {
    try {
      newDiscussion.ID = problem.ID;
      const response = await createDiscussion(newDiscussion);
      setDiscussions([...discussions, { ...newDiscussion, ID: response.data.id }]);
    } catch (error) {
      console.error("Error creating discussion", error);
    }
  };

  const handleUpdateDiscussion = async (updatedDiscussion: Discussion) => {
    try {
      updatedDiscussion.ID = problem.ID;
      await updateDiscussion(updatedDiscussion);
      setDiscussions(discussions.map(d =>
        d.ID === updatedDiscussion.ID ? updatedDiscussion : d
      ));
    } catch (error) {
      console.error("Error updating discussion", error);
    }
  };

  const handleVoteDiscussion = async (discussionId: number, voteType: 'up' | 'down') => {
    const payload: AddVotePayload = {
      discussionId,
      voteType
    };

    try {
      await voteDiscussion(payload);
      setDiscussions(discussions.map(d =>
        d.ID === discussionId ? { ...d, Votes: d.Votes + (voteType === 'up' ? 1 : -1) } : d
      ));
    } catch (error) {
      console.error("Error voting discussion", error);
    }
  };

  const handleAddDiscussion = (discussion: Discussion) => {
    createDiscussion(discussion)
      .then(() => {
        setDiscussions((prevDiscussions) => [discussion, ...prevDiscussions]);
      })
      .catch(() => {
        alert('Error creating discussion.');
      });
  };

  if (loading) return <p className="p-6">Loading problem...</p>;
  if (!problem) return <p className="p-6 text-red-500">Problem not found.</p>;

  return (
    <div className="max-w-5xl mx-auto px-4 py-6 space-y-6">
      <div className="flex space-x-4 border-b-2 pb-2">
        <button
          className={clsx(
            'tab px-6 py-2 font-medium text-sm rounded-lg transition-all duration-300',
            activeTab === 0
              ? 'border-b-4 border-blue-500 text-blue-500'
              : 'hover:bg-gray-100 text-gray-700'
          )}
          onClick={() => setActiveTab(0)}
        >
          Problem
        </button>
        <button
          className={clsx(
            'tab px-6 py-2 font-medium text-sm rounded-lg transition-all duration-300',
            activeTab === 1
              ? 'border-b-4 border-blue-500 text-blue-500'
              : 'hover:bg-gray-100 text-gray-700'
          )}
          onClick={() => setActiveTab(1)}
        >
          Discussions
        </button>
      </div>

      {activeTab === 0 && (
        <ProblemDetails
          problem={problem}
          code={code}
          setCode={setCode}
          language={language}
          setLanguage={setLanguage}
          testCases={testCases}
          setTestCases={setTestCases}
          handleRun={handleRun}
          handleSubmit={handleSubmit}
          running={running}
          activeTab={activeTab}
          setActiveTab={setActiveTab}
          submitting={submitting}
          output={output}
        />
      )}

      {activeTab === 1 && (
        <div>
          {isDiscussionsLoading ? (
            <p>Loading discussions...</p>
          ) : (
            <DiscussionsTab
              discussions={discussions}
              onCreateDiscussion={handleCreateDiscussion}
              onUpdateDiscussion={handleUpdateDiscussion}
              onVoteDiscussion={handleVoteDiscussion}
            />
          )}
        </div>

      )}
    </div>
  );
};

export default ProblemDetailPage;
