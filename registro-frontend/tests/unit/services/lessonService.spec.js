import { describe, it, expect, vi, beforeEach } from 'vitest'
import { lessonService } from 'src/services/lessonService'

const mockApi = vi.hoisted(() => ({
    get: vi.fn(),
    post: vi.fn()
}))

vi.mock('@/services/api', () => ({
    default: mockApi
}))

describe('Lesson Service', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('fetches lessons by class', async () => {
        mockApi.get.mockResolvedValue({ data: [] })
        await lessonService.getLessons('class-1', 'sub-1', '2026-07-20')
        expect(mockApi.get).toHaveBeenCalledWith('/lessons/class/class-1', {
            params: { subject_id: 'sub-1', date: '2026-07-20' }
        })
    })

    it('fetches teacher diary across classes', async () => {
        mockApi.get.mockResolvedValue({ data: [{ id: 'l1', topic: 'Math' }] })
        const res = await lessonService.getTeacherDiary('2026-07-01', '2026-07-20')
        expect(mockApi.get).toHaveBeenCalledWith('/lessons/my-diary', {
            params: { from: '2026-07-01', to: '2026-07-20' }
        })
        expect(res.data.length).toBe(1)
    })
})
