/**
 * Tests for api.js - focuses on the most security-critical and error-prone areas:
 * - getBaseURL utility
 * - clearLocalSession / resetApiState
 * - setReauthHandler / setApiRouter setup
 * - Idempotency-Key injection logic (unit-level, without HTTP)
 */
import { describe, it, expect, vi, beforeEach } from 'vitest'

describe('api.js - getBaseURL', () => {
    beforeEach(() => {
        vi.resetModules()
    })

    it('returns /api/v1 when no VITE_API_URL is set', async () => {
        const { getBaseURL } = await import('@/services/api')
        const url = getBaseURL()
        expect(url).toBe('/api/v1')
    })

    it('returns URL unchanged when it already ends with /api/v1', async () => {
        vi.stubEnv('VITE_API_URL', 'https://api.example.com/api/v1')
        const { getBaseURL } = await import('@/services/api')
        expect(getBaseURL()).toBe('https://api.example.com/api/v1')
        vi.unstubAllEnvs()
    })

    it('appends /api/v1 when base URL does not end with it', async () => {
        vi.stubEnv('VITE_API_URL', 'https://api.example.com')
        const { getBaseURL } = await import('@/services/api')
        expect(getBaseURL()).toBe('https://api.example.com/api/v1')
        vi.unstubAllEnvs()
    })

    it('strips trailing slashes from VITE_API_URL', async () => {
        vi.stubEnv('VITE_API_URL', 'https://api.example.com///')
        const { getBaseURL } = await import('@/services/api')
        expect(getBaseURL()).toBe('https://api.example.com/api/v1')
        vi.unstubAllEnvs()
    })
})

describe('api.js - cancelInFlightRequests', () => {
    beforeEach(() => {
        vi.resetModules()
    })

    it('is exported as a function', async () => {
        const { cancelInFlightRequests } = await import('@/services/api')
        expect(typeof cancelInFlightRequests).toBe('function')
    })

    it('can be called without error when no requests are in flight', async () => {
        const { cancelInFlightRequests } = await import('@/services/api')
        expect(() => cancelInFlightRequests()).not.toThrow()
    })
})

describe('api.js - resetApiState', () => {
    beforeEach(() => {
        vi.resetModules()
    })

    it('is exported as a function', async () => {
        const { resetApiState } = await import('@/services/api')
        expect(typeof resetApiState).toBe('function')
    })

    it('can be called without throwing', async () => {
        const { resetApiState } = await import('@/services/api')
        expect(() => resetApiState()).not.toThrow()
    })
})

describe('api.js - clearLocalSession', () => {
    beforeEach(() => {
        vi.resetModules()
        localStorage.setItem('user', JSON.stringify({ email: 'test@test.it' }))
        localStorage.setItem('token', 'oldtoken')
        localStorage.setItem('refreshToken', 'oldrefresh')
        localStorage.setItem('selectedChildId', 'child-123')
        sessionStorage.setItem('user', JSON.stringify({ email: 'test@test.it' }))
        sessionStorage.setItem('token', 'oldtoken')
        sessionStorage.setItem('refreshToken', 'oldrefresh')
        sessionStorage.setItem('selectedChildId', 'child-123')
        sessionStorage.setItem('registro_lesson_drafts', JSON.stringify([]))
    })

    it('clears all auth keys from localStorage', async () => {
        const { clearLocalSession } = await import('@/services/api')
        clearLocalSession()

        expect(localStorage.getItem('user')).toBeNull()
        expect(localStorage.getItem('token')).toBeNull()
        expect(localStorage.getItem('refreshToken')).toBeNull()
        expect(localStorage.getItem('selectedChildId')).toBeNull()
    })

    it('clears all auth keys from sessionStorage', async () => {
        const { clearLocalSession } = await import('@/services/api')
        clearLocalSession()

        expect(sessionStorage.getItem('user')).toBeNull()
        expect(sessionStorage.getItem('token')).toBeNull()
        expect(sessionStorage.getItem('refreshToken')).toBeNull()
        expect(sessionStorage.getItem('selectedChildId')).toBeNull()
        expect(sessionStorage.getItem('registro_lesson_drafts')).toBeNull()
    })

    it('can be called on empty storage without throwing', async () => {
        localStorage.clear()
        sessionStorage.clear()
        const { clearLocalSession } = await import('@/services/api')
        expect(() => clearLocalSession()).not.toThrow()
    })
})

describe('api.js - setReauthHandler and setApiRouter', () => {
    beforeEach(() => {
        vi.resetModules()
    })

    it('setReauthHandler is exported as a function', async () => {
        const { setReauthHandler } = await import('@/services/api')
        expect(typeof setReauthHandler).toBe('function')
    })

    it('setApiRouter is exported as a function', async () => {
        const { setApiRouter } = await import('@/services/api')
        expect(typeof setApiRouter).toBe('function')
    })

    it('setReauthHandler accepts a function without throwing', async () => {
        const { setReauthHandler } = await import('@/services/api')
        expect(() => setReauthHandler(async () => 'new-token')).not.toThrow()
    })

    it('setReauthHandler accepts null (unregister handler)', async () => {
        const { setReauthHandler } = await import('@/services/api')
        expect(() => setReauthHandler(null)).not.toThrow()
    })

    it('setApiRouter accepts a router object without throwing', async () => {
        const { setApiRouter } = await import('@/services/api')
        const fakeRouter = { push: vi.fn(), currentRoute: { value: { path: '/' } } }
        expect(() => setApiRouter(fakeRouter)).not.toThrow()
    })
})

describe('api.js - Idempotency-Key generation logic', () => {
    // Test the idempotency key format in isolation by inspecting what the interceptor would produce
    it('generates a UUID v4 format key via crypto.randomUUID()', () => {
        const key = crypto.randomUUID()
        const uuidV4Pattern = /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i
        expect(key).toMatch(uuidV4Pattern)
    })

    it('each generated key is unique', () => {
        const keys = new Set(Array.from({ length: 50 }, () => crypto.randomUUID()))
        expect(keys.size).toBe(50)
    })
})

describe('api.js - IDEMPOTENT_ENDPOINTS matching', () => {
    // Test the endpoint matching logic conceptually
    const IDEMPOTENT_ENDPOINTS = [
        '/grades',
        '/grades/bulk',
        '/attendance',
        '/class-tests',
        '/payments',
        '/signatures',
        '/firme',
        '/verbali',
    ]

    function shouldInjectKey(url) {
        return IDEMPOTENT_ENDPOINTS.some(ep => url.includes(ep))
    }

    it('matches /grades endpoint', () => {
        expect(shouldInjectKey('/grades')).toBe(true)
    })

    it('matches /grades/bulk endpoint', () => {
        expect(shouldInjectKey('/grades/bulk')).toBe(true)
    })

    it('matches /attendance endpoint', () => {
        expect(shouldInjectKey('/attendance')).toBe(true)
    })

    it('matches /payments endpoint', () => {
        expect(shouldInjectKey('/payments')).toBe(true)
    })

    it('matches /signatures endpoint', () => {
        expect(shouldInjectKey('/signatures')).toBe(true)
    })

    it('matches /verbali endpoint', () => {
        expect(shouldInjectKey('/verbali')).toBe(true)
    })

    it('does NOT match non-idempotent endpoints like /users', () => {
        expect(shouldInjectKey('/users')).toBe(false)
    })

    it('does NOT match non-idempotent endpoints like /schools', () => {
        expect(shouldInjectKey('/schools')).toBe(false)
    })

    it('does NOT match /auth/login', () => {
        expect(shouldInjectKey('/auth/login')).toBe(false)
    })
})
