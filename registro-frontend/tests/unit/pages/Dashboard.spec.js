import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { setActivePinia } from 'pinia'
import { Quasar } from 'quasar'
import Dashboard from '@/pages/Dashboard.vue'
import { useAuthStore } from '@/stores/auth'

vi.mock('vue-router', () => ({
    useRouter: () => ({
        push: vi.fn(),
        replace: vi.fn()
    })
}))

vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn().mockResolvedValue({ data: [] }),
        post: vi.fn().mockResolvedValue({ data: {} })
    }
}))

vi.mock('@/services/dashboardService', () => ({
    default: {
        getDashboardStats: vi.fn().mockImplementation((role) => {
            if (role === 'teacher') {
                return Promise.resolve({
                    classes_count: 3,
                    students_count: 75,
                    lessons_today_count: 2,
                    grades_pending_count: 5
                })
            }
            return Promise.resolve({
                average_grade: 8.2,
                attendance_rate: 96,
                homework_count: 3,
                documents_count: 1
            })
        })
    }
}))

describe('Dashboard.vue', () => {
    let wrapper

    beforeEach(() => {
        const pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { first_name: 'TestUser', role: 'student' },
                    userRole: 'student'
                }
            },
            stubActions: false
        })
        setActivePinia(pinia)

        wrapper = mount(Dashboard, {
            global: {
                plugins: [
                    [Quasar, {}],
                    pinia
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-icon': true,
                    'q-btn': true,
                    'q-list': { template: '<ul><slot /></ul>' },
                    'q-item': { template: '<li><slot /></li>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<span><slot /></span>' },
                    'q-chip': true
                }
            }
        })
    })

    it('renders greeting with user name', () => {
        expect(wrapper.text()).toContain('TestUser')
    })

    it('renders student stats by default', async () => {
        await flushPromises()
        expect(wrapper.text()).toContain('Media Voti')
        expect(wrapper.text()).toContain('Presenze')
    })


    it('renders teacher stats when role is teacher', async () => {
        const pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { first_name: 'Teacher', role: 'teacher' }
                }
            },
            stubActions: false
        })
        setActivePinia(pinia)
        const store = useAuthStore(pinia)
        store.userRole = 'teacher'

        const teacherWrapper = mount(Dashboard, {
            global: {
                plugins: [
                    [Quasar, {}],
                    pinia
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-icon': true,
                    'q-btn': true,
                    'q-list': { template: '<ul><slot /></ul>' },
                    'q-item': { template: '<li><slot /></li>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<span><slot /></span>' },
                    'q-chip': true
                }
            }
        })

        await flushPromises()

        expect(teacherWrapper.text()).toContain('Le Mie Classi')
        expect(teacherWrapper.text()).toContain('Lezioni Oggi')
    })
})
