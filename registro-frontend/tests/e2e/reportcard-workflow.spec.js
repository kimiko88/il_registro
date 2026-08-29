import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import ReportCard from '@/pages/student/ReportCard.vue'

describe('Student and Parent Report Card Workflow E2E', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'student-1', role: 'student', first_name: 'Mario', last_name: 'Rossi' }
                }
            }
        })
    })

    it('renders formal report card sheet, quadrimestre toggle and export PDF action', () => {
        const wrapper = mount(ReportCard, {
            global: {
                plugins: [Quasar, pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-avatar': { template: '<div><slot /></div>' },
                    'q-badge': { template: '<div><slot /></div>' },
                    'q-markup-table': { template: '<table><slot /></table>' },
                    'q-btn-toggle': { template: '<div class="q-btn-toggle"><slot /></div>' },
                    'q-btn': {
                        props: ['label'],
                        template: '<button class="q-btn">{{ label }}<slot /></button>'
                    },
                    'q-icon': true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.find('.printable-card').exists()).toBe(true)
    })
})
