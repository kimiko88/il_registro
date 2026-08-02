import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSchedulingStore } from '@/stores/scheduling'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
    default: {
        post: vi.fn()
    }
}))

describe('Scheduling Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
    })

    it('generates schedule successfully', async () => {
        const store = useSchedulingStore()
        const constraints = {
            school_id: 'S1',
            days: [1], // Monday
            teachers: ['T1'],
            classes: ['1A']
        }
        api.post.mockResolvedValue({ data: [{ day: 1, hour: 8 }, { day: 1, hour: 9 }] })

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
        api.post.mockResolvedValue({ data: { conflicts: [{ type: 'TEACHER_OVERLAP' }] } })

        const conflicts = await store.validateSchedule(badSlots)

        expect(conflicts).toHaveLength(1)
        expect(store.conflicts).toHaveLength(1)
        expect(conflicts[0].type).toBe('TEACHER_OVERLAP')
    })
})
