import axios from 'axios';
import { useAuthStore } from '@/stores/auth';
import { useErrorStore } from '@/stores/error';

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

/**
 * Genera un UUID v4 per le chiavi di idempotenza.
 * Usa crypto.randomUUID() se disponibile, con fallback manuale.
 */
function _generateIdempotencyKey() {
    if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
        return crypto.randomUUID();
    }
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
        const r = (Math.random() * 16) | 0;
        const v = c === 'x' ? r : (r & 0x3) | 0x8;
        return v.toString(16);
    });
}

/**
 * Endpoint critici su cui viene aggiunto automaticamente l'header Idempotency-Key.
 * Il backend può usare questa chiave per deduplicare richieste duplicate
 * (doppio click, retry di rete) e restituire la risposta già elaborata.
 */
const IDEMPOTENT_ENDPOINTS = [
    '/grades',
    '/grades/bulk',
    '/attendance',
    '/class-tests',
    '/payments',
    '/signatures',
    '/firme',
    '/verbali',
];

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

        // ── Idempotency-Key automatico su POST/PUT/PATCH critici ────────────
        // Aggiunto solo se l'header non è già presente (es. impostato da useIdempotency).
        // Previene doppi inserimenti su retry di rete (timeout, 502/503) senza
        // che ogni service call debba gestirlo manualmente.
        const method = (config.method || '').toLowerCase();
        if (['post', 'put', 'patch'].includes(method) && !config.headers['Idempotency-Key']) {
            const url = config.url || '';
            const isIdempotentEndpoint = IDEMPOTENT_ENDPOINTS.some(ep => url.includes(ep));
            if (isIdempotentEndpoint) {
                config.headers['Idempotency-Key'] = _generateIdempotencyKey();
            }
        }
        // ─────────────────────────────────────────────────────────────────────

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

        // ── Automatic retry for transient server errors ──────────────────────
        // Retry up to 2 times on 502/503/504 and network timeouts.
        // Skip: auth endpoints (avoid loops), 4xx (client errors), already-retried requests.
        if (originalRequest && !originalRequest._retryCount) {
            originalRequest._retryCount = 0;
        }
        const isRetryable =
            originalRequest &&
            !originalRequest._isRetryRequest &&
            (originalRequest._retryCount || 0) < 2 &&
            !originalRequest.url?.includes('/auth/') &&
            (
                error.code === 'ECONNABORTED' ||
                (error.response?.status >= 502 && error.response?.status <= 504)
            );

        if (isRetryable) {
            originalRequest._retryCount = (originalRequest._retryCount || 0) + 1;
            originalRequest._isRetryRequest = true;
            const delay = 500 * originalRequest._retryCount; // 500ms, then 1000ms
            await new Promise(resolve => setTimeout(resolve, delay));
            return api(originalRequest);
        }
        // ─────────────────────────────────────────────────────────────────────

        if (!error.response) {
            error.userMessage = appI18n?.global?.t ? appI18n.global.t('errors.connectionError') : 'Errore di connessione al server. Verifica la tua connessione e riprova.';
            // Report network errors to error store
            try { useErrorStore().reportError(error, error.userMessage); } catch { /* Pinia not ready */ }
            return Promise.reject(error);
        }

        const serverMsg = error.response?.data?.message || error.response?.data?.error;

        if (error.response.status === 403) {
            error.userMessage = serverMsg || (appI18n?.global?.t ? appI18n.global.t('errors.forbidden') : 'Non disponi dei permessi necessari per completare questa operazione.');
            try { useErrorStore().reportError(error, error.userMessage); } catch { /* Pinia not ready */ }
        } else if (error.response.status === 429) {
            error.userMessage = serverMsg || (appI18n?.global?.t ? appI18n.global.t('errors.rateLimit') : 'Troppi tentativi di accesso. Riprova tra un minuto.');
            try { useErrorStore().reportError(error, error.userMessage); } catch { /* Pinia not ready */ }
        } else if (error.response.status >= 500) {
            error.userMessage = serverMsg || (appI18n?.global?.t ? appI18n.global.t('errors.serverError') : 'Si è verificato un errore sul server. Riprova più tardi.');
            try { useErrorStore().reportError(error, error.userMessage); } catch { /* Pinia not ready */ }
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

