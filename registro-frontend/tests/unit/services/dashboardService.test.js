/**
 * @file dashboardService.test.js
 * Unit tests for dashboardService: role-based endpoint routing, principal/vice_principal/coordinator/auditor, dual export
 */
import { describe, it, expect, vi, beforeEach } from 'vitest'
import dashboardServiceDefault, { dashboardService } from '@/services/dashboardService'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn()
  }
}))

describe('dashboardService — Role-based Stats Routing', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    api.get.mockResolvedValue({ data: { count: 42 } })
  })

  it('exports both default and named export', () => {
    expect(dashboardService).toBeDefined()
    expect(dashboardServiceDefault).toBeDefined()
    expect(dashboardService.getDashboardStats).toBe(dashboardServiceDefault.getDashboardStats)
  })

  it.each([
    'admin',
    'superadmin',
    'secretary',
    'principal',
    'vice_principal',
    'system_auditor'
  ])('routes role %s to /admin/dashboard/stats', async (role) => {
    const res = await dashboardService.getDashboardStats(role)
    expect(api.get).toHaveBeenCalledWith('/admin/dashboard/stats')
    expect(res).toEqual({ count: 42 })
  })

  it.each([
    'teacher',
    'coordinator'
  ])('routes role %s to /teachers/dashboard/stats', async (role) => {
    const res = await dashboardService.getDashboardStats(role)
    expect(api.get).toHaveBeenCalledWith('/teachers/dashboard/stats')
    expect(res).toEqual({ count: 42 })
  })

  it('routes student role to /students/dashboard/stats', async () => {
    const res = await dashboardService.getDashboardStats('student')
    expect(api.get).toHaveBeenCalledWith('/students/dashboard/stats')
    expect(res).toEqual({ count: 42 })
  })

  it('routes parent role to /parents/dashboard/stats', async () => {
    const res = await dashboardService.getDashboardStats('parent')
    expect(api.get).toHaveBeenCalledWith('/parents/dashboard/stats')
    expect(res).toEqual({ count: 42 })
  })

  it('returns empty object for unknown role', async () => {
    const res = await dashboardService.getDashboardStats('unknown_role')
    expect(api.get).not.toHaveBeenCalled()
    expect(res).toEqual({})
  })
})
