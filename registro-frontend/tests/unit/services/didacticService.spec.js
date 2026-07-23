import { describe, it, expect, vi, beforeEach } from 'vitest'
import didacticService from 'src/services/didacticService'

const mockApi = vi.hoisted(() => ({
    get: vi.fn(),
    post: vi.fn(),
    delete: vi.fn()
}))

vi.mock('@/services/api', () => ({
    default: mockApi
}))

describe('Didactic Service', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('fetches didactic materials by class id', async () => {
        mockApi.get.mockResolvedValue({ data: [] })
        await didacticService.getMaterials('class-1')
        expect(mockApi.get).toHaveBeenCalledWith('/didactic-materials/class/class-1')
    })

    it('creates a didactic material', async () => {
        const payload = { title: 'Slide' }
        await didacticService.createMaterial(payload)
        expect(mockApi.post).toHaveBeenCalledWith('/didactic-materials', payload)
    })

    it('deletes a didactic material by id', async () => {
        await didacticService.deleteMaterial('mat-1')
        expect(mockApi.delete).toHaveBeenCalledWith('/didactic-materials/mat-1')
    })
})
