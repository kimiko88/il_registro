import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { resetApiState, clearLocalSession } from '../services/api'

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
            // Only allow non-JWT mock strings in test environment if explicitly prefixed
            if (typeof import.meta !== 'undefined' && import.meta.env && import.meta.env.MODE === 'test' && tokenStr.startsWith('mock-')) {
                return false
            }
            return true
        }
        const base64Url = parts[1]
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)).join(''))
        const payload = JSON.parse(jsonPayload)
        if (payload && payload.exp) {
            return Date.now() >= payload.exp * 1000
        }
        return false
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
    const token = ref(sessionStorage.getItem('token') || localStorage.getItem('token') || null)
    const refreshToken = ref(sessionStorage.getItem('refreshToken') || localStorage.getItem('refreshToken') || null)

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

    function login(userData, tokenData, refreshTokenData, rememberMe = false) {
        const sanitized = sanitizeUserData(userData)
        user.value = sanitized
        token.value = tokenData
        refreshToken.value = refreshTokenData

        const role = sanitized?.role
        const isPrivileged = role === 'admin' || role === 'superadmin' || role === 'teacher' || role === 'secretary' || role === 'principal' || role === 'vice_principal'
        const shouldRemember = rememberMe && !isPrivileged

        if (shouldRemember) {
            localStorage.setItem('user', JSON.stringify(sanitized))
            localStorage.setItem('token', tokenData)
            if (refreshTokenData) {
                localStorage.setItem('refreshToken', refreshTokenData)
            } else {
                localStorage.removeItem('refreshToken')
            }
            sessionStorage.removeItem('user')
            sessionStorage.removeItem('token')
            sessionStorage.removeItem('refreshToken')
        } else {
            sessionStorage.setItem('user', JSON.stringify(sanitized))
            sessionStorage.setItem('token', tokenData)
            if (refreshTokenData) {
                sessionStorage.setItem('refreshToken', refreshTokenData)
            } else {
                sessionStorage.removeItem('refreshToken')
            }
            localStorage.removeItem('user')
            localStorage.removeItem('token')
            localStorage.removeItem('refreshToken')
        }
    }

    function updateTokens(newTokenData, newRefreshTokenData = undefined, newUserData = null) {
        if (!newTokenData) return
        token.value = newTokenData
        if (newRefreshTokenData !== undefined) {
            refreshToken.value = newRefreshTokenData
        }
        if (newUserData) {
            user.value = sanitizeUserData(newUserData)
        }

        const isLocal = !!(localStorage.getItem('token') || localStorage.getItem('user'))
        const storage = isLocal ? localStorage : sessionStorage

        storage.setItem('token', newTokenData)
        if (newRefreshTokenData !== undefined) {
            if (newRefreshTokenData) {
                storage.setItem('refreshToken', newRefreshTokenData)
            } else {
                storage.removeItem('refreshToken')
            }
        }
        if (newUserData && user.value) {
            storage.setItem('user', JSON.stringify(user.value))
        }
    }

    function logout() {
        user.value = null
        token.value = null
        refreshToken.value = null

        clearLocalSession()
        resetApiState()
    }

    function updateUser(userData) {
        const sanitized = sanitizeUserData(userData)
        user.value = sanitized
        if (localStorage.getItem('token')) {
            localStorage.setItem('user', JSON.stringify(sanitized))
        } else if (sessionStorage.getItem('token')) {
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
