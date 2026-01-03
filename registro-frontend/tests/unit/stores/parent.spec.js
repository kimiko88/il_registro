import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useParentStore } from '@/stores/parent'
import { api } from '@/boot/axios'

vi.mock('@/boot/axios', () => ({
    api: {
        get: vi.fn()
    }
}))

describe('Parent Store', () => {
    let store

    beforeEach(() => {
        setActivePinia(createPinia())
        store = useParentStore()
        vi.clearAllMocks()
        localStorage.clear()
    })

    it('initializes with default state', () => {
        expect(store.children).toEqual([])
        expect(store.selectedChildId).toBeNull()
        expect(store.loading).toBe(false)
        expect(store.error).toBeNull()
    })

    it('fetches children successfully', async () => {
        const mockChildren = [
            { id: '1', firstName: 'Child', lastName: 'One' },
            { id: '2', firstName: 'Child', lastName: 'Two' }
        ]
        api.get.mockResolvedValue({ data: mockChildren })

        await store.fetchChildren()

        expect(store.children).toEqual(mockChildren)
        expect(store.loading).toBe(false)
        expect(store.error).toBeNull()
        // Should auto-select first child
        expect(store.selectedChildId).toBe('1')
        expect(localStorage.getItem('selectedChildId')).toBe('1')
    })

    it('preserves existing selection if valid', async () => {
        // Reset Pinia to ensure we start fresh with new localStorage
        setActivePinia(createPinia())
        localStorage.setItem('selectedChildId', '2')
        store = useParentStore() // Now it should read '2'

        const mockChildren = [
            { id: '1', firstName: 'Child', lastName: 'One' },
            { id: '2', firstName: 'Child', lastName: 'Two' }
        ]
        api.get.mockResolvedValue({ data: mockChildren })

        await store.fetchChildren()

        expect(store.selectedChildId).toBe('2')
    })

    it('handles fetch error', async () => {
        const error = new Error('Network Error')
        api.get.mockRejectedValue(error)

        await store.fetchChildren()

        expect(store.children).toEqual([])
        expect(store.loading).toBe(false)
        expect(store.error).toBe('Failed to load children')
    })

    it('selects child manually', () => {
        store.children = [
            { id: '1', firstName: 'Child', lastName: 'One' },
            { id: '2', firstName: 'Child', lastName: 'Two' }
        ]

        store.selectChild('2')
        expect(store.selectedChildId).toBe('2')
        expect(localStorage.getItem('selectedChildId')).toBe('2')
    })

    it('ignores selection of invalid child id', () => {
        store.children = [{ id: '1' }]
        store.selectedChildId = '1'

        store.selectChild('999')
        expect(store.selectedChildId).toBe('1')
    })

    it('computes selectedChild correctly', () => {
        store.children = [
            { id: '1', name: 'Huey' },
            { id: '2', name: 'Dewey' }
        ]

        store.selectedChildId = '2'
        expect(store.selectedChild).toEqual({ id: '2', name: 'Dewey' })

        store.selectedChildId = '3' // Invalid
        expect(store.selectedChild).toEqual({ id: '1', name: 'Huey' }) // Fallback to first

        store.children = []
        expect(store.selectedChild).toBeNull()
    })
})
