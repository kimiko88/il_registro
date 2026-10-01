import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import Scrutiny from '@/pages/teacher/Scrutiny.vue'
import { createPinia, setActivePinia } from 'pinia'
import { scrutinyService } from '@/services/scrutinyService'
import { useAuthStore } from '@/stores/auth'
import { useClassesStore } from '@/stores/classes'

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

describe('Teacher Scrutiny - Deferred Scrutiny & Normative Period 3 (O.M. 92/2007)', () => {
  let pinia
  let authStore
  let classesStore
  let wrapper

  const mockStudents = [
    {
      student_id: 'stud-deferred-1',
      student_name: 'Giulia Bianchi',
      subject_data: { 'sub-math': { average: 4.5 } },
      attendance_stats: { absences: 5, lates: 1, early_exits: 0 },
      record: {
        conduct_grade: 8,
        final_decision: 'Sospeso',
        grades: [{ subject_id: 'sub-math', final_grade: 4.5 }]
      }
    },
    {
      student_id: 'stud-regular-2',
      student_name: 'Marco Verdi',
      subject_data: { 'sub-math': { average: 7.5 } },
      attendance_stats: { absences: 2, lates: 0, early_exits: 0 },
      record: {
        conduct_grade: 9,
        final_decision: 'Ammesso',
        grades: [{ subject_id: 'sub-math', final_grade: 7.5 }]
      }
    }
  ]

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)

    authStore = useAuthStore()
    authStore.user = { id: 'teacher-coord-1', name: 'Prof. Coordinatore', role: 'teacher' }

    classesStore = useClassesStore()
    classesStore.classes = [
      { id: 'class-4b', name: '4', section: 'B', coordinator_id: 'teacher-coord-1', academic_year: '2025/2026' }
    ]

    vi.clearAllMocks()

    scrutinyService.getMatrix.mockResolvedValue({
      data: {
        subjects: [{ id: 'sub-math', name: 'Matematica', is_religion: false }],
        students: mockStudents
      }
    })

    scrutinyService.getStudentDeficiencies.mockResolvedValue({
      data: [
        {
          id: 'def-101',
          subject_id: 'sub-math',
          subject_name: 'Matematica',
          topics: 'Trigonometria ed equazioni goniometriche',
          status: 'da_recuperare',
          recovery_grade: null
        }
      ]
    })

    scrutinyService.saveDeferredScrutiny.mockResolvedValue({
      data: { message: 'deferred scrutiny saved successfully' }
    })

    wrapper = mount(Scrutiny, {
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
  })

  it('guarantees Period 3 (Scrutinio Differito) is always present in periodOptions', () => {
    const options = wrapper.vm.periodOptions
    expect(options).toBeDefined()
    const period3 = options.find(p => p.value === 3)
    expect(period3).toBeDefined()
    expect(period3.label).toBeTruthy()
  })

  it('renders deferred scrutiny banner and notice when period is 3', async () => {
    wrapper.vm.period = 3
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.period).toBe(3)
    const banner = wrapper.find('.q-banner')
    expect(banner.exists()).toBe(true)
    expect(banner.text()).toContain('O.M. 92/2007')
  })

  it('queries 2nd semester matrix when period is 3 (deferred scrutiny operates on June evaluations)', async () => {
    wrapper.vm.selectedClassId = 'class-4b'
    wrapper.vm.period = 3
    await wrapper.vm.fetchMatrix()

    expect(scrutinyService.getMatrix).toHaveBeenCalledWith('class-4b', 2)
  })

  it('opens deferred scrutiny modal and loads student deficiencies on openDeferredModal', async () => {
    wrapper.vm.selectedClassId = 'class-4b'
    await wrapper.vm.openDeferredModal(mockStudents[0])

    expect(wrapper.vm.showDeferredModal).toBe(true)
    expect(wrapper.vm.selectedStudent).toEqual(mockStudents[0])
    expect(wrapper.vm.deferredForm.student_id).toBe('stud-deferred-1')
    expect(wrapper.vm.deferredForm.class_id).toBe('class-4b')
    expect(wrapper.vm.deferredForm.final_decision).toBe('promosso_con_debiti_saldati')
    expect(scrutinyService.getStudentDeficiencies).toHaveBeenCalledWith('stud-deferred-1')
    expect(wrapper.vm.studentDeficienciesList).toHaveLength(1)
  })

  it('saves deferred scrutiny resolution and refreshes matrix', async () => {
    wrapper.vm.selectedClassId = 'class-4b'
    await wrapper.vm.openDeferredModal(mockStudents[0])

    wrapper.vm.deferredForm.final_decision = 'promosso_con_debiti_saldati'
    wrapper.vm.deferredForm.notes = 'Debito di Matematica superato con prova scritta di recupero il 28/08.'
    wrapper.vm.studentDeficienciesList[0].status = 'recuperato'
    wrapper.vm.studentDeficienciesList[0].recovery_grade = 6

    await wrapper.vm.saveDeferredScrutiny()

    expect(scrutinyService.saveDeferredScrutiny).toHaveBeenCalledWith({
      student_id: 'stud-deferred-1',
      class_id: 'class-4b',
      final_decision: 'promosso_con_debiti_saldati',
      notes: 'Debito di Matematica superato con prova scritta di recupero il 28/08.',
      deficiencies: [
        {
          deficiency_id: 'def-101',
          status: 'recuperato',
          recovery_grade: 6
        }
      ]
    })
    expect(wrapper.vm.showDeferredModal).toBe(false)
  })
})
