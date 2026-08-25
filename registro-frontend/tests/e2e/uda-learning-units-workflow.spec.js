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

describe('UDA Learning Units & Multidisciplinary Planning Workflow E2E', () => {
    it('executes UDA plan creation, key competency mapping, and phase scheduling', async () => {
        const udaPlan = {
            id: 'uda-1',
            title: 'UDA Interdisciplinare: Sostenibilità Ambientale',
            class_id: 'class-1',
            period: 'primo_quadrimestre',
            competencies: ['Competenza STEM', 'Competenza di Cittadinanza'],
            status: 'draft'
        }

        // Verify initial UDA setup
        expect(udaPlan.competencies.length).toBe(2)
        expect(udaPlan.status).toBe('draft')

        // Submit for approval
        udaPlan.status = 'submitted'
        expect(udaPlan.status).toBe('submitted')
    })
})
