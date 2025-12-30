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
                            teacher: { name: 'Prof. Rossi' }
                        }
                    })
                ]
            }
        })

        expect(wrapper.text()).toContain('Prof. Rossi')
        // Simulate navigation or interaction
    })
})
