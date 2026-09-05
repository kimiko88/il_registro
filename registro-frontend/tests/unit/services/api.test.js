import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import api, {
  getBaseURL,
  setApiRouter,
  resetApiState,
  clearLocalSession,
  setReauthHandler,
  setApiI18n,
} from '@/services/api'

// Mock auth store
const mockAuthStore = {
  token: 'mock-jwt-token',
  user: { language: 'it-IT', email: 'test@school.it' },
  updateTokens: vi.fn(),
  logout: vi.fn(),
}

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => mockAuthStore,
}))


describe('api.js — getBaseURL', () => {
  const originalEnv = import.meta.env.VITE_API_URL

  afterEach(() => {
    if (originalEnv === undefined) {
      delete import.meta.env.VITE_API_URL
    } else {
      import.meta.env.VITE_API_URL = originalEnv
    }
  })

  it('returns default /api/v1 when VITE_API_URL is undefined or empty', () => {
    delete import.meta.env.VITE_API_URL
    expect(getBaseURL()).toBe('/api/v1')
  })

  it('handles URL without trailing slash', () => {
    import.meta.env.VITE_API_URL = 'http://localhost:3000'
    expect(getBaseURL()).toBe('http://localhost:3000/api/v1')
  })

  it('handles URL with trailing slash without creating double slash', () => {
    import.meta.env.VITE_API_URL = 'http://localhost:3000/'
    expect(getBaseURL()).toBe('http://localhost:3000/api/v1')
  })

  it('preserves URL that already ends with /api/v1', () => {
    import.meta.env.VITE_API_URL = 'https://api.school.it/api/v1'
    expect(getBaseURL()).toBe('https://api.school.it/api/v1')
  })

  it('removes trailing slash from URL that ends with /api/v1/', () => {
    import.meta.env.VITE_API_URL = 'https://api.school.it/api/v1/'
    expect(getBaseURL()).toBe('https://api.school.it/api/v1')
  })
})

describe('api.js — Session and State Management', () => {
  beforeEach(() => {
    localStorage.clear()
    sessionStorage.clear()
    vi.clearAllMocks()
  })

  it('clearLocalSession removes all auth and session keys', () => {
    localStorage.setItem('user', '{"id":1}')
    localStorage.setItem('token', 'abc')
    localStorage.setItem('refreshToken', 'def')
    localStorage.setItem('selectedChildId', 'child1')
    sessionStorage.setItem('token', 'abc')
    sessionStorage.setItem('registro_lesson_drafts', '[]')

    clearLocalSession()

    expect(localStorage.getItem('user')).toBeNull()
    expect(localStorage.getItem('token')).toBeNull()
    expect(localStorage.getItem('refreshToken')).toBeNull()
    expect(localStorage.getItem('selectedChildId')).toBeNull()
    expect(sessionStorage.getItem('token')).toBeNull()
    expect(sessionStorage.getItem('registro_lesson_drafts')).toBeNull()
  })

  it('resetApiState cleans up in-flight refresh state without crashing', () => {
    expect(() => resetApiState()).not.toThrow()
  })
})

describe('api.js — Request Interceptor', () => {
  it('adds Authorization Bearer header when authStore has token', () => {
    const handler = api.interceptors.request.handlers[0].fulfilled
    const config = { headers: {} }
    const result = handler(config)
    expect(result.headers.Authorization).toBe('Bearer mock-jwt-token')
  })

  it('sets Accept-Language header from user language or fallback', () => {
    const handler = api.interceptors.request.handlers[0].fulfilled
    const config = { headers: {} }
    const result = handler(config)
    expect(result.headers['Accept-Language']).toBe('it-IT')
  })

  it('strips redundant /api/v1/ prefix from request url', () => {
    const handler = api.interceptors.request.handlers[0].fulfilled
    const config1 = { headers: {}, url: '/api/v1/users' }
    const result1 = handler(config1)
    expect(result1.url).toBe('/users')

    const config2 = { headers: {}, url: 'api/v1/classes' }
    const result2 = handler(config2)
    expect(result2.url).toBe('/classes')
  })
})

describe('api.js — Response Interceptor Error Messages', () => {
  const errorHandler = api.interceptors.response.handlers[0].rejected

  it('handles network error (no response) with userMessage', async () => {
    const error = {}
    await expect(errorHandler(error)).rejects.toBe(error)
    expect(error.userMessage).toContain('connessione')
  })

  it('handles 403 Forbidden with userMessage', async () => {
    const error = {
      response: { status: 403, data: {} },
      config: { url: '/admin/users' }
    }
    await expect(errorHandler(error)).rejects.toBe(error)
    expect(error.userMessage).toContain('permessi')
  })

  it('handles 429 Rate Limit with server message if present', async () => {
    const error = {
      response: { status: 429, data: { error: 'Troppe richieste!' } },
      config: { url: '/auth/login' }
    }
    await expect(errorHandler(error)).rejects.toBe(error)
    expect(error.userMessage).toBe('Troppe richieste!')
  })

  it('handles 500 Internal Server Error with userMessage', async () => {
    const error = {
      response: { status: 500, data: {} },
      config: { url: '/classes' }
    }
    await expect(errorHandler(error)).rejects.toBe(error)
    expect(error.userMessage).toContain('server')
  })
})

describe('api.js — Router, I18n and Reauth Handlers Setup', () => {
  it('setApiRouter registers router without errors', () => {
    const mockRouter = { push: vi.fn(), currentRoute: { value: { path: '/home' } } }
    expect(() => setApiRouter(mockRouter)).not.toThrow()
  })

  it('setReauthHandler registers reauth callback without errors', () => {
    const mockHandler = vi.fn().mockResolvedValue('new-token')
    expect(() => setReauthHandler(mockHandler)).not.toThrow()
  })

  it('setApiI18n registers i18n instance used for error formatting', () => {
    const mockI18n = {
      global: {
        t: vi.fn((key) => `translated:${key}`)
      }
    }
    expect(() => setApiI18n(mockI18n)).not.toThrow()
  })
})

