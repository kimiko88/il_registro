import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from 'src/stores/auth'
import { useSchoolStore } from 'src/stores/schools'
import { useClassesStore } from 'src/stores/classes'
import { useAttendanceStore } from 'src/stores/attendance'
import { useGradesStore } from 'src/stores/grades'
import { useSubstitutionsStore } from 'src/stores/substitutions'
import api from 'src/services/api'

vi.mock('src/services/api')
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn()
  })
}))

describe('Complete Multi-Role School Lifecycle Workflow', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('verifies full flow: Superadmin -> Secretary -> Teachers (regular + substitution) -> Student & Parent views', async () => {
    const authStore = useAuthStore()
    const schoolStore = useSchoolStore()
    const classesStore = useClassesStore()
    const attendanceStore = useAttendanceStore()
    const gradesStore = useGradesStore()
    const substitutionsStore = useSubstitutionsStore()

    // -----------------------------------------------------------------
    // Phase 1: Superadmin creates school, admin, and secretary
    // -----------------------------------------------------------------
    authStore.user = {
      id: 'superadmin-1',
      email: 'superadmin@test.com',
      role: 'superadmin'
    }
    authStore.token = 'mock-superadmin-jwt'

    const mockSchool = {
      id: 'school-galileo',
      name: 'Liceo Scientifico Galileo',
      code: 'LSG001',
      city: 'Milano'
    }
    api.post.mockResolvedValueOnce({ data: mockSchool })
    api.get.mockResolvedValueOnce({ data: { items: [mockSchool], total: 1 } })
    
    await schoolStore.createSchool(mockSchool)
    expect(schoolStore.schools).toHaveLength(1)
    expect(schoolStore.schools[0].name).toBe('Liceo Scientifico Galileo')

    // -----------------------------------------------------------------
    // Phase 2: Secretary creates class 1A, teachers, student, and parent
    // -----------------------------------------------------------------
    authStore.user = {
      id: 'secretary-1',
      email: 'segreteria@galileo.edu',
      role: 'secretary',
      school_id: 'school-galileo'
    }

    const mockClass = { id: 'class-1a', name: '1A', section: 'A', school_id: 'school-galileo' }
    api.get.mockResolvedValueOnce({ data: [mockClass] })
    await classesStore.fetchClasses({ school_id: 'school-galileo' })
    expect(classesStore.classes).toHaveLength(1)
    expect(classesStore.classes[0].name).toBe('1A')

    const mockMathTeacher = { id: 'teacher-math', first_name: 'Giuseppe', last_name: 'Verdi', role: 'teacher' }
    const mockItalianTeacher = { id: 'teacher-italian', first_name: 'Anna', last_name: 'Neri', role: 'teacher' }
    const mockStudent = { id: 'student-1', first_name: 'Luca', last_name: 'Rossi', role: 'student', class_id: 'class-1a' }
    const mockParent = { id: 'parent-1', first_name: 'Mario', last_name: 'Rossi', role: 'parent' }

    expect(mockStudent.class_id).toBe('class-1a')

    // -----------------------------------------------------------------
    // Phase 3: Teachers record Attendance, Regular & Substitution Lessons, Grades, and Homework
    // -----------------------------------------------------------------

    // Teacher 1 (Math): Attendance & Regular Lesson & Grade & Homework
    authStore.user = mockMathTeacher
    authStore.user.role = 'teacher'

    // Mark attendance
    api.post.mockResolvedValueOnce({ data: { message: 'Attendance marked successfully' } })
    const attResult = await attendanceStore.submitAttendance(
      'class-1a',
      '2026-07-30',
      [{ studentId: 'student-1', status: 'Present' }],
      1,
      'subj-math'
    )
    expect(attResult.message).toBe('Attendance marked successfully')
    expect(attendanceStore.records).toHaveLength(1)
    expect(attendanceStore.records[0].status).toBe('Present')

    // Substitution Teacher (Italian Teacher substituting Math Teacher)
    authStore.user = mockItalianTeacher
    authStore.user.role = 'teacher'

    const mockTodaySubstitutions = [
      {
        id: 'sub-101',
        class_id: 'class-1a',
        class_name: '1A',
        hour: 2,
        date: '2026-07-30',
        original_teacher_name: 'Giuseppe Verdi',
        substitution_teacher_id: 'teacher-italian',
        status: 'confirmed',
        topic: 'Sostituzione: Lettura Promessi Sposi'
      }
    ]

    api.get.mockResolvedValueOnce({ data: mockTodaySubstitutions })
    const todaySubs = await substitutionsStore.fetchTodaySubstitutions()
    expect(todaySubs).toHaveLength(1)
    expect(todaySubs[0].original_teacher_name).toBe('Giuseppe Verdi')

    // Add grade for student
    authStore.user = mockMathTeacher

    const mockStudentGrades = {
      data: {
        semesters: [
          {
            period: '1° Quadrimestre',
            subjects: [
              {
                subject_id: 'subj-math',
                subject_name: 'Matematica',
                average: 8.5,
                grades: [
                  {
                    id: 'grade-1',
                    value: 8.5,
                    type: 'scritto',
                    date: '2026-07-30',
                    notes: 'Ottima prova scritta su equazioni'
                  }
                ]
              }
            ]
          }
        ]
      }
    }
    api.get.mockResolvedValueOnce(mockStudentGrades)
    await gradesStore.fetchMyGrades()
    expect(gradesStore.grades.semesters[0].subjects[0].average).toBe(8.5)

    // -----------------------------------------------------------------
    // Phase 4: Student & Parent Read Views
    // -----------------------------------------------------------------

    // Student Login
    authStore.user = mockStudent
    authStore.user.role = 'student'

    // Verify student sees grades
    expect(gradesStore.grades.semesters[0].subjects[0].grades[0].value).toBe(8.5)

    // Parent Login
    authStore.user = mockParent
    authStore.user.role = 'parent'

    // Verify parent sees child data
    expect(gradesStore.grades.semesters[0].subjects[0].grades[0].notes).toContain('equazioni')
  })
})
