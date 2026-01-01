import { describe, it, expect, vi, beforeEach } from 'vitest'
import schoolService from '@/services/schoolService'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn(),
        post: vi.fn(),
        patch: vi.fn(),
        delete: vi.fn()
    }
}))

describe('schoolService', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('getSchools calls api.get', () => {
        schoolService.getSchools({})
        expect(api.get).toHaveBeenCalledWith('/schools', { params: {} })
    })

    it('getSchool calls api.get', () => {
        schoolService.getSchool('1')
        expect(api.get).toHaveBeenCalledWith('/schools/1')
    })

    it('createSchool calls api.post', () => {
        schoolService.createSchool({})
        expect(api.post).toHaveBeenCalledWith('/schools', {})
    })

    it('updateSchool calls api.patch', () => {
        schoolService.updateSchool('1', {})
        expect(api.patch).toHaveBeenCalledWith('/schools/1', {})
    })

    it('deleteSchool calls api.delete', () => {
        schoolService.deleteSchool('1')
        expect(api.delete).toHaveBeenCalledWith('/schools/1')
    })

    it('importSchools calls api.post with FormData', () => {
        const file = new File([''], 'test.csv')
        schoolService.importSchools(file)
        expect(api.post).toHaveBeenCalledWith('/schools/import', expect.any(FormData), expect.objectContaining({ headers: expect.anything() }))
    })
})
