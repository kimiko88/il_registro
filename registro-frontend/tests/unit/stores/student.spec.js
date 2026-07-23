import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useStudentStore } from '@/stores/student'
import authService from '@/services/authService'

// Mock authService
vi.mock('@/services/authService', () => ({
    default: {
        getCurrentUser: vi.fn()
    }
}))

describe('Student Store', () => {
    let store

    beforeEach(() => {
        setActivePinia(createPinia())
        store = useStudentStore()
        vi.useFakeTimers()
    })

    afterEach(() => {
        vi.useRealTimers()
    })

    it('initializes correct state', () => {
        expect(store.profile).toBeNull()
        expect(store.loading).toBe(false)
        expect(store.notifications).toEqual([])
    })

    it('fetches profile', async () => {
        authService.getCurrentUser.mockResolvedValue({ firstName: 'Marco', lastName: 'Rossi', className: '5A Scientifico' })
        const promise = store.fetchProfile()
        expect(store.loading).toBe(true)
        await promise
        expect(store.loading).toBe(false)
        expect(store.profile.firstName).toBe('Marco')
    })

    it('getters work', async () => {
        authService.getCurrentUser.mockResolvedValue({ firstName: 'Marco', lastName: 'Rossi', className: '5A Scientifico' })
        await store.fetchProfile()
        expect(store.fullName).toBe('Marco Rossi')
        expect(store.className).toBe('5A Scientifico')
        expect(store.isAuthenticated).toBe(true)
    })

    it('fetching notifications', async () => {
        await store.fetchNotifications()
        expect(store.notifications.length).toBe(0)
    })
})
