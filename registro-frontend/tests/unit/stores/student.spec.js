import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useStudentStore } from '@/stores/student'

describe('Student Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('initializes', () => {
        const store = useStudentStore()
        expect(store.profile).toBeNull()
    })

    it('fetches profile', async () => {
        const store = useStudentStore()
        await store.fetchProfile()
        expect(store.profile).not.toBeNull()
        expect(store.fullName).toContain('Marco')
    })
})
