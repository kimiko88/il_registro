import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import TeacherDashboard from '@/pages/teacher/Index.vue' // Adjust path if needed
import { createTestingPinia } from '@pinia/testing'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
    history: createWebHistory(),
    routes: [{ path: '/', component: TeacherDashboard }]
})

describe('Teacher Workflow', () => {
    it('loads dashboard and navigates', async () => {
        const wrapper = mount(TeacherDashboard, {
            global: {
                plugins: [
                    router,
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            teacher: {
                                profile: {
                                    firstName: 'Prof.',
                                    lastName: 'Rossi'
                                }
                            }
                        }
                    })
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-avatar': true,
                    'q-list': { template: '<ul><slot /></ul>' },
                    'q-item': { template: '<li><slot /></li>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<span><slot /></span>' },
                    'q-icon': true,
                    'q-btn': true
                }
            }
        })

        expect(wrapper.text()).toContain('Prof. Rossi')
        // Simulate navigation or interaction
    })
})
