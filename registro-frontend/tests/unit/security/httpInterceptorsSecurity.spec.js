import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'

describe('HTTP Interceptors & Token Security', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
        sessionStorage.clear()
    })

    it('injects Bearer token into outgoing request headers', async () => {
        const authStore = useAuthStore()
        authStore.token = 'my-secret-jwt-token'

        const interceptor = api.interceptors.request.handlers[0]
        const config = { headers: {} }

        const result = await interceptor.fulfilled(config)

        expect(result.headers.Authorization).toBe('Bearer my-secret-jwt-token')
    })

    it('injects saved Accept-Language header into outgoing requests', async () => {
        localStorage.setItem('superadmin_language', 'en-US')

        const interceptor = api.interceptors.request.handlers[0]
        const config = { headers: {} }

        const result = await interceptor.fulfilled(config)

        expect(result.headers['Accept-Language']).toBe('en-US')
    })

    it('attaches user-friendly error message for 403 Forbidden response', async () => {
        const responseInterceptor = api.interceptors.response.handlers[0]
        const error = {
            response: { status: 403, data: { error: 'forbidden' } },
            config: { url: '/admin/users' }
        }

        await expect(responseInterceptor.rejected(error)).rejects.toHaveProperty('userMessage')
    })

    it('attaches user-friendly error message for 500 Server Error response', async () => {
        const responseInterceptor = api.interceptors.response.handlers[0]
        const error = {
            response: { status: 500, data: { error: 'internal error' } },
            config: { url: '/reports' }
        }

        await expect(responseInterceptor.rejected(error)).rejects.toHaveProperty('userMessage')
    })
})
