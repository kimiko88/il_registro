import { useAuthStore } from 'src/stores/auth'

export const authGuard = async (to, from, next) => {
    const authStore = useAuthStore()

    // If initAuth is currently running or token is missing with stored user session, wait for initAuth to resolve
    if (authStore.isInitializing) {
        await authStore.initAuth()
    } else if (!authStore.token && (localStorage.getItem('user') || sessionStorage.getItem('user'))) {
        await authStore.initAuth()
    }

    const publicRoutes = ['/login', '/register', '/forgot-password']

    const getUserDashboard = (role) => {
        if (role === 'admin' || role === 'superadmin' || role === 'system_auditor') return '/admin/dashboard'
        if (role === 'teacher' || role === 'coordinator') return '/teacher'
        if (role === 'student') return '/student'
        if (role === 'parent') return '/parent'
        if (role === 'secretary' || role === 'principal' || role === 'vice_principal') return '/secretary'
        return '/login'
    }

    const currentRole = authStore.userRole

    const isPublic = to.meta?.requiresAuth === false ||
        publicRoutes.includes(to.path) ||
        to.matched?.some(record => record.meta?.requiresAuth === false)

    // If route is public
    if (isPublic) {
        // Only redirect logged in users away if navigating to auth-entry routes (/login, /register, /forgot-password)
        const isAuthEntry = publicRoutes.includes(to.path) ||
            to.matched?.some(record => publicRoutes.includes(record.path))

        if (authStore.isAuthenticated && isAuthEntry) {
            next(getUserDashboard(currentRole))
            return
        }
        next()
        return
    }

    // If not authenticated or no valid role present from JWT, redirect to login
    if (!authStore.isAuthenticated || !currentRole) {
        next({ path: '/login', query: { reason: 'session_expired' } })
        return
    }

    // Role-based access control (deny-by-default for specified role lists)
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
