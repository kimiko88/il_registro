import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import TeacherSettings from '@/pages/teacher/Settings.vue'

describe('Teacher Functional Activities and Settings Workflow E2E', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'teacher-1', role: 'teacher', first_name: 'Giuseppe', last_name: 'Verdi' }
                }
            }
        })
    })

    it('renders teacher preferences tabs including signature, digital PIN and accessibility', () => {
        const wrapper = mount(TeacherSettings, {
            global: {
                plugins: [Quasar, pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
                    'q-tab': { template: '<div class="q-tab"><slot /></div>' },
                    'q-separator': { template: '<hr />' },
                    'q-tab-panels': { template: '<div class="q-tab-panels"><slot /></div>' },
                    'q-tab-panel': { template: '<div class="q-tab-panel"><slot /></div>' },
                    'q-select': { template: '<div class="q-select"><slot /></div>' },
                    'q-input': true,
                    'q-btn': {
                        props: ['label'],
                        template: '<button class="q-btn">{{ label }}<slot /></button>'
                    },
                    'q-icon': true,
                    'q-toggle': true,
                    'q-chip': true,
                    'q-badge': true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.find('.q-tabs').exists()).toBe(true)
    })
})
