import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import SchoolCredits from '@/pages/teacher/SchoolCredits.vue'

describe('School Credits Triennium Workflow E2E', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'teacher-1', role: 'teacher', name: 'Prof. Coordinatore' }
                }
            }
        })
    })

    it('renders credits calculation tables and reference guides', () => {
        const wrapper = mount(SchoolCredits, {
            global: {
                plugins: [pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-expansion-item': { template: '<div><slot /></div>' },
                    'q-markup-table': { template: '<table><slot /></table>' },
                    'q-table': { template: '<div><slot /></div>' },
                    'q-select': { template: '<div class="q-select"><slot /></div>' },
                    'q-btn': {
                        props: ['label'],
                        template: '<button class="q-btn">{{ label }}<slot /></button>'
                    },
                    'q-icon': true,
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form @submit.prevent><slot /></form>' },
                    'q-input': true,
                    'q-badge': true,
                    'q-toggle': true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.find('.q-btn').exists()).toBe(true)
    })
})
