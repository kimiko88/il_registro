import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import CoordinatorView from '@/pages/teacher/CoordinatorView.vue'

describe('Teacher Class Coordinator Dashboard Workflow E2E', () => {
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

    it('renders coordinator dashboard and empty or populated class view', () => {
        const wrapper = mount(CoordinatorView, {
            global: {
                plugins: [Quasar, pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-space': true,
                    'q-select': true,
                    'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
                    'q-tab': { template: '<div class="q-tab"><slot /></div>' },
                    'q-separator': { template: '<hr />' },
                    'q-tab-panels': { template: '<div class="q-tab-panels"><slot /></div>' },
                    'q-tab-panel': { template: '<div class="q-tab-panel"><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-spinner': true,
                    'q-markup-table': { template: '<table><slot /></table>' },
                    'q-icon': true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.text()).toContain('coordinatore')
    })
})
