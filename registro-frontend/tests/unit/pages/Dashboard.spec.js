import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { setActivePinia } from 'pinia'
import { Quasar } from 'quasar'
import Dashboard from '@/pages/Dashboard.vue'
import { useAuthStore } from '@/stores/auth'

vi.mock('@/services/dashboardService', () => ({
    default: {
        getDashboardStats: vi.fn().mockResolvedValue(null)
    }
}))

describe('Dashboard.vue', () => {
    let wrapper

    beforeEach(() => {
        wrapper = mount(Dashboard, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            auth: {
                                user: { first_name: 'TestUser', role: 'student' },
                                userRole: 'student'
                            }
                        }
                    })
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

    it('renders student stats by default', () => {
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
