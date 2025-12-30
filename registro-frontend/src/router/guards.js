import { useAuthStore } from '@/stores/auth'

export function authGuard(to, from, next) {
    const authStore = useAuthStore()
    const isAuthenticated = authStore.isAuthenticated
    const userRole = authStore.user?.role

    if (to.path === '/login' && isAuthenticated) {
        return next('/')
    }

    if (to.meta.role && !isAuthenticated) {
        return next('/login')
    }

    if (to.meta.role && to.meta.role !== userRole && userRole !== 'superadmin') {
        // Basic role check - superadmin can access everything
        return next('/')
    }

    next()
}
