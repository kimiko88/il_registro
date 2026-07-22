import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
    const user = ref(JSON.parse(sessionStorage.getItem('user') || localStorage.getItem('user')) || null)
    const token = ref(sessionStorage.getItem('token') || localStorage.getItem('token') || null)
    const refreshToken = ref(sessionStorage.getItem('refreshToken') || localStorage.getItem('refreshToken') || null)

    const isAuthenticated = computed(() => !!token.value)
    const userRole = computed(() => user.value?.role || null)
    const userName = computed(() => {
        if (!user.value) return 'User'
        return `${user.value.first_name || user.value.firstName || ''} ${user.value.last_name || user.value.lastName || ''}`.trim() || 'User'
    })

    function login(userData, tokenData, refreshTokenData, rememberMe = false) {
        user.value = userData
        token.value = tokenData
        refreshToken.value = refreshTokenData

        const storage = rememberMe ? localStorage : sessionStorage
        storage.setItem('user', JSON.stringify(userData))
        storage.setItem('token', tokenData)
        if (refreshTokenData) {
            storage.setItem('refreshToken', refreshTokenData)
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
        if (localStorage.getItem('user')) {
            localStorage.setItem('user', JSON.stringify(userData))
        }
        if (sessionStorage.getItem('user')) {
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
