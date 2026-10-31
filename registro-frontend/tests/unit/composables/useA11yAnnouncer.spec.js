import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

describe('useA11yAnnouncer Composable', () => {
    beforeEach(() => {
        vi.useFakeTimers()
        vi.resetModules()
    })

    afterEach(() => {
        vi.useRealTimers()
    })

    it('exposes politeMessage and assertiveMessage refs and announce function', async () => {
        const { useA11yAnnouncer } = await import('@/composables/useA11yAnnouncer')
        const announcer = useA11yAnnouncer()

        expect(typeof announcer.announce).toBe('function')
        expect(announcer.politeMessage).toBeDefined()
        expect(announcer.assertiveMessage).toBeDefined()
    })

    it('announces a polite message after the 50ms debounce delay', async () => {
        const { useA11yAnnouncer } = await import('@/composables/useA11yAnnouncer')
        const { announce, politeMessage } = useA11yAnnouncer()

        announce('Pagina caricata con successo', 'polite')
        // Before timer advances
        expect(politeMessage.value).toBe('')

        vi.advanceTimersByTime(50)
        expect(politeMessage.value).toBe('Pagina caricata con successo')
    })

    it('announces an assertive message after the 50ms delay and clears previous', async () => {
        const { useA11yAnnouncer } = await import('@/composables/useA11yAnnouncer')
        const { announce, assertiveMessage } = useA11yAnnouncer()

        announce('Attenzione: sessione in scadenza', 'assertive')
        expect(assertiveMessage.value).toBe('')

        vi.advanceTimersByTime(50)
        expect(assertiveMessage.value).toBe('Attenzione: sessione in scadenza')
    })

    it('defaults politeness to polite if omitted', async () => {
        const { useA11yAnnouncer } = await import('@/composables/useA11yAnnouncer')
        const { announce, politeMessage, assertiveMessage } = useA11yAnnouncer()

        announce('Operazione completata')
        vi.advanceTimersByTime(50)

        expect(politeMessage.value).toBe('Operazione completata')
        expect(assertiveMessage.value).not.toBe('Operazione completata')
    })

    it('ignores empty or null messages', async () => {
        const { useA11yAnnouncer } = await import('@/composables/useA11yAnnouncer')
        const { announce, politeMessage } = useA11yAnnouncer()

        politeMessage.value = 'Existing message'
        announce('')
        announce(null)
        announce(undefined)

        vi.advanceTimersByTime(100)
        expect(politeMessage.value).toBe('Existing message')
    })
})

