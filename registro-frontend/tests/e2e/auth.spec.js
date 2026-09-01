import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import Login from '@/pages/Login.vue'
import { useAuthStore } from '@/stores/auth'
import { authGuard } from '@/router/guards'

function createFakeJwt(role, expiresInSec = 3600) {
    const header = btoa(JSON.stringify({ alg: 'HS256', typ: 'JWT' }))
    const payload = btoa(JSON.stringify({
        sub: 'user-1',
        role: role,
        user_role: role,
        exp: Math.floor(Date.now() / 1000) + expiresInSec
    }))
    return `${header}.${payload}.signature`
}

describe('Auth Workflow E2E', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        sessionStorage.clear()
        localStorage.clear()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: null,
                    token: null,
                    isAuthenticated: false,
                    isInitializing: false
                }
            }
        })
    })

    it('E01 — renders login form elements properly', () => {
        const wrapper = mount(Login, {
            global: {
                plugins: [pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-btn-dropdown': { template: '<div><slot /></div>' },
                    'q-list': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form @submit.prevent><slot /></form>' },
                    'q-input': {
                        props: ['modelValue', 'label', 'type'],
                        template: '<input :type="type || \'text\'" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />'
                    },
                    'q-btn': {
                        props: ['label', 'type', 'loading'],
                        template: '<button :type="type || \'button\'">{{ label }}<slot /></button>'
                    },
                    'q-icon': true,
                    'q-banner': { template: '<div class="q-banner"><slot /></div>' },
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-checkbox': true
                }
            }
        })

        expect(wrapper.find('form').exists()).toBe(true)
        expect(wrapper.findAll('input').length).toBeGreaterThanOrEqual(2)
    })

    it('E02 — authGuard redirects unauthenticated users to /login?reason=session_expired', async () => {
        const authStore = useAuthStore()
        authStore.isAuthenticated = false
        authStore.token = null

        const next = vi.fn()
        const to = { path: '/teacher/grades', meta: { requiresAuth: true } }
        const from = { path: '/' }

        await authGuard(to, from, next)

        expect(next).toHaveBeenCalledWith({
            path: '/login',
            query: { reason: 'session_expired' }
        })
    })

    it('E03 — authGuard redirects authenticated teacher to /teacher on public auth routes', async () => {
        const authStore = useAuthStore()
        authStore.token = createFakeJwt('teacher')
        authStore.user = { role: 'teacher', email: 'teacher@school.it' }

        const next = vi.fn()
        const to = { path: '/login', meta: { requiresAuth: false } }
        const from = { path: '/' }

        await authGuard(to, from, next)

        expect(next).toHaveBeenCalledWith('/teacher')
    })

    it('E04 — authGuard redirects authenticated admin to /admin/dashboard on login route', async () => {
        const authStore = useAuthStore()
        authStore.token = createFakeJwt('admin')
        authStore.user = { role: 'admin', email: 'admin@school.it' }

        const next = vi.fn()
        const to = { path: '/login', meta: { requiresAuth: false } }
        const from = { path: '/' }

        await authGuard(to, from, next)

        expect(next).toHaveBeenCalledWith('/admin/dashboard')
    })

    it('E05 — authGuard allows navigation for authorized roles on protected routes', async () => {
        const authStore = useAuthStore()
        authStore.token = createFakeJwt('teacher')
        authStore.user = { role: 'teacher', email: 'teacher@school.it' }

        const next = vi.fn()
        const to = {
            path: '/teacher/grades',
            meta: { requiresAuth: true, roles: ['teacher', 'coordinator'] }
        }
        const from = { path: '/teacher' }

        await authGuard(to, from, next)

        expect(next).toHaveBeenCalledWith()
    })
})
