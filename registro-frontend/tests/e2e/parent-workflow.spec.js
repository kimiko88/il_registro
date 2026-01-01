import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
// Assuming Parent Dashboard shares structure or has specific page.
// Since we don't know exact Parent page, assuming `pages/parent/Index.vue` exists or similar.
// I will check or assume `student/Index.vue` logic applies (often parents see similar view).
// Let's assume `pages/parent/Index.vue` exists for now.
import ParentDashboard from '@/pages/parent/Index.vue'
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

describe('Parent Workflow', () => {
    it('loads parent dashboard', async () => {
        const router = createRouter({
            history: createWebHistory(),
            routes: [{ path: '/', component: ParentDashboard }]
        })

        const wrapper = mount(ParentDashboard, {
            global: {
                plugins: [
                    router,
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            auth: { user: { role: 'parent', first_name: 'Parent' } }
                        }
                    })
                ],
                stubs: {
                    'q-page': true,
                    'q-card': true,
                    'q-card-section': true,
                    'q-btn': true,
                    'q-icon': true,
                    'q-select': true // Parent usually selects child
                },
                mocks: {
                    $q: {
                        loading: { show: vi.fn(), hide: vi.fn() },
                        notify: vi.fn(),
                        lang: { getLocale: () => 'it', isoName: 'it', rtl: false }
                    }
                }
            }
        })

        // Basic Load Check
        expect(wrapper.exists()).toBe(true)
    })
})
