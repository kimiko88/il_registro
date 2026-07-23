import { describe, it, expect, vi, beforeEach } from 'vitest'
import { pctoService } from 'src/services/pctoService'

const mockApi = vi.hoisted(() => ({
    get: vi.fn(),
    post: vi.fn()
}))

vi.mock('@/services/api', () => ({
    default: mockApi
}))

describe('PCTO Service', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('fetches student PCTO projects', async () => {
        mockApi.get.mockResolvedValue({ data: [] })
        await pctoService.getMyProjects()
        expect(mockApi.get).toHaveBeenCalledWith('/pcto/my-projects')
    })

    it('logs hours for PCTO project', async () => {
        const payload = { project_id: 'p1', hours: 4, date: '2026-07-20' }
        mockApi.post.mockResolvedValue({ data: {} })
        await pctoService.logHours(payload)
        expect(mockApi.post).toHaveBeenCalledWith('/pcto/hours', payload)
    })

    it('approves host company hours log', async () => {
        mockApi.post.mockResolvedValue({ data: { message: 'hours log status updated' } })
        await pctoService.approveHours('hour-123', true)
        expect(mockApi.post).toHaveBeenCalledWith('/pcto/hours/hour-123/approve', { approved: true })
    })
})
