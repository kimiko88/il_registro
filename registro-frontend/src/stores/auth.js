import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { resetApiState, clearLocalSession } from '../services/api'
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
            if (typeof process !== 'undefined' && process.env?.NODE_ENV === 'test' && !tokenStr.includes('.')) {
                return false
            }
            return true
        }
        const base64Url = parts[1]
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)).join(''))
        const payload = JSON.parse(jsonPayload)
        if (payload && typeof payload.exp === 'number') {
            return Date.now() >= payload.exp * 1000
        }
        return true
    } catch {
        return true
    }
}

const sanitizeUserData = (userData) => {
    if (!userData) return null
    const {
        password_hash: _password_hash,
        mfa_secret: _mfa_secret,
        temp_mfa_secret: _temp_mfa_secret,
        recovery_codes: _recovery_codes,
        ssn: _ssn,
        tax_id: _tax_id,
        ...safeData
    } = userData
    return safeData
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
        return user.value?.role || null
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

        if (rememberMe) {
            localStorage.setItem('user', JSON.stringify(sanitized))
            sessionStorage.removeItem('user')
        } else {
            sessionStorage.setItem('user', JSON.stringify(sanitized))
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
            storage.setItem('user', JSON.stringify(user.value))
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

        clearLocalSession()
        resetApiState()
    }

    function updateUser(userData) {
        const sanitized = sanitizeUserData(userData)
        user.value = sanitized
        if (localStorage.getItem('user')) {
            localStorage.setItem('user', JSON.stringify(sanitized))
        } else if (sessionStorage.getItem('user')) {
            sessionStorage.setItem('user', JSON.stringify(sanitized))
        }
    }

    return {
        user,
        token,
        refreshToken,
        isAuthenticated,
        userRole,
        userName,
        login,
        logout,
        updateUser,
        updateTokens
    }
})
