import { describe, it, expect, vi, beforeEach } from 'vitest'
import { gradeService } from 'src/services/gradeService'

// Hoist mockApi variable so it's available in vi.mock factory
const mockApi = vi.hoisted(() => ({
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn()
}))

vi.mock('@/services/api', () => ({
    default: mockApi
}))

describe('Grade Service', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('fetches my grades', async () => {
        mockApi.get.mockResolvedValue({ data: [] })
        await gradeService.getMyGrades()
        expect(mockApi.get).toHaveBeenCalledWith('/grades/my-grades')
    })

    it('saves a grade', async () => {
        const payload = { value: 8 }
        await gradeService.saveGrade(payload)
        expect(mockApi.post).toHaveBeenCalledWith('/grades', payload)
    })
})
