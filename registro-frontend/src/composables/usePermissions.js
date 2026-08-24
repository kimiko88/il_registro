import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

/**
 * Composable for checking user permissions based on role
 * Handles SuperAdmin (global access) vs Admin (single school access)
 */
export function usePermissions() {
    const authStore = useAuthStore()

    // Role checks
    const isSuperAdmin = computed(() => authStore.userRole === 'superadmin')
    const isAdmin = computed(() => authStore.userRole === 'admin')
    const isAdminOrSuperAdmin = computed(() => isSuperAdmin.value || isAdmin.value)

    // User's assigned school ID
    const userSchoolId = computed(() => authStore.user?.school_id || null)

    // Global permissions (SuperAdmin only)
    const canManageAllSchools = computed(() => isSuperAdmin.value)
    const canCreateSchools = computed(() => isSuperAdmin.value)
    const canDeleteSchools = computed(() => isSuperAdmin.value)
    const canManageAdminUsers = computed(() => isSuperAdmin.value)
    const canAccessMonitoring = computed(() => isSuperAdmin.value)
    const canAccessGlobalAnalytics = computed(() => isSuperAdmin.value)
    const canAccessSystemSettings = computed(() => isSuperAdmin.value)

    // School-level permissions (evaluate against reactive computed getters)
    const canViewSchool = (schoolId) => {
        if (isSuperAdmin.value) return true
        if (isAdmin.value) return schoolId === userSchoolId.value
        return false
    }

    const canEditSchool = (schoolId) => {
        if (isSuperAdmin.value) return true
        if (isAdmin.value) return schoolId === userSchoolId.value
        return false
    }

    const canViewSchoolUsers = (schoolId) => {
        if (isSuperAdmin.value) return true
        if (isAdmin.value) return schoolId === userSchoolId.value
        return false
    }

    // Analytics permissions
    const canViewGlobalAnalytics = computed(() => isSuperAdmin.value)
    const canViewSchoolAnalytics = computed(() => isAdminOrSuperAdmin.value)

    // Feature visibility
    const showSchoolFilter = computed(() => isSuperAdmin.value)
    const showCreateSchoolButton = computed(() => isSuperAdmin.value)
    const showDeleteSchoolButton = computed(() => isSuperAdmin.value)
    const showAdminUsersMenu = computed(() => isSuperAdmin.value)
    const showMonitoringMenu = computed(() => isSuperAdmin.value)
    const showSystemSettingsMenu = computed(() => isSuperAdmin.value)

    return {
        // Role checks
        isSuperAdmin,
        isAdmin,
        isAdminOrSuperAdmin,
        userSchoolId,

        // Global permissions
        canManageAllSchools,
        canCreateSchools,
        canDeleteSchools,
        canManageAdminUsers,
        canAccessMonitoring,
        canAccessGlobalAnalytics,
        canAccessSystemSettings,

        // School-level permissions
        canViewSchool,
        canEditSchool,
        canViewSchoolUsers,

        // Analytics permissions
        canViewGlobalAnalytics,
        canViewSchoolAnalytics,

        // UI visibility
        showSchoolFilter,
        showCreateSchoolButton,
        showDeleteSchoolButton,
        showAdminUsersMenu,
        showMonitoringMenu,
        showSystemSettingsMenu
    }
}
