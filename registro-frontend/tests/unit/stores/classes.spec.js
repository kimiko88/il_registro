import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useClassesStore } from '@/stores/classes'
import api from '@/services/api'

// Mock API
vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        delete: vi.fn()
    }
}))

describe('Classes Store', () => {
    let store

    beforeEach(() => {
        setActivePinia(createPinia())
        store = useClassesStore()
        vi.clearAllMocks()
    })

    it('fetches classes successfully', async () => {
        const mockClasses = [{ id: '1', name: '1A' }]
        api.get.mockResolvedValue({ data: mockClasses })

        await store.fetchClasses()

        expect(api.get).toHaveBeenCalledWith('/classes', { params: {} })
        expect(store.classes[0].id).toBe('1')
        expect(store.classes[0].name).toBe('1A')
        expect(store.loading).toBe(false)
        expect(store.error).toBe(null)
    })

    it('fetches assigned classes successfully', async () => {
        const mockClasses = [{ id: '1', name: '1A' }]
        api.get.mockResolvedValue({ data: mockClasses })

        await store.fetchAssignedClasses()

        expect(api.get).toHaveBeenCalledWith('/teacher/classes')
        expect(store.classes[0].id).toBe('1')
        expect(store.classes[0].name).toBe('1A')
    })

    it('creates class successfully', async () => {
        const newClass = { name: '1B' }
        const createdClass = { id: '2', ...newClass }
        api.post.mockResolvedValue({ data: createdClass })

        await store.createClass(newClass)

        expect(api.post).toHaveBeenCalledWith('/classes', newClass)
        expect(store.classes).toContainEqual(createdClass)
    })

    it('updates class successfully', async () => {
        store.classes = [{ id: '1', name: '1A' }]
        const updateData = { name: '1A Updated' }
        api.put.mockResolvedValue({ data: { id: '1', ...updateData } })

        await store.updateClass('1', updateData)

        expect(api.put).toHaveBeenCalledWith('/classes/1', updateData)
        expect(store.classes[0].name).toBe('1A Updated')
    })

    it('deletes class successfully', async () => {
        store.classes = [{ id: '1', name: '1A' }, { id: '2', name: '1B' }]
        api.delete.mockResolvedValue({})

        await store.deleteClass('1')

        expect(api.delete).toHaveBeenCalledWith('/classes/1')
        expect(store.classes).toHaveLength(1)
        expect(store.classes[0].id).toBe('2')
    })

    it('handles fetch error', async () => {
        api.get.mockRejectedValue(new Error('Network Error'))

        await store.fetchClasses()

        expect(store.error).toBe('Failed to fetch classes')
        expect(store.loading).toBe(false)
        expect(store.classes).toEqual([])
    })
})
