import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import authService from '@/services/authService'

export function useAuth() {
    const authStore = useAuthStore()
    const router = useRouter()
    const { user, isAuthenticated, refreshToken } = storeToRefs(authStore)

    async function login(email, password, rememberMe = false) {
        try {
            const response = await authService.login(email, password)
            const { user: userData, access_token: token, refresh_token: refreshTokenValue } = response
            authStore.login(userData, token, refreshTokenValue, rememberMe)

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

    async function logout() {
        try {
            // Call backend logout API if we have a refresh token
            if (refreshToken.value) {
                await authService.logout(refreshToken.value)
            }
        } catch (error) {
            console.error('Logout API error:', error)
            // Continue with local logout even if API call fails
        } finally {
            // Clear local state and redirect
            authStore.logout()
            router.push('/login')
        }
    }

    return {
        user,
        isAuthenticated,
        login,
        logout
    }
}
