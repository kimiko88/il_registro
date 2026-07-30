import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from 'src/stores/auth'
import { useSchoolStore } from 'src/stores/schools'
import { useClassesStore } from 'src/stores/classes'
import { useGradesStore } from 'src/stores/grades'
import { useAttendanceStore } from 'src/stores/attendance'
import api from 'src/services/api'

vi.mock('src/services/api')

describe('Student and Parent RBAC Access Control Tests', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  const restrictedRoles = ['student', 'parent']

  restrictedRoles.forEach(role => {
    describe(`Role: ${role}`, () => {
      it('rejects administrative school creation with HTTP 403', async () => {
        const authStore = useAuthStore()
        const schoolStore = useSchoolStore()
        authStore.user = { id: 'u1', role }

        api.post.mockRejectedValueOnce({
          response: { status: 403, data: { error: 'forbidden' } }
        })

        await expect(schoolStore.createSchool({ name: 'Unauthorized School' })).rejects.toBeDefined()
      })

      it('rejects class creation with HTTP 403', async () => {
        const authStore = useAuthStore()
        const classesStore = useClassesStore()
        authStore.user = { id: 'u1', role }

        api.post.mockRejectedValueOnce({
          response: { status: 403, data: { error: 'forbidden' } }
        })

        await expect(classesStore.createClass({ name: '5Z' })).rejects.toBeDefined()
      })

      it('rejects grade addition by student or parent with HTTP 403', async () => {
        const authStore = useAuthStore()
        const gradesStore = useGradesStore()
        authStore.user = { id: 'u1', role }

        api.post.mockRejectedValueOnce({
          response: { status: 403, data: { error: 'forbidden' } }
        })

        await expect(
          gradesStore.addGrade({
            student_id: 'student-1',
            subject_id: 'sub-1',
            grade_value: 10
          })
        ).rejects.toBeDefined()
      })

      it('rejects attendance submit for class with HTTP 403', async () => {
        const authStore = useAuthStore()
        const attendanceStore = useAttendanceStore()
        authStore.user = { id: 'u1', role }

        api.post.mockRejectedValueOnce({
          response: { status: 403, data: { error: 'forbidden' } }
        })

        await expect(
          attendanceStore.submitAttendance('class-1a', '2026-07-30', [{ studentId: 'student-1', status: 'Present' }])
        ).rejects.toBeDefined()
      })
    })
  })
})
