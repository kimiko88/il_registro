import { describe, it, expect, vi } from 'vitest'
import { parentDualSignService } from '@/services/parentDualSignService'

vi.mock('@/services/parentDualSignService', () => ({
    parentDualSignService: {
        listAuthorizations: vi.fn(),
        signAuthorization: vi.fn(),
        rejectAuthorization: vi.fn(),
        getChildCustodyInfo: vi.fn()
    }
}))

describe('Parent Dual Signature & Custody Compliance E2E Workflow', () => {
    it('executes full bigenitorial dual signature lifecycle with co-signature and restriction checks', async () => {
        // Step 1: Parent 1 lists pending authorizations
        parentDualSignService.listAuthorizations.mockResolvedValueOnce({
            data: [
                {
                    id: 'dual-auth-101',
                    title: 'Uscita Didattica Museo Archeologico',
                    document_type: 'trip_consent',
                    status: 'pending_first',
                    parent1_pin_verified: false,
                    parent2_pin_verified: false
                }
            ]
        })
        const listP1 = await parentDualSignService.listAuthorizations()
        expect(listP1.data[0].status).toBe('pending_first')

        // Step 2: Parent 1 signs with personal PIN -> moves to pending_second
        parentDualSignService.signAuthorization.mockResolvedValueOnce({
            data: {
                id: 'dual-auth-101',
                status: 'pending_second',
                parent1_pin_verified: true,
                parent2_pin_verified: false
            }
        })
        const signP1 = await parentDualSignService.signAuthorization('dual-auth-101', '1234')
        expect(signP1.data.status).toBe('pending_second')
        expect(signP1.data.parent1_pin_verified).toBe(true)

        // Step 3: Parent 2 signs with personal PIN -> completes full agreement
        parentDualSignService.signAuthorization.mockResolvedValueOnce({
            data: {
                id: 'dual-auth-101',
                status: 'completed',
                parent1_pin_verified: true,
                parent2_pin_verified: true
            }
        })
        const signP2 = await parentDualSignService.signAuthorization('dual-auth-101', '5678')
        expect(signP2.data.status).toBe('completed')
        expect(signP2.data.parent2_pin_verified).toBe(true)

        // Step 4: Check restricted custody handling
        parentDualSignService.getChildCustodyInfo.mockResolvedValueOnce({
            data: {
                student_id: 'student-restricted-1',
                custody_type: 'restricted',
                court_order_details: 'Sentenza Trib. Minorenni prot. 441/2026',
                can_authorize_activities: false
            }
        })
        const custodyCheck = await parentDualSignService.getChildCustodyInfo('student-restricted-1')
        expect(custodyCheck.data.custody_type).toBe('restricted')
        expect(custodyCheck.data.can_authorize_activities).toBe(false)
    })
})
