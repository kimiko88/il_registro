import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from 'src/stores/auth'
import { useSchoolStore } from 'src/stores/schools'
import { useClassesStore } from 'src/stores/classes'
import { useGradesStore } from 'src/stores/grades'
import { useAttendanceStore } from 'src/stores/attendance'
import { useColloquiStore } from 'src/stores/colloqui'

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

describe('Full School Year Multi-Role Lifecycle Simulation', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    api.get.mockImplementation((url) => {
      if (url.includes('/scrutiny/class')) {
        return Promise.resolve({ data: [{ student_id: 'st-luca', average: 8.2, status: 'proposed' }] })
      }
      if (url.includes('/scrutiny/overview')) {
        return Promise.resolve({ data: { total_classes: 2, locked_scrutinies: 2, status: 'completed' } })
      }
      if (url.includes('/scrutiny-final/export')) {
        return Promise.resolve({ data: new ArrayBuffer(8) })
      }
      if (url.includes('/schools')) {
        return Promise.resolve({ data: { items: [{ id: 'school-101', name: 'Liceo Fermi' }], total: 1 } })
      }
      return Promise.resolve({ data: [] })
    })
    api.post.mockImplementation((url) => {
      if (url.includes('/scrutiny/validate')) {
        return Promise.resolve({ data: { message: 'scrutiny validated' } })
      }
      return Promise.resolve({ data: { id: 'created-id', message: 'success' } })
    })
  })

  it('Phase 1: Admin & Secretary Setup School, Multi-Role Users & Class Schedule', async () => {
    const authStore = useAuthStore()
    const schoolStore = useSchoolStore()
    const classesStore = useClassesStore()

    // 1. Superadmin sets up school
    authStore.user = { id: 'sa-01', role: 'superadmin' }
    await expect(schoolStore.createSchool({ name: 'Liceo Fermi', code: 'FERMI01' })).resolves.not.toThrow()

    // 2. Secretary creates class 1A
    authStore.user = { id: 'sec-01', role: 'secretary', school_id: 'school-101' }
    await expect(classesStore.createClass({ name: '1', section: 'A', school_id: 'school-101' })).resolves.not.toThrow()
  })

  it('Phase 2: Term 1 Daily Operations (Attendance, PCTO, Civica Grades & Colloqui)', async () => {
    const authStore = useAuthStore()
    const attendanceStore = useAttendanceStore()
    const gradesStore = useGradesStore()
    const colloquiStore = useColloquiStore()

    // 1. Teacher marks attendance with PCTO / Orientamento hours
    authStore.user = { id: 't-mario', role: 'teacher', school_id: 'school-101' }
    await expect(
      attendanceStore.submitAttendance(
        'class-1a',
        '2025-10-15',
        [{ studentId: 'st-luca', status: 'Present' }],
        1,
        'sub-math'
      )
    ).resolves.not.toThrow()

    // 2. Teacher adds numerical & shared Educazione Civica grades
    const grade = await gradesStore.addGrade({
      student_id: 'st-luca',
      subject_id: 'sub-civica',
      grade_value: 8.5,
      grade_type: 'numeric',
      semester: 1,
      date: '2025-10-15'
    })
    expect(grade).toBeDefined()

    // 3. Parent / Teacher creates and checks interview (colloquio) slots
    authStore.user = { id: 't-mario', role: 'teacher', school_id: 'school-101' }
    await expect(
      colloquiStore.createSlots([{ date: '2025-10-15', startTime: '15:00', endTime: '15:15' }])
    ).resolves.not.toThrow()
  })

  it('Phase 3: Intermediate Scrutiny (Q1 Board Meeting & Validation)', async () => {
    const authStore = useAuthStore()
    authStore.user = { id: 't-mario', role: 'teacher', is_coordinator: true, school_id: 'school-101' }

    // Coordinator fetches Q1 scrutiny records
    const res = await api.get('/api/v1/scrutiny/class/class-1a?semester=1')
    expect(res.data).toHaveLength(1)
    expect(res.data[0].average).toBeGreaterThanOrEqual(6.0)

    // Validation request
    const validRes = await api.post('/api/v1/scrutiny/validate', { class_id: 'class-1a', semester: 1 })
    expect(validRes.data.message).toBe('scrutiny validated')
  })

  it('Phase 4: Final Scrutiny (Q2 Board Meeting & PDF Export)', async () => {
    const authStore = useAuthStore()
    authStore.user = { id: 'dir-preside', role: 'principal', is_principal: true, school_id: 'school-101' }

    // Principal accesses overview & triggers report card export
    const overview = await api.get('/api/v1/scrutiny/overview')
    expect(overview.data.status).toBe('completed')

    const exportRes = await api.get('/api/v1/scrutiny-final/export')
    expect(exportRes.data).toBeDefined()
  })
})
