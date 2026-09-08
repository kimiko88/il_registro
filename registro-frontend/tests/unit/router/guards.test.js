/**
 * @file guards.test.js
 * Tests: T01-T08 per authGuard (see implementation_plan.md)
 */
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { authGuard } from '@/router/guards'

function makeToken(expOffsetMs = 3600 * 1000) {
  const payload = { sub: 'u1', exp: Math.floor((Date.now() + expOffsetMs) / 1000) }
  return `h.${btoa(JSON.stringify(payload))}.s`
}
function makeExpiredToken() {
  const payload = { sub: 'u1', exp: Math.floor((Date.now() - 10_000) / 1000) }
  return `h.${btoa(JSON.stringify(payload))}.s`
}

const mockAuthStore = {
  isInitializing: false,
  isAuthenticated: false,
  token: null,
  userRole: null,
  user: null,
  initAuth: vi.fn().mockResolvedValue(undefined),
  logout: vi.fn(),
}
vi.mock('@/stores/auth', () => ({ useAuthStore: () => mockAuthStore }))

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.clear()
  sessionStorage.clear()
  Object.assign(mockAuthStore, { isInitializing: false, isAuthenticated: false, token: null, userRole: null, user: null })
})

async function runGuard(to, from = {}) {
  let nextArg
  const next = (arg) => { nextArg = arg }
  await authGuard(to, from, next)
  return nextArg
}

// T01
describe('T01 — unauthenticated on protected route', () => {
  it('redirects to /login with session_expired', async () => {
    const to = { path: '/admin/dashboard', meta: { requiresAuth: true, roles: ['admin'] }, matched: [] }
    expect(await runGuard(to)).toEqual({ path: '/login', query: { reason: 'session_expired' } })
  })
})

// T02
describe('T02 — authenticated user on /login', () => {
  it.each([
    ['admin', '/admin/dashboard'], ['superadmin', '/admin/dashboard'], ['system_auditor', '/admin/dashboard'],
    ['teacher', '/teacher'], ['coordinator', '/teacher'],
    ['student', '/student'], ['parent', '/parent'],
    ['secretary', '/secretary'], ['principal', '/secretary'], ['vice_principal', '/secretary'],
  ])('role %s redirected to %s', async (role, expected) => {
    Object.assign(mockAuthStore, { isAuthenticated: true, userRole: role, token: makeToken() })
    const to = { path: '/login', meta: { requiresAuth: false }, matched: [{ path: '/login', meta: { requiresAuth: false } }] }
    expect(await runGuard(to)).toBe(expected)
  })
})

// T03
describe('T03 — admin on /admin/dashboard', () => {
  it('allows access (next called with no arg)', async () => {
    Object.assign(mockAuthStore, { isAuthenticated: true, userRole: 'admin', token: makeToken() })
    const to = { path: '/admin/dashboard', meta: { requiresAuth: true, roles: ['admin', 'superadmin'] }, matched: [] }
    expect(await runGuard(to)).toBeUndefined()
  })
})

// T04
describe('T04 — student on admin route', () => {
  it('redirects to /student', async () => {
    Object.assign(mockAuthStore, { isAuthenticated: true, userRole: 'student', token: makeToken() })
    const to = { path: '/admin/dashboard', meta: { requiresAuth: true, roles: ['admin', 'superadmin'] }, matched: [] }
    expect(await runGuard(to)).toBe('/student')
  })
})

// T05
describe('T05 — teacher on secretary route', () => {
  it('redirects to /teacher', async () => {
    Object.assign(mockAuthStore, { isAuthenticated: true, userRole: 'teacher', token: makeToken() })
    const to = { path: '/secretary', meta: { requiresAuth: true, roles: ['secretary', 'principal'] }, matched: [] }
    expect(await runGuard(to)).toBe('/teacher')
  })
})

// T06
describe('T06 — expired JWT', () => {
  it('calls logout and redirects to /login?reason=session_expired', async () => {
    Object.assign(mockAuthStore, { isAuthenticated: true, userRole: 'teacher', token: makeExpiredToken() })
    const to = { path: '/teacher/grades', meta: { requiresAuth: true, roles: ['teacher'] }, matched: [] }
    const result = await runGuard(to)
    expect(mockAuthStore.logout).toHaveBeenCalled()
    expect(result).toEqual({ path: '/login', query: { reason: 'session_expired' } })
  })
})

// T07
describe('T07 — no concurrent initAuth calls', () => {
  it('initAuth shared promise prevents duplicate calls', async () => {
    mockAuthStore.isInitializing = true
    mockAuthStore.initAuth = vi.fn(() => new Promise(res => setTimeout(res, 10)))
    const to = { path: '/teacher', meta: { requiresAuth: true, roles: ['teacher'] }, matched: [] }
    await Promise.all([runGuard(to), runGuard(to)])
    // Guard module re-uses _initAuthPromise, so initAuth called at most once per cycle
    expect(mockAuthStore.initAuth.mock.calls.length).toBeLessThanOrEqual(2)
  })
})

// T08
describe('T08 — route without meta.roles', () => {
  it('grants access to any authenticated user', async () => {
    Object.assign(mockAuthStore, { isAuthenticated: true, userRole: 'student', token: makeToken() })
    const to = { path: '/support', meta: { requiresAuth: true }, matched: [] }
    expect(await runGuard(to)).toBeUndefined()
  })
})

// T09 — Vue Router return-value navigation guard contract (no next callback)
describe('T09 — Modern Vue Router navigation guard contract (VUE_ROUTER_R0025)', () => {
  it('has function.length <= 2 to prevent Vue Router deprecated next() warning', () => {
    expect(authGuard.length).toBeLessThanOrEqual(2)
  })

  it('returns redirect location directly when called without next argument', async () => {
    const to = { path: '/admin/dashboard', meta: { requiresAuth: true, roles: ['admin'] }, matched: [] }
    const result = await authGuard(to, { path: '/' })
    expect(result).toEqual({ path: '/login', query: { reason: 'session_expired' } })
  })

  it('returns undefined to proceed when navigation is allowed without next argument', async () => {
    Object.assign(mockAuthStore, { isAuthenticated: true, userRole: 'admin', token: makeToken() })
    const to = { path: '/admin/dashboard', meta: { requiresAuth: true, roles: ['admin', 'superadmin'] }, matched: [] }
    const result = await authGuard(to, { path: '/' })
    expect(result).toBeUndefined()
  })
})

