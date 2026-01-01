import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useColloquiStore } from '@/stores/colloqui'

describe('Colloqui Store', () => {
    let store

    beforeEach(() => {
        vi.useFakeTimers()
        setActivePinia(createPinia())
        store = useColloquiStore()
    })

    afterEach(() => {
        vi.useRealTimers()
    })

    it('fetches slots', async () => {
        const promise = store.fetchSlots('start', 'end')

        expect(store.loading).toBe(true)
        await vi.runAllTimersAsync()
        await promise

        expect(store.loading).toBe(false)
        expect(store.slots).toHaveLength(2)
        expect(store.bookings).toHaveLength(1)
    })

    it('creates slots', async () => {
        store.slots = []
        const newSlots = [{ date: '2025-01-21', startTime: '10:00', endTime: '10:15' }]
        const promise = store.createSlots(newSlots)

        await vi.runAllTimersAsync()
        await promise

        expect(store.slots).toHaveLength(1)
        expect(store.slots[0].date).toBe('2025-01-21')
    })

    it('deletes slot', async () => {
        store.slots = [{ id: 1 }, { id: 2 }]
        await store.deleteSlot(1)

        expect(store.slots).toHaveLength(1)
        expect(store.slots[0].id).toBe(2)
    })
})
