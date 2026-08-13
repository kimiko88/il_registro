import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { scrutinyService } from '@/services/scrutinyService'
import StudentReportCard from '@/pages/Student/ReportCard.vue'
import ParentReportCard from '@/pages/Parent/ReportCard.vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn()
  }
}))

describe('Deficiencies & Deferred Scrutiny Test Suite', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('scrutinyService calls correct API endpoints for deficiencies', async () => {
    api.post.mockResolvedValueOnce({ data: { message: 'ok' } })
    api.get.mockResolvedValueOnce({ data: [{ id: 'def-1' }] })
    api.get.mockResolvedValueOnce({ data: [{ id: 'def-2' }] })
    api.post.mockResolvedValueOnce({ data: { message: 'deferred ok' } })

    await scrutinyService.saveDeficiency({ student_id: 's1', topics: 'Equazioni' })
    expect(api.post).toHaveBeenCalledWith('/scrutiny/deficiencies', { student_id: 's1', topics: 'Equazioni' })

    await scrutinyService.getStudentDeficiencies('s1')
    expect(api.get).toHaveBeenCalledWith('/scrutiny/deficiencies/student/s1')

    await scrutinyService.getClassDeficiencies('c1', 1)
    expect(api.get).toHaveBeenCalledWith('/scrutiny/deficiencies/class/c1', { params: { semester: 1 } })

    await scrutinyService.saveDeferredScrutiny({ student_id: 's1', final_decision: 'promosso_con_debiti_saldati' })
    expect(api.post).toHaveBeenCalledWith('/scrutiny/deferred', { student_id: 's1', final_decision: 'promosso_con_debiti_saldati' })
  })

  it('renders student deficiency topics card in Student ReportCard view', async () => {
    const authStore = useAuthStore()
    authStore.user = { id: 's1', first_name: 'Mario', last_name: 'Rossi' }

    api.get.mockImplementation((url) => {
      if (url.includes('/scrutiny/deficiencies/student/')) {
        return Promise.resolve({
          data: [
            {
              id: 'def-101',
              subject_name: 'Matematica',
              topics: 'Equazioni di 2° grado e Trigonometria',
              recovery_mode: 'corso_recupero',
              status: 'da_recuperare',
              recovery_grade: null
            }
          ]
        })
      }
      if (url.includes('/attendance/my-attendance/summary')) {
        return Promise.resolve({ data: { absences: 3, lates: 1, early_exits: 0 } })
      }
      return Promise.resolve({ data: {} })
    })

    const wrapper = mount(StudentReportCard, {
      global: {
        stubs: {
          'q-page': { template: '<div class="q-page"><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-avatar': true,
          'q-badge': { template: '<span><slot /></span>' },
          'q-btn': true,
          'q-btn-toggle': true,
          'q-table': true,
          'q-td': true,
          'q-chip': true,
          'q-icon': true
        }
      }
    })

    await flushPromises()

    expect(wrapper.text()).toContain('Carenze Formative & Argomenti da Recuperare')
    expect(wrapper.text()).toContain('Matematica')
    expect(wrapper.text()).toContain('Equazioni di 2° grado e Trigonometria')
    expect(wrapper.text()).toContain('DA RECUPERARE')
  })

  it('renders child deficiency topics card in Parent ReportCard view', async () => {
    api.get.mockImplementation((url) => {
      if (url.includes('/users/me/children')) {
        return Promise.resolve({ data: [{ id: 'child-1', first_name: 'Marco', last_name: 'Rossi' }] })
      }
      if (url.includes('/grades/child-grades/child-1')) {
        return Promise.resolve({
          data: {
            student_name: 'Marco Rossi',
            class_name: '3A',
            subject_grades: [{ subject: 'Fisica', final_grade: 5 }]
          }
        })
      }
      if (url.includes('/scrutiny/deficiencies/student/child-1')) {
        return Promise.resolve({
          data: [
            {
              id: 'def-102',
              subject_name: 'Fisica',
              topics: 'Cinematica e Vettori',
              recovery_mode: 'studio_individuale',
              status: 'da_recuperare',
              recovery_grade: null
            }
          ]
        })
      }
      return Promise.resolve({ data: {} })
    })

    const wrapper = mount(ParentReportCard, {
      global: {
        stubs: {
          'q-page': { template: '<div class="q-page"><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-avatar': true,
          'q-badge': { template: '<span><slot /></span>' },
          'q-btn': true,
          'q-btn-toggle': true,
          'q-select': true,
          'q-banner': true,
          'q-table': true,
          'q-td': true,
          'q-chip': true,
          'q-icon': true
        }
      }
    })

    await flushPromises()

    expect(wrapper.text()).toContain('Carenze Formative & Argomenti da Recuperare (Figlio/a)')
    expect(wrapper.text()).toContain('Fisica')
    expect(wrapper.text()).toContain('Cinematica e Vettori')
  })
})
