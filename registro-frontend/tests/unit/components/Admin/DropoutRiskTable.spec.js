import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import DropoutRiskTable from '@/components/Admin/DropoutRiskTable.vue'

const mockApiGet = vi.fn()

vi.mock('@/services/api', () => ({
  default: {
    get: (...args) => mockApiGet(...args)
  }
}))

describe('DropoutRiskTable.vue - Early Warning Analytics', () => {
  let wrapper

  beforeEach(() => {
    vi.clearAllMocks()

    mockApiGet.mockImplementation((url) => {
      if (url === '/reports/dropout-risk') {
        return Promise.resolve({
          data: {
            count: 2,
            data: [
              {
                student_id: 's1-uuid-1234',
                first_name: 'Marco',
                last_name: 'Bianchi',
                class_id: 'c1',
                class_name: '3B',
                absence_rate: 26.5,
                total_hours: 100,
                absence_hours: 26,
                failing_subjects_count: 3,
                failing_subjects: ['Matematica (4.2)', 'Fisica (4.0)', 'Inglese (4.5)'],
                lates_count: 5,
                early_exits_count: 2,
                risk_score: 75.0,
                risk_level: 'Critico',
                recommended_action: 'Convocazione urgente famiglia e tutor dedicato'
              },
              {
                student_id: 's2-uuid-5678',
                first_name: 'Giulia',
                last_name: 'Neri',
                class_id: 'c1',
                class_name: '3B',
                absence_rate: 18.0,
                total_hours: 100,
                absence_hours: 18,
                failing_subjects_count: 1,
                failing_subjects: ['Latino (4.8)'],
                lates_count: 1,
                early_exits_count: 0,
                risk_score: 32.0,
                risk_level: 'Moderato',
                recommended_action: 'Colloquio con coordinatore'
              }
            ]
          }
        })
      }
      return Promise.resolve({ data: [] })
    })

    wrapper = mount(DropoutRiskTable, {
      global: {
        plugins: [
          [Quasar, {}],
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              classes: {
                classes: [{ id: 'c1', name: '3B', section: '' }]
              }
            },
            stubActions: true
          })
        ],
        stubs: {
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-select': true,
          'q-btn': true,
          'q-chip': true,
          'q-icon': true,
          'q-badge': true,
          'q-tooltip': true,
          'q-table': {
            props: ['rows', 'columns'],
            template: '<div class="q-table-stub"><div v-for="r in rows" :key="r.student_id" class="row-item">{{ r.last_name }} - {{ r.risk_level }} - {{ r.absence_rate }}%</div></div>'
          }
        }
      }
    })
  })

  it('renders table and calls dropout-risk API on mount', async () => {
    expect(wrapper.exists()).toBe(true)
    expect(mockApiGet).toHaveBeenCalledWith('/reports/dropout-risk', expect.any(Object))
  })

  it('populates rows with student risk data', async () => {
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.items.length).toBe(2)
    expect(wrapper.vm.countCritical).toBe(1)
    expect(wrapper.vm.countMedium).toBe(1)
    expect(wrapper.text()).toContain('Bianchi - Critico - 26.5%')
  })

  it('calls exportSupportPlan when requested', async () => {
    mockApiGet.mockResolvedValueOnce({
      data: 'ID Studente;Cognome;Nome\ns1;Bianchi;Marco'
    })

    // Mock URL.createObjectURL
    window.URL.createObjectURL = vi.fn().mockReturnValue('blob:mock-url')
    window.URL.revokeObjectURL = vi.fn()

    await wrapper.vm.exportSupportPlan()
    expect(mockApiGet).toHaveBeenCalledWith(
      '/reports/dropout-risk/export',
      expect.objectContaining({ responseType: 'blob' })
    )
  })
})
