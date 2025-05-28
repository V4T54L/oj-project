import axios from "./axios";
import type { ProblemDetail, ProblemInfo } from "../types";

// Fetch all problems
export const getProblems = async (): Promise<ProblemInfo[]> => {
    try {
        const response = await axios.get('/api/problems');
        return response.data;
    } catch (error: any) {
        throw new Error(error.response?.data?.message || 'An error occurred while fetching problems');
    }
};

// Get problem details by ID
export const getProblemDetail = async (problemId: number): Promise<ProblemDetail> => {
    try {
        const response = await axios.get(`/api/problems/${problemId}`);
        return response.data;
    } catch (error: any) {
        throw new Error(error.response?.data?.message || 'An error occurred while fetching problem details');
    }
};
