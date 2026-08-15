import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from 'src/stores/auth'
import { describe, it, expect, beforeEach, afterEach } from 'vitest'

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

            localStorage.setItem('user', JSON.stringify(mockUser))

            const store = useAuthStore()

            expect(store.user).toEqual(mockUser)
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

            store.login(mockUser, 'access-token', 'refresh-token')

            expect(store.user).toEqual(mockUser)
            expect(store.token).toBe('access-token')
            expect(store.refreshToken).toBeNull()
            expect(store.isAuthenticated).toBe(true)
        })

        it('should persist user profile to localStorage while keeping tokens in memory', () => {
            const store = useAuthStore()
            const mockUser = {
                id: '123',
                email: 'test@example.com',
                first_name: 'John',
                last_name: 'Doe',
                role: 'teacher'
            }

            store.login(mockUser, 'access-token', 'refresh-token')

            expect(localStorage.getItem('user')).toBe(JSON.stringify(mockUser))
            expect(localStorage.getItem('token')).toBeNull()
            expect(localStorage.getItem('refreshToken')).toBeNull()
        })
    })

    describe('logout', () => {
        it('should clear user, token, and refreshToken', () => {
            const store = useAuthStore()
            const mockUser = { id: '123', email: 'test@example.com', role: 'teacher' }

            store.login(mockUser, 'access-token', 'refresh-token')
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

            store.login(mockUser, 'access-token', 'refresh-token')
            store.logout()

            expect(localStorage.getItem('user')).toBeNull()
            expect(localStorage.getItem('token')).toBeNull()
            expect(localStorage.getItem('refreshToken')).toBeNull()
        })
    })

    describe('updateUser', () => {
        it('should update user data', () => {
            const store = useAuthStore()
            const initialUser = { id: '123', first_name: 'John', last_name: 'Doe', role: 'teacher' }
            const updatedUser = { id: '123', first_name: 'Jane', last_name: 'Smith', role: 'admin' }

            store.login(initialUser, 'token', 'refresh')
            store.updateUser(updatedUser)

            expect(store.user).toEqual(updatedUser)
            expect(localStorage.getItem('user')).toBe(JSON.stringify(updatedUser))
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

                store.login(mockUser, 'token', 'refresh')

                expect(store.userRole).toBe(role)
                expect(store.isAuthenticated).toBe(true)
            })
        })
    })
})
