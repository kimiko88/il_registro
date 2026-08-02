import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

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
        if (parts.length !== 3) return false // Opaque or mock token in test/dev environment, assume valid
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

    function login(userData, tokenData, refreshTokenData, rememberMe = true) {
        const sanitized = sanitizeUserData(userData)
        user.value = sanitized
        token.value = tokenData
        refreshToken.value = refreshTokenData

        if (rememberMe) {
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

    function updateTokens(newTokenData, newRefreshTokenData = null) {
        if (!newTokenData) return
        token.value = newTokenData
        if (newRefreshTokenData) {
            refreshToken.value = newRefreshTokenData
        }

        if (localStorage.getItem('token') || localStorage.getItem('user')) {
            localStorage.setItem('token', newTokenData)
            if (newRefreshTokenData) {
                localStorage.setItem('refreshToken', newRefreshTokenData)
            }
        } else {
            sessionStorage.setItem('token', newTokenData)
            if (newRefreshTokenData) {
                sessionStorage.setItem('refreshToken', newRefreshTokenData)
            }
        }
    }

    function logout() {
        user.value = null
        token.value = null
        refreshToken.value = null

        localStorage.removeItem('user')
        localStorage.removeItem('token')
        localStorage.removeItem('refreshToken')
        localStorage.removeItem('selectedChildId')

        sessionStorage.removeItem('user')
        sessionStorage.removeItem('token')
        sessionStorage.removeItem('refreshToken')
        sessionStorage.removeItem('selectedChildId')
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
