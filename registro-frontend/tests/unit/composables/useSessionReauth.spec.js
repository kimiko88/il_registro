import { describe, it, expect, vi, beforeEach } from 'vitest'

// Reset singleton state between tests by re-importing with resetModules
describe('useSessionReauth', () => {
    beforeEach(() => {
        vi.resetModules()
    })

    it('opens the dialog and sets the email when triggerReauth is called', async () => {
        const { useSessionReauth } = await import('@/composables/useSessionReauth')
        const { showDialog, userEmail, triggerReauth } = useSessionReauth()

        expect(showDialog.value).toBe(false)
        triggerReauth('docente@scuola.it')
        expect(showDialog.value).toBe(true)
        expect(userEmail.value).toBe('docente@scuola.it')
    })

    it('triggerReauth returns a Promise that resolves with resolveReauth', async () => {
        const { useSessionReauth } = await import('@/composables/useSessionReauth')
        const { triggerReauth, resolveReauth } = useSessionReauth()

        const prom = triggerReauth('a@b.it')
        resolveReauth('new-access-token')
        await expect(prom).resolves.toBe('new-access-token')
    })

    it('triggerReauth returns a Promise that rejects with cancelReauth', async () => {
        const { useSessionReauth } = await import('@/composables/useSessionReauth')
        const { triggerReauth, cancelReauth } = useSessionReauth()

        const prom = triggerReauth('a@b.it')
        cancelReauth()
        await expect(prom).rejects.toThrow('reauth_cancelled')
    })

    it('resolveReauth closes the dialog', async () => {
        const { useSessionReauth } = await import('@/composables/useSessionReauth')
        const { showDialog, triggerReauth, resolveReauth } = useSessionReauth()

        const p = triggerReauth('user@school.it')
        expect(showDialog.value).toBe(true)

        resolveReauth('token-xyz')
        await p // consume the promise to avoid unhandled rejection
        expect(showDialog.value).toBe(false)
    })

    it('cancelReauth closes the dialog', async () => {
        const { useSessionReauth } = await import('@/composables/useSessionReauth')
        const { showDialog, triggerReauth, cancelReauth } = useSessionReauth()

        const p = triggerReauth('user@school.it')
        cancelReauth()
        await expect(p).rejects.toThrow() // consume the promise
        expect(showDialog.value).toBe(false)
    })

    it('resolves all queued promises (concurrent triggers)', async () => {
        const { useSessionReauth } = await import('@/composables/useSessionReauth')
        const { triggerReauth, resolveReauth } = useSessionReauth()

        const p1 = triggerReauth('a@b.it')
        const p2 = triggerReauth('a@b.it') // second concurrent trigger
        resolveReauth('shared-token')

        await expect(p1).resolves.toBe('shared-token')
        await expect(p2).resolves.toBe('shared-token')
    })

    it('rejects all queued promises on cancel', async () => {
        const { useSessionReauth } = await import('@/composables/useSessionReauth')
        const { triggerReauth, cancelReauth } = useSessionReauth()

        const p1 = triggerReauth('a@b.it')
        const p2 = triggerReauth('a@b.it')
        cancelReauth()

        await expect(p1).rejects.toThrow('reauth_cancelled')
        await expect(p2).rejects.toThrow('reauth_cancelled')
    })

    it('works with empty email (no email provided)', async () => {
        const { useSessionReauth } = await import('@/composables/useSessionReauth')
        const { userEmail, triggerReauth, cancelReauth } = useSessionReauth()

        const prom = triggerReauth()
        expect(userEmail.value).toBe('')
        cancelReauth()
        await expect(prom).rejects.toThrow()
    })
})
