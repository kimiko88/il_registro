import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useCommunicationsStore } from '@/stores/communications'

// Mock api from boot/axios
const { mockGet, mockPost } = vi.hoisted(() => ({
    mockGet: vi.fn(),
    mockPost: vi.fn()
}))

vi.mock('src/boot/axios', () => ({
    api: {
        get: mockGet,
        post: mockPost
    }
}))

describe('Communications Store', () => {
    let store

    beforeEach(() => {
        setActivePinia(createPinia())
        store = useCommunicationsStore()
        vi.clearAllMocks()
    })

    it('fetches messages successfully', async () => {
        const mockMsgs = [{ id: 1, subject: 'Hello' }]
        mockGet.mockResolvedValue({ data: mockMsgs })

        await store.fetchMessages()

        expect(store.loading).toBe(false)
        expect(store.messages).toEqual(mockMsgs)
        expect(store.error).toBeNull()
        expect(mockGet).toHaveBeenCalledWith('/communications')
    })

    it('handles fetch messages error', async () => {
        mockGet.mockRejectedValue({ response: { data: { error: 'Server Error' } } })

        await store.fetchMessages()

        expect(store.loading).toBe(false)
        expect(store.error).toBe('Server Error')
        expect(store.messages).toEqual([])
    })

    it('sends message successfully', async () => {
        const payload = { subject: 'Hi', body: 'Test' }
        const responseData = { id: 2, ...payload }
        mockPost.mockResolvedValue({ data: responseData })

        const result = await store.sendMessage(payload)

        expect(result).toEqual(responseData)
        expect(store.messages[0]).toEqual(responseData) // Unshift check
        expect(mockPost).toHaveBeenCalledWith('/communications', payload)
    })

    it('handles send message error', async () => {
        mockPost.mockRejectedValue(new Error('Send Failed'))

        await expect(store.sendMessage({})).rejects.toThrow('Send Failed')
        expect(store.loading).toBe(false)
    })
})
