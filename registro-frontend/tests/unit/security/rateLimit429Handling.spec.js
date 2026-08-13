import { describe, it, expect } from 'vitest'
import api from '@/services/api'

describe('Frontend Rate Limit (HTTP 429) Security Response Handling', () => {
    it('handles HTTP 429 response error and provides user-friendly retry message', async () => {
        const responseInterceptor = api.interceptors.response.handlers[0]
        const error = {
            response: {
                status: 429,
                data: {
                    code: 'AUTH_RATE_LIMIT_EXCEEDED',
                    error: 'Troppi tentativi di accesso. Riprova tra un minuto.'
                }
            },
            config: { url: '/api/v1/auth/login' }
        }

        try {
            await responseInterceptor.rejected(error)
        } catch (err) {
            expect(err).toHaveProperty('userMessage')
            expect(err.userMessage).toContain('Troppi tentativi')
        }
    })
})
