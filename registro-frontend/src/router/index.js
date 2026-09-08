import { createRouter, createWebHistory } from 'vue-router'
import routes from './routes'
import { authGuard } from './guards'
import { setApiRouter } from '@/services/api'
import { i18n } from '@/i18n'

const router = createRouter({
    history: createWebHistory(),
    routes
})

if (typeof setApiRouter === 'function') {
    setApiRouter(router)
}

router.beforeEach((to, from) => authGuard(to, from))

router.afterEach((to) => {
    const base = 'Registro Elettronico'
    let resolvedTitle = ''

    if (to.meta?.titleKey && typeof i18n?.global?.t === 'function') {
        const translated = i18n.global.t(to.meta.titleKey)
        if (translated && translated !== to.meta.titleKey) {
            resolvedTitle = translated
        }
    }

    if (!resolvedTitle && to.meta?.title) {
        resolvedTitle = to.meta.title
    }

    if (resolvedTitle) {
        document.title = `${resolvedTitle} — ${base}`
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
        const retryKey = `chunk_retry_${to?.fullPath || 'unknown'}`
        if (!sessionStorage.getItem(retryKey)) {
            // First attempt: mark the retry and reload once
            sessionStorage.setItem(retryKey, '1')
            if (to?.fullPath) {
                window.location.assign(to.fullPath)
            } else {
                window.location.reload()
            }
        } else {
            // Already retried: clear the flag and do not loop
            sessionStorage.removeItem(retryKey)
            console.error('Chunk load definitively failed after retry. The user may need to clear cache.', error)
        }
    }
})

export default router
