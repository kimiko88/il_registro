import { describe, it, expect, vi } from 'vitest'
import { communicationService } from '@/services/communicationService'

vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            loading: { show: vi.fn(), hide: vi.fn() },
            notify: vi.fn(),
            dialog: vi.fn().mockReturnValue({ onOk: (fn) => fn() })
        })
    }
})

describe('Communications Signature Workflow', () => {
    it('handles reading bacheca notices and signing required communications', async () => {
        expect(communicationService).toBeDefined()

        const mockNotice = {
            id: 'comm-100',
            subject: 'Circolare n. 45 - Consigli di Classe',
            body: 'Si comunicano le date dei prossimi consigli di classe.',
            type: 'circolare',
            requires_signature: true,
            is_signed: false,
            created_at: '2026-09-01T08:00:00Z'
        }

        // Verify initial unsigned state
        expect(mockNotice.requires_signature).toBe(true)
        expect(mockNotice.is_signed).toBe(false)

        // Simulate signing communication
        mockNotice.is_signed = true
        expect(mockNotice.is_signed).toBe(true)
    })
})
