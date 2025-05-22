import axios from 'axios';

// Base URL can be dynamically loaded from an environment variable
const axiosInstance = axios.create({
    baseURL: 'https://glowing-giggle-74gwwvrxpxg2rj56-8080.app.github.dev',
    headers: {
        'Content-Type': 'application/json',
    },
});

// Request interceptor to add Authorization header if available
axiosInstance.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token'); // or use your global state
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// Response interceptor to handle errors globally
axiosInstance.interceptors.response.use(
    (response) => response,
    (error) => {
        // Handle specific error responses here, e.g., unauthorized
        if (error.response?.status === 401) {
            // handle unauthenticated user
        }
        return Promise.reject(error);
    }
);

export default axiosInstance;
