import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
    const user = ref(JSON.parse(localStorage.getItem('user')) || null)
    const token = ref(localStorage.getItem('token') || null)

    const isAuthenticated = computed(() => !!token.value)

    function login(userData, tokenData) {
        user.value = userData
        token.value = tokenData
        localStorage.setItem('user', JSON.stringify(userData))
        localStorage.setItem('token', tokenData)
    }

    function logout() {
        user.value = null
        token.value = null
        localStorage.removeItem('user')
        localStorage.removeItem('token')
    }

    return { user, token, isAuthenticated, login, logout }
})
