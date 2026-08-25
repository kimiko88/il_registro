import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from 'src/stores/auth'
import { usePermissions } from 'src/composables/usePermissions'

describe('usePermissions composable', () => {
    let authStore

    beforeEach(() => {
        setActivePinia(createPinia())
        authStore = useAuthStore()
    })

    describe('SuperAdmin Role', () => {
        beforeEach(() => {
            authStore.user = {
                id: 'sa-1',
                email: 'superadmin@example.com',
                role: 'superadmin',
                school_id: null
            }
        })

        it('should identify superadmin and give global privileges', () => {
            const {
                isSuperAdmin,
                isAdmin,
                isAdminOrSuperAdmin,
                canManageAllSchools,
                canCreateSchools,
                canAccessMonitoring
            } = usePermissions()

            expect(isSuperAdmin.value).toBe(true)
            expect(isAdmin.value).toBe(false)
            expect(isAdminOrSuperAdmin.value).toBe(true)
            expect(canManageAllSchools.value).toBe(true)
            expect(canCreateSchools.value).toBe(true)
            expect(canAccessMonitoring.value).toBe(true)
        })

        it('should allow viewing and editing any school ID', () => {
            const { canViewSchool, canEditSchool, canViewSchoolUsers } = usePermissions()

            expect(canViewSchool('school-123')).toBe(true)
            expect(canEditSchool('school-999')).toBe(true)
            expect(canViewSchoolUsers('school-456')).toBe(true)
        })

        it('should enable UI elements for SuperAdmin', () => {
            const { showSchoolFilter, showCreateSchoolButton, showMonitoringMenu } = usePermissions()

            expect(showSchoolFilter.value).toBe(true)
            expect(showCreateSchoolButton.value).toBe(true)
            expect(showMonitoringMenu.value).toBe(true)
        })
    })

    describe('Admin Role (School Level)', () => {
        beforeEach(() => {
            authStore.user = {
                id: 'admin-1',
                email: 'admin@school1.it',
                role: 'admin',
                school_id: 'school-100'
            }
        })

        it('should identify admin with single school scope', () => {
            const {
                isSuperAdmin,
                isAdmin,
                isAdminOrSuperAdmin,
                userSchoolId,
                canManageAllSchools,
                canAccessMonitoring
            } = usePermissions()

            expect(isSuperAdmin.value).toBe(false)
            expect(isAdmin.value).toBe(true)
            expect(isAdminOrSuperAdmin.value).toBe(true)
            expect(userSchoolId.value).toBe('school-100')
            expect(canManageAllSchools.value).toBe(false)
            expect(canAccessMonitoring.value).toBe(false)
        })

        it('should restrict school viewing/editing to assigned school ID', () => {
            const { canViewSchool, canEditSchool } = usePermissions()

            expect(canViewSchool('school-100')).toBe(true) // Assigned school
            expect(canViewSchool('school-200')).toBe(false) // Different school
            expect(canEditSchool('school-100')).toBe(true)
            expect(canEditSchool('school-200')).toBe(false)
        })

        it('should disable superadmin-only UI elements', () => {
            const { showSchoolFilter, showCreateSchoolButton, showMonitoringMenu } = usePermissions()

            expect(showSchoolFilter.value).toBe(false)
            expect(showCreateSchoolButton.value).toBe(false)
            expect(showMonitoringMenu.value).toBe(false)
        })
    })

    describe('Teacher / Student / Parent Role', () => {
        beforeEach(() => {
            authStore.user = {
                id: 'teacher-1',
                email: 'teacher@school1.it',
                role: 'teacher',
                school_id: 'school-100'
            }
        })

        it('should deny admin privileges for standard user roles', () => {
            const {
                isSuperAdmin,
                isAdmin,
                isAdminOrSuperAdmin,
                canViewSchool,
                canViewSchoolAnalytics
            } = usePermissions()

            expect(isSuperAdmin.value).toBe(false)
            expect(isAdmin.value).toBe(false)
            expect(isAdminOrSuperAdmin.value).toBe(false)
            expect(canViewSchool('school-100')).toBe(false)
            expect(canViewSchoolAnalytics.value).toBe(false)
        })
    })
})
