import { createRouter, createWebHistory } from 'vue-router'
import routes from './routes'
import { authGuard } from './guards'

const router = createRouter({
    history: createWebHistory(),
    routes
})

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

export default router
