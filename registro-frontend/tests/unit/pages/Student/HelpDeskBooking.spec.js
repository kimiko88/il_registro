import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import HelpDeskBooking from '@/pages/student/HelpDeskBooking.vue'
import { helpDeskService } from '@/services/helpDeskService'

vi.mock('@/services/helpDeskService', () => ({
    helpDeskService: {
        getSlots: vi.fn(),
        bookSlot: vi.fn()
    }
}))

describe('HelpDeskBooking.vue', () => {
    it('renders available help desk slots with teacher and subject info', async () => {
        helpDeskService.getSlots.mockResolvedValueOnce({
            data: [
                {
                    id: 'slot-1',
                    teacher_name: 'Prof. Carlo Rossi',
                    subject_name: 'Fisica',
                    slot_date: '2026-10-10',
                    start_time: '15:00',
                    end_time: '16:30',
                    max_capacity: 4,
                    bookings_count: 2,
                    status: 'open'
                }
            ]
        })

        const wrapper = mount(HelpDeskBooking, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-separator': true,
                    'q-icon': true,
                    'q-btn': true,
                    'q-badge': true,
                    'q-spinner-dots': true,
                    'q-dialog': true
                }
            }
        })

        await flushPromises()
        expect(helpDeskService.getSlots).toHaveBeenCalled()
        expect(wrapper.vm.slots).toHaveLength(1)
        expect(wrapper.vm.slots[0].teacher_name).toBe('Prof. Carlo Rossi')
        expect(wrapper.vm.slots[0].subject_name).toBe('Fisica')
    })
})
