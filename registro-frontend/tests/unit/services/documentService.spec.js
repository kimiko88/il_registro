import { describe, it, expect, vi, beforeEach } from 'vitest'
import documentService from '@/services/documentService'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn(),
        post: vi.fn()
    }
}))

describe('documentService', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('getInbox calls api.get', () => {
        documentService.getInbox({})
        expect(api.get).toHaveBeenCalledWith('/documents/inbox', { params: {} })
    })

    it('getDocument calls api.get', () => {
        documentService.getDocument('1')
        expect(api.get).toHaveBeenCalledWith('/documents/1')
    })

    it('reviewDocument calls api.post', () => {
        documentService.reviewDocument('1', 'ok', 'note')
        expect(api.post).toHaveBeenCalledWith('/documents/1/review', { decision: 'ok', notes: 'note' })
    })

    it('archiveDocument calls api.post', () => {
        documentService.archiveDocument('1')
        expect(api.post).toHaveBeenCalledWith('/documents/1/archive')
    })

    it('getDocumentVersions calls api.get', () => {
        documentService.getDocumentVersions('1')
        expect(api.get).toHaveBeenCalledWith('/documents/1/versions')
    })
})
