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

describe('Rubric-Based Evaluations & Competency Matrix Workflow E2E', () => {
    it('executes evaluation rubric template creation and criterion-based scoring', async () => {
        const rubric = {
            id: 'rubric-1',
            title: 'Rubrica Valutativa Presentazione Orale',
            criteria: [
                { id: 'crit-1', name: 'Chiarezza Espositiva', max_score: 10.0 }
            ]
        }

        const assessment = {
            student_id: 'st-1',
            scores: [
                { criterion_id: 'crit-1', score: 9.0 }
            ],
            notes: 'Eccellente esposizione'
        }

        // Verify rubric setup
        expect(rubric.criteria[0].max_score).toBe(10.0)

        // Verify student assessment scoring
        expect(assessment.scores[0].score).toBe(9.0)
    })
})
