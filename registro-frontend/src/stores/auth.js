import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'
import { resetApiState, clearLocalSession, getBaseURL } from '../services/api'
import { useWebSocketStore } from './websocket'

const parseUser = (val) => {
    if (!val) return null
    try {
        return JSON.parse(val)
    } catch {
        return null
    }
}

const isTokenExpired = (tokenStr) => {
    if (!tokenStr) return true
    try {
        const parts = tokenStr.split('.')
        if (parts.length !== 3) {
            return true
        }
        const base64Url = parts[1]
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)).join(''))
        const payload = JSON.parse(jsonPayload)
        if (payload && typeof payload.exp === 'number') {
            return Date.now() >= payload.exp * 1000
        }
        return false
    } catch {
        return true
    }
}

import { useGradesStore } from './grades'

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

// Minimal opaque storage payload to prevent PII leakage (email, names, tax_id) in web storage (localStorage / sessionStorage)
const toStorageUser = (userData) => {
    if (!userData) return null
    return {
        id: userData.id,
        role: userData.role || userData.user_role || null
    }
}

const getRoleFromToken = (tokenStr) => {
    if (!tokenStr) return null
    try {
        const parts = tokenStr.split('.')
        if (parts.length !== 3) return null
        const base64Url = parts[1]
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)).join(''))
        const payload = JSON.parse(jsonPayload)
        return payload ? (payload.role || payload.user_role || null) : null
    } catch {
        return null
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
            const wsStore = useWebSocketStore()
            wsStore.disconnect(true)
        } catch {
            // WebSocket store not initialized or already closed
        }

        try {
            const gradesStore = useGradesStore()
            gradesStore.clearCache()
        } catch {
            // Grades store not initialized
        }

        clearLocalSession()
        resetApiState()
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
        if (token.value && !isTokenExpired(token.value)) return

        const hasSavedUser = !!(localStorage.getItem('user') || sessionStorage.getItem('user'))
        if (!hasSavedUser) return

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
                console.warn('Initial session restore failed:', err)
                logout()
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

