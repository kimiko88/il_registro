import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import ParentAttendance from '@/pages/parent/Attendance.vue'

// Mock services
const { mockGetChildAttendance, mockJustify } = vi.hoisted(() => ({
    mockGetChildAttendance: vi.fn(),
    mockJustify: vi.fn()
}))

vi.mock('src/services/attendanceService', () => ({
    attendanceService: {
        getChildAttendance: mockGetChildAttendance,
        justify: mockJustify
    }
}))

describe('Parent/Attendance.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        mockGetChildAttendance.mockResolvedValue({
            data: [
                { id: 1, date: '2023-01-01', status: 'absent', is_justified: false },
                { id: 2, date: '2023-01-02', status: 'absent', is_justified: true },
                { id: 3, date: '2023-01-03', status: 'present' }
            ]
        })

        wrapper = mount(ParentAttendance, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            parent: {
                                children: [{ id: 'child1', firstName: 'Mario' }],
                                selectedChildId: 'child1'
                            }
                        },
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

    it('fetches attendance on mount', () => {
        expect(mockGetChildAttendance).toHaveBeenCalledWith('child1')
    })

    it('filters and displays absences', async () => {
        await wrapper.vm.$nextTick()
        await new Promise(resolve => setTimeout(resolve, 10))
        await wrapper.vm.$nextTick()

        // Original data has 3 records: 2 absent, 1 present.
        // Logic filters out present. So length 2.
        // events.value
        expect(wrapper.vm.events).toHaveLength(2)
    })

    it('opens justify dialog', async () => {
        await wrapper.vm.$nextTick()
        await new Promise(resolve => setTimeout(resolve, 10))
        await wrapper.vm.$nextTick()

        // Trigger openJustifyDialog directly or find button
        // The list items are stubbed.
        // Let's call the method on VM
        wrapper.vm.openJustifyDialog({ id: 1, date: '2023-01-01' })
        expect(wrapper.vm.justifyDialog).toBe(true)
        expect(wrapper.vm.selectedEvent.id).toBe(1)
    })
})
