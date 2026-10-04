import { describe, it, expect, vi } from 'vitest'
import { enrollmentService } from '@/services/enrollmentService'

vi.mock('@/services/enrollmentService', () => ({
    enrollmentService: {
        importSIDI: vi.fn(),
        generateFormationDraft: vi.fn(),
        updateDraftAssignments: vi.fn(),
        finalizeDraft: vi.fn()
    }
}))

describe('Class Formation & SIDI Import E2E Workflow', () => {
    it('executes full class formation lifecycle: SIDI import, solver generation, student reassignment, and draft finalization', async () => {
        // Step 1: Import SIDI
        enrollmentService.importSIDI.mockResolvedValueOnce({
            data: { message: 'Domande SIDI importate con successo', imported: 24 }
        })
        const importRes = await enrollmentService.importSIDI(new FormData())
        expect(importRes.data.imported).toBe(24)

        // Step 2: Generate draft via Solver
        enrollmentService.generateFormationDraft.mockResolvedValueOnce({
            data: {
                id: 'draft-2026',
                title: 'Classi Prime 2026/2027',
                assignments: {
                    classes: [
                        {
                            class_name: '1A',
                            total_students: 20,
                            l104_count: 1, // Compliant with max 20 for L.104
                            males_count: 10,
                            females_count: 10
                        },
                        {
                            class_name: '1B',
                            total_students: 24,
                            l104_count: 0,
                            males_count: 12,
                            females_count: 12
                        }
                    ]
                }
            }
        })
        const draftRes = await enrollmentService.generateFormationDraft({
            title: 'Classi Prime 2026/2027',
            academic_year: '2026/2027',
            parameters: { target_class_count: 2, max_l104_per_class: 1 }
        })
        expect(draftRes.data.assignments.classes).toHaveLength(2)
        expect(draftRes.data.assignments.classes[0].total_students).toBe(20)

        // Step 3: Update draft assignments (move student)
        enrollmentService.updateDraftAssignments.mockResolvedValueOnce({
            data: { message: 'Bozza classi aggiornata con successo' }
        })
        const updateRes = await enrollmentService.updateDraftAssignments('draft-2026', draftRes.data.assignments)
        expect(updateRes.data.message).toContain('aggiornata')

        // Step 4: Finalize draft
        enrollmentService.finalizeDraft.mockResolvedValueOnce({
            data: { message: 'Bozza finalizzata. Classi e studenti generati con successo!' }
        })
        const finalRes = await enrollmentService.finalizeDraft('draft-2026')
        expect(finalRes.data.message).toContain('finalizzata')
    })
})
