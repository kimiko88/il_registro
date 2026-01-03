import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useChildAttendance } from '@/composables/useChildAttendance'
import { useChildrenStore } from '@/stores/children'
import { useAttendanceStore } from '@/stores/attendance'

// Mock Stores
// Since we use createTestingPinia usually, but here composable uses stores directly.
// We can use real pinia and spy on store actions, or mock global stores.
// Let's use real pinia and mock useMyAttendance.

const { mockUseMyAttendance } = vi.hoisted(() => ({
    mockUseMyAttendance: vi.fn(() => ({
        stats: { present: 10 },
        records: [],
        loading: false
    }))
}))

vi.mock('src/composables/useMyAttendance', () => ({
    useMyAttendance: mockUseMyAttendance
}))

describe('useChildAttendance', () => {
    let childrenStore
    let attendanceStore

    beforeEach(() => {
        setActivePinia(createPinia())
        childrenStore = useChildrenStore()
        attendanceStore = useAttendanceStore()

        // Mock fetchMyAttendance
        attendanceStore.fetchMyAttendance = vi.fn()
        childrenStore.selectedChildId = null
    })

    it('returns stats and records from useMyAttendance', () => {
        const { stats, loading } = useChildAttendance()
        expect(stats.present).toBe(10)
        expect(loading).toBe(false)
    })

    it('fetches attendance when selected child changes', async () => {
        useChildAttendance()

        // Trigger watch
        childrenStore.selectedChildId = 'child1'
        await new Promise(resolve => setTimeout(resolve, 10)) // wait for watch

        expect(attendanceStore.fetchMyAttendance).toHaveBeenCalled()
    })

    it('does not fetch if no child selected', async () => {
        useChildAttendance()
        childrenStore.selectedChildId = null
        await new Promise(resolve => setTimeout(resolve, 10))

        expect(attendanceStore.fetchMyAttendance).not.toHaveBeenCalled()
    })
})
