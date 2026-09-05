import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useTeacherStore, TEACHER_CACHE_TTL } from '@/stores/teacher'
import { useStudentStore, STUDENT_CACHE_TTL } from '@/stores/student'
import { useParentStore, PARENT_CACHE_TTL } from '@/stores/parent'
import authService from '@/services/authService'
import api from '@/services/api'

vi.mock('@/services/authService', () => ({
    default: {
        getCurrentUser: vi.fn()
    }
}))

vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn()
    }
}))

describe('Pinia Stores TTL & Cache Invalidation', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
    })

    describe('useTeacherStore TTL Caching', () => {
        it('caches profile fetch and bypasses with force: true or invalidateCache()', async () => {
            const store = useTeacherStore()
            authService.getCurrentUser.mockResolvedValue({ id: 't-1', first_name: 'Mario', last_name: 'Rossi' })

            await store.fetchProfile()
            expect(authService.getCurrentUser).toHaveBeenCalledTimes(1)

            // Second call within TTL should return cached profile without API call
            await store.fetchProfile()
            expect(authService.getCurrentUser).toHaveBeenCalledTimes(1)

            // Force refresh should trigger API call
            await store.fetchProfile({ force: true })
            expect(authService.getCurrentUser).toHaveBeenCalledTimes(2)

            // Invalidate cache should trigger API call on next fetch
            store.invalidateCache()
            await store.fetchProfile()
            expect(authService.getCurrentUser).toHaveBeenCalledTimes(3)
        })

        it('caches notifications and upcoming colloqui', async () => {
            const store = useTeacherStore()
            api.get.mockResolvedValue({ data: [{ id: 1 }] })

            await store.fetchNotifications()
            await store.fetchNotifications()
            expect(api.get).toHaveBeenCalledWith('/notifications')
            expect(api.get).toHaveBeenCalledTimes(1)

            api.get.mockResolvedValue({ data: [{ id: 1, booked: true }] })
            await store.fetchUpcomingColloqui()
            await store.fetchUpcomingColloqui()
            expect(api.get).toHaveBeenCalledWith('/colloqui/slots/my?upcoming=true')
            expect(api.get).toHaveBeenCalledTimes(2) // 1 notification + 1 colloqui
        })
    })

    describe('useStudentStore TTL Caching', () => {
        it('caches student profile fetch and invalidates on demand', async () => {
            const store = useStudentStore()
            authService.getCurrentUser.mockResolvedValue({ id: 's-1', first_name: 'Luigi', last_name: 'Verdi' })

            await store.fetchProfile()
            await store.fetchProfile()
            expect(authService.getCurrentUser).toHaveBeenCalledTimes(1)

            store.invalidateCache()
            await store.fetchProfile()
            expect(authService.getCurrentUser).toHaveBeenCalledTimes(2)
        })
    })

    describe('useParentStore TTL Caching', () => {
        it('caches parent children list and invalidates on demand', async () => {
            const store = useParentStore()
            api.get.mockResolvedValue({ data: [{ id: 'child-1', first_name: 'Paolo', last_name: 'Rossi' }] })

            await store.fetchChildren()
            await store.fetchChildren()
            expect(api.get).toHaveBeenCalledTimes(1)

            await store.fetchChildren({ force: true })
            expect(api.get).toHaveBeenCalledTimes(2)

            store.invalidateCache()
            await store.fetchChildren()
            expect(api.get).toHaveBeenCalledTimes(3)
        })
    })
})
