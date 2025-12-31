import authService from 'src/services/authService'
import api from 'src/services/api'
import { describe, it, expect, vi, beforeEach } from 'vitest'

vi.mock('src/services/api')

describe('Auth Service', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    describe('login', () => {
        it('should call POST /auth/login with credentials', async () => {
            const mockResponse = {
                data: {
                    user: {
                        id: '123',
                        email: 'test@example.com',
                        first_name: 'John',
                        last_name: 'Doe',
                        role: 'teacher'
                    },
                    access_token: 'access-token',
                    refresh_token: 'refresh-token'
                }
            }

            api.post.mockResolvedValue(mockResponse)

            const result = await authService.login('test@example.com', 'password123')

            expect(api.post).toHaveBeenCalledWith('/auth/login', {
                email: 'test@example.com',
                password: 'password123'
            })
            expect(result).toEqual(mockResponse.data)
        })

        it('should throw error on failed login', async () => {
            api.post.mockRejectedValue(new Error('Invalid credentials'))

            await expect(authService.login('wrong@example.com', 'wrong'))
                .rejects
                .toThrow('Invalid credentials')
        })
    })

    describe('logout', () => {
        it('should call POST /auth/logout with refresh token', async () => {
            const mockResponse = {
                data: { message: 'logged out successfully' }
            }

            api.post.mockResolvedValue(mockResponse)

            const result = await authService.logout('refresh-token-123')

            expect(api.post).toHaveBeenCalledWith('/auth/logout', {
                refresh_token: 'refresh-token-123'
            })
            expect(result).toEqual(mockResponse.data)
        })

        it('should handle logout error', async () => {
            api.post.mockRejectedValue(new Error('Token invalid'))

            await expect(authService.logout('invalid-token'))
                .rejects
                .toThrow('Token invalid')
        })
    })

    describe('getCurrentUser', () => {
        it('should call GET /auth/me', async () => {
            const mockUser = {
                id: '123',
                email: 'test@example.com',
                first_name: 'John',
                last_name: 'Doe',
                role: 'teacher',
                school_id: 'school-1'
            }

            api.get.mockResolvedValue({ data: mockUser })

            const result = await authService.getCurrentUser()

            expect(api.get).toHaveBeenCalledWith('/auth/me')
            expect(result).toEqual(mockUser)
        })

        it('should handle unauthorized error', async () => {
            api.get.mockRejectedValue(new Error('Unauthorized'))

            await expect(authService.getCurrentUser())
                .rejects
                .toThrow('Unauthorized')
        })
    })

    describe('refreshToken', () => {
        it('should call POST /auth/refresh-token', async () => {
            const mockResponse = {
                data: {
                    access_token: 'new-access-token',
                    refresh_token: 'new-refresh-token',
                    expires_in: 3600
                }
            }

            api.post.mockResolvedValue(mockResponse)

            const result = await authService.refreshToken('old-refresh-token')

            expect(api.post).toHaveBeenCalledWith('/auth/refresh-token', {
                refresh_token: 'old-refresh-token'
            })
            expect(result).toEqual(mockResponse.data)
        })
    })

    describe('register', () => {
        it('should call POST /auth/register with user data', async () => {
            const userData = {
                email: 'new@example.com',
                password: 'Password123!',
                first_name: 'New',
                last_name: 'User',
                role: 'student'
            }

            const mockResponse = {
                data: {
                    id: '456',
                    email: 'new@example.com',
                    first_name: 'New',
                    last_name: 'User',
                    role: 'student'
                }
            }

            api.post.mockResolvedValue(mockResponse)

            const result = await authService.register(userData)

            expect(api.post).toHaveBeenCalledWith('/auth/register', userData)
            expect(result).toEqual(mockResponse.data)
        })
    })
})
