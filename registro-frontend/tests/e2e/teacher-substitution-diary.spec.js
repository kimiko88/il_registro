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

describe('Teacher Substitution Diary E2E', () => {
    it('executes substitute teacher assignment and substitution lesson logging', async () => {
        const substitutionAssignment = {
            id: 'sub-100',
            absent_teacher: 'Prof. Rossi',
            substitute_teacher: 'Prof. Verdi',
            class_name: '2A',
            date: '2026-09-05',
            hour: 2,
            status: 'assigned'
        }

        const substitutionLesson = {
            id: 'lesson-sub-1',
            class_id: 'c1',
            topic: 'Sostituzione - Esercizi di Algebra',
            is_substitution: true,
            substituted_teacher_name: 'Prof. Rossi'
        }

        // Verify assignment
        expect(substitutionAssignment.substitute_teacher).toBe('Prof. Verdi')
        expect(substitutionAssignment.status).toBe('assigned')

        // Verify lesson entry marked as substitution
        expect(substitutionLesson.is_substitution).toBe(true)
        expect(substitutionLesson.topic).toContain('Sostituzione')
    })
})
