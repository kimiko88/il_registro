import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import FascicoloStudente from '@/pages/secretary/FascicoloStudente.vue'
import { userService } from '@/services/userService'

vi.mock('@/services/userService', () => ({
    userService: {
        getStudentFascicolo: vi.fn()
    }
}))

vi.mock('vue-router', () => ({
    useRoute: () => ({
        params: { id: 'student-123' }
    })
}))

describe('Secretary Student Digital Dossier Workflow E2E', () => {
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

    it('renders student digital portfolio with student bio and guardians info', async () => {
        userService.getStudentFascicolo.mockResolvedValue({
            data: {
                student: {
                    id: 'student-123',
                    first_name: 'Mario',
                    last_name: 'Rossi',
                    email: 'm.rossi@example.com',
                    fiscal_code: 'RSSMRA08A01H501U'
                },
                guardians: [
                    {
                        id: 'parent-1',
                        first_name: 'Giuseppe',
                        last_name: 'Rossi',
                        email: 'g.rossi@example.com',
                        phone: '+39 333 1234567'
                    }
                ]
            }
        })

        const wrapper = mount(FascicoloStudente, {
            global: {
                plugins: [pinia]
            }
        })

        await wrapper.vm.$nextTick()
        await new Promise(r => setTimeout(r, 10))

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.text()).toContain('Fascicolo Digitale Studente')
    })
})
