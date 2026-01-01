import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useStudentStore } from '@/stores/student'

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
        const promise = store.fetchProfile()
        expect(store.loading).toBe(true)
        await vi.advanceTimersByTimeAsync(1000)
        await promise
        expect(store.loading).toBe(false)
        expect(store.profile.firstName).toBe('Marco')
    })

    it('getters work', async () => {
        const promise = store.fetchProfile()
        await vi.advanceTimersByTimeAsync(1000)
        await promise
        expect(store.fullName).toBe('Marco Rossi')
        expect(store.className).toBe('5A Scientifico')
        expect(store.isAuthenticated).toBe(true)
    })

    it('fetching notifications', async () => {
        await store.fetchNotifications()
        expect(store.notifications.length).toBe(2)
    })
})
