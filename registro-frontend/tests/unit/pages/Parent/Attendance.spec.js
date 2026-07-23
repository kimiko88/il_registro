import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import ParentAttendance from '@/pages/parent/Attendance.vue'
import api from 'src/services/api'

vi.mock('src/services/api', () => ({
    default: {
        get: vi.fn(),
        post: vi.fn()
    }
}))

describe('Parent/Attendance.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        api.get.mockImplementation((url) => {
            if (url.includes('/users/me/children')) {
                return Promise.resolve({ data: [{ id: 'child1', first_name: 'Mario' }] })
            }
            if (url.includes('/attendance/child-attendance/child1')) {
                return Promise.resolve({
                    data: [
                        { id: 1, date: '2023-01-01', status: 'Absent', justified: false },
                        { id: 2, date: '2023-01-02', status: 'Absent', justified: true },
                        { id: 3, date: '2023-01-03', status: 'Present' }
                    ]
                })
            }
            return Promise.resolve({ data: [] })
        })

        wrapper = mount(ParentAttendance, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        stubActions: false
                    })
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-list': { template: '<ul><slot /></ul>' },
                    'q-item': { template: '<li><slot /></li>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-icon': true,
                    'q-btn': true,
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-select': true,
                    'q-input': true,
                    'q-card-actions': { template: '<div><slot /></div>' }
                }
            }
        })
    })

    it('fetches attendance on mount', async () => {
        await wrapper.vm.$nextTick()
        expect(api.get).toHaveBeenCalledWith('/users/me/children')
    })

    it('filters and displays absences', async () => {
        await wrapper.vm.$nextTick()
        await new Promise(resolve => setTimeout(resolve, 10))
        await wrapper.vm.$nextTick()

        expect(wrapper.vm.events.length).toBeGreaterThan(0)
    })

    it('opens justify dialog', async () => {
        await wrapper.vm.$nextTick()
        await new Promise(resolve => setTimeout(resolve, 10))
        await wrapper.vm.$nextTick()

        wrapper.vm.openJustifyModal({ id: 1, date: '2023-01-01' })
        expect(wrapper.vm.justifyModal).toBe(true)
        expect(wrapper.vm.selectedAttendance.id).toBe(1)
    })
})
