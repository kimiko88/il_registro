import { describe, it, expect, vi, beforeEach } from 'vitest'
import authService from '@/services/authService'
import api from '@/services/api'

// Mock API
vi.mock('@/services/api', () => ({
    default: {
        post: vi.fn(),
        get: vi.fn()
    }
}))

describe('Auth Service', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('login', async () => {
        const email = 'test@test.com'
        const password = 'password'
        api.post.mockResolvedValue({ data: { token: '123' } })

        const res = await authService.login(email, password)

        expect(api.post).toHaveBeenCalledWith('/auth/login', { email, password })
        expect(res.token).toBe('123')
    })

    it('register', async () => {
        const user = { name: 'Test' }
        api.post.mockResolvedValue({ data: { id: 1 } })

        const res = await authService.register(user)

        expect(api.post).toHaveBeenCalledWith('/auth/register', user)
        expect(res.id).toBe(1)
    })

    it('logout', async () => {
        api.post.mockResolvedValue({})

        await authService.logout('refresh_token_123')

        expect(api.post).toHaveBeenCalledWith('/auth/logout', { refresh_token: 'refresh_token_123' })
    })

    it('get current user', async () => {
        api.get.mockResolvedValue({ data: { id: 1, name: 'User' } })
        const res = await authService.getCurrentUser()
        expect(api.get).toHaveBeenCalledWith('/auth/me')
        expect(res.id).toBe(1)
    })
})
