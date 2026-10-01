import { describe, it, expect, vi, beforeEach } from 'vitest'
import staffAttendanceService from '@/services/staffAttendanceService'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn()
  }
}))

describe('staffAttendanceServiceExtended.spec.js — Complete ATA & Staff Attendance Service Suite', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Attendance and Strike Management Endpoints', () => {
    it('getDailySummary calls GET /staff-attendance/summary with date param', async () => {
      const mockSummary = { date: '2026-10-15', total_staff: { total: 40, present: 35 } }
      api.get.mockResolvedValue({ data: mockSummary })

      const res = await staffAttendanceService.getDailySummary('2026-10-15')
      expect(api.get).toHaveBeenCalledWith('/staff-attendance/summary', { params: { date: '2026-10-15' } })
      expect(res).toEqual(mockSummary)
    })

    it('getList calls GET /staff-attendance with date param', async () => {
      const mockList = [{ id: 'att-1', user_id: 'u-1', status: 'present' }]
      api.get.mockResolvedValue({ data: mockList })

      const res = await staffAttendanceService.getList('2026-10-15')
      expect(api.get).toHaveBeenCalledWith('/staff-attendance', { params: { date: '2026-10-15' } })
      expect(res).toEqual(mockList)
    })

    it('setStrikeMode calls POST /staff-attendance/strike-mode', async () => {
      const payload = { date: '2026-10-15', is_strike_day: true }
      api.post.mockResolvedValue({ data: { success: true } })

      const res = await staffAttendanceService.setStrikeMode(payload)
      expect(api.post).toHaveBeenCalledWith('/staff-attendance/strike-mode', payload)
      expect(res.success).toBe(true)
    })

    it('recordAttendance calls POST /staff-attendance', async () => {
      const payload = { user_id: 'u-1', date: '2026-10-15', status: 'present' }
      api.post.mockResolvedValue({ data: { id: 'att-new', ...payload } })

      const res = await staffAttendanceService.recordAttendance(payload)
      expect(api.post).toHaveBeenCalledWith('/staff-attendance', payload)
      expect(res.id).toBe('att-new')
    })

    it('recordBulk calls POST /staff-attendance/bulk', async () => {
      const payload = {
        date: '2026-10-15',
        attendances: [{ user_id: 'u-1', status: 'present' }]
      }
      api.post.mockResolvedValue({ data: { count: 1 } })

      const res = await staffAttendanceService.recordBulk(payload)
      expect(api.post).toHaveBeenCalledWith('/staff-attendance/bulk', payload)
      expect(res.count).toBe(1)
    })

    it('deleteAttendance calls DELETE /staff-attendance/:id', async () => {
      api.delete.mockResolvedValue({ data: { success: true } })

      const res = await staffAttendanceService.deleteAttendance('att-99')
      expect(api.delete).toHaveBeenCalledWith('/staff-attendance/att-99')
      expect(res.success).toBe(true)
    })
  })

  describe('Badge Swipes and Assignment Endpoints', () => {
    it('registerBadgeSwipe calls POST /staff-attendance/badge-swipe', async () => {
      const swipe = { badge_code: 'B-101', device_id: 'DEV-1', swipe_type: 'in' }
      api.post.mockResolvedValue({ data: { id: 'swipe-1' } })

      const res = await staffAttendanceService.registerBadgeSwipe(swipe)
      expect(api.post).toHaveBeenCalledWith('/staff-attendance/badge-swipe', swipe)
      expect(res.id).toBe('swipe-1')
    })

    it('processBadgeSwipes calls POST /staff-attendance/badge-swipe/process', async () => {
      api.post.mockResolvedValue({ data: { processed: 5 } })

      const res = await staffAttendanceService.processBadgeSwipes()
      expect(api.post).toHaveBeenCalledWith('/staff-attendance/badge-swipe/process')
      expect(res.processed).toBe(5)
    })

    it('assignBadge calls POST /staff-attendance/badges', async () => {
      const badgeData = { user_id: 'u-ata', badge_code: 'B-202' }
      api.post.mockResolvedValue({ data: { success: true } })

      const res = await staffAttendanceService.assignBadge(badgeData)
      expect(api.post).toHaveBeenCalledWith('/staff-attendance/badges', badgeData)
      expect(res.success).toBe(true)
    })
  })

  describe('Timecard and Leaves Endpoints', () => {
    it('getTimecard calls GET /staff-attendance/timecard with params', async () => {
      const mockTimecard = { month: '2026-10', total_hours: 154 }
      api.get.mockResolvedValue({ data: mockTimecard })

      const res = await staffAttendanceService.getTimecard({ month: '2026-10' })
      expect(api.get).toHaveBeenCalledWith('/staff-attendance/timecard', { params: { month: '2026-10' } })
      expect(res).toEqual(mockTimecard)
    })

    it('exportTimecard calls GET with responseType blob', async () => {
      const mockBlob = new Blob(['csv-content'], { type: 'text/csv' })
      api.get.mockResolvedValue({ data: mockBlob })

      const res = await staffAttendanceService.exportTimecard({ month: '2026-10' })
      expect(api.get).toHaveBeenCalledWith('/staff-attendance/timecard/export', {
        params: { month: '2026-10' },
        responseType: 'blob'
      })
      expect(res).toEqual(mockBlob)
    })

    it('listLeaves calls GET /staff-attendance/leaves with filter params', async () => {
      const mockLeaves = [{ id: 'l-1', type: 'ferie', status: 'approved' }]
      api.get.mockResolvedValue({ data: mockLeaves })

      const res = await staffAttendanceService.listLeaves({ status: 'approved' })
      expect(api.get).toHaveBeenCalledWith('/staff-attendance/leaves', { params: { status: 'approved' } })
      expect(res).toEqual(mockLeaves)
    })

    it('createLeave calls POST /staff-attendance/leaves', async () => {
      const newLeave = { type: 'ferie', start_date: '2026-11-01', end_date: '2026-11-03' }
      api.post.mockResolvedValue({ data: { id: 'l-new', ...newLeave } })

      const res = await staffAttendanceService.createLeave(newLeave)
      expect(api.post).toHaveBeenCalledWith('/staff-attendance/leaves', newLeave)
      expect(res.id).toBe('l-new')
    })

    it('approveLeave calls PATCH /staff-attendance/leaves/:id/approve', async () => {
      api.patch.mockResolvedValue({ data: { success: true } })

      const res = await staffAttendanceService.approveLeave('l-5', { notes: 'Approvato' })
      expect(api.patch).toHaveBeenCalledWith('/staff-attendance/leaves/l-5/approve', { notes: 'Approvato' })
      expect(res.success).toBe(true)
    })

    it('rejectLeave calls PATCH /staff-attendance/leaves/:id/reject', async () => {
      api.patch.mockResolvedValue({ data: { success: true } })

      const res = await staffAttendanceService.rejectLeave('l-5', { reason: 'Esigenze di servizio' })
      expect(api.patch).toHaveBeenCalledWith('/staff-attendance/leaves/l-5/reject', { reason: 'Esigenze di servizio' })
      expect(res.success).toBe(true)
    })

    it('deleteLeave calls DELETE /staff-attendance/leaves/:id', async () => {
      api.delete.mockResolvedValue({ data: { success: true } })

      const res = await staffAttendanceService.deleteLeave('l-10')
      expect(api.delete).toHaveBeenCalledWith('/staff-attendance/leaves/l-10')
      expect(res.success).toBe(true)
    })
  })
})
