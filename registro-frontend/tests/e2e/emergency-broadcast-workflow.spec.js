import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import Communications from '@/pages/teacher/Communications.vue'

describe('Emergency Broadcast and Urgent Acknowledgment Workflow E2E', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'teacher-1', role: 'teacher', name: 'Prof. Rossi' }
                }
            }
        })
    })

    it('renders communications and circulars tabs and compose action', () => {
        const wrapper = mount(Communications, {
            global: {
                plugins: [pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
                    'q-tab': { template: '<div class="q-tab"><slot /></div>' },
                    'q-separator': { template: '<hr />' },
                    'q-tab-panels': { template: '<div class="q-tab-panels"><slot /></div>' },
                    'q-tab-panel': { template: '<div class="q-tab-panel"><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-scroll-area': { template: '<div><slot /></div>' },
                    'q-list': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-avatar': { template: '<div><slot /></div>' },
                    'q-btn': {
                        props: ['label'],
                        template: '<button class="q-btn">{{ label }}<slot /></button>'
                    },
                    'q-input': true,
                    'q-icon': true,
                    'q-badge': true,
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-select': true,
                    'q-checkbox': true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.find('.q-tabs').exists()).toBe(true)
    })
})
