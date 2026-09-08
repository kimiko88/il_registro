import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from 'src/stores/auth'
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'

const createMockJWT = (role = 'teacher', expInSeconds = 3600) => {
    const header = btoa(JSON.stringify({ alg: 'HS256', typ: 'JWT' }))
    const payload = btoa(JSON.stringify({
        sub: 'user-123',
        role: role,
        exp: Math.floor(Date.now() / 1000) + expInSeconds
    }))
    return `${header}.${payload}.signature`
}

describe('Auth Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
    })

    afterEach(() => {
        localStorage.clear()
    })

    describe('initial state', () => {
        it('should have null user and token when not logged in', () => {
            const store = useAuthStore()

            expect(store.user).toBeNull()
            expect(store.token).toBeNull()
            expect(store.refreshToken).toBeNull()
            expect(store.isAuthenticated).toBe(false)
        })

        it('should load user from storage while keeping token in memory', () => {
            const mockUser = {
                id: '123',
                email: 'test@example.com',
                first_name: 'John',
                last_name: 'Doe',
                role: 'teacher'
            }

            localStorage.setItem('user', JSON.stringify({ id: '123', role: 'teacher' }))

            const store = useAuthStore()

            expect(store.user).toEqual({ id: '123', role: 'teacher' })
            expect(store.token).toBeNull()
            expect(store.refreshToken).toBeNull()
        })
    })

    describe('computed properties', () => {
        it('should compute userRole correctly', () => {
            const store = useAuthStore()

            expect(store.userRole).toBeNull()

            store.user = { role: 'teacher' }
            expect(store.userRole).toBe('teacher')
        })

        it('should compute userName correctly', () => {
            const store = useAuthStore()

            expect(store.userName).toBe('User')

            store.user = { first_name: 'John', last_name: 'Doe' }
            expect(store.userName).toBe('John Doe')
        })

        it('should return "User" when user is null', () => {
            const store = useAuthStore()
            store.user = null

            expect(store.userName).toBe('User')
        })
    })

    describe('login', () => {
        it('should set user, token, and refreshToken', () => {
            const store = useAuthStore()
            const mockUser = {
                id: '123',
                email: 'test@example.com',
                first_name: 'John',
                last_name: 'Doe',
                role: 'teacher'
            }
            const token = createMockJWT('teacher')

            store.login(mockUser, token, 'refresh-token')

            expect(store.user).toEqual(mockUser)
            expect(store.token).toBe(token)
            expect(store.refreshToken).toBeNull()
            expect(store.isAuthenticated).toBe(true)
        })

        it('should persist minimal user profile to localStorage while keeping tokens in memory', () => {
            const store = useAuthStore()
            const mockUser = {
                id: '123',
                email: 'test@example.com',
                first_name: 'John',
                last_name: 'Doe',
                role: 'teacher'
            }
            const token = createMockJWT('teacher')

            store.login(mockUser, token, 'refresh-token')

            expect(localStorage.getItem('user')).toBe(JSON.stringify({ id: '123', role: 'teacher' }))
            expect(localStorage.getItem('token')).toBeNull()
            expect(localStorage.getItem('refreshToken')).toBeNull()
        })
    })

    describe('logout', () => {
        it('should clear user, token, and refreshToken', () => {
            const store = useAuthStore()
            const mockUser = { id: '123', email: 'test@example.com', role: 'teacher' }
            const token = createMockJWT('teacher')

            store.login(mockUser, token, 'refresh-token')
            expect(store.isAuthenticated).toBe(true)

            store.logout()

            expect(store.user).toBeNull()
            expect(store.token).toBeNull()
            expect(store.refreshToken).toBeNull()
            expect(store.isAuthenticated).toBe(false)
        })

        it('should clear localStorage', () => {
            const store = useAuthStore()
            const mockUser = { id: '123', email: 'test@example.com', role: 'teacher' }
            const token = createMockJWT('teacher')

            store.login(mockUser, token, 'refresh-token')
            store.logout()

            expect(localStorage.getItem('user')).toBeNull()
            expect(localStorage.getItem('token')).toBeNull()
            expect(localStorage.getItem('refreshToken')).toBeNull()
        })

        it('should invalidate api-static-lists cache on logout if caches API is available', () => {
            const deleteMock = vi.fn().mockReturnValue(Promise.resolve(true))
            const originalCaches = window.caches
            Object.defineProperty(window, 'caches', {
                value: { delete: deleteMock },
                configurable: true,
                writable: true
            })

            const store = useAuthStore()
            store.logout()

            expect(deleteMock).toHaveBeenCalledWith('api-static-lists')

            if (originalCaches !== undefined) {
                Object.defineProperty(window, 'caches', {
                    value: originalCaches,
                    configurable: true,
                    writable: true
                })
            } else {
                delete window.caches
            }
        })
    })

    describe('updateUser', () => {
        it('should update user data', () => {
            const store = useAuthStore()
            const initialUser = { id: '123', first_name: 'John', last_name: 'Doe', role: 'teacher' }
            const updatedUser = { id: '123', first_name: 'Jane', last_name: 'Smith', role: 'admin' }
            const token = createMockJWT('teacher')

            store.login(initialUser, token, 'refresh')
            store.updateUser(updatedUser)

            expect(store.user).toEqual(updatedUser)
            expect(localStorage.getItem('user')).toBe(JSON.stringify({ id: '123', role: 'admin' }))
        })
    })

    describe('role-based functionality', () => {
        const roles = ['admin', 'secretary', 'teacher', 'student', 'parent']

        roles.forEach(role => {
            it(`should handle ${role} role correctly`, () => {
                const store = useAuthStore()
                const mockUser = {
                    id: '123',
                    email: 'test@example.com',
                    first_name: 'Test',
                    last_name: 'User',
                    role: role
                }
                const token = createMockJWT(role)

                store.login(mockUser, token, 'refresh')

                expect(store.userRole).toBe(role)
                expect(store.isAuthenticated).toBe(true)
            })
        })
    })
})

