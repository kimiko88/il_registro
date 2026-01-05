import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useDocumentsStore } from './documents'
import { api } from 'src/boot/axios'

// Mock api
vi.mock('src/boot/axios', () => ({
    api: {
        post: vi.fn(),
        get: vi.fn()
    }
}))

describe('Documents Store - Signing', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
    })

    it('signs a document successfully', async () => {
        const store = useDocumentsStore()
        const docId = '123'
        const pin = '1234'
        const mockResponse = { data: { id: 'sig-1', document_id: docId, signature_hash: 'abc' } }

        api.post.mockResolvedValue(mockResponse)

        const result = await store.signDocument(docId, pin)

        expect(api.post).toHaveBeenCalledWith('/signatures/', { document_id: docId, pin: pin })
        expect(result).toEqual(mockResponse.data)
        expect(store.loading).toBe(false)
    })

    it('handles signing errors', async () => {
        const store = useDocumentsStore()
        const error = new Error('Invalid PIN')
        api.post.mockRejectedValue(error)

        await expect(store.signDocument('123', '0000')).rejects.toThrow('Invalid PIN')
        expect(store.loading).toBe(false)
    })
})
