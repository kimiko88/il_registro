import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import StudentIndex from '@/pages/student/Index.vue'

// Mock services
const { mockGetMyGrades, mockGetMyAttendance, mockGetMyProjects, mockGetSubjects } = vi.hoisted(() => ({
    mockGetMyGrades: vi.fn(),
    mockGetMyAttendance: vi.fn(),
    mockGetMyProjects: vi.fn(),
    mockGetSubjects: vi.fn()
}))

vi.mock('@/services/gradeService', () => ({
    gradeService: {
        getMyGrades: mockGetMyGrades
    }
}))
vi.mock('@/services/attendanceService', () => ({
    attendanceService: {
        getMyAttendance: mockGetMyAttendance
    }
}))
vi.mock('@/services/pctoService', () => ({
    pctoService: {
        getMyProjects: mockGetMyProjects
    }
}))
vi.mock('@/services/adminService', () => ({
    default: {
        getSubjects: mockGetSubjects
    }
}))
// Communications might be used too, assuming mockGetActiveNotifications if needed

describe('Student/Index.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        mockGetMyGrades.mockResolvedValue({
            data: {
                semesters: [
                    {
                        grades: [
                            { id: 1, subject_id: 'Math', grade_value: 8, date: '2023-01-01', grade_type: 'oral', description: 'Good' }
                        ]
                    }
                ]
            }
        })
        mockGetMyAttendance.mockResolvedValue({
            data: [
                { id: 1, date: '2023-01-01', status: 'present' }
            ]
        })
        mockGetMyProjects.mockResolvedValue({
            data: [] // Empty list for now
        })
        mockGetSubjects.mockResolvedValue({
            data: [
                { id: 'Math', name: 'Matematica' }
            ]
        })

        wrapper = mount(StudentIndex, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            student: {
                                profile: { first_name: 'Mario', last_name: 'Rossi', class_id: '5A', school_id: '1' }
                            }
                        },
                        stubActions: true
                    })
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-btn': true,
                    'q-knob': true,
                    'q-linear-progress': true,
                    'q-list': true,
                    'q-item': true,
                    'q-item-section': true,
                    'q-item-label': true,
                    'q-badge': true,
                    'q-icon': true,
                    'q-avatar': true,
                    'router-link': true
                }
            }
        })
    })

    it('renders student name', () => {
        expect(wrapper.text()).toContain('Mario')
    })

    it('fetches and displays dashboard data', async () => {
        // Wait for all async calls in onMounted
        await wrapper.vm.$nextTick()
        await new Promise(resolve => setTimeout(resolve, 100))
        await wrapper.vm.$nextTick()
        await new Promise(resolve => setTimeout(resolve, 100))
        await wrapper.vm.$nextTick()

        // Verify service calls
        expect(mockGetMyGrades).toHaveBeenCalled()
        expect(mockGetMyAttendance).toHaveBeenCalled()

        // Check if name is still there
        expect(wrapper.text()).toContain('Mario')
    })
})
