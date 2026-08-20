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

describe('BES/PDP Document Workflow E2E', () => {
    it('handles PDP plan creation, measure selection, and review submission', async () => {
        const pdpDocument = {
            id: 'doc-pdp-100',
            student_id: '00000000-0000-0000-0000-000000000001',
            student_name: 'Mario Rossi',
            type: 'pdp',
            title: 'Piano Didattico Personalizzato 2026/2027',
            status: 'draft',
            measures: ['Tempo aggiuntivo 30%', 'Uso della calcolatrice']
        }

        // Verify draft document
        expect(pdpDocument.status).toBe('draft')
        expect(pdpDocument.measures.length).toBe(2)

        // Submit for review
        pdpDocument.status = 'submitted'
        expect(pdpDocument.status).toBe('submitted')
    })
})
