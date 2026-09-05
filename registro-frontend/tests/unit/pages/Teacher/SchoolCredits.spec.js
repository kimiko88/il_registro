import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import SchoolCredits from '@/pages/teacher/SchoolCredits.vue'
import creditService from '@/services/creditService'
import api from '@/services/api'

vi.mock('@/services/creditService', () => ({
  default: {
    listClassCredits: vi.fn(),
    calculateSuggested: vi.fn(),
    assignCredit: vi.fn(),
    getStudentSummary: vi.fn()
  }
}))

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn()
  }
}))

describe('SchoolCredits.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    creditService.listClassCredits.mockResolvedValue({
      data: [
        { id: 'sc-1', student_name: 'Luigi Bianchi', grade_level: 5, grade_average: 8.7, conduct_grade: 9, base_credit_range_min: 12, base_credit_range_max: 13, assigned_credit: 13 }
      ]
    })
    creditService.calculateSuggested.mockResolvedValue({
      data: {
        base_credit_range_min: 12,
        base_credit_range_max: 13,
        suggested_credit: 13,
        motivation: 'Attribuzione fascia massima D.Lgs 62/2017'
      }
    })
    api.get.mockImplementation((url) => {
      if (url.includes('/classes')) return Promise.resolve({ data: [{ id: 'cls-5A', name: '5A' }] })
      if (url.includes('/users')) return Promise.resolve({ data: [{ id: 's-1', first_name: 'Luigi', last_name: 'Bianchi' }] })
      return Promise.resolve({ data: [] })
    })
  })

  it('renders school credits page and loads class credits', async () => {
    const wrapper = mount(SchoolCredits, {
      global: {
        plugins: [createTestingPinia()],
        mocks: {
          $t: (msg) => msg
        },
        stubs: {
          'q-page': { template: '<div class="q-page"><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-table': { template: '<div class="q-table"><slot /></div>' },
          'q-expansion-item': { template: '<div class="q-expansion-item"><slot /></div>' },
          'q-markup-table': { template: '<table class="q-markup-table"><slot /></table>' },
          'q-btn': { template: '<button class="q-btn"><slot /></button>' },
          'q-badge': { template: '<span class="q-badge"><slot /></span>' },
          'q-chip': { template: '<span class="q-chip"><slot /></span>' },
          'q-icon': true,
          'q-tooltip': true,
          'q-dialog': true,
          'q-form': true,
          'q-input': true,
          'q-select': { template: '<div class="q-select-stub"><slot /></div>' },
          'q-checkbox': true,
          'q-space': true,
          'q-skeleton': true
        }
      }
    })

    await flushPromises()
    expect(wrapper.exists()).toBe(true)
    expect(creditService.listClassCredits).toHaveBeenCalled()
  })
})
