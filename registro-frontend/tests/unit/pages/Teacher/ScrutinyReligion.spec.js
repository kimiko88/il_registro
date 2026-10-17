import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import Scrutiny from '@/pages/teacher/Scrutiny.vue'
import { createPinia, setActivePinia } from 'pinia'

const { mockGetMatrix, mockSaveScrutiny } = vi.hoisted(() => ({
  mockGetMatrix: vi.fn(),
  mockSaveScrutiny: vi.fn()
}))

vi.mock('@/services/scrutinyService', () => ({
  scrutinyService: {
    getMatrix: mockGetMatrix,
    save: mockSaveScrutiny,
    saveDeficiency: vi.fn(() => Promise.resolve({ data: { message: 'ok' } })),
    getStudentDeficiencies: vi.fn(() => Promise.resolve({ data: [] })),
    saveDeferredScrutiny: vi.fn(() => Promise.resolve({ data: { message: 'ok' } }))
  }
}))

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn((url) => {
      if (url.includes('/school-calendar/periods')) {
        return Promise.resolve({ data: [{ period: 1, name: '1° Quadrimestre' }, { period: 2, name: '2° Quadrimestre' }] })
      }
      return Promise.resolve({ data: [] })
    })
  }
}))

describe('Teacher Scrutiny - Religion & Avvalimento Matrix Workflow', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()

    mockGetMatrix.mockResolvedValue({
      data: {
        subjects: [
          { id: 'sub-ita', name: 'Italiano', is_religion: false, is_judgment_only: false },
          { id: 'sub-rel', name: 'Religione Cattolica', is_religion: true, is_judgment_only: true }
        ],
        students: [
          {
            student_id: 'stud-1',
            student_name: 'Mario Rossi',
            religion_choice: 'avvalente',
            subject_data: {
              'sub-ita': { average: 7.5 },
              'sub-rel': { proposed_judgment: 'Ottimo' }
            },
            attendance_stats: { absences: 2, lates: 0, early_exits: 0 },
            record: { conduct_grade: 8, final_decision: 'Promosso' }
          },
          {
            student_id: 'stud-2',
            student_name: 'Luigi Verdi',
            religion_choice: 'non_avvalente',
            subject_data: {
              'sub-ita': { average: 6.0 },
              'sub-rel': { proposed_judgment: '-' }
            },
            attendance_stats: { absences: 5, lates: 1, early_exits: 0 },
            record: { conduct_grade: 7, final_decision: 'Promosso' }
          },
          {
            student_id: 'stud-3',
            student_name: 'Chiara Neri',
            religion_choice: 'attivita_alternativa',
            subject_data: {
              'sub-ita': { average: 8.0 },
              'sub-rel': { proposed_judgment: '-' }
            },
            attendance_stats: { absences: 1, lates: 0, early_exits: 0 },
            record: { conduct_grade: 9, final_decision: 'Promosso' }
          }
        ]
      }
    })

    mockSaveScrutiny.mockResolvedValue({ data: { message: 'Salvataggio riuscito' } })
  })

  function createWrapper() {
    return mount(Scrutiny, {
      global: {
        stubs: {
          'q-page': { template: '<div class="q-page"><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-btn': true,
          'q-select': true,
          'q-btn-toggle': true,
          'q-table': true,
          'q-tr': true,
          'q-th': true,
          'q-td': true,
          'q-input': true,
          'q-tooltip': true,
          'q-icon': true,
          'q-separator': true,
          'q-dialog': true,
          'q-badge': true
        },
        mocks: {
          $t: (key) => key
        }
      }
    })
  }

  it('loads matrix data and populates scrutinyData with judgments for avvalente students', async () => {
    const wrapper = createWrapper()
    wrapper.vm.selectedClassId = 'cls-1'
    await wrapper.vm.fetchMatrix()

    expect(mockGetMatrix).toHaveBeenCalledWith('cls-1', 1)
    expect(wrapper.vm.matrix.students).toHaveLength(3)

    // Check scrutinyData populated
    const s1Data = wrapper.vm.scrutinyData['stud-1']
    expect(s1Data).toBeDefined()
    expect(s1Data.conduct_grade).toBe(8)
    expect(s1Data.grades['sub-rel']).toBe('Ottimo')
  })

  it('excludes religion grades when saving scrutiny for non_avvalente student', async () => {
    const wrapper = createWrapper()
    wrapper.vm.selectedClassId = 'cls-1'
    await wrapper.vm.fetchMatrix()

    // student 2 is non_avvalente
    const success = await wrapper.vm.saveStudentScrutiny('stud-2', true)
    expect(success).toBe(true)

    expect(mockSaveScrutiny).toHaveBeenCalledWith(expect.objectContaining({
      student_id: 'stud-2',
      class_id: 'cls-1',
      conduct_grade: 7,
      grades: [
        { subject_id: 'sub-ita', final_grade: 6 }
      ]
    }))

    // Ensure sub-rel is NOT included in the grades payload for stud-2
    const lastCall = mockSaveScrutiny.mock.calls[0][0]
    const relGrade = lastCall.grades.find(g => g.subject_id === 'sub-rel')
    expect(relGrade).toBeUndefined()
  })

  it('excludes religion grades when saving scrutiny for attivita_alternativa student', async () => {
    const wrapper = createWrapper()
    wrapper.vm.selectedClassId = 'cls-1'
    await wrapper.vm.fetchMatrix()

    // student 3 is attivita_alternativa
    const success = await wrapper.vm.saveStudentScrutiny('stud-3', true)
    expect(success).toBe(true)

    expect(mockSaveScrutiny).toHaveBeenCalledWith(expect.objectContaining({
      student_id: 'stud-3',
      class_id: 'cls-1',
      conduct_grade: 9,
      grades: [
        { subject_id: 'sub-ita', final_grade: 8 }
      ]
    }))

    const lastCall = mockSaveScrutiny.mock.calls[0][0]
    const relGrade = lastCall.grades.find(g => g.subject_id === 'sub-rel')
    expect(relGrade).toBeUndefined()
  })

  it('maps ministerial judgment to numerical value when saving scrutiny for avvalente student', async () => {
    const wrapper = createWrapper()
    wrapper.vm.selectedClassId = 'cls-1'
    await wrapper.vm.fetchMatrix()

    // student 1 is avvalente with judgment 'Ottimo' (maps to 10)
    const success = await wrapper.vm.saveStudentScrutiny('stud-1', true)
    expect(success).toBe(true)

    expect(mockSaveScrutiny).toHaveBeenCalledWith(expect.objectContaining({
      student_id: 'stud-1',
      class_id: 'cls-1',
      conduct_grade: 8,
      grades: expect.arrayContaining([
        { subject_id: 'sub-ita', final_grade: 8 },
        { subject_id: 'sub-rel', final_grade: 10 }
      ])
    }))
  })

  it('handles "Non classificabile" judgment mapped to 0 for avvalente student', async () => {
    const wrapper = createWrapper()
    wrapper.vm.selectedClassId = 'cls-1'
    await wrapper.vm.fetchMatrix()

    // Change judgment to Non classificabile
    wrapper.vm.scrutinyData['stud-1'].grades['sub-rel'] = 'Non classificabile'

    const success = await wrapper.vm.saveStudentScrutiny('stud-1', true)
    expect(success).toBe(true)

    expect(mockSaveScrutiny).toHaveBeenCalledWith(expect.objectContaining({
      student_id: 'stud-1',
      grades: expect.arrayContaining([
        { subject_id: 'sub-rel', final_grade: 0 }
      ])
    }))
  })
})
