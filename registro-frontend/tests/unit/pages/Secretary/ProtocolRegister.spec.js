import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import ProtocolRegister from '@/pages/secretary/ProtocolRegister.vue'
import { protocolService } from '@/services/protocolService'

vi.mock('@/services/protocolService', () => ({
    protocolService: {
        getEntries: vi.fn(),
        protocolDocument: vi.fn()
    }
}))

describe('ProtocolRegister.vue', () => {
    it('renders official AgID protocol register entries', async () => {
        protocolService.getEntries.mockResolvedValueOnce({
            data: [
                {
                    id: 'prot-1',
                    protocol_number: 14,
                    protocol_year: 2026,
                    protocol_date: '2026-10-04T09:30:00Z',
                    flow_direction: 'in',
                    subject: 'Richiesta accesso agli atti L. 241/90',
                    sender: 'Avv. Bianchi',
                    recipient: 'Dirigente Scolastico',
                    classification_title: 1,
                    classification_class: '3'
                }
            ]
        })

        const wrapper = mount(ProtocolRegister, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-tabs': true,
                    'q-tab': true,
                    'q-icon': true,
                    'q-btn': true,
                    'q-badge': true,
                    'q-table': true,
                    'q-dialog': true
                }
            }
        })

        await flushPromises()
        expect(protocolService.getEntries).toHaveBeenCalled()
        expect(wrapper.vm.entries).toHaveLength(1)
        expect(wrapper.vm.entries[0].protocol_number).toBe(14)
        expect(wrapper.vm.entries[0].flow_direction).toBe('in')
    })
})
