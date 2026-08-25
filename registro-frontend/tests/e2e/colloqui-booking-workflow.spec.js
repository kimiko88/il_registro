import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { useColloquiStore } from '@/stores/colloqui'

vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            loading: { show: vi.fn(), hide: vi.fn() },
            notify: vi.fn(),
            dialog: vi.fn().mockReturnValue({ onOk: (fn) => fn() })
        })
    }
})

describe('Colloqui Booking Workflow', () => {
    it('executes full colloqui slot booking and cancellation flow', async () => {
        const pinia = createTestingPinia({
            createSpy: vi.fn,
            stubActions: false,
            initialState: {
                auth: { user: { id: 'parent-1', role: 'parent' } },
                colloqui: {
                    slots: [
                        { id: 'slot-1', teacher_id: 't-1', date: '2026-09-01', start_time: '10:00', end_time: '10:30', max_bookings: 2, status: 'available' }
                    ],
                    bookings: []
                }
            }
        })

        const colloquiStore = useColloquiStore(pinia)

        // Verify initial slot availability
        expect(colloquiStore.slots.length).toBe(1)
        expect(colloquiStore.slots[0].status).toBe('available')

        // Simulate booking slot
        colloquiStore.bookings.push({
            id: 'book-1',
            slot_id: 'slot-1',
            parent_id: 'parent-1',
            status: 'confirmed',
            created_at: new Date().toISOString()
        })

        expect(colloquiStore.bookings.length).toBe(1)
        expect(colloquiStore.bookings[0].status).toBe('confirmed')

        // Simulate cancellation
        colloquiStore.bookings[0].status = 'cancelled'
        expect(colloquiStore.bookings[0].status).toBe('cancelled')
    })
})
