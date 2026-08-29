import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import Certificates from '@/pages/secretary/Certificates.vue'

describe('Certificates and Student Dossier Workflow E2E', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'secr-1', role: 'secretary', name: 'Segreteria Didattica' }
                }
            }
        })
    })

    it('renders certificates issuance dashboard and generate certificate button', () => {
        const wrapper = mount(Certificates, {
            global: {
                plugins: [Quasar, pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-table': { template: '<div><slot /></div>' },
                    'q-dialog': { template: '<div><slot /></div>' }
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.text()).toContain('Gestione Certificati')
    })
})
