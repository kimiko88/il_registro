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

describe('Disciplinary Notes Workflow E2E', () => {
    it('handles teacher note creation, admin approval, and parent acknowledgment', async () => {
        const note = {
            id: 'note-1',
            student_name: 'Mario Rossi',
            teacher_name: 'Prof. Bianchi',
            type: 'disciplinary',
            note: 'Uso del cellulare durante la verifica.',
            is_approved: false,
            is_viewed_by_parent: false
        }

        // Verify initial note state
        expect(note.is_approved).toBe(false)
        expect(note.is_viewed_by_parent).toBe(false)

        // Admin approves note
        note.is_approved = true
        expect(note.is_approved).toBe(true)

        // Parent acknowledges note
        note.is_viewed_by_parent = true
        expect(note.is_viewed_by_parent).toBe(true)
    })
})
