import axios from './axios';
import type { SubmissionPayload, SubmissionResult, TestCaseResult } from '../types';

// Submit code for a problem
export const submitCode = async (payload: SubmissionPayload): Promise<{ submission_id: number }> => {
    try {
        const response = await axios.post(`/api/submit`, payload);
        return response.data;
    } catch (error: any) {
        throw new Error(error.response?.data?.message || 'An error occurred while submitting the code');
    }
};

// Run code for a problem's test cases
export const runCode = async (problemId: number, payload: SubmissionPayload): Promise<{ run_id: number }> => {
    try {
        const response = await axios.post(`/api/problems/${problemId}/run`, payload);
        return response.data;
    } catch (error: any) {
        throw new Error(error.response?.data?.message || 'An error occurred while running the code');
    }
};

// Get the result of a submission
export const getSubmissionResult = async (submissionId: number): Promise<SubmissionResult> => {
    try {
        const response = await axios.get(`/api/submissions/${submissionId}`);
        return response.data;
    } catch (error: any) {
        throw new Error(error.response?.data?.message || 'An error occurred while fetching the submission result');
    }
};

// Get the result of a specific run
export const getRunResult = async (runId: number): Promise<TestCaseResult[]> => {
    try {
        const response = await axios.get(`/a[i/runs/${runId}`);
        return response.data;
    } catch (error: any) {
        throw new Error(error.response?.data?.message || 'An error occurred while fetching the run result');
    }
};
