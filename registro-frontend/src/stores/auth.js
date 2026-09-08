import { defineStore, getActivePinia } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'
import { resetApiState, clearLocalSession, getBaseURL } from '../services/api'
import { isTokenExpired, getRoleFromToken } from '../utils/jwt'

const parseUser = (val) => {
    if (!val) return null
    try {
        return JSON.parse(val)
    } catch (_err) {
        return null
    }
}

const ALLOWED_USER_FIELDS = new Set([
    'id', 'first_name', 'last_name', 'email', 'role', 'user_role', 'school_id', 'class_id', 'is_staff', 'created_at', 'updated_at', 'avatar'
])

const sanitizeUserData = (userData) => {
    if (!userData) return null
    const clean = {}
    for (const key of Object.keys(userData)) {
        if (ALLOWED_USER_FIELDS.has(key)) {
            clean[key] = userData[key]
        }
    }
    return clean
}

const toStorageUser = (userData) => {
    if (!userData) return null
    return {
        id: userData.id,
        role: userData.role || userData.user_role || null
    }
}

export const useAuthStore = defineStore('auth', () => {
    const user = ref(parseUser(sessionStorage.getItem('user')) || parseUser(localStorage.getItem('user')) || null)
    const isInitializing = ref(false)
    // SECURITY: Access token is kept strictly in memory (Pinia ref) to prevent theft via XSS.
    // Refresh tokens are handled exclusively via HttpOnly cookies by backend; refreshToken ref always evaluates to null in memory.
    const token = ref(null)
    const refreshToken = computed(() => null)

    const isAuthenticated = computed(() => {
        if (!token.value) return false
        return !isTokenExpired(token.value)
    })

    const userRole = computed(() => {
        const jwtRole = getRoleFromToken(token.value)
        if (jwtRole) return jwtRole
        return user.value?.role || user.value?.user_role || null
    })

    const userName = computed(() => {
        if (!user.value) return 'User'
        return `${user.value.first_name || user.value.firstName || ''} ${user.value.last_name || user.value.lastName || ''}`.trim() || 'User'
    })

    function login(userData, tokenData, _refreshTokenData = null, rememberMe = true) {
        const sanitized = sanitizeUserData(userData)
        user.value = sanitized
        token.value = tokenData

        // Clear legacy token items from storage
        localStorage.removeItem('token')
        localStorage.removeItem('refreshToken')
        sessionStorage.removeItem('token')
        sessionStorage.removeItem('refreshToken')

        const storageUser = toStorageUser(sanitized)
        if (rememberMe) {
            localStorage.setItem('user', JSON.stringify(storageUser))
            sessionStorage.removeItem('user')
        } else {
            sessionStorage.setItem('user', JSON.stringify(storageUser))
            localStorage.removeItem('user')
        }
    }

    function updateTokens(newTokenData, _newRefreshTokenData = undefined, newUserData = null) {
        if (!newTokenData) return
        token.value = newTokenData
        if (newUserData) {
            user.value = sanitizeUserData(newUserData)
        }

        // Ensure legacy tokens are removed from storage
        localStorage.removeItem('token')
        localStorage.removeItem('refreshToken')
        sessionStorage.removeItem('token')
        sessionStorage.removeItem('refreshToken')

        if (newUserData && user.value) {
            const isLocal = !!localStorage.getItem('user')
            const storage = isLocal ? localStorage : sessionStorage
            storage.setItem('user', JSON.stringify(toStorageUser(user.value)))
        }
    }

    function logout() {
        user.value = null
        token.value = null

        try {
            const pinia = getActivePinia()
            if (pinia && pinia._s) {
                for (const [storeId, store] of pinia._s.entries()) {
                    if (storeId === 'auth') continue
                    try {
                        if (storeId === 'websocket' && typeof store.disconnect === 'function') {
                            store.disconnect(true)
                        }
                        if (typeof store.clearCache === 'function') {
                            store.clearCache()
                        }
                        if (typeof store.invalidateCache === 'function') {
                            store.invalidateCache()
                        }
                        if (typeof store.reset === 'function') {
                            store.reset()
                        }
                        if (typeof store.$reset === 'function') {
                            store.$reset()
                        }
                    } catch {
                        // Store cleanup guard
                    }
                }
            }
        } catch {
            // Pinia context not active or already disposed
        }

        clearLocalSession()
        resetApiState()

        if (typeof window !== 'undefined' && 'caches' in window) {
            caches.delete('api-static-lists').catch(() => {})
        }
    }

    function updateUser(userData) {
        const sanitized = sanitizeUserData(userData)
        user.value = sanitized
        const storageUser = toStorageUser(sanitized)
        if (localStorage.getItem('user')) {
            localStorage.setItem('user', JSON.stringify(storageUser))
        } else if (sessionStorage.getItem('user')) {
            sessionStorage.setItem('user', JSON.stringify(storageUser))
        }
    }

    const initPromise = ref(null)

    async function initAuth() {
        if (initPromise.value) return initPromise.value
        if (token.value && !isTokenExpired(token.value)) return Promise.resolve(true)

        const hasSavedUser = !!(localStorage.getItem('user') || sessionStorage.getItem('user'))
        if (!hasSavedUser) return Promise.resolve(false)

        isInitializing.value = true
        initPromise.value = (async () => {
            try {
                const refreshResponse = await axios.post(
                    `${getBaseURL()}/auth/refresh-token`,
                    {},
                    { withCredentials: true }
                )
                const { access_token, user: userData } = refreshResponse.data || {}
                if (access_token) {
                    let fullUser = userData
                    if (!fullUser || !fullUser.first_name) {
                        try {
                            const meRes = await axios.get(`${getBaseURL()}/auth/me`, {
                                headers: { Authorization: `Bearer ${access_token}` }
                            })
                            fullUser = meRes.data || fullUser
                        } catch (meErr) {
                            console.warn('Initial session restore: /auth/me fetch failed, falling back to minimal profile', meErr)
                        }
                    }
                    updateTokens(access_token, null, fullUser || user.value)
                }
            } catch (err) {
                if (err.response && (err.response.status === 400 || err.response.status === 401 || err.response.status === 403)) {
                    logout()
                } else {
                    console.warn('Initial session restore failed:', err)
                }
            } finally {
                isInitializing.value = false
                initPromise.value = null
            }
        })()

        return initPromise.value
    }

    return {
        user,
        token,
        refreshToken,
        isAuthenticated,
        userRole,
        userName,
        isInitializing,
        initAuth,
        login,
        logout,
        updateUser,
        updateTokens
    }
})

