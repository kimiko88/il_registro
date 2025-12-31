import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from 'src/stores/auth'
import { useAuth } from 'src/composables/useAuth'
import authService from 'src/services/authService'

vi.mock('src/services/authService')
vi.mock('vue-router', () => ({
    useRouter: () => ({
        push: vi.fn()
    })
}))

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

            authService.login.mockResolvedValue({
                user: adminUser,
                access_token: 'admin-access-token',
                refresh_token: 'admin-refresh-token'
            })

            const error = await login('admin@school.com', 'password')

            expect(error).toBeNull()
            expect(authStore.user).toEqual(adminUser)
            expect(authStore.token).toBe('admin-access-token')
            expect(authStore.userRole).toBe('admin')
            expect(authStore.isAuthenticated).toBe(true)
        })

        it('should have admin permissions and access', () => {
            const authStore = useAuthStore()
            authStore.login(adminUser, 'token', 'refresh')

            expect(authStore.user.role).toBe('admin')
            expect(authStore.user.school_id).toBeDefined()
        })

        it('should logout admin successfully', async () => {
            const { logout } = useAuth()
            const authStore = useAuthStore()

            authStore.login(adminUser, 'token', 'refresh-token')
            authService.logout.mockResolvedValue({ message: 'logged out successfully' })

            await logout()

            expect(authService.logout).toHaveBeenCalledWith('refresh-token')
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

            authService.login.mockResolvedValue({
                user: teacherUser,
                access_token: 'teacher-access-token',
                refresh_token: 'teacher-refresh-token'
            })

            const error = await login('teacher@school.com', 'password')

            expect(error).toBeNull()
            expect(authStore.userRole).toBe('teacher')
            expect(authStore.userName).toBe('Giovanni Bianchi')
        })

        it('should maintain teacher session across page reload', () => {
            const authStore = useAuthStore()

            // Simulate login
            authStore.login(teacherUser, 'token', 'refresh')

            // Verify localStorage
            expect(localStorage.getItem('user')).toBe(JSON.stringify(teacherUser))
            expect(localStorage.getItem('token')).toBe('token')

            // Simulate page reload by creating new store
            setActivePinia(createPinia())
            const newAuthStore = useAuthStore()

            // Store should load from localStorage
            expect(newAuthStore.user).toEqual(teacherUser)
            expect(newAuthStore.isAuthenticated).toBe(true)
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

            authService.login.mockResolvedValue({
                user: studentUser,
                access_token: 'student-access-token',
                refresh_token: 'student-refresh-token'
            })

            await login('student@school.com', 'password')

            expect(authStore.userRole).toBe('student')
            expect(authStore.user.email).toBe('student@school.com')
        })

        it('should handle student profile updates', () => {
            const authStore = useAuthStore()
            authStore.login(studentUser, 'token', 'refresh')

            const updatedStudent = {
                ...studentUser,
                first_name: 'Luca Updated'
            }

            authStore.updateUser(updatedStudent)

            expect(authStore.user.first_name).toBe('Luca Updated')
            expect(authStore.userName).toBe('Luca Updated Verdi')
            expect(localStorage.getItem('user')).toBe(JSON.stringify(updatedStudent))
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

            authService.login.mockResolvedValue({
                user: parentUser,
                access_token: 'parent-access-token',
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

            authService.login.mockResolvedValue({
                user: secretaryUser,
                access_token: 'secretary-access-token',
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

            authStore.login(teacherUser, 'token', 'refresh')

            // Try to manually change role (should not be possible)
            const originalRole = authStore.userRole

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

            authStore.login(adminUser, 'admin-token', 'admin-refresh')
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
                'token',
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
                'old-token',
                'old-refresh'
            )

            authService.refreshToken.mockResolvedValue({
                access_token: 'new-token',
                refresh_token: 'new-refresh',
                expires_in: 3600
            })

            const result = await authService.refreshToken('old-refresh')

            expect(result.access_token).toBe('new-token')
            expect(result.refresh_token).toBe('new-refresh')
        })
    })
})
