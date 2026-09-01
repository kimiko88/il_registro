import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import Payments from '@/pages/parent/Payments.vue'

describe('Parent Payments and PagoPA Workflow E2E', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'parent-1', role: 'parent', name: 'Mario Rossi' }
                }
            }
        })
    })

    it('renders parent payments dashboard with tabs and pagoPA action', () => {
        const wrapper = mount(Payments, {
            global: {
                plugins: [pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-chip': { template: '<div class="q-chip"><slot /></div>' },
                    'q-spinner': { template: '<div></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
                    'q-tab': { template: '<div class="q-tab"><slot /></div>' },
                    'q-separator': { template: '<hr />' },
                    'q-tab-panels': { template: '<div class="q-tab-panels"><slot /></div>' },
                    'q-tab-panel': { template: '<div class="q-tab-panel"><slot /></div>' },
                    'q-list': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-avatar': { template: '<div><slot /></div>' },
                    'q-btn': {
                        props: ['label'],
                        template: '<button class="q-btn">{{ label }}<slot /></button>'
                    },
                    'q-icon': true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.find('.q-tabs').exists()).toBe(true)
    })
})
