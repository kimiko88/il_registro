import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

// We mock quasar at module level — Notify.create is shared across tests
const mockNotifyCreate = vi.fn()

vi.mock('quasar', async () => {
    const actual = await vi.importActual('quasar')
    return {
        ...actual,
        Notify: {
            create: mockNotifyCreate
        }
    }
})

describe('useErrorHandler', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        mockNotifyCreate.mockClear()
        vi.resetModules()
    })

    async function setup() {
        const { useErrorHandler } = await import('@/composables/useErrorHandler')
        return useErrorHandler()
    }

    it('handleError returns the error message', async () => {
        const { handleError } = await setup()
        const err = new Error('Something went wrong')
        const msg = handleError(err, null, false)
        expect(msg).toBeTruthy()
    })

    it('handleError prefers err.userMessage over generic message', async () => {
        const { handleError } = await setup()
        const err = { userMessage: 'Custom user message', message: 'Internal' }
        const msg = handleError(err, null, false)
        expect(msg).toBe('Custom user message')
    })

    it('handleError uses response.data.error when available', async () => {
        const { handleError } = await setup()
        const err = { response: { data: { error: 'Server-side error text' } } }
        const msg = handleError(err, null, false)
        expect(msg).toBe('Server-side error text')
    })

    it('handleError falls back to fallbackMessage when error has no message', async () => {
        const { handleError } = await setup()
        const err = {}
        const msg = handleError(err, 'Fallback message', false)
        expect(msg).toBe('Fallback message')
    })

    it('handleError shows a Notify toast by default (showToast=true)', async () => {
        const { handleError } = await setup()

        handleError(new Error('test'), null, true)

        expect(mockNotifyCreate).toHaveBeenCalledOnce()
        const call = mockNotifyCreate.mock.calls[0][0]
        expect(call.color).toBe('negative')
        expect(call.icon).toBe('error')
    })

    it('handleError does NOT show Notify when showToast=false', async () => {
        const { handleError } = await setup()

        handleError(new Error('silent error'), null, false)
        expect(mockNotifyCreate).not.toHaveBeenCalled()
    })

    it('handleError Notify toast has negative color', async () => {
        const { handleError } = await setup()

        handleError(new Error('a11y test'), null, true)
        const callArg = mockNotifyCreate.mock.calls[0][0]
        expect(callArg.color).toBe('negative')
        expect(callArg.icon).toBe('error')
        expect(callArg.timeout).toBeGreaterThan(0)
    })

    it('handleError reports error to errorStore', async () => {
        const { useErrorStore } = await import('@/stores/error')
        const errorStore = useErrorStore()
        const reportSpy = vi.spyOn(errorStore, 'reportError')

        const { handleError } = await setup()
        const err = new Error('store error')
        handleError(err, null, false)

        expect(reportSpy).toHaveBeenCalledWith(err, expect.any(String), false)
    })
})
