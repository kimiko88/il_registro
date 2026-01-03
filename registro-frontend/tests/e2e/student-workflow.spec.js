import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import StudentDashboard from '@/pages/student/Index.vue'
import { createTestingPinia } from '@pinia/testing'
import { createRouter, createWebHistory } from 'vue-router'

// Mock Quasar
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            loading: { show: vi.fn(), hide: vi.fn() },
            notify: vi.fn(),
            lang: { getLocale: () => 'it', isoName: 'it' }
        })
    }
})

// Mock Services
const { mockGradeService, mockAttendanceService, mockPctoService } = vi.hoisted(() => {
    return {
        mockGradeService: { getMyGrades: vi.fn().mockResolvedValue({ data: { semesters: [] } }) },
        mockAttendanceService: { getMyAttendance: vi.fn().mockResolvedValue({ data: [] }) },
        mockPctoService: { getMyProjects: vi.fn().mockResolvedValue({ data: [] }) }
    }
})

vi.mock('@/services/gradeService', () => ({ gradeService: mockGradeService }))
vi.mock('@/services/attendanceService', () => ({ attendanceService: mockAttendanceService }))
vi.mock('@/services/pctoService', () => ({ pctoService: mockPctoService }))

describe('Student Workflow', () => {
    it('loads student dashboard and navigation', async () => {
        const router = createRouter({
            history: createWebHistory(),
            routes: [{ path: '/', component: StudentDashboard }]
        })

        const wrapper = mount(StudentDashboard, {
            global: {
                plugins: [
                    router,
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            auth: { user: { role: 'student', first_name: 'Student', id: 1 } },
                            student: {
                                // Fix: Provide profile data used in template
                                profile: { firstName: 'Student', lastName: 'User' },
                                notifications: [],
                                stats: { average: 7.5, absences: 2 },
                                recentGrades: [],
                                upcomingEvents: []
                            }
                        }
                    })
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-btn': true,
                    'q-icon': true,
                    'q-list': true,
                    'q-item': true,
                    'q-item-section': true,
                    'q-item-label': true,
                    'q-circular-progress': true,
                    'q-linear-progress': true,
                    'q-avatar': true,
                    'q-separator': true,
                    'q-menu': true,
                    'q-badge': true,
                    'q-chip': true
                },
                mocks: {
                    $q: {
                        loading: { show: vi.fn(), hide: vi.fn() },
                        notify: vi.fn(),
                        lang: { getLocale: () => 'it', isoName: 'it', rtl: false },
                        screen: { lt: { md: false } } // sometimes used in responsive layouts
                    }
                }
            }
        })

        // Wait for async mounts if any (flushPromises? mostly safe with mockResolvedValue)
        await new Promise(resolve => setTimeout(resolve, 0))

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.text()).toContain('Student') // Should now match "Bentornato, Student"
        expect(wrapper.text()).toContain('Media Voti')
    })
})
