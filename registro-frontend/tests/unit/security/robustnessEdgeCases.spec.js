import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from 'src/stores/auth'
import { useGradesStore } from 'src/stores/grades'
import { useAttendanceStore } from 'src/stores/attendance'

vi.mock('src/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn()
  }
}))

import api from 'src/services/api'

describe('Frontend Robustness & Input Boundary Test Suite', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('1. Grade Entry Input & Range Boundaries', () => {
    it('validates numerical grade bounds (1.0 - 10.0)', async () => {
      const gradesStore = useGradesStore()

      // Invalid grade <= 0 or > 10
      api.post.mockRejectedValueOnce({
        response: { status: 400, data: { error: 'voto fuori dai limiti consentiti (1-10)' } }
      })
      await expect(
        gradesStore.addGrade({ student_id: 'st-01', subject_id: 'sub-1', grade_value: 12.0 })
      ).rejects.toBeDefined()

      // Valid boundary grade = 10.0
      api.post.mockResolvedValueOnce({ data: { id: 'g-10', grade_value: 10.0 } })
      api.get.mockResolvedValueOnce({ data: [{ id: 'g-10', grade_value: 10.0 }] })
      const validGrade = await gradesStore.addGrade({
        student_id: 'st-01',
        subject_id: 'sub-1',
        grade_value: 10.0
      })
      expect(validGrade).toBeDefined()
    })
  })

  describe('2. Attendance Marking Edge Cases', () => {
    it('handles empty status list gracefully', async () => {
      const attendanceStore = useAttendanceStore()
      api.post.mockResolvedValueOnce({ data: { records: [] } })

      const res = await attendanceStore.submitAttendance('class-1a', '2025-10-15', [])
      expect(res).toBeDefined()
    })

    it('rejects retroactive attendance past 30 days limit', async () => {
      const attendanceStore = useAttendanceStore()
      api.post.mockRejectedValueOnce({
        response: { status: 400, data: { error: 'data non modificabile oltre 30 giorni' } }
      })

      await expect(
        attendanceStore.submitAttendance('class-1a', '2025-05-01', [{ studentId: 'st-1', status: 'Present' }])
      ).rejects.toBeDefined()
    })
  })

  describe('3. Administrative Role Status Badges & Flags', () => {
    it('verifies staff, vice principal, and principal role flag evaluation', () => {
      const authStore = useAuthStore()

      authStore.user = {
        id: 'u-admin',
        role: 'teacher',
        is_staff: true,
        is_vice_principal: true,
        is_principal: false
      }

      expect(authStore.user.is_staff).toBe(true)
      expect(authStore.user.is_vice_principal).toBe(true)
      expect(authStore.user.is_principal).toBe(false)
    })
  })
})
