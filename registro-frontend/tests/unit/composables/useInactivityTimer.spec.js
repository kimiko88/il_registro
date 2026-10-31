import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { nextTick } from 'vue'
import { createPinia, setActivePinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'

describe('useInactivityTimer', () => {
    let authStore

    beforeEach(() => {
        vi.useFakeTimers()
        setActivePinia(createPinia())
        authStore = useAuthStore()
        // Simulate authenticated state
        authStore.login({ email: 'test@scuola.it', role: 'teacher' }, 'fake-token', 'fake-refresh', false)
        vi.resetModules()
    })

    afterEach(() => {
        vi.useRealTimers()
        vi.clearAllMocks()
    })

    it('initialises with warning dialog hidden', async () => {
        const { useInactivityTimer } = await import('@/composables/useInactivityTimer')
        const { showWarningDialog, secondsRemaining } = useInactivityTimer(3600)
        expect(showWarningDialog.value).toBe(false)
        expect(secondsRemaining.value).toBe(60)
    })

    it('stayLoggedIn hides the warning dialog and resets countdown', async () => {
        const { useInactivityTimer } = await import('@/composables/useInactivityTimer')
        const { showWarningDialog, secondsRemaining, stayLoggedIn } = useInactivityTimer(3600)

        showWarningDialog.value = true
        secondsRemaining.value = 10
        stayLoggedIn()

        expect(showWarningDialog.value).toBe(false)
        expect(secondsRemaining.value).toBe(60)
    })

    it('logoutNow hides dialog and calls authStore.logout()', async () => {
        const { useInactivityTimer } = await import('@/composables/useInactivityTimer')
        const logoutSpy = vi.spyOn(authStore, 'logout')
        const { showWarningDialog, logoutNow } = useInactivityTimer(3600)

        showWarningDialog.value = true
        logoutNow()

        expect(showWarningDialog.value).toBe(false)
        expect(logoutSpy).toHaveBeenCalledOnce()
    })

    it('resetActivity resets lastActivity and closes dialog if open', async () => {
        const { useInactivityTimer } = await import('@/composables/useInactivityTimer')
        const { showWarningDialog, resetActivity } = useInactivityTimer(3600)

        showWarningDialog.value = true
        resetActivity()

        expect(showWarningDialog.value).toBe(false)
    })

    it('does not show warning for unauthenticated users', async () => {
        authStore.logout()
        const { useInactivityTimer } = await import('@/composables/useInactivityTimer')
        const { showWarningDialog } = useInactivityTimer(10)

        vi.advanceTimersByTime(60_000)
        await nextTick()

        expect(showWarningDialog.value).toBe(false)
    })
})
