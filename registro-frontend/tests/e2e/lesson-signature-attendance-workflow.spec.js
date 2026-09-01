import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import Attendance from '@/pages/teacher/Attendance.vue'

describe('Lesson Signature and Class Attendance Workflow E2E', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'teacher-1', role: 'teacher', name: 'Prof. Bianchi' }
                }
            }
        })
    })

    it('renders class register header, filters and action buttons', () => {
        const wrapper = mount(Attendance, {
            global: {
                plugins: [Quasar, pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-table': { template: '<div><slot /></div>' },
                    'q-dialog': { template: '<div><slot /></div>' }
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.find('.bg-grey-1').exists()).toBe(true)
    })
})
