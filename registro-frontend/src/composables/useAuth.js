import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import api from '@/services/api'

export function useAuth() {
    const authStore = useAuthStore()
    const router = useRouter()
    const { user, isAuthenticated } = storeToRefs(authStore)

    async function login(email, password) {
        try {
            const response = await api.post('/auth/login', { email, password })
            const { user: userData, token } = response.data
            authStore.login(userData, token)
            router.push('/')
            return null
        } catch (error) {
            return error.response?.data?.error || 'Login failed'
        }
    }

    function logout() {
        authStore.logout()
        router.push('/login')
    }

    return {
        user,
        isAuthenticated,
        login,
        logout
    }
}
