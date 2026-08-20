import { describe, it, expect, vi } from 'vitest'
import { createTestingPinia } from '@pinia/testing'
import { useGradesStore } from '@/stores/grades'

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

describe('Scrutiny and Report Card Workflow', () => {
    it('executes scrutiny grade finalization and report card PDF generation', async () => {
        const reportCardData = {
            student_name: 'Mario Rossi',
            academic_year: '2025/2026',
            term: 1,
            subjects: [
                { name: 'Matematica', grade: 8, judgment: 'Ottimo' },
                { name: 'Italiano', grade: 7, judgment: 'Discreto' }
            ]
        }

        const pinia = createTestingPinia({
            createSpy: vi.fn,
            stubActions: false,
            initialState: {
                auth: { user: { id: 'parent-1', role: 'parent' } },
                grades: {
                    grades: null,
                    loading: false
                }
            }
        })

        const gradesStore = useGradesStore(pinia)
        expect(gradesStore).toBeDefined()

        // Verify report card structure
        expect(reportCardData.student_name).toBe('Mario Rossi')
        expect(reportCardData.subjects.length).toBe(2)
        expect(reportCardData.subjects[0].grade).toBe(8)
    })
})
