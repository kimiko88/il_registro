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
    withCredentials: true,
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
            // Store not ready yet
        }
        // Note: Tokens are stored exclusively in Pinia memory to prevent XSS.
        // Any client-side decoding is used strictly for UX optimization (e.g. routing/UI state);
        // cryptographical signature verification and authorization are strictly enforced server-side.
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        const lang = localStorage.getItem('superadmin_language') || 'it-IT';
        const sanitizedLang = /^[a-zA-Z0-9_-]{2,10}$/.test(lang) ? lang : 'it-IT';
        config.headers['Accept-Language'] = sanitizedLang;

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

let appRouter = null;

export const setApiRouter = (router) => {
    appRouter = router;
};

export const resetApiState = () => {
    isRefreshing = false;
    failedQueue = [];
};

export const clearLocalSession = () => {
    localStorage.removeItem('user');
    localStorage.removeItem('token');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('selectedChildId');
    sessionStorage.removeItem('user');
    sessionStorage.removeItem('token');
    sessionStorage.removeItem('refreshToken');
    sessionStorage.removeItem('selectedChildId');
};

const handleSessionExpired = () => {
    if (appRouter && typeof appRouter.push === 'function') {
        if (appRouter.currentRoute?.value?.path !== '/login') {
            appRouter.push('/login?reason=session_expired');
        }
    } else if (window.location.pathname !== '/login') {
        window.location.href = '/login?reason=session_expired';
    }
};

let appI18n = null;

export const setApiI18n = (i18nInstance) => {
    appI18n = i18nInstance;
};

api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error?.config;

        if (!error.response) {
            error.userMessage = appI18n?.global?.t ? appI18n.global.t('errors.connectionError') : 'Errore di connessione al server. Verifica la tua connessione e riprova.';
            return Promise.reject(error);
        }

        if (error.response.status === 403) {
            error.userMessage = appI18n?.global?.t ? appI18n.global.t('errors.forbidden') : 'Non disponi dei permessi necessari per completare questa operazione.';
        } else if (error.response.status === 429) {
            error.userMessage = error.response.data?.error || 'Troppi tentativi di accesso. Riprova tra un minuto.';
        } else if (error.response.status >= 500) {
            error.userMessage = appI18n?.global?.t ? appI18n.global.t('errors.serverError') : 'Si è verificato un errore sul server. Riprova più tardi.';
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
                    const refreshResponse = await axios.post(
                        `${getBaseURL()}/auth/refresh-token`,
                        {},
                        { withCredentials: true }
                    );
                    const access_token = refreshResponse.data?.access_token;
                    const newRefreshToken = refreshResponse.data?.refresh_token;

                    if (authStore && authStore.updateTokens) {
                        authStore.updateTokens(access_token, newRefreshToken);
                    }

                    processQueue(null, access_token);
                    originalRequest.headers.Authorization = `Bearer ${access_token}`;
                    return api(originalRequest);
                } catch (refreshErr) {
                    processQueue(refreshErr, null);
                    if (authStore) {
                        authStore.logout();
                    } else {
                        clearLocalSession();
                    }
                    handleSessionExpired();
                    return Promise.reject(refreshErr);
                } finally {
                    isRefreshing = false;
                }
            }
        }
        return Promise.reject(error);
    }
);

export default api;
