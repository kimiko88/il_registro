import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import authService from '@/services/authService'

import { i18n } from '@/i18n'

function getTranslation(key, fallback) {
    if (i18n?.global?.t && i18n.global.te && i18n.global.te(key)) {
        return i18n.global.t(key)
    }
    return fallback
}

export function useAuth() {
    const authStore = useAuthStore()
    const router = useRouter()
    const { user, isAuthenticated } = storeToRefs(authStore)

    async function login(email, password, rememberMe = false) {
        try {
            const response = await authService.login(email, password)
            const { user: userData, access_token: token, refresh_token: refreshTokenValue } = response
            authStore.login(userData, token, refreshTokenValue, rememberMe)

            // Redirect based on role
            switch (userData.role) {
                case 'superadmin':
                case 'admin':
                case 'system_auditor':
                    router.push('/admin')
                    break
                case 'secretary':
                case 'principal':
                case 'vice_principal':
                case 'staff':
                    router.push('/secretary')
                    break
                case 'teacher':
                case 'coordinator':
                case 'docente':
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
        if (!err) return getTranslation('errors.serverError', 'Errore durante l\'accesso. Riprova.')
        const errData = err.response?.data
        const code = errData?.code
        if (code && i18n?.global?.te && i18n.global.te(`errors.${code}`)) {
            return i18n.global.t(`errors.${code}`)
        }

        const raw = String(errData?.error || err.message || err).toLowerCase()

        if (raw.includes('invalid credentials') || raw.includes('invalid email') || raw.includes('password')) {
            return getTranslation('errors.invalidCredentials', 'Email o password non corrette.')
        }
        if (raw.includes('too many') || raw.includes('rate limit') || raw.includes('429')) {
            return getTranslation('errors.rateLimit', 'Troppi tentativi di accesso. Riprova tra qualche minuto.')
        }
        if (raw.includes('unauthorized')) {
            return getTranslation('errors.unauthorized', 'Sessione non valida o scaduta.')
        }
        if (raw.includes('not found')) {
            return getTranslation('errors.userNotFound', 'Utente non trovato.')
        }
        if (raw.includes('disabled') || raw.includes('suspended') || raw.includes('blocked')) {
            return getTranslation('errors.accountDisabled', 'Account disabilitato o sospeso. Contatta la Segreteria.')
        }
        if (raw.includes('network') || raw.includes('failed to fetch') || raw.includes('timeout')) {
            return getTranslation('errors.connectionError', 'Errore di connessione al server. Verifica la tua connessione e riprova.')
        }
        return errData?.error || errData?.message || getTranslation('errors.serverError', 'Si è verificato un errore sul server. Riprova più tardi.')
    }

    async function logout() {
        try {
            await authService.logout()
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
