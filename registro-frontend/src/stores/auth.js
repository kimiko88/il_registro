import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
    const user = ref(JSON.parse(localStorage.getItem('user')) || null)
    const token = ref(localStorage.getItem('token') || null)
    const refreshToken = ref(localStorage.getItem('refreshToken') || null)

    const isAuthenticated = computed(() => !!token.value)
    const userRole = computed(() => user.value?.role || null)
    const userName = computed(() => {
        if (!user.value) return 'User'
        return `${user.value.first_name || user.value.firstName || ''} ${user.value.last_name || user.value.lastName || ''}`.trim() || 'User'
    })

    function login(userData, tokenData, refreshTokenData) {
        user.value = userData
        token.value = tokenData
        refreshToken.value = refreshTokenData

        localStorage.setItem('user', JSON.stringify(userData))
        localStorage.setItem('token', tokenData)
        localStorage.setItem('refreshToken', refreshTokenData)
    }

    function logout() {
        user.value = null
        token.value = null
        refreshToken.value = null

        localStorage.removeItem('user')
        localStorage.removeItem('token')
        localStorage.removeItem('refreshToken')
    }

    function updateUser(userData) {
        user.value = userData
        localStorage.setItem('user', JSON.stringify(userData))
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
