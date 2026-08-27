import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import SupportRegister from '@/pages/teacher/SupportRegister.vue'
import supportService from '@/services/supportService'
import api from '@/services/api'

vi.mock('@/services/supportService', () => ({
  default: {
    listDiaryEntries: vi.fn(),
    listPeiGoals: vi.fn(),
    createDiaryEntry: vi.fn(),
    createPeiGoal: vi.fn(),
    updateGoalProgress: vi.fn()
  }
}))

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn()
  }
}))

describe('SupportRegister.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    supportService.listDiaryEntries.mockResolvedValue({
      data: [
        { id: 'd-1', student_name: 'Paolo Neri', entry_date: '2025-05-10', time_slot: '1ª Ora', activity_type: 'in_classe', topic_and_activities: 'Supporto Geometria' }
      ]
    })
    supportService.listPeiGoals.mockResolvedValue({
      data: [
        { id: 'g-1', student_name: 'Paolo Neri', pei_type: 'equipollente', axis: 'cognitiva', title: 'Geometria Solida', progress_status: 'intermedio' }
      ]
    })
    api.get.mockImplementation((url) => {
      if (url.includes('/classes')) return Promise.resolve({ data: [{ id: 'cls-1', name: '1A' }] })
      if (url.includes('/users')) return Promise.resolve({ data: [{ id: 's-1', first_name: 'Paolo', last_name: 'Neri' }] })
      return Promise.resolve({ data: [] })
    })
  })

  it('renders support register and tabs correctly', async () => {
    const wrapper = mount(SupportRegister, {
      global: {
        plugins: [createTestingPinia()],
        mocks: {
          $t: (msg) => msg
        },
        stubs: {
          'q-page': { template: '<div class="q-page"><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
          'q-tab': { template: '<div class="q-tab"><slot /></div>' },
          'q-tab-panels': { template: '<div class="q-tab-panels"><slot /></div>' },
          'q-tab-panel': { template: '<div class="q-tab-panel"><slot /></div>' },
          'q-table': { template: '<div class="q-table"><slot /></div>' },
          'q-btn': { template: '<button class="q-btn"><slot /></button>' },
          'q-badge': { template: '<span class="q-badge"><slot /></span>' },
          'q-chip': { template: '<span class="q-chip"><slot /></span>' },
          'q-separator': { template: '<hr />' },
          'q-icon': true,
          'q-tooltip': true,
          'q-dialog': true,
          'q-form': true,
          'q-input': true,
          'q-select': true,
          'q-checkbox': true,
          'q-space': true
        }
      }
    })

    await flushPromises()
    expect(wrapper.exists()).toBe(true)
    expect(supportService.listDiaryEntries).toHaveBeenCalled()
    expect(supportService.listPeiGoals).toHaveBeenCalled()
  })
})
