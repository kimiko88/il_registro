import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import Scrutiny from '@/pages/teacher/Scrutiny.vue'
import { createPinia, setActivePinia } from 'pinia'

vi.mock('@/services/scrutinyService', () => ({
  scrutinyService: {
    getMatrix: vi.fn(() => Promise.resolve({
      data: {
        subjects: [{ id: 'sub-math', name: 'Matematica' }],
        students: [
          {
            student_id: 'stud-1',
            student_name: 'Mario Rossi',
            subject_data: { 'sub-math': { average: 5.0 } },
            attendance_stats: { absences: 2, lates: 1, early_exits: 0 },
            record: { conduct_grade: 8, final_decision: 'Sospeso' }
          }
        ]
      }
    })),
    saveDeficiency: vi.fn(() => Promise.resolve({ data: { message: 'ok' } })),
    getStudentDeficiencies: vi.fn(() => Promise.resolve({
      data: [
        { id: 'def-1', subject_name: 'Matematica', topics: 'Equazioni di 2° grado', status: 'da_recuperare' }
      ]
    })),
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

describe('Teacher Scrutiny & Deficiency Component', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders scrutiny title correctly', () => {
    const wrapper = mount(Scrutiny, {
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
        }
      }
    })

    expect(wrapper.text()).toContain('Scrutinio Accademico & Differito')
  })
})
