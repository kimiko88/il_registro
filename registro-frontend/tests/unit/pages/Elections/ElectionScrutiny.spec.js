import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import ElectionScrutiny from '@/pages/admin/ElectionScrutiny.vue'
import { electionsService } from '@/services/electionsService'

vi.mock('vue-router', () => ({
    useRoute: () => ({
        params: { id: 'elec-test-101' },
        query: {}
    })
}))

vi.mock('@/services/electionsService', () => ({
    electionsService: {
        getScrutiny: vi.fn(),
        getElections: vi.fn(),
        closeElection: vi.fn()
    }
}))

describe('ElectionScrutiny.vue', () => {
    it('renders d\'Hondt list seats allocation and proclaimed candidates', async () => {
        electionsService.getScrutiny.mockResolvedValueOnce({
            scrutiny: {
                election_id: 'elec-test-101',
                total_voters: 45,
                total_votes_cast: 45,
                blank_votes: 2,
                lists_results: [
                    { list_id: 'l-1', list_number: 1, motto: 'Lista Insieme', total_votes: 30, seats_won: 3 },
                    { list_id: 'l-2', list_number: 2, motto: 'Lista Futuro', total_votes: 13, seats_won: 1 }
                ],
                elected_members: [
                    { candidate_id: 'cand-1', first_name: 'Sara', last_name: 'Romano', list_motto: 'Lista Insieme', votes: 18 }
                ]
            }
        })

        const wrapper = mount(ElectionScrutiny, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-separator': true,
                    'q-icon': true,
                    'q-btn': true,
                    'q-badge': true,
                    'q-table': true,
                    'q-list': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-avatar': true,
                    'q-chip': true
                }
            }
        })

        await flushPromises()
        expect(electionsService.getScrutiny).toHaveBeenCalledWith('elec-test-101', 4)
        expect(wrapper.vm.scrutiny.total_voters).toBe(45)
        expect(wrapper.vm.scrutiny.lists_results).toHaveLength(2)
        expect(wrapper.vm.scrutiny.lists_results[0].seats_won).toBe(3)
        expect(wrapper.vm.scrutiny.elected_members[0].first_name).toBe('Sara')
    })
})
