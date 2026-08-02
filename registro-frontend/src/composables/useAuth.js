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
                case 'superadmin':
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
            return translateLoginError(error)
        }
    }

    function translateLoginError(err) {
        if (!err) return 'Errore durante l\'accesso. Riprova.'
        const raw = String(err.response?.data?.error || err.message || err).toLowerCase()

        if (raw.includes('invalid credentials') || raw.includes('invalid email') || raw.includes('password')) {
            return 'Email o password non corrette'
        }
        if (raw.includes('too many') || raw.includes('rate limit') || raw.includes('429')) {
            return 'Troppi tentativi di accesso. Riprova tra qualche minuto.'
        }
        if (raw.includes('unauthorized')) {
            return 'Credenziali non autorizzate'
        }
        if (raw.includes('not found')) {
            return 'Utente non trovato'
        }
        if (raw.includes('disabled') || raw.includes('suspended') || raw.includes('blocked')) {
            return 'Account disabilitato o sospeso. Contatta la Segreteria.'
        }
        if (raw.includes('network') || raw.includes('failed to fetch') || raw.includes('timeout')) {
            return 'Impossibile connettersi al server. Verifica la connessione di rete.'
        }
        return err.response?.data?.error || 'Errore durante l\'accesso. Riprova.'
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
