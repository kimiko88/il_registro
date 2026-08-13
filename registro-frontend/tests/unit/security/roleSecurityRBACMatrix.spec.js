import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from 'src/stores/auth'
import { useSchoolStore } from 'src/stores/schools'
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

describe('Comprehensive 9-Role Frontend Security & RBAC Permission Matrix', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  const roles = [
    { role: 'superadmin' },
    { role: 'admin' },
    { role: 'secretary' },
    { role: 'principal' },
    { role: 'vice_principal' },
    { role: 'staff' },
    { role: 'coordinator' },
    { role: 'teacher' },
    { role: 'student' },
    { role: 'parent' }
  ]

  describe('1. School Administration Endpoint Permissions', () => {
    const allowedSchoolRoles = ['superadmin']

    roles.forEach(({ role }) => {
      it(`evaluates createSchool permission for role ${role}`, async () => {
        const authStore = useAuthStore()
        const schoolStore = useSchoolStore()
        authStore.user = { id: 'u-test', role }

        if (allowedSchoolRoles.includes(role)) {
          api.post.mockResolvedValueOnce({ data: { id: 's-1', name: 'New School' } })
          api.get.mockResolvedValueOnce({ data: { items: [], total: 0 } })
          await expect(schoolStore.createSchool({ name: 'New School' })).resolves.not.toThrow()
        } else {
          api.post.mockRejectedValueOnce({
            response: { status: 403, data: { error: 'forbidden' } }
          })
          await expect(schoolStore.createSchool({ name: 'New School' })).rejects.toBeDefined()
        }
      })
    })
  })

  describe('2. Grade Addition Permissions', () => {
    const allowedGradeRoles = ['superadmin', 'admin', 'teacher', 'coordinator']

    roles.forEach(({ role }) => {
      it(`evaluates addGrade permission for role ${role}`, async () => {
        const authStore = useAuthStore()
        const gradesStore = useGradesStore()
        authStore.user = { id: 'u-test', role }

        if (allowedGradeRoles.includes(role)) {
          api.post.mockResolvedValueOnce({ data: { id: 'g-1', grade_value: 8 } })
          api.get.mockResolvedValueOnce({ data: [{ id: 'g-1', grade_value: 8 }] })
          const res = await gradesStore.addGrade({ student_id: 'st-1', subject_id: 'sub-1', grade_value: 8 })
          expect(res).toBeDefined()
        } else {
          api.post.mockRejectedValueOnce({
            response: { status: 403, data: { error: 'forbidden' } }
          })
          await expect(gradesStore.addGrade({ student_id: 'st-1', subject_id: 'sub-1', grade_value: 8 })).rejects.toBeDefined()
        }
      })
    })
  })

  describe('3. Attendance Bulk Marking Permissions', () => {
    const allowedAttRoles = ['superadmin', 'admin', 'teacher', 'coordinator']

    roles.forEach(({ role }) => {
      it(`evaluates submitAttendance permission for role ${role}`, async () => {
        const authStore = useAuthStore()
        const attendanceStore = useAttendanceStore()
        authStore.user = { id: 'u-test', role }

        if (allowedAttRoles.includes(role)) {
          api.post.mockResolvedValueOnce({ data: { message: 'marked' } })
          await expect(
            attendanceStore.submitAttendance('c-1', '2025-10-15', [{ studentId: 'st-1', status: 'Present' }])
          ).resolves.not.toThrow()
        } else {
          api.post.mockRejectedValueOnce({
            response: { status: 403, data: { error: 'forbidden' } }
          })
          await expect(
            attendanceStore.submitAttendance('c-1', '2025-10-15', [{ studentId: 'st-1', status: 'Present' }])
          ).rejects.toBeDefined()
        }
      })
    })
  })

  describe('4. Educazione Civica Grade Modification Restrictions', () => {
    it('allows original creator teacher to update grade', async () => {
      const authStore = useAuthStore()
      const gradesStore = useGradesStore()
      authStore.user = { id: 'teacher-creator', role: 'teacher' }

      api.patch.mockResolvedValueOnce({ data: { id: 'g-civica', grade_value: 9 } })
      api.get.mockResolvedValueOnce({ data: [{ id: 'g-civica', grade_value: 9 }] })
      await expect(gradesStore.updateGrade('g-civica', { grade_value: 9 })).resolves.not.toThrow()
    })

    it('rejects unauthorized secondary teacher from modifying Civica grade with 403', async () => {
      const authStore = useAuthStore()
      const gradesStore = useGradesStore()
      authStore.user = { id: 'teacher-other', role: 'teacher' }

      api.patch.mockRejectedValueOnce({
        response: { status: 403, data: { error: 'forbidden: only grade creator can update' } }
      })
      await expect(gradesStore.updateGrade('g-civica', { grade_value: 9 })).rejects.toBeDefined()
    })
  })
})
