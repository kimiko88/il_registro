import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuth } from '@/composables/useAuth'
import authService from '@/services/authService'

const mockPush = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush
  })
}))

vi.mock('@/services/authService', () => ({
  default: {
    login: vi.fn(),
    logout: vi.fn()
  }
}))

describe('useAuth Role Redirection & Routing Suite', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  const testCases = [
    { role: 'superadmin', expectedRoute: '/admin' },
    { role: 'admin', expectedRoute: '/admin' },
    { role: 'system_auditor', expectedRoute: '/admin' },
    { role: 'secretary', expectedRoute: '/secretary' },
    { role: 'principal', expectedRoute: '/secretary' },
    { role: 'vice_principal', expectedRoute: '/secretary' },
    { role: 'staff', expectedRoute: '/secretary' },
    { role: 'teacher', expectedRoute: '/teacher' },
    { role: 'coordinator', expectedRoute: '/teacher' },
    { role: 'docente', expectedRoute: '/teacher' },
    { role: 'student', expectedRoute: '/student' },
    { role: 'parent', expectedRoute: '/parent' },
    { role: 'unknown_role', expectedRoute: '/' }
  ]

  testCases.forEach(({ role, expectedRoute }) => {
    it(`routes role "${role}" to "${expectedRoute}" upon successful login`, async () => {
      authService.login.mockResolvedValueOnce({
        user: { id: 'u1', email: 'test@school.it', role },
        access_token: 'fake-jwt-token',
        refresh_token: 'fake-refresh-token'
      })

      const { login } = useAuth()
      const error = await login('test@school.it', 'Secret123!')

      expect(error).toBeNull()
      expect(mockPush).toHaveBeenCalledWith(expectedRoute)
    })
  })

  it('handles and translates login credentials errors gracefully', async () => {
    authService.login.mockRejectedValueOnce({
      response: {
        data: { error: 'Invalid credentials' }
      }
    })

    const { login } = useAuth()
    const error = await login('wrong@school.it', 'WrongPass')

    expect(typeof error).toBe('string')
    expect(error.length).toBeGreaterThan(0)
  })
})
