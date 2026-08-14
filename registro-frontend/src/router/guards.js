import { useAuthStore } from 'src/stores/auth'

export const authGuard = (to, from, next) => {
    const authStore = useAuthStore()

    const publicRoutes = ['/login', '/register', '/forgot-password']

    const getUserDashboard = (role) => {
        if (role === 'admin' || role === 'superadmin') return '/admin/dashboard'
        if (role === 'system_auditor') return '/admin/audit-log'
        if (role === 'teacher' || role === 'coordinator') return '/teacher'
        if (role === 'student') return '/student'
        if (role === 'parent') return '/parent'
        if (role === 'secretary' || role === 'principal' || role === 'vice_principal') return '/secretary'
        return '/'
    }

    const currentRole = authStore.userRole

    const isPublic = to.meta?.requiresAuth === false || publicRoutes.includes(to.path)

    // If route is public
    if (isPublic) {
        // If already logged in, redirect to user's dashboard
        if (authStore.isAuthenticated) {
            next(getUserDashboard(currentRole))
            return
        }
        next()
        return
    }

    // If not authenticated, redirect to login
    if (!authStore.isAuthenticated) {
        next('/login')
        return
    }

    // Role-based access control
    if (to.meta) {
        const requiredRoles = Array.isArray(to.meta.roles)
            ? to.meta.roles
            : (to.meta.role ? [to.meta.role] : null)

        if (requiredRoles && requiredRoles.length > 0 && !requiredRoles.includes(currentRole)) {
            if (import.meta.env.DEV) {
                console.warn(`Access denied: role '${currentRole}' is not allowed for path '${to.path}'`)
            }
            next(getUserDashboard(currentRole))
            return
        }
    }

    next()
}
