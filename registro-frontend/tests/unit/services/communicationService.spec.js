import { describe, it, expect, vi, beforeEach } from 'vitest'
import { communicationService } from 'src/services/communicationService'

const mockApi = vi.hoisted(() => ({
    get: vi.fn(),
    post: vi.fn()
}))

vi.mock('@/services/api', () => ({
    default: mockApi
}))

describe('Communication Service', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('fetches communications list', async () => {
        mockApi.get.mockResolvedValue({ data: [] })
        await communicationService.getMessages()
        expect(mockApi.get).toHaveBeenCalledWith('/communications')
    })

    it('sends a communication message', async () => {
        const payload = { subject: 'Alert' }
        mockApi.post.mockResolvedValue({ data: {} })
        await communicationService.sendMessage(payload)
        expect(mockApi.post).toHaveBeenCalledWith('/communications', payload)
    })

    it('signs a communication message', async () => {
        mockApi.post.mockResolvedValue({ data: {} })
        await communicationService.signMessage('msg-1')
        expect(mockApi.post).toHaveBeenCalledWith('/communications/msg-1/sign')
    })

    it('gets signatures for a message', async () => {
        mockApi.get.mockResolvedValue({ data: [] })
        await communicationService.getSignatures('msg-1')
        expect(mockApi.get).toHaveBeenCalledWith('/communications/msg-1/signatures')
    })
})
