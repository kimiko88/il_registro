import axios from 'axios';
import { useAuthStore } from '@/stores/auth';

export const getBaseURL = () => {
    const rawUrl = import.meta.env.VITE_API_URL;
    if (!rawUrl) {
        return '/api/v1';
    }
    const cleaned = String(rawUrl).trim().replace(/\/+$/, '');
    if (cleaned.endsWith('/api/v1')) {
        return cleaned;
    }
    return `${cleaned}/api/v1`;
};

const api = axios.create({
    baseURL: getBaseURL(),
    timeout: 15000,
    withCredentials: true,
    headers: {
        'Content-Type': 'application/json',
    },
});

const getEffectiveLanguage = (userLang) => {
    if (userLang && typeof userLang === 'string') return userLang;
    try {
        const lang = (typeof localStorage !== 'undefined' && (localStorage.getItem('app_language') || localStorage.getItem('user_locale'))) || 'it-IT';
        return typeof lang === 'string' && /^[a-zA-Z0-9_-]{2,10}$/.test(lang) ? lang : 'it-IT';
    } catch {
        return 'it-IT';
    }
};

const activeAbortControllers = new Set();

export const cancelInFlightRequests = () => {
    for (const controller of activeAbortControllers) {
        try {
            controller.abort();
        } catch (_e) {
            // Ignore errors if the request was already completed or aborted
        }
    }
    activeAbortControllers.clear();
};

api.interceptors.request.use(
    (config) => {
        let token = null;
        let userLang = null;
        try {
            const authStore = useAuthStore();
            token = authStore.token;
            userLang = authStore.user?.language;
        } catch {
            // Store not ready yet
        }
        // Note: Tokens are stored exclusively in Pinia memory to prevent XSS.
        // Any client-side decoding is used strictly for UX optimization (e.g. routing/UI state);
        // cryptographical signature verification and authorization are strictly enforced server-side.
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        config.headers['Accept-Language'] = getEffectiveLanguage(userLang);

        if (config.url) {
            if (config.url.startsWith('/api/v1/')) {
                config.url = config.url.substring('/api/v1'.length);
            } else if (config.url.startsWith('api/v1/')) {
                config.url = '/' + config.url.substring('api/v1/'.length);
            }
        }

        return config;
    },
    (error) => Promise.reject(error)
);

let isRefreshing = false;
let failedQueue = [];

const processQueue = (error, token = null) => {
    const queue = failedQueue;
    failedQueue = [];
    queue.forEach((prom) => {
        try {
            if (error) {
                prom.reject(error);
            } else {
                prom.resolve(token);
            }
        } catch (_e) {
            // Ignore errors if promise was already settled
        }
    });
};

let appRouter = null;

export const setApiRouter = (router) => {
    appRouter = router;
};

export const resetApiState = () => {
    isRefreshing = false;
    processQueue(new Error('Session reset or logged out'), null);
    cancelInFlightRequests();
};

export const clearLocalSession = () => {
    resetApiState();
    localStorage.removeItem('user');
    localStorage.removeItem('token');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('selectedChildId');
    sessionStorage.removeItem('user');
    sessionStorage.removeItem('token');
    sessionStorage.removeItem('refreshToken');
    sessionStorage.removeItem('selectedChildId');
    sessionStorage.removeItem('registro_lesson_drafts');
};

// Handler per la re-autenticazione in-page (registrato da MainLayout).
// Se presente, invece di navigare a /login, mostra un dialog modale.
// Firma: (email: string) => Promise<string>  (risolve col nuovo access_token)
let _reauthHandler = null;

export const setReauthHandler = (fn) => {
    _reauthHandler = fn;
};

const handleSessionExpired = () => {
    if (appRouter && typeof appRouter.push === 'function') {
        if (appRouter.currentRoute?.value?.path !== '/login') {
            appRouter.push('/login?reason=session_expired');
        }
    } else if (typeof window !== 'undefined' && window.location?.pathname !== '/login') {
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

        const serverMsg = error.response?.data?.message || error.response?.data?.error;

        if (error.response.status === 403) {
            error.userMessage = serverMsg || (appI18n?.global?.t ? appI18n.global.t('errors.forbidden') : 'Non disponi dei permessi necessari per completare questa operazione.');
        } else if (error.response.status === 429) {
            error.userMessage = serverMsg || (appI18n?.global?.t ? appI18n.global.t('errors.rateLimit') : 'Troppi tentativi di accesso. Riprova tra un minuto.');
        } else if (error.response.status >= 500) {
            error.userMessage = serverMsg || (appI18n?.global?.t ? appI18n.global.t('errors.serverError') : 'Si è verificato un errore sul server. Riprova più tardi.');
        } else if (serverMsg && typeof serverMsg === 'string') {
            error.userMessage = serverMsg;
        }

        const isAuthUrl = originalRequest?.url && (
            originalRequest.url.endsWith('/auth/login') ||
            originalRequest.url.endsWith('/auth/refresh-token') ||
            originalRequest.url.endsWith('/auth/refresh')
        );

        if (
            error.response.status === 401 &&
            originalRequest &&
            !isAuthUrl
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
                    // Prova la re-autenticazione in-page se il handler è registrato
                    if (_reauthHandler) {
                        try {
                            const email = authStore?.user?.email || '';
                            const newToken = await _reauthHandler(email);
                            if (authStore && authStore.updateTokens) {
                                authStore.updateTokens(newToken);
                            }
                            // Successo: risolvi la coda e riprova la request originale col nuovo token
                            processQueue(null, newToken);
                            originalRequest.headers.Authorization = `Bearer ${newToken}`;
                            return api(originalRequest);
                        } catch (reauthErr) {
                            // L'utente ha annullato il dialog o reauth fallito -> rifiuta la coda
                            processQueue(reauthErr, null);
                            if (authStore) {
                                authStore.logout();
                            } else {
                                clearLocalSession();
                            }
                            handleSessionExpired();
                            return Promise.reject(reauthErr);
                        }
                    } else {
                        // Fallback: nessun handler registrato → svuota coda con errore + logout + redirect
                        processQueue(refreshErr, null);
                        if (authStore) {
                            authStore.logout();
                        } else {
                            clearLocalSession();
                        }
                        handleSessionExpired();
                        return Promise.reject(refreshErr);
                    }
                } finally {
                    isRefreshing = false;
                }
            }
        }
        return Promise.reject(error);
    }
);

export default api;

