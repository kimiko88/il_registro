import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSchedulingStore } from './scheduling'

describe('Scheduling Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('generates schedule successfully', async () => {
        const store = useSchedulingStore()
        const constraints = {
            school_id: 'S1',
            days: [1], // Monday
            teachers: ['T1'],
            classes: ['1A']
        }

        const result = await store.generateSchedule(constraints)

        expect(result).toHaveLength(2)
        expect(store.generatedSchedule).toHaveLength(2)
        expect(store.isLoading).toBe(false)
    })

    it('validates schedule with conflicts', async () => {
        const store = useSchedulingStore()
        const badSlots = [
            { day: 1, hour: 8, teacher_id: 'T1' },
            { day: 1, hour: 8, teacher_id: 'T1' } // Conflict
        ]

        const conflicts = await store.validateSchedule(badSlots)

        expect(conflicts).toHaveLength(1)
        expect(store.conflicts).toHaveLength(1)
        expect(conflicts[0].type).toBe('TEACHER_OVERLAP')
    })
})
