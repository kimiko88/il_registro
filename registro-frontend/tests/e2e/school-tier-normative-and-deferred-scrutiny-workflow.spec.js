import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import SchoolManagement from '@/pages/admin/SchoolManagement.vue'
import Scrutiny from '@/pages/teacher/Scrutiny.vue'
import adminService from '@/services/adminService'
import { scrutinyService } from '@/services/scrutinyService'
import { createPinia, setActivePinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useClassesStore } from '@/stores/classes'

// Mock services
vi.mock('@/services/adminService', () => ({
  default: {
    getSchools: vi.fn(),
    createSchool: vi.fn(),
    updateSchool: vi.fn(),
    deleteSchool: vi.fn(),
    getTiers: vi.fn()
  }
}))

vi.mock('@/services/scrutinyService', () => ({
  scrutinyService: {
    getMatrix: vi.fn(),
    save: vi.fn(),
    saveDeficiency: vi.fn(),
    getStudentDeficiencies: vi.fn(),
    saveDeferredScrutiny: vi.fn(),
    lock: vi.fn(),
    generatePagellaPdf: vi.fn(),
    exportClassZip: vi.fn()
  }
}))

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn((url) => {
      if (url.includes('/school-calendar/periods')) {
        return Promise.resolve({
          data: [
            { period: 1, name: '1° Quadrimestre' },
            { period: 2, name: '2° Quadrimestre' }
          ]
        })
      }
      return Promise.resolve({ data: [] })
    })
  }
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: vi.fn() })
}))

vi.mock('@/composables/usePermissions', () => ({
  usePermissions: () => ({
    isSuperAdmin: true,
    canCreateSchools: true,
    canDeleteSchools: true,
    canEditSchool: () => true
  })
}))

describe('E2E Workflow: School Tier Normative Configuration & Deferred Scrutiny (O.M. 92/2007)', () => {
  let pinia

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)
    vi.clearAllMocks()
  })

  // ─── Step 1: Admin configures schools with normative tiers ────────────────
  it('E2E Step 1: SuperAdmin manages schools with normative levels (Infanzia, Primaria, Secondaria)', async () => {
    const mockSchoolsList = [
      {
        id: 'school-inf-1',
        name: 'Polo Infanzia Munari',
        code: 'RMAA100',
        school_level: 'infanzia',
        city: 'Roma',
        province: 'RM',
        student_count: 60,
        teacher_count: 6,
        is_active: true
      },
      {
        id: 'school-pri-1',
        name: 'Scuola Primaria Rodari',
        code: 'RMEE100',
        school_level: 'primaria',
        city: 'Roma',
        province: 'RM',
        student_count: 240,
        teacher_count: 20,
        is_active: true
      },
      {
        id: 'school-sec2-1',
        name: 'Liceo Statale Righi',
        code: 'RMPS100',
        school_level: 'secondaria_secondo_grado',
        city: 'Roma',
        province: 'RM',
        student_count: 900,
        teacher_count: 70,
        is_active: true
      }
    ]

    adminService.getSchools.mockResolvedValue({
      data: {
        items: mockSchoolsList,
        total: mockSchoolsList.length
      }
    })
    adminService.createSchool.mockResolvedValue({ data: { id: 'new-school-primaria-id' } })

    const wrapper = mount(SchoolManagement, {
      global: {
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-table': true,
          'q-badge': true,
          'q-btn': true,
          'q-select': true,
          'q-input': true,
          'q-dialog': true,
          'q-toggle': true,
          'q-icon': true,
          'q-tooltip': true,
          'q-space': true
        }
      }
    })

    // Verify page initialization
    expect(wrapper.text()).toContain('Gestione Scuole')

    // Verify tier options cover all Italian levels
    const tierOptions = wrapper.vm.tierSelectOptions
    expect(tierOptions.map(t => t.value)).toEqual([
      'infanzia',
      'primaria',
      'secondaria_primo_grado',
      'secondaria_secondo_grado',
      'comprensivo',
      'omnicomprensivo'
    ])

    // Admin creates a new Primaria school
    wrapper.vm.openCreate()
    expect(wrapper.vm.showCreateDialog).toBe(true)

    wrapper.vm.schoolForm.name = 'Scuola Primaria Collodi'
    wrapper.vm.schoolForm.code = 'RMEE20000X'
    wrapper.vm.schoolForm.school_level = 'primaria'
    wrapper.vm.schoolForm.address = 'Via Pinocchio 1'
    wrapper.vm.schoolForm.city = 'Roma'
    wrapper.vm.schoolForm.province = 'RM'
    wrapper.vm.schoolForm.zip_code = '00100'

    await wrapper.vm.saveSchool()

    expect(adminService.createSchool).toHaveBeenCalledWith(
      expect.objectContaining({
        name: 'Scuola Primaria Collodi',
        code: 'RMEE20000X',
        school_level: 'primaria'
      })
    )
  })

  // ─── Step 2: Teacher performs deferred scrutiny for student with suspended judgment
  it('E2E Step 2: Teacher accesses Period 3, inspects deficiencies, and resolves suspended judgment (O.M. 92/2007)', async () => {
    const authStore = useAuthStore()
    authStore.user = { id: 'teacher-sec2', name: 'Prof. Bianchi', role: 'teacher' }

    const classesStore = useClassesStore()
    classesStore.classes = [
      { id: 'class-5a', name: '5', section: 'A', coordinator_id: 'teacher-sec2', academic_year: '2025/2026' }
    ]

    const studentWithDebt = {
      student_id: 'student-debt-42',
      student_name: 'Lorenzo Neri',
      subject_data: { 'sub-fisica': { average: 5.0 } },
      attendance_stats: { absences: 8, lates: 2, early_exits: 0 },
      record: {
        conduct_grade: 8,
        final_decision: 'Sospeso',
        grades: [{ subject_id: 'sub-fisica', final_grade: 5.0 }]
      }
    }

    scrutinyService.getMatrix.mockResolvedValue({
      data: {
        subjects: [{ id: 'sub-fisica', name: 'Fisica', is_religion: false }],
        students: [studentWithDebt]
      }
    })

    scrutinyService.getStudentDeficiencies.mockResolvedValue({
      data: [
        {
          id: 'def-phys-1',
          subject_id: 'sub-fisica',
          subject_name: 'Fisica',
          topics: 'Elettromagnetismo e circuiti RLC',
          status: 'da_recuperare',
          recovery_grade: null
        }
      ]
    })

    scrutinyService.saveDeferredScrutiny.mockResolvedValue({
      data: { message: 'deferred scrutiny saved successfully' }
    })

    const wrapper = mount(Scrutiny, {
      global: {
        stubs: {
          'q-page': { template: '<div class="q-page"><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-banner': { template: '<div class="q-banner"><slot /><slot name="avatar" /></div>' },
          'q-table': true,
          'q-tr': true,
          'q-th': true,
          'q-td': true,
          'q-btn': true,
          'q-select': true,
          'q-btn-toggle': true,
          'q-input': true,
          'q-tooltip': true,
          'q-icon': true,
          'q-separator': true,
          'q-dialog': true,
          'q-badge': true
        }
      }
    })

    // 1. Verify period 3 exists in options
    expect(wrapper.vm.periodOptions.some(p => p.value === 3)).toBe(true)

    // 2. Select Period 3 (Scrutinio Differito)
    wrapper.vm.period = 3
    wrapper.vm.selectedClassId = 'class-5a'
    await wrapper.vm.$nextTick()

    // 3. Deferred scrutiny banner is visible
    const banner = wrapper.find('.q-banner')
    expect(banner.exists()).toBe(true)
    expect(banner.text()).toContain('O.M. 92/2007')

    // 4. Open deferred scrutiny modal for Lorenzo Neri
    await wrapper.vm.openDeferredModal(studentWithDebt)

    expect(wrapper.vm.showDeferredModal).toBe(true)
    expect(wrapper.vm.selectedStudent.student_id).toBe('student-debt-42')
    expect(wrapper.vm.studentDeficienciesList).toHaveLength(1)
    expect(wrapper.vm.studentDeficienciesList[0].topics).toContain('Elettromagnetismo')

    // 5. Teacher evaluates exam result and promotes student
    wrapper.vm.deferredForm.final_decision = 'promosso_con_debiti_saldati'
    wrapper.vm.deferredForm.notes = 'Debito di Fisica saldato con successo in data 01/09.'
    wrapper.vm.studentDeficienciesList[0].status = 'recuperato'
    wrapper.vm.studentDeficienciesList[0].recovery_grade = 6.5

    await wrapper.vm.saveDeferredScrutiny()

    // 6. Verify deferred scrutiny submission to backend
    expect(scrutinyService.saveDeferredScrutiny).toHaveBeenCalledWith({
      student_id: 'student-debt-42',
      class_id: 'class-5a',
      final_decision: 'promosso_con_debiti_saldati',
      notes: 'Debito di Fisica saldato con successo in data 01/09.',
      deficiencies: [
        {
          deficiency_id: 'def-phys-1',
          status: 'recuperato',
          recovery_grade: 6.5
        }
      ]
    })
    expect(wrapper.vm.showDeferredModal).toBe(false)
  })
})
