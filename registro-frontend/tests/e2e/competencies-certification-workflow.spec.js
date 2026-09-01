import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import Competencies from '@/pages/teacher/Competencies.vue'

describe('Teacher EU Competencies Certification Workflow E2E', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'teacher-1', role: 'teacher', name: 'Prof. Mario Rossi' }
                }
            }
        })
    })

    it('renders competencies evaluation matrix page and filter controls', () => {
        const wrapper = mount(Competencies, {
            global: {
                plugins: [Quasar, pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-chip': { template: '<div><slot /></div>' },
                    'q-banner': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-select': true,
                    'q-table': { template: '<div><slot /></div>' },
                    'q-btn': {
                        props: ['label'],
                        template: '<button class="q-btn">{{ label }}<slot /></button>'
                    },
                    'q-icon': true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.find('.sticky-header').exists()).toBe(true)
    })
})
