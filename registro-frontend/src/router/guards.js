import { useAuthStore } from 'src/stores/auth'

export const authGuard = (to, from, next) => {
    const authStore = useAuthStore()

    // Define public routes
    const publicRoutes = ['/login', '/register', '/forgot-password']

    // Check if route is public
    if (publicRoutes.includes(to.path)) {
        // If logged in and trying to access login, redirect to dashboard? 
        // Usually yes, but user prompt "If not logged in, always redirect to login".
        // It implies protection of protected routes.
        next()
        return
    }

    // If not authenticated, redirect to login
    if (!authStore.token) {
        next('/login')
    } else {
        // Check role permissions if needed (optional enhancement)
        // if (to.meta.role && to.meta.role !== authStore.role) ...
        next()
    }
}
