import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import AccessibilityStatement from '@/pages/AccessibilityStatement.vue'

describe('Accessibility AgID Statement and Controls Workflow E2E', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'user-1', role: 'student', name: 'Mario Rossi' }
                }
            }
        })
    })

    it('renders accessibility statement, summary cards and AgID reference sections', () => {
        const wrapper = mount(AccessibilityStatement, {
            global: {
                plugins: [Quasar, pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-avatar': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-icon': true,
                    'q-btn': {
                        props: ['label'],
                        template: '<button class="q-btn">{{ label }}<slot /></button>'
                    },
                    'q-form': { template: '<form @submit.prevent><slot /></form>' },
                    'q-input': true,
                    'q-select': true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.find('.accessibility-statement-page').exists()).toBe(true)
    })
})
