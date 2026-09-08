import { useAuthStore } from 'src/stores/auth'
import { isTokenExpired } from 'src/utils/jwt'
import { Notify } from 'quasar'

// Shared in-flight promise to prevent concurrent initAuth calls
let _initAuthPromise = null

export const authGuard = async (to, from, ...rest) => {
    const next = typeof rest[0] === 'function' ? rest[0] : null
    const proceed = (target) => {
        if (next) {
            if (target !== undefined) {
                next(target)
            } else {
                next()
            }
        }
        return target
    }

    const authStore = useAuthStore()

    // If initAuth is currently running or token is missing with stored user session,
    // wait for the shared promise to avoid concurrent executions.
    if (authStore.isInitializing || (!authStore.token && (localStorage.getItem('user') || sessionStorage.getItem('user')))) {
        if (!_initAuthPromise) {
            _initAuthPromise = authStore.initAuth().finally(() => { _initAuthPromise = null })
        }
        await _initAuthPromise
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
            return proceed(getUserDashboard(currentRole))
        }
        return proceed()
    }

    // Client-side JWT expiry check (second layer — backend is the primary authority)
    if (authStore.token && isTokenExpired(authStore.token)) {
        authStore.logout?.()
        return proceed({ path: '/login', query: { reason: 'session_expired' } })
    }

    // If not authenticated or no valid role present from JWT, redirect to login
    if (!authStore.isAuthenticated || !currentRole) {
        return proceed({ path: '/login', query: { reason: 'session_expired' } })
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
            try {
                if (typeof Notify !== 'undefined' && typeof Notify.create === 'function') {
                    Notify.create({
                        type: 'warning',
                        message: 'Accesso negato: non disponi dei permessi necessari per questa sezione.',
                        icon: 'lock',
                        position: 'top',
                        timeout: 3000
                    })
                }
            } catch { /* ignore notification failure in test/headless */ }
            return proceed(getUserDashboard(currentRole))
        }
    }

    return proceed()
}
