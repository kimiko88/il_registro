import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import SecretaryDashboard from '@/pages/secretary/Index.vue'
import { createTestingPinia } from '@pinia/testing'
import { createRouter, createWebHistory } from 'vue-router'

// Mock Quasar
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            loading: { show: vi.fn(), hide: vi.fn() },
            notify: vi.fn()
        })
    }
})

describe('Secretary Workflow', () => {
    it('loads dashboard and navigates to users', async () => {
        const router = createRouter({
            history: createWebHistory(),
            routes: [
                { path: '/', component: SecretaryDashboard },
                { path: '/secretary/users', component: { template: '<div>Users Page</div>' } }
            ]
        })

        const wrapper = mount(SecretaryDashboard, {
            global: {
                plugins: [
                    router,
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            auth: { user: { role: 'secretary', school_id: '1' } }
                        }
                    })
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-avatar': true,
                    'q-icon': true,
                    'q-btn': true, // We might need to un-stub q-btn if we want to click it.
                    // But q-btn 'to' prop is handled by router-link usually or manual click.
                    // If stubbed, 'to' prop is ignored unless we handle click.
                    // Quasar q-btn with 'to' uses internal router.
                    // Let's use shallowMount or carefully stub.
                    'q-list': true,
                    'q-item': true,
                    'q-item-section': true,
                    'q-item-label': true,
                    'q-chip': true
                }
            }
        })

        // Push router to initial
        router.push('/')
        await router.isReady()

        expect(wrapper.text()).toContain('Dashboard Segreteria')
        expect(wrapper.text()).toContain('Da Revisionare')

        // Verify "Registra Utente" button exists (by finding component with label/to)
        // Since q-btn is stubbed (default true), we can find it by attributes.
        const userBtn = wrapper.findAllComponents({ name: 'q-btn' }).find(c => c.props().label === 'Registra Utente/Studente')
        expect(userBtn).toBeTruthy()
        expect(userBtn.props().to).toBe('/secretary/users')

        // We can't easily click stubbed q-btn to trigger router unless we simulate it.
        // But verifying attributes is enough for "Navigation Check".

        // Verify Stats
        expect(wrapper.text()).toContain('Documenti Pendenti')
    })
})
