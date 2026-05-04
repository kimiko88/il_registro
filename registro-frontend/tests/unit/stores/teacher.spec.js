import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useTeacherStore } from '@/stores/teacher'
import authService from '@/services/authService'

// Mock authService
vi.mock('@/services/authService', () => ({
    default: {
        getCurrentUser: vi.fn()
    }
}))

describe('Teacher Store', () => {
    let store

    beforeEach(() => {
        setActivePinia(createPinia())
        store = useTeacherStore()
    })

    it('initializes correctly', () => {
        expect(store.profile).toBeNull()
        expect(store.loading).toBe(false)
        expect(store.error).toBeNull()
        expect(store.isAuthenticated).toBe(false)
    })

    it('fetches profile successfully', async () => {
        authService.getCurrentUser.mockResolvedValue({ firstName: 'Mario', lastName: 'Rossi', is_coordinator: true })
        await store.fetchProfile()

        expect(store.profile).not.toBeNull()
        expect(store.profile.firstName).toBe('Mario')
        expect(store.isAuthenticated).toBe(true)
        expect(store.fullName).toBe('Mario Rossi')
    })

    it('identifies coordinator role', async () => {
        authService.getCurrentUser.mockResolvedValue({ firstName: 'Mario', lastName: 'Rossi', is_coordinator: true })
        await store.fetchProfile()
        // Default mock profile has isCoordinator: true
        expect(store.isCoordinator).toBe(true)
    })

    it('fetches notifications', async () => {
        await store.fetchNotifications()
        expect(store.notifications.length).toBe(0)
    })

    it('calculates fullName correctly', () => {
        store.profile = { firstName: 'John', lastName: 'Doe' }
        expect(store.fullName).toBe('John Doe')

        store.profile = null
        expect(store.fullName).toBe('')
    })
})
