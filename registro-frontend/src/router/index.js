import { createRouter, createWebHistory } from 'vue-router'
import routes from './routes'
import { authGuard } from './guards'

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach(authGuard)

router.afterEach((to) => {
    let title = 'Registro Elettronico'
    if (to.path.startsWith('/admin')) {
        title = 'Registro Elettronico - Admin'
    } else if (to.path.startsWith('/teacher')) {
        title = 'Registro Elettronico - Docente'
    } else if (to.path.startsWith('/student')) {
        title = 'Registro Elettronico - Studente'
    } else if (to.path.startsWith('/parent')) {
        title = 'Registro Elettronico - Famiglie'
    } else if (to.path.startsWith('/secretary')) {
        title = 'Registro Elettronico - Segreteria'
    }
    document.title = title
})

export default router
