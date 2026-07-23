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

export const useAuthStore = defineStore('auth', () => {
    const user = ref(parseUser(sessionStorage.getItem('user')) || parseUser(localStorage.getItem('user')) || null)
    const token = ref(sessionStorage.getItem('token') || localStorage.getItem('token') || null)
    const refreshToken = ref(sessionStorage.getItem('refreshToken') || localStorage.getItem('refreshToken') || null)

    const isAuthenticated = computed(() => !!token.value)
    const userRole = computed(() => user.value?.role || null)
    const userName = computed(() => {
        if (!user.value) return 'User'
        return `${user.value.first_name || user.value.firstName || ''} ${user.value.last_name || user.value.lastName || ''}`.trim() || 'User'
    })

    function login(userData, tokenData, refreshTokenData, rememberMe = true) {
        user.value = userData
        token.value = tokenData
        refreshToken.value = refreshTokenData

        if (rememberMe) {
            localStorage.setItem('user', JSON.stringify(userData))
            localStorage.setItem('token', tokenData)
            if (refreshTokenData) {
                localStorage.setItem('refreshToken', refreshTokenData)
            }
            sessionStorage.removeItem('user')
            sessionStorage.removeItem('token')
            sessionStorage.removeItem('refreshToken')
        } else {
            sessionStorage.setItem('user', JSON.stringify(userData))
            sessionStorage.setItem('token', tokenData)
            if (refreshTokenData) {
                sessionStorage.setItem('refreshToken', refreshTokenData)
            }
            localStorage.removeItem('user')
            localStorage.removeItem('token')
            localStorage.removeItem('refreshToken')
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
        user.value = userData
        if (localStorage.getItem('token')) {
            localStorage.setItem('user', JSON.stringify(userData))
        } else if (sessionStorage.getItem('token')) {
            sessionStorage.setItem('user', JSON.stringify(userData))
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
        updateUser
    }
})
