import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSchoolStore } from '@/stores/schools'

// Mock Service
const mockSchoolService = vi.hoisted(() => ({
    getSchools: vi.fn(),
    createSchool: vi.fn(),
    updateSchool: vi.fn(),
    deleteSchool: vi.fn()
}))

vi.mock('@/services/schoolService', () => ({
    default: mockSchoolService
}))

describe('Schools Store', () => {
    let store

    beforeEach(() => {
        setActivePinia(createPinia())
        store = useSchoolStore()
        vi.clearAllMocks()
    })

    it('fetches schools', async () => {
        mockSchoolService.getSchools.mockResolvedValue({
            data: { items: [{ id: 1, name: 'S1' }], total: 1 }
        })

        await store.fetchSchools()

        expect(store.loading).toBe(false)
        expect(store.schools).toHaveLength(1)
        expect(store.schools[0].name).toBe('S1')
        expect(store.pagination.rowsNumber).toBe(1)
    })

    it('creates school and refreshes', async () => {
        mockSchoolService.createSchool.mockResolvedValue({})
        store.fetchSchools = vi.fn() // Spy internal action

        await store.createSchool({ name: 'New' })

        expect(mockSchoolService.createSchool).toHaveBeenCalledWith({ name: 'New' })
        expect(store.fetchSchools).toHaveBeenCalled()
    })

    it('updates school and refreshes', async () => {
        mockSchoolService.updateSchool.mockResolvedValue({})
        store.fetchSchools = vi.fn()

        await store.updateSchool(1, { name: 'Up' })

        expect(mockSchoolService.updateSchool).toHaveBeenCalledWith(1, { name: 'Up' })
        expect(store.fetchSchools).toHaveBeenCalled()
    })

    it('deletes school and refreshes', async () => {
        mockSchoolService.deleteSchool.mockResolvedValue({})
        store.fetchSchools = vi.fn()

        await store.deleteSchool(1)

        expect(mockSchoolService.deleteSchool).toHaveBeenCalledWith(1)
        expect(store.fetchSchools).toHaveBeenCalled()
    })

    it('handles fetch error', async () => {
        mockSchoolService.getSchools.mockRejectedValue(new Error('Fail'))
        await store.fetchSchools()
        expect(store.error).toBe('Fail')
        expect(store.loading).toBe(false)
    })

    it('uses cached schools when called within TTL and bypasses cache with force: true', async () => {
        mockSchoolService.getSchools.mockResolvedValue({
            data: { items: [{ id: 1, name: 'S1' }], total: 1 }
        })

        // Call 1: calls API
        await store.fetchSchools()
        expect(mockSchoolService.getSchools).toHaveBeenCalledTimes(1)

        // Call 2 within TTL: no API call
        await store.fetchSchools()
        expect(mockSchoolService.getSchools).toHaveBeenCalledTimes(1)

        // Call 3 with force: true: calls API
        await store.fetchSchools({}, { force: true })
        expect(mockSchoolService.getSchools).toHaveBeenCalledTimes(2)
    })
})
