import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import CoordinatorDashboard from '@/components/Teacher/CoordinatorDashboard.vue'

// Mock Composable
const mockGetProblemStudents = vi.fn()
vi.mock('src/composables/useCoordination', () => ({
    useCoordination: () => ({
        getProblemStudents: mockGetProblemStudents
    })
}))

describe('Teacher/CoordinatorDashboard.vue', () => {
    it('renders correctly with problem students', () => {
        mockGetProblemStudents.mockReturnValue([
            { id: 1, name: 'Alice', issue: 'Low Attendance' }
        ])

        const wrapper = mount(CoordinatorDashboard, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-list': { template: '<ul><slot /></ul>' },
                    'q-item': { template: '<li><slot /></li>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<span><slot /></span>' },
                    'q-icon': true
                }
            },
            props: {
                classId: 'c1',
                className: 'Class 5A'
            }
        })

        expect(wrapper.text()).toMatch(/Coordinator Dashboard|Pannello Coordinatore/)
        expect(wrapper.text()).toContain('Alice')
        expect(wrapper.text()).toContain('Low Attendance')
        expect(mockGetProblemStudents).toHaveBeenCalledWith('c1')
    })
})
