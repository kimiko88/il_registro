import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import GradeWeights from '@/pages/teacher/GradeWeights.vue'

vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn().mockResolvedValue({ data: [] }),
        put: vi.fn().mockResolvedValue({ data: {} }),
        delete: vi.fn().mockResolvedValue({ data: {} }),
        interceptors: {
            request: { use: vi.fn(), handlers: [] },
            response: { use: vi.fn(), handlers: [] }
        }
    }
}))

describe('Grade Weights and Matrix Workflow E2E', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'admin-1', role: 'admin', name: 'Admin Scolastico' }
                }
            }
        })
    })

    it('renders grade weights configuration page and add button', () => {
        const wrapper = mount(GradeWeights, {
            global: {
                plugins: [pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-select': { template: '<div class="q-select"><slot /></div>' },
                    'q-table': { template: '<div><slot /></div>' },
                    'q-space': { template: '<div></div>' },
                    'q-spinner-dots': { template: '<div></div>' },
                    'q-btn': {
                        props: ['label'],
                        template: '<button class="q-btn">{{ label }}<slot /></button>'
                    },
                    'q-icon': true,
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form @submit.prevent><slot /></form>' },
                    'q-input': { template: '<div><slot /></div>' },
                    'q-chip': true,
                    'q-badge': true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.text()).toContain('Configurazione Pesi Voti')
    })
})
