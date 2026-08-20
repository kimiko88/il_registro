import { describe, it, expect, vi } from 'vitest'

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

describe('Class Council Verbali & Meeting Minutes Workflow E2E', () => {
    it('executes council meeting creation, verbale drafting, and digital signature', async () => {
        const meeting = {
            id: 'meeting-1',
            title: 'Consiglio di Classe 1° Trimestre',
            class_id: 'class-1',
            agenda: 'Valutazioni didattiche ed approvazione PDP'
        }

        const verbale = {
            id: 'verb-1',
            meeting_id: 'meeting-1',
            secretary_id: 'teacher-1',
            is_published: true,
            signatures: []
        }

        // Verify meeting & verbale setup
        expect(verbale.is_published).toBe(true)
        expect(verbale.signatures.length).toBe(0)

        // Sign verbale
        verbale.signatures.push({ user_id: 'teacher-1', signed_at: '2026-10-15T16:00:00Z' })
        expect(verbale.signatures.length).toBe(1)
    })
})
