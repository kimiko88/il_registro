import { createRouter, createWebHistory } from 'vue-router'
import routes from './routes'
import { authGuard } from './guards'
import { setApiRouter } from '@/services/api'

const router = createRouter({
    history: createWebHistory(),
    routes
})

if (typeof setApiRouter === 'function') {
    setApiRouter(router)
}

router.beforeEach(authGuard)

router.afterEach((to) => {
    const base = 'Registro Elettronico'
    const pageTitle = to.meta?.title
    if (pageTitle) {
        document.title = `${pageTitle} — ${base}`
        return
    }

    let section = ''
    if (to.path.startsWith('/admin')) {
        section = 'Admin'
    } else if (to.path.startsWith('/teacher')) {
        section = 'Docente'
    } else if (to.path.startsWith('/student')) {
        section = 'Studente'
    } else if (to.path.startsWith('/parent')) {
        section = 'Famiglie'
    } else if (to.path.startsWith('/secretary')) {
        section = 'Segreteria'
    }
    
    document.title = section ? `${section} — ${base}` : base
})

router.onError((error, to) => {
    if (error.message && /loading chunk|failed to fetch dynamically imported module/i.test(error.message)) {
        console.error('Lazy-load chunk failure detected:', error)
        if (to?.fullPath) {
            window.location.assign(to.fullPath)
        } else {
            window.location.reload()
        }
    }
})

export default router
