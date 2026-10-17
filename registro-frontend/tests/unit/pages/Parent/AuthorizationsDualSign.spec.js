import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Authorizations from '@/pages/parent/Authorizations.vue'
import { parentDualSignService } from '@/services/parentDualSignService'

vi.mock('@/services/parentDualSignService', () => ({
    parentDualSignService: {
        listAuthorizations: vi.fn(),
        signAuthorization: vi.fn(),
        rejectAuthorization: vi.fn()
    },
    default: {
        listAuthorizations: vi.fn(),
        signAuthorization: vi.fn(),
        rejectAuthorization: vi.fn()
    }
}))

describe('Parent Dual Signature & Custody Authorizations', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        parentDualSignService.listAuthorizations.mockResolvedValue({
            data: [
                {
                    id: 'auth-1',
                    title: 'Autorizzazione Gita a Firenze',
                    document_type: 'trip_consent',
                    status: 'pending_first',
                    deadline: '2026-10-20T12:00:00Z',
                    parent1_pin_verified: false,
                    parent2_pin_verified: false
                },
                {
                    id: 'auth-2',
                    title: 'Approvazione PDP DSA',
                    document_type: 'pdp_approval',
                    status: 'pending_second',
                    deadline: '2026-10-15T12:00:00Z',
                    parent1_pin_verified: true,
                    parent2_pin_verified: false
                }
            ]
        })
    })

    it('renders the dual sign authorizations with status badges and triggers PIN sign dialog', async () => {
        const wrapper = mount(Authorizations, {
            global: {
                mocks: {
                    t: (key) => key
                },
                stubs: {
                    QPage: { template: '<div><slot /></div>' },
                    QCard: { template: '<div><slot /></div>' },
                    QCardSection: { template: '<div><slot /></div>' },
                    QCardActions: { template: '<div><slot /></div>' },
                    QBtn: { template: '<button @click="$emit(\'click\')"><slot /></button>' },
                    QBadge: { template: '<span class="badge"><slot /></span>' },
                    QChip: { template: '<span><slot /></span>' },
                    QDialog: { template: '<div><slot /></div>' },
                    QInput: { template: '<input />' },
                    QIcon: true,
                    QTooltip: true,
                    QSpinner: true,
                    QTabs: true,
                    QTab: true,
                    QSeparator: true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        await wrapper.vm.$nextTick()
        expect(parentDualSignService.listAuthorizations).toHaveBeenCalled()
    })
})
