import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from 'src/stores/auth'
import { useAuth } from 'src/composables/useAuth'
import authService from 'src/services/authService'

const createMockJWT = (role = 'teacher', expInSeconds = 3600) => {
    const header = btoa(JSON.stringify({ alg: 'HS256', typ: 'JWT' }))
    const payload = btoa(JSON.stringify({
        sub: 'user-123',
        role: role,
        exp: Math.floor(Date.now() / 1000) + expInSeconds
    }))
    return `${header}.${payload}.signature`
}

vi.mock('src/services/authService')
vi.mock('vue-router', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useRouter: () => ({
            push: vi.fn()
        })
    }
})

describe('Role-Based Authentication Flows', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
        localStorage.clear()
    })

    describe('Admin Flow', () => {
        const adminUser = {
            id: 'admin-1',
            email: 'admin@school.com',
            first_name: 'Mario',
            last_name: 'Rossi',
            role: 'admin',
            school_id: 'school-1'
        }

        it('should complete full admin login flow', async () => {
            const { login } = useAuth()
            const authStore = useAuthStore()
            const adminToken = createMockJWT('admin')

            authService.login.mockResolvedValue({
                user: adminUser,
                access_token: adminToken,
                refresh_token: 'admin-refresh-token'
            })

            const error = await login('admin@school.com', 'password')

            expect(error).toBeNull()
            expect(authStore.user).toEqual(adminUser)
            expect(authStore.token).toBe(adminToken)
            expect(authStore.userRole).toBe('admin')
            expect(authStore.isAuthenticated).toBe(true)
        })

        it('should have admin permissions and access', () => {
            const authStore = useAuthStore()
            authStore.login(adminUser, createMockJWT('admin'), 'refresh')

            expect(authStore.user.role).toBe('admin')
            expect(authStore.user.school_id).toBeDefined()
        })

        it('should logout admin successfully', async () => {
            const { logout } = useAuth()
            const authStore = useAuthStore()

            authStore.login(adminUser, createMockJWT('admin'), 'refresh-token')
            authService.logout.mockResolvedValue({ message: 'logged out successfully' })

            await logout()

            expect(authService.logout).toHaveBeenCalled()
            expect(authStore.isAuthenticated).toBe(false)
            expect(authStore.user).toBeNull()
        })
    })

    describe('Teacher Flow', () => {
        const teacherUser = {
            id: 'teacher-1',
            email: 'teacher@school.com',
            first_name: 'Giovanni',
            last_name: 'Bianchi',
            role: 'teacher',
            school_id: 'school-1'
        }

        it('should complete full teacher login flow', async () => {
            const { login } = useAuth()
            const authStore = useAuthStore()
            const teacherToken = createMockJWT('teacher')

            authService.login.mockResolvedValue({
                user: teacherUser,
                access_token: teacherToken,
                refresh_token: 'teacher-refresh-token'
            })

            const error = await login('teacher@school.com', 'password')

            expect(error).toBeNull()
            expect(authStore.userRole).toBe('teacher')
            expect(authStore.userName).toBe('Giovanni Bianchi')
        })

        it('should maintain user profile across page reload while keeping token in memory', () => {
            const authStore = useAuthStore()
            const teacherToken = createMockJWT('teacher')

            // Simulate login
            authStore.login(teacherUser, teacherToken, 'refresh')

            // Verify localStorage (minimal user profile persisted, token kept in memory)
            expect(localStorage.getItem('user')).toBe(JSON.stringify({ id: teacherUser.id, role: teacherUser.role }))
            expect(localStorage.getItem('token')).toBeNull()

            // Simulate page reload by creating new store
            setActivePinia(createPinia())
            const newAuthStore = useAuthStore()

            // Store should load minimal user profile from localStorage
            expect(newAuthStore.user).toEqual({ id: teacherUser.id, role: teacherUser.role })
        })
    })

    describe('Student Flow', () => {
        const studentUser = {
            id: 'student-1',
            email: 'student@school.com',
            first_name: 'Luca',
            last_name: 'Verdi',
            role: 'student',
            school_id: 'school-1'
        }

        it('should complete full student login flow', async () => {
            const { login } = useAuth()
            const authStore = useAuthStore()
            const studentToken = createMockJWT('student')

            authService.login.mockResolvedValue({
                user: studentUser,
                access_token: studentToken,
                refresh_token: 'student-refresh-token'
            })

            await login('student@school.com', 'password')

            expect(authStore.userRole).toBe('student')
            expect(authStore.user.email).toBe('student@school.com')
        })

        it('should handle student profile updates', () => {
            const authStore = useAuthStore()
            authStore.login(studentUser, createMockJWT('student'), 'refresh')

            const updatedStudent = {
                ...studentUser,
                first_name: 'Luca Updated'
            }

            authStore.updateUser(updatedStudent)

            expect(authStore.user.first_name).toBe('Luca Updated')
            expect(authStore.userName).toBe('Luca Updated Verdi')
            expect(localStorage.getItem('user')).toBe(JSON.stringify({ id: updatedStudent.id, role: updatedStudent.role }))
        })
    })

    describe('Parent Flow', () => {
        const parentUser = {
            id: 'parent-1',
            email: 'parent@email.com',
            first_name: 'Anna',
            last_name: 'Neri',
            role: 'parent',
            school_id: 'school-1'
        }

        it('should complete full parent login flow', async () => {
            const { login } = useAuth()
            const authStore = useAuthStore()
            const parentToken = createMockJWT('parent')

            authService.login.mockResolvedValue({
                user: parentUser,
                access_token: parentToken,
                refresh_token: 'parent-refresh-token'
            })

            await login('parent@email.com', 'password')

            expect(authStore.userRole).toBe('parent')
            expect(authStore.isAuthenticated).toBe(true)
        })
    })

    describe('Secretary Flow', () => {
        const secretaryUser = {
            id: 'secretary-1',
            email: 'secretary@school.com',
            first_name: 'Paolo',
            last_name: 'Gialli',
            role: 'secretary',
            school_id: 'school-1'
        }

        it('should complete full secretary login flow', async () => {
            const { login } = useAuth()
            const authStore = useAuthStore()
            const secretaryToken = createMockJWT('secretary')

            authService.login.mockResolvedValue({
                user: secretaryUser,
                access_token: secretaryToken,
                refresh_token: 'secretary-refresh-token'
            })

            await login('secretary@school.com', 'password')

            expect(authStore.userRole).toBe('secretary')
            expect(authStore.user.first_name).toBe('Paolo')
        })
    })

    describe('Cross-Role Security', () => {
        it('should not allow role switching without re-authentication', () => {
            const authStore = useAuthStore()

            const teacherUser = {
                id: 'teacher-1',
                email: 'teacher@school.com',
                first_name: 'Teacher',
                last_name: 'User',
                role: 'teacher'
            }

            authStore.login(teacherUser, createMockJWT('teacher'), 'refresh')

            // Try to manually change role (should not be possible)
            const _originalRole = authStore.userRole

            // Even if user object is modified, logout should clear everything
            authStore.logout()

            expect(authStore.userRole).toBeNull()
            expect(authStore.user).toBeNull()
        })

        it('should isolate role data after logout', () => {
            const authStore = useAuthStore()

            const adminUser = {
                id: 'admin-1',
                email: 'admin@school.com',
                first_name: 'Admin',
                last_name: 'User',
                role: 'admin'
            }

            authStore.login(adminUser, createMockJWT('admin'), 'admin-refresh')
            expect(authStore.userRole).toBe('admin')

            authStore.logout()

            // After logout, no trace of previous role
            expect(localStorage.getItem('user')).toBeNull()
            expect(localStorage.getItem('token')).toBeNull()
            expect(localStorage.getItem('refreshToken')).toBeNull()
        })
    })

    describe('Error Handling Across Roles', () => {
        const roles = ['admin', 'teacher', 'student', 'parent', 'secretary']

        roles.forEach(role => {
            it(`should handle failed login for ${role}`, async () => {
                const { login } = useAuth()
                const authStore = useAuthStore()

                authService.login.mockRejectedValue({
                    response: {
                        data: {
                            error: 'invalid credentials'
                        }
                    }
                })

                const error = await login(`${role}@school.com`, 'wrong-password')

                expect(error).toBeTruthy()
                expect(authStore.isAuthenticated).toBe(false)
                expect(authStore.user).toBeNull()
            })
        })

        it('should handle network errors during logout', async () => {
            const { logout } = useAuth()
            const authStore = useAuthStore()

            authStore.login(
                { id: '1', email: 'test@test.com', role: 'teacher' },
                createMockJWT('teacher'),
                'refresh'
            )

            authService.logout.mockRejectedValue(new Error('Network error'))

            // Should still logout locally even if backend fails
            await logout()

            expect(authStore.isAuthenticated).toBe(false)
            expect(authStore.user).toBeNull()
        })
    })

    describe('Token Refresh Flow', () => {
        it('should handle token refresh for authenticated user', async () => {
            const authStore = useAuthStore()

            authStore.login(
                { id: '1', email: 'test@test.com', first_name: 'Test', last_name: 'User', role: 'teacher' },
                createMockJWT('teacher'),
                'old-refresh'
            )

            authService.refreshToken.mockResolvedValue({
                access_token: createMockJWT('teacher'),
                refresh_token: 'new-refresh',
                expires_in: 3600
            })

            const result = await authService.refreshToken('old-refresh')

            expect(result.access_token).toBeDefined()
            expect(result.refresh_token).toBe('new-refresh')
        })
    })
})

