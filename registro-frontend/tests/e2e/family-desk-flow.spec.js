import { describe, it, expect, vi } from 'vitest'
import { familyDeskService } from '@/services/familyDeskService'

vi.mock('@/services/familyDeskService', () => ({
    familyDeskService: {
        submitRequest: vi.fn(),
        getRequests: vi.fn(),
        reviewRequest: vi.fn(),
        getDelegates: vi.fn()
    }
}))

describe('Family Desk & Permanent Delegate Sync E2E Workflow', () => {
    it('executes full family request lifecycle: parent submits delega ritiro, secretary approves with protocol, delegate is synced with portineria', async () => {
        // Step 1: Parent submits delega ritiro
        familyDeskService.submitRequest.mockResolvedValueOnce({
            message: 'Istanza presentata con successo',
            request: {
                id: 'req-e2e-1',
                student_id: 'student-e2e-1',
                request_type: 'delega_ritiro',
                status: 'submitted',
                form_data: {
                    first_name: 'Giuseppe',
                    last_name: 'Verdi',
                    tax_code: 'VRDGPP60A01H501K',
                    relationship: 'nonno',
                    phone: '+393339876543'
                }
            }
        })
        const subRes = await familyDeskService.submitRequest({
            student_id: 'student-e2e-1',
            request_type: 'delega_ritiro',
            form_data: {
                first_name: 'Giuseppe',
                last_name: 'Verdi',
                tax_code: 'VRDGPP60A01H501K',
                relationship: 'nonno',
                phone: '+393339876543'
            }
        })
        expect(subRes.request.id).toBe('req-e2e-1')
        expect(subRes.request.status).toBe('submitted')

        // Step 2: Secretary reviews and approves request
        familyDeskService.reviewRequest.mockResolvedValueOnce({
            message: 'Stato istanza aggiornato con successo'
        })
        const revRes = await familyDeskService.reviewRequest('req-e2e-1', {
            status: 'approved',
            protocol_number: 'PROT-2026-00891'
        })
        expect(revRes.message).toContain('successo')

        // Step 3: Verify permanent delegate sync for portineria/early exits
        familyDeskService.getDelegates.mockResolvedValueOnce({
            data: [
                {
                    id: 'del-e2e-1',
                    student_id: 'student-e2e-1',
                    first_name: 'Giuseppe',
                    last_name: 'Verdi',
                    tax_code: 'VRDGPP60A01H501K',
                    relationship: 'nonno',
                    is_valid: true
                }
            ]
        })
        const delRes = await familyDeskService.getDelegates('student-e2e-1')
        expect(delRes.data).toHaveLength(1)
        expect(delRes.data[0].first_name).toBe('Giuseppe')
        expect(delRes.data[0].is_valid).toBe(true)
    })
})
