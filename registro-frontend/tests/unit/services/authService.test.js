/**
 * @file authService.test.js
 * Unit tests for authService: trimming, login, logout, refresh, forgot/reset password, dual export
 */
import { describe, it, expect, vi, beforeEach } from 'vitest'
import authServiceDefault, { authService } from '@/services/authService'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    post: vi.fn(),
    get: vi.fn()
  }
}))

describe('authService — Trimming and Endpoints', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('exports both default and named export referencing the same object', () => {
    expect(authService).toBeDefined()
    expect(authServiceDefault).toBeDefined()
    expect(authService.login).toBe(authServiceDefault.login)
  })

  it('trims email when logging in', async () => {
    api.post.mockResolvedValue({ data: { user: { id: 1 }, access_token: 'tok' } })

    const res = await authService.login('   teacher@school.it   ', 'Password123!')

    expect(api.post).toHaveBeenCalledWith('/auth/login', {
      email: 'teacher@school.it',
      password: 'Password123!'
    })
    expect(res.access_token).toBe('tok')
  })

  it('trims email when registering', async () => {
    api.post.mockResolvedValue({ data: { id: 'u1', email: 'student@school.it' } })

    const res = await authService.register({ email: '  student@school.it  ', first_name: 'Mario' })

    expect(api.post).toHaveBeenCalledWith('/auth/register', {
      email: 'student@school.it',
      first_name: 'Mario'
    })
    expect(res.id).toBe('u1')
  })

  it('logout sends refresh_token if provided', async () => {
    api.post.mockResolvedValue({ data: { message: 'Logged out' } })

    const res = await authService.logout('ref-tok-123')

    expect(api.post).toHaveBeenCalledWith('/auth/logout', { refresh_token: 'ref-tok-123' })
    expect(res.message).toBe('Logged out')
  })

  it('logout sends empty object when refresh_token is not provided', async () => {
    api.post.mockResolvedValue({ data: { message: 'Logged out' } })

    await authService.logout()

    expect(api.post).toHaveBeenCalledWith('/auth/logout', {})
  })

  it('forgotPassword trims email and calls /auth/forgot-password', async () => {
    api.post.mockResolvedValue({ data: { message: 'Email sent' } })

    const res = await authService.forgotPassword('   user@school.it   ')

    expect(api.post).toHaveBeenCalledWith('/auth/forgot-password', { email: 'user@school.it' })
    expect(res.message).toBe('Email sent')
  })

  it('resetPassword calls /auth/reset-password with token and new_password', async () => {
    api.post.mockResolvedValue({ data: { message: 'Password reset' } })

    const res = await authService.resetPassword('reset-tok-xyz', 'NewPass123!')

    expect(api.post).toHaveBeenCalledWith('/auth/reset-password', {
      token: 'reset-tok-xyz',
      new_password: 'NewPass123!'
    })
    expect(res.message).toBe('Password reset')
  })

  it('getCurrentUser calls /auth/me', async () => {
    api.get.mockResolvedValue({ data: { id: 'u1', email: 'test@school.it' } })

    const res = await authService.getCurrentUser()

    expect(api.get).toHaveBeenCalledWith('/auth/me')
    expect(res.id).toBe('u1')
  })

  it('changePassword calls /auth/change-password with payload', async () => {
    api.post.mockResolvedValue({ data: { message: 'Password changed' } })

    const res = await authService.changePassword('OldPass!', 'NewPass!')

    expect(api.post).toHaveBeenCalledWith('/auth/change-password', {
      current_password: 'OldPass!',
      new_password: 'NewPass!'
    })
    expect(res.message).toBe('Password changed')
  })
})
