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
            const { user: userData, access_token: token } = response.data
            authStore.login(userData, token)

            // Redirect based on role
            switch (userData.role) {
                case 'admin':
                    router.push('/admin')
                    break
                case 'secretary':
                    router.push('/secretary')
                    break
                case 'teacher':
                    router.push('/teacher')
                    break
                case 'student':
                    router.push('/student')
                    break
                case 'parent':
                    router.push('/parent')
                    break
                default:
                    router.push('/')
            }
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
