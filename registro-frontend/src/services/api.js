import axios from 'axios';

const getBaseURL = () => {
    const rawUrl = import.meta.env.VITE_API_URL;
    if (!rawUrl) {
        return '/api/v1';
    }
    if (rawUrl.endsWith('/api/v1') || rawUrl.endsWith('/api/v1/')) {
        return rawUrl;
    }
    return rawUrl.endsWith('/') ? `${rawUrl}api/v1` : `${rawUrl}/api/v1`;
};

const api = axios.create({
    baseURL: getBaseURL(),
    headers: {
        'Content-Type': 'application/json',
    },
});

api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

api.interceptors.response.use(
    (response) => response,
    (error) => {
        // Handle 401 Unauthorized globally, but ignore for login requests
        if (error.response && error.response.status === 401 && !error.config.url.includes('/auth/login')) {
            localStorage.removeItem('token');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

export default api;
