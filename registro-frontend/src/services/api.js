import axios from 'axios';
import { useAuthStore } from '@/stores/auth';

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
    timeout: 15000,
    headers: {
        'Content-Type': 'application/json',
    },
});

api.interceptors.request.use(
    (config) => {
        let token = null;
        try {
            const authStore = useAuthStore();
            token = authStore.token;
        } catch {
            // fallback if store not ready
        }
        if (!token) {
            token = localStorage.getItem('token') || sessionStorage.getItem('token');
        }
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        const lang = localStorage.getItem('superadmin_language') || 'it-IT';
        config.headers['Accept-Language'] = lang;

        return config;
    },
    (error) => Promise.reject(error)
);

let isRefreshing = false;
let failedQueue = [];

const processQueue = (error, token = null) => {
    failedQueue.forEach((prom) => {
        if (error) {
            prom.reject(error);
        } else {
            prom.resolve(token);
        }
    });
    failedQueue = [];
};

api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error?.config;

        if (!error.response) {
            error.userMessage = 'Errore di connessione al server. Verifica la tua connessione e riprova.';
            return Promise.reject(error);
        }

        if (error.response.status === 403) {
            error.userMessage = 'Non disponi dei permessi necessari per completare questa operazione.';
        } else if (error.response.status >= 500) {
            error.userMessage = 'Si è verificato un errore sul server. Riprova più tardi.';
        }

        if (
            error.response.status === 401 &&
            originalRequest &&
            !originalRequest.url.includes('/auth/login') &&
            !originalRequest.url.includes('/auth/refresh')
        ) {
            if (!originalRequest._retry) {
                originalRequest._retry = true;
                let authStore = null;
                try {
                    authStore = useAuthStore();
                } catch {
                    // Store not initialized
                }

                const refreshToken = authStore?.refreshToken || localStorage.getItem('refreshToken') || sessionStorage.getItem('refreshToken');

                if (refreshToken) {
                    if (isRefreshing) {
                        return new Promise((resolve, reject) => {
                            failedQueue.push({ resolve, reject });
                        })
                            .then((token) => {
                                originalRequest.headers.Authorization = `Bearer ${token}`;
                                return api(originalRequest);
                            })
                            .catch((err) => Promise.reject(err));
                    }

                    isRefreshing = true;

                    try {
                        const refreshResponse = await axios.post(`${getBaseURL()}/auth/refresh-token`, {
                            refresh_token: refreshToken,
                        });
                        const access_token = refreshResponse.data?.access_token;
                        const newRefreshToken = refreshResponse.data?.refresh_token;

                        if (authStore && authStore.updateTokens) {
                            authStore.updateTokens(access_token, newRefreshToken);
                        } else {
                            if (localStorage.getItem('token')) {
                                localStorage.setItem('token', access_token);
                                if (newRefreshToken) localStorage.setItem('refreshToken', newRefreshToken);
                            } else if (sessionStorage.getItem('token')) {
                                sessionStorage.setItem('token', access_token);
                                if (newRefreshToken) sessionStorage.setItem('refreshToken', newRefreshToken);
                            }
                        }

                        isRefreshing = false;
                        processQueue(null, access_token);

                        originalRequest.headers.Authorization = `Bearer ${access_token}`;
                        return api(originalRequest);
                    } catch (refreshErr) {
                        isRefreshing = false;
                        processQueue(refreshErr, null);
                        if (authStore) {
                            authStore.logout();
                        } else {
                            localStorage.removeItem('user');
                            localStorage.removeItem('token');
                            localStorage.removeItem('refreshToken');
                            sessionStorage.removeItem('user');
                            sessionStorage.removeItem('token');
                            sessionStorage.removeItem('refreshToken');
                        }
                        if (window.location.pathname !== '/login') {
                            window.location.href = '/login?reason=session_expired';
                        }
                        return Promise.reject(refreshErr);
                    }
                } else {
                    if (authStore) {
                        authStore.logout();
                    } else {
                        localStorage.removeItem('user');
                        localStorage.removeItem('token');
                        localStorage.removeItem('refreshToken');
                        sessionStorage.removeItem('user');
                        sessionStorage.removeItem('token');
                        sessionStorage.removeItem('refreshToken');
                    }
                    if (window.location.pathname !== '/login') {
                        window.location.href = '/login?reason=session_expired';
                    }
                    return Promise.reject(error);
                }
            }
        }
        return Promise.reject(error);
    }
);

export default api;
