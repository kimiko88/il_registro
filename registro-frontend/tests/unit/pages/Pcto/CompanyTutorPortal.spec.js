import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import CompanyTutorPortal from '@/pages/pcto/CompanyTutorPortal.vue'
import { pctoCompanyTutorService } from '@/services/pctoCompanyTutorService'

vi.mock('vue-router', () => ({
    useRoute: () => ({
        query: { token: 'valid-test-token-12345' }
    })
}))

vi.mock('@/services/pctoCompanyTutorService', () => ({
    pctoCompanyTutorService: {
        getSession: vi.fn(),
        getAssignedStudents: vi.fn(),
        verifyTimesheet: vi.fn(),
        submitEvaluation: vi.fn()
    }
}))

describe('CompanyTutorPortal.vue', () => {
    it('authenticates with token and loads assigned students', async () => {
        pctoCompanyTutorService.getSession.mockResolvedValueOnce({
            tutor: {
                id: 'tut-1',
                company_name: 'Acme Robotics Srl',
                tutor_first_name: 'Laura',
                tutor_last_name: 'Neri'
            }
        })
        pctoCompanyTutorService.getAssignedStudents.mockResolvedValueOnce({
            data: [
                {
                    student_id: 'stud-1',
                    student_name: 'Andrea Costa',
                    project_title: 'Automazione Industriale',
                    total_hours: 60,
                    completed_hours: 32,
                    is_evaluated: false
                }
            ]
        })

        const wrapper = mount(CompanyTutorPortal, {
            global: {
                stubs: {
                    'q-layout': { template: '<div><slot /></div>' },
                    'q-header': { template: '<div><slot /></div>' },
                    'q-toolbar': { template: '<div><slot /></div>' },
                    'q-toolbar-title': { template: '<div><slot /></div>' },
                    'q-page-container': { template: '<div><slot /></div>' },
                    'q-page': { template: '<div><slot /></div>' },
                    'q-icon': true,
                    'q-btn': true,
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-banner': { template: '<div><slot /></div>' },
                    'q-spinner-dots': true,
                    'q-dialog': true
                }
            }
        })

        await flushPromises()
        expect(pctoCompanyTutorService.getSession).toHaveBeenCalledWith('valid-test-token-12345')
        expect(pctoCompanyTutorService.getAssignedStudents).toHaveBeenCalledWith('valid-test-token-12345')
        expect(wrapper.vm.tutor.company_name).toBe('Acme Robotics Srl')
        expect(wrapper.vm.students).toHaveLength(1)
        expect(wrapper.vm.students[0].student_name).toBe('Andrea Costa')
    })
})
