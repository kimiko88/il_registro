import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useDocumentsStore } from '@/stores/documents'
import documentService from '@/services/documentService'

// Mock Service
vi.mock('@/services/documentService', () => ({
    default: {
        getInbox: vi.fn(),
        getDocument: vi.fn(),
        reviewDocument: vi.fn()
    }
}))

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

describe('Documents Store', () => {
    let store

    beforeEach(() => {
        vi.useFakeTimers()
        setActivePinia(createPinia())
        store = useDocumentsStore()
        vi.clearAllMocks()
    })

    afterEach(() => {
        vi.useRealTimers()
    })

    it('fetches inbox successfully', async () => {
        const mockResponse = { data: { items: [{ id: 1 }], total: 1 } }
        documentService.getInbox.mockResolvedValue(mockResponse)

        await store.fetchInbox({ page: 1 })

        expect(documentService.getInbox).toHaveBeenCalled()
        expect(store.inbox).toEqual([{ id: 1 }])
        expect(store.pagination.rowsNumber).toBe(1)
        expect(store.loading).toBe(false)
    })

    it('handles fetch inbox error', async () => {
        documentService.getInbox.mockRejectedValue(new Error('Fetch Error'))

        await store.fetchInbox()

        expect(store.error).toBe('Fetch Error')
        expect(store.loading).toBe(false)
    })

    it('fetches document successfully', async () => {
        const mockDoc = { id: 1, title: 'Test' }
        documentService.getDocument.mockResolvedValue({ data: mockDoc })

        await store.fetchDocument(1)

        expect(documentService.getDocument).toHaveBeenCalledWith(1)
        expect(store.currentDocument).toEqual(mockDoc)
    })

    it('reviews document and refreshes inbox', async () => {
        documentService.reviewDocument.mockResolvedValue({})
        // Mock getInbox for the refresh call
        documentService.getInbox.mockResolvedValue({ data: { items: [], total: 0 } })

        await store.reviewDocument(1, 'approve', 'Good job')

        expect(documentService.reviewDocument).toHaveBeenCalledWith(1, 'approve', 'Good job')
        expect(documentService.getInbox).toHaveBeenCalled()
    })

    it('fetches my documents (teacher) with simulated delay', async () => {
        documentService.getInbox.mockResolvedValue({ data: { items: [{ id: 1 }, { id: 2 }], total: 2 } })
        const promise = store.fetchMyDocuments()

        expect(store.loading).toBe(true)
        await vi.runAllTimersAsync()
        await promise

        expect(store.inbox).toHaveLength(2) // Hardcoded 2 items
        expect(store.loading).toBe(false)
    })

    it('fetches templates', async () => {
        mockGet.mockResolvedValue({ data: [{ id: 1, type: 'PDP' }, { id: 2, type: 'PEI' }] })
        const templates = await store.fetchTemplates()
        expect(templates).toHaveLength(2)
        expect(templates[0].type).toBe('PDP')
    })

    it('creates document', async () => {
        const docData = { title: 'New Doc' }
        mockPost.mockResolvedValue({ data: { id: 100, ...docData } })
        const promise = store.createDocument(docData)

        const result = await promise

        expect(result.title).toBe('New Doc')
        expect(store.inbox[0].title).toBe('New Doc')
    })

    it('fetches my files (student)', async () => {
        mockGet.mockResolvedValue({ data: [{ id: 1 }, { id: 2 }] })
        const promise = store.fetchMyFiles()

        expect(store.loading).toBe(true)
        await promise

        expect(store.inbox).toHaveLength(2)
        expect(store.loading).toBe(false)
    })
})
