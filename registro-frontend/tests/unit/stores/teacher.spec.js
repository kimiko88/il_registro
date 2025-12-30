import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useTeacherStore } from '@/stores/teacher'

describe('Teacher Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('initializes with default state', () => {
        const store = useTeacherStore()
        expect(store.profile).toBeNull()
        expect(store.loading).toBe(false)
    })

    it('can fetch profile (mock)', async () => {
        const store = useTeacherStore()
        const p = store.fetchProfile() // async
        expect(store.loading).toBe(true)
        await p
        expect(store.profile).not.toBeNull()
        expect(store.profile.firstName).toBe('Mario')
        expect(store.loading).toBe(false)
    })
})
