import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import FamilyDesk from '@/pages/parent/FamilyDesk.vue'
import { familyDeskService } from '@/services/familyDeskService'

vi.mock('@/services/familyDeskService', () => ({
    familyDeskService: {
        getRequests: vi.fn(),
        getDelegates: vi.fn(),
        submitRequest: vi.fn()
    }
}))

describe('FamilyDesk.vue', () => {
    it('renders family requests and active delegates', async () => {
        familyDeskService.getRequests.mockResolvedValueOnce({
            data: [
                {
                    id: 'req-1',
                    request_type: 'delega_ritiro',
                    status: 'approved',
                    protocol_number: 'PROT-2026-0001',
                    created_at: '2026-10-04T08:00:00Z'
                }
            ]
        })
        familyDeskService.getDelegates.mockResolvedValueOnce({
            data: [
                {
                    id: 'del-1',
                    first_name: 'Giovanni',
                    last_name: 'Bianchi',
                    tax_code: 'BNCGNN60A01H501U',
                    relationship: 'nonno',
                    phone: '+393331122334',
                    id_card_details: 'CI AB12345'
                }
            ]
        })

        const wrapper = mount(FamilyDesk, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-icon': true,
                    'q-btn': true,
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-separator': true,
                    'q-badge': true,
                    'q-chip': true,
                    'q-list': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-avatar': true,
                    'q-tabs': true,
                    'q-tab': true,
                    'q-table': { template: '<div><slot name="body-cell-request_type" :row="{ request_type: \'delega_ritiro\' }" /></div>' },
                    'q-dialog': true
                }
            }
        })

        await flushPromises()
        expect(familyDeskService.getRequests).toHaveBeenCalled()
        expect(familyDeskService.getDelegates).toHaveBeenCalled()
        expect(wrapper.vm.requests).toHaveLength(1)
        expect(wrapper.vm.delegates).toHaveLength(1)
        expect(wrapper.vm.delegates[0].first_name).toBe('Giovanni')
    })
})
