import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useCommunicationsStore } from '@/stores/communications'
import api from '@/services/api'

// Mock api
vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn(),
        post: vi.fn(),
        delete: vi.fn()
    }
}))

describe('Communications Store', () => {
    let store

    beforeEach(() => {
        setActivePinia(createPinia())
        store = useCommunicationsStore()
        vi.clearAllMocks()
    })

    it('fetches communications successfully', async () => {
        const mockMsgs = [{ id: 1, subject: 'Hello' }]
        api.get.mockResolvedValue({ data: mockMsgs })

        await store.fetchCommunications()

        expect(store.loading).toBe(false)
        expect(store.communications).toEqual(mockMsgs)
        expect(store.error).toBeNull()
        expect(api.get).toHaveBeenCalledWith('/communications')
    })

    it('handles fetch communications error', async () => {
        api.get.mockRejectedValue({ response: { data: { error: 'Server Error' } } })

        await store.fetchCommunications()

        expect(store.loading).toBe(false)
        expect(store.error).toBe('Server Error')
        expect(store.communications).toEqual([])
    })

    it('sends message successfully', async () => {
        const payload = { title: 'Hi', content: 'Test' }
        const responseData = { id: 2, ...payload }
        api.post.mockResolvedValue({ data: responseData })

        const result = await store.sendMessage(payload)

        expect(result).toEqual(responseData)
        expect(store.communications[0]).toEqual(responseData) // Unshift check
        expect(api.post).toHaveBeenCalledWith('/communications', payload)
    })

    it('handles send message error', async () => {
        api.post.mockRejectedValue(new Error('Send Failed'))

        await expect(store.sendMessage({})).rejects.toThrow('Send Failed')
        expect(store.loading).toBe(false)
    })
})
