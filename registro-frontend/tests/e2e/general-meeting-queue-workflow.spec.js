import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import GeneralMeetingLiveQueue from '@/pages/teacher/GeneralMeetingLiveQueue.vue'

describe('General Meeting Live Queue Workflow E2E', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'teacher-1', role: 'teacher', name: 'Prof. Bianchi' }
                }
            }
        })
    })

    it('renders live queue spotlight and call next parent controls', () => {
        const wrapper = mount(GeneralMeetingLiveQueue, {
            global: {
                plugins: [pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-select': { template: '<div class="q-select"><slot /></div>' },
                    'q-tooltip': { template: '<div></div>' },
                    'q-table': { template: '<div><slot /></div>' },
                    'q-btn': {
                        props: ['label'],
                        template: '<button class="q-btn">{{ label }}<slot /></button>'
                    },
                    'q-icon': true,
                    'q-badge': true,
                    'q-chip': true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.find('.q-btn').exists()).toBe(true)
    })
})
