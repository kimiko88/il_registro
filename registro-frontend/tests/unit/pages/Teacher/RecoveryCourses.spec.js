import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import RecoveryCourses from '@/pages/teacher/RecoveryCourses.vue'
import recoveryService from '@/services/recoveryService'
import api from '@/services/api'

vi.mock('@/services/recoveryService', () => ({
  default: {
    listCourses: vi.fn(),
    listTests: vi.fn(),
    createCourse: vi.fn(),
    recordTest: vi.fn()
  }
}))

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn()
  }
}))

describe('RecoveryCourses.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    recoveryService.listCourses.mockResolvedValue({
      data: [
        { id: 'c-1', title: 'Recupero Matematica 2B', subject_name: 'Matematica', total_hours: 10, status: 'scheduled' }
      ]
    })
    recoveryService.listTests.mockResolvedValue({
      data: [
        { id: 't-1', student_name: 'Mario Rossi', class_name: '2B', grade: 6.5, outcome: 'recuperato', final_deliberation: 'Ammesso' }
      ]
    })
    api.get.mockImplementation((url) => {
      if (url.includes('/classes')) return Promise.resolve({ data: [{ id: 'cls-1', name: '2B' }] })
      if (url.includes('/subjects')) return Promise.resolve({ data: [{ id: 'sub-1', name: 'Matematica' }] })
      if (url.includes('/users')) return Promise.resolve({ data: [{ id: 's-1', first_name: 'Mario', last_name: 'Rossi' }] })
      return Promise.resolve({ data: [] })
    })
  })

  it('renders recovery courses and metrics correctly', async () => {
    const wrapper = mount(RecoveryCourses, {
      global: {
        plugins: [createTestingPinia()],
        mocks: {
          $t: (msg) => msg
        },
        stubs: {
          'q-page': { template: '<div class="q-page"><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
          'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
          'q-tab': { template: '<div class="q-tab"><slot /></div>' },
          'q-tab-panels': { template: '<div class="q-tab-panels"><slot /></div>' },
          'q-tab-panel': { template: '<div class="q-tab-panel"><slot /></div>' },
          'q-table': { template: '<div class="q-table"><slot /></div>' },
          'q-btn': { template: '<button class="q-btn"><slot /></button>' },
          'q-badge': { template: '<span class="q-badge"><slot /></span>' },
          'q-avatar': { template: '<div class="q-avatar"><slot /></div>' },
          'q-separator': { template: '<hr />' },
          'q-icon': true,
          'q-tooltip': true,
          'q-dialog': true,
          'q-form': true,
          'q-input': true,
          'q-select': true,
          'q-banner': true,
          'q-chip': true,
          'q-space': true
        }
      }
    })

    await flushPromises()
    expect(wrapper.exists()).toBe(true)
    expect(recoveryService.listCourses).toHaveBeenCalled()
    expect(recoveryService.listTests).toHaveBeenCalled()
  })
})
