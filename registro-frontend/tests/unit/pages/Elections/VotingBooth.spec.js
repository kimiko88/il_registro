import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import VotingBooth from '@/pages/elections/VotingBooth.vue'
import { electionsService } from '@/services/electionsService'

vi.mock('vue-router', () => ({
    useRoute: () => ({
        params: { id: 'elec-test-101' },
        query: {}
    })
}))

vi.mock('@/services/electionsService', () => ({
    electionsService: {
        getElectionDetails: vi.fn(),
        getElections: vi.fn(),
        castVote: vi.fn()
    }
}))

describe('VotingBooth.vue', () => {
    it('renders election title, lists, and allows selecting list', async () => {
        electionsService.getElectionDetails.mockResolvedValueOnce({
            election: {
                id: 'elec-test-101',
                title: 'Rinnovo Rappresentanti di Classe 2026',
                max_preferences: 1,
                is_closed: false
            },
            lists: [
                {
                    id: 'list-1',
                    list_number: 1,
                    motto: 'Impegno e Studio',
                    candidates: [
                        { id: 'c-1', first_name: 'Elena', last_name: 'Galli' }
                    ]
                }
            ]
        })

        const wrapper = mount(VotingBooth, {
            global: {
                directives: {
                    ripple: {}
                },
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-icon': true,
                    'q-btn': true,
                    'q-list': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-radio': true,
                    'q-checkbox': true,
                    'q-banner': { template: '<div><slot /></div>' },
                    'q-spinner-dots': true
                }
            }
        })

        await flushPromises()
        expect(electionsService.getElectionDetails).toHaveBeenCalledWith('elec-test-101')
        expect(wrapper.vm.election.title).toBe('Rinnovo Rappresentanti di Classe 2026')
        expect(wrapper.vm.lists).toHaveLength(1)
    })
})
