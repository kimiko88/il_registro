import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import Dashboard from '@/pages/Dashboard.vue'
import { useAuthStore } from '@/stores/auth'

describe('Dashboard.vue', () => {
    let wrapper
    let store

    beforeEach(() => {
        wrapper = mount(Dashboard, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            auth: {
                                user: { first_name: 'TestUser' },
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
        store = useAuthStore()
    })

    it('renders greeting with user name', () => {
        expect(wrapper.text()).toContain('TestUser')
    })

    it('renders student stats by default', () => {
        expect(wrapper.text()).toContain('Media Voti')
        expect(wrapper.text()).toContain('Presenze')
    })

    it('renders teacher stats when role is teacher', async () => {
        // We use createTestingPinia which mocks the store.
        // userRole is a computed property (getter). getters are writable in mocked stores.

        wrapper = mount(Dashboard, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            auth: {
                                user: { first_name: 'Teacher' }
                            }
                        },
                        stubActions: false
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

        const store = useAuthStore()
        store.userRole = 'teacher' // Override the getter value
        await wrapper.vm.$nextTick()

        expect(wrapper.text()).toContain('Le Mie Classi')
        expect(wrapper.text()).toContain('Lezioni Oggi')
    })
})
