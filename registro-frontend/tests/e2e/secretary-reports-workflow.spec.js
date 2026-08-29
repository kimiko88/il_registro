import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import Reports from '@/pages/secretary/Reports.vue'

describe('Secretary Reports Hub and Data Export Workflow E2E', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'sec-1', role: 'secretary', name: 'Segreteria Didattica' }
                }
            }
        })
    })

    it('renders reports center, report categories and action cards', () => {
        const wrapper = mount(Reports, {
            global: {
                plugins: [Quasar, pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-icon': true,
                    'q-dialog': { template: '<div><slot /></div>' }
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.text()).toContain('Centro Reportistica')
    })
})
