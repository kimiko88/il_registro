import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useSessionReauth } from '@/composables/useSessionReauth'

describe('Session Reauth In-Page Flow', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
        sessionStorage.clear()
        vi.clearAllMocks()
    })

    it('triggers reauth without navigating away when token expires', async () => {
        const { showDialog, userEmail, triggerReauth, resolveReauth } = useSessionReauth()

        expect(showDialog.value).toBe(false)

        // Simulate API interceptor triggering reauth
        const reauthPromise = triggerReauth('teacher@school.it')

        expect(showDialog.value).toBe(true)
        expect(userEmail.value).toBe('teacher@school.it')

        // Simulate user successfully typing password and logging in
        resolveReauth('new-refreshed-jwt-token')

        const resultToken = await reauthPromise
        expect(resultToken).toBe('new-refreshed-jwt-token')
        expect(showDialog.value).toBe(false)
    })

    it('rejects reauth promise when user cancels reauth dialog', async () => {
        const { showDialog, triggerReauth, cancelReauth } = useSessionReauth()

        const reauthPromise = triggerReauth('admin@school.it')
        expect(showDialog.value).toBe(true)

        // User clicks "Cancel / Exit"
        cancelReauth()

        await expect(reauthPromise).rejects.toThrow('reauth_cancelled')
        expect(showDialog.value).toBe(false)
    })

    it('updates auth store token on reauth without resetting user state', () => {
        const authStore = useAuthStore()
        authStore.user = { id: 'u1', email: 'teacher@school.it', first_name: 'Prof', role: 'teacher' }
        authStore.token = 'old-token'

        // Reauth success updates token without wiping out user
        authStore.updateTokens('new-valid-token')

        expect(authStore.token).toBe('new-valid-token')
        expect(authStore.user.email).toBe('teacher@school.it')
        expect(authStore.user.role).toBe('teacher')
    })
})
