import axios from "./axios";

export const login = async (): Promise<any> => {
    try {
        const response = await axios.post('/login',{});
        return response.data;
    } catch (error: any) {
        throw new Error(error.response?.data?.message || 'An error occurred while fetching problems');
    }
};