import { createRouter, createWebHistory } from 'vue-router'
import routes from './routes'
import { authGuard } from './guards'

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach(authGuard)

export default router
