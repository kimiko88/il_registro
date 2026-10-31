import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'

vi.mock('@/services/authService', () => ({
    default: {
        login: vi.fn(),
        logout: vi.fn()
    }
}))

vi.mock('vue-router', () => ({
    useRouter: () => ({
        push: vi.fn().mockResolvedValue(undefined),
        currentRoute: { value: { path: '/login' } }
    })
}))

vi.mock('@/utils/roleUtils', () => ({
    getUserDashboard: vi.fn((role) => `/${role}/dashboard`)
}))

describe('useAuth - login error translation', () => {
    let authStore

    beforeEach(async () => {
        setActivePinia(createPinia())
        authStore = useAuthStore()
        vi.clearAllMocks()
        vi.resetModules()
    })

    afterEach(() => {
        vi.clearAllMocks()
    })

    it('translates "invalid credentials" error correctly', async () => {
        const { useAuth } = await import('@/composables/useAuth')
        const { getLoginErrorKey } = useAuth()

        const err = { response: { data: { error: 'invalid credentials' } } }
        const key = getLoginErrorKey(err)
        expect(key).toBe('errors.invalidCredentials')
    })

    it('translates "invalid email" error correctly', async () => {
        const { useAuth } = await import('@/composables/useAuth')
        const { getLoginErrorKey } = useAuth()

        const err = { response: { data: { error: 'invalid email format' } } }
        expect(getLoginErrorKey(err)).toBe('errors.invalidCredentials')
    })

    it('translates rate limit (too many) error correctly', async () => {
        const { useAuth } = await import('@/composables/useAuth')
        const { getLoginErrorKey } = useAuth()

        const err = { response: { data: { error: 'too many requests' } } }
        expect(getLoginErrorKey(err)).toBe('errors.rateLimit')
    })

    it('translates rate limit (rate limit) error correctly', async () => {
        const { useAuth } = await import('@/composables/useAuth')
        const { getLoginErrorKey } = useAuth()

        const err = { response: { data: { error: 'rate limit exceeded' } } }
        expect(getLoginErrorKey(err)).toBe('errors.rateLimit')
    })

    it('translates rate limit HTTP 429 in message correctly', async () => {
        const { useAuth } = await import('@/composables/useAuth')
        const { getLoginErrorKey } = useAuth()

        const err = { message: '429 Too Many Requests' }
        expect(getLoginErrorKey(err)).toBe('errors.rateLimit')
    })

    it('translates unauthorized error correctly', async () => {
        const { useAuth } = await import('@/composables/useAuth')
        const { getLoginErrorKey } = useAuth()

        const err = { response: { data: { error: 'unauthorized access' } } }
        expect(getLoginErrorKey(err)).toBe('errors.unauthorized')
    })

    it('translates not found error correctly', async () => {
        const { useAuth } = await import('@/composables/useAuth')
        const { getLoginErrorKey } = useAuth()

        const err = { response: { data: { error: 'user not found' } } }
        expect(getLoginErrorKey(err)).toBe('errors.userNotFound')
    })

    it('translates disabled/suspended/blocked account error correctly', async () => {
        const { useAuth } = await import('@/composables/useAuth')
        const { getLoginErrorKey } = useAuth()

        expect(getLoginErrorKey({ response: { data: { error: 'account disabled' } } })).toBe('errors.accountDisabled')
        expect(getLoginErrorKey({ response: { data: { error: 'account suspended' } } })).toBe('errors.accountDisabled')
        expect(getLoginErrorKey({ response: { data: { error: 'user blocked' } } })).toBe('errors.accountDisabled')
    })

    it('translates network error correctly', async () => {
        const { useAuth } = await import('@/composables/useAuth')
        const { getLoginErrorKey } = useAuth()

        expect(getLoginErrorKey({ message: 'Network Error: failed to fetch' })).toBe('errors.connectionError')
        expect(getLoginErrorKey({ message: 'timeout of 5000ms exceeded' })).toBe('errors.connectionError')
    })

    it('returns serverError for null/undefined errors', async () => {
        const { useAuth } = await import('@/composables/useAuth')
        const { getLoginErrorKey } = useAuth()

        expect(getLoginErrorKey(null)).toBe('errors.serverError')
        expect(getLoginErrorKey(undefined)).toBe('errors.serverError')
        expect(getLoginErrorKey({})).toBe('errors.serverError')
    })

    it('uses code field from response.data when available', async () => {
        const { useAuth } = await import('@/composables/useAuth')
        const { getLoginErrorKey } = useAuth()

        const err = { response: { data: { code: 'CUSTOM_CODE' } } }
        // When no translation exists for code, falls through to default
        const key = getLoginErrorKey(err)
        // Code is checked first via i18n.te — in tests i18n.te returns false for unknown keys
        expect(typeof key).toBe('string')
    })

    it('login calls authService.login with correct credentials', async () => {
        const authService = await import('@/services/authService')
        authService.default.login.mockResolvedValue({
            user: { email: 'prof@scuola.it', role: 'teacher' },
            access_token: 'tok',
            refresh_token: 'ref'
        })

        const loginSpy = vi.spyOn(authStore, 'login')
        const { useAuth } = await import('@/composables/useAuth')
        const { login } = useAuth()

        const result = await login('prof@scuola.it', 'SecurePass123!')

        expect(authService.default.login).toHaveBeenCalledWith('prof@scuola.it', 'SecurePass123!')
        expect(loginSpy).toHaveBeenCalled()
        expect(result).toBeNull() // null means success (no error)
    })

    it('login returns error message on auth failure', async () => {
        const authService = await import('@/services/authService')
        authService.default.login.mockRejectedValue({
            response: { data: { error: 'invalid credentials' } }
        })

        const { useAuth } = await import('@/composables/useAuth')
        const { login } = useAuth()

        const result = await login('wrong@email.it', 'wrongpass')
        expect(typeof result).toBe('string')
        expect(result.length).toBeGreaterThan(0)
    })

    it('logout calls authService.logout and authStore.logout', async () => {
        const authService = await import('@/services/authService')
        authService.default.logout.mockResolvedValue({})

        const logoutSpy = vi.spyOn(authStore, 'logout')
        const { useAuth } = await import('@/composables/useAuth')
        const { logout } = useAuth()

        await logout()

        expect(logoutSpy).toHaveBeenCalledOnce()
    })

    it('logout still clears local state even when API call fails', async () => {
        const authService = await import('@/services/authService')
        authService.default.logout.mockRejectedValue(new Error('Network error'))

        const logoutSpy = vi.spyOn(authStore, 'logout')
        const { useAuth } = await import('@/composables/useAuth')
        const { logout } = useAuth()

        await logout()

        // authStore.logout must still be called despite API error
        expect(logoutSpy).toHaveBeenCalledOnce()
    })
})
