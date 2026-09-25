import { mount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import SchedulePreferences from '@/pages/teacher/SchedulePreferences.vue'
import timetableGenService from '@/services/timetableGenService'

vi.mock('@/services/timetableGenService', () => ({
  default: {
    getPreferences: vi.fn(),
    savePreferences: vi.fn()
  }
}))

vi.mock('quasar', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useQuasar: () => ({
      dark: { isActive: false },
      notify: vi.fn()
    })
  }
})

describe('Teacher SchedulePreferences.vue — Teacher Desiderata Matrix', () => {
  let wrapper

  beforeEach(() => {
    vi.clearAllMocks()
    timetableGenService.getPreferences.mockResolvedValue({
      data: [
        { day_of_week: 1, hour_index: 1, preference_type: 'preferred' },
        { day_of_week: 5, hour_index: 8, preference_type: 'unavailable' }
      ]
    })

    wrapper = mount(SchedulePreferences, {
      global: {
        stubs: {
          'q-page': { template: '<div class="q-page"><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-banner': { template: '<div class="q-banner"><slot /></div>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-icon': true,
          'q-chip': { template: '<span><slot /></span>' },
          'q-spinner-dots': true
        }
      }
    })
  })

  it('renders header, seniority alert banner, and loads existing preferences', async () => {
    await flushPromises()
    expect(wrapper.text()).toContain('Desiderata Orario Scolastico')
    expect(wrapper.text()).toContain('Informativa Generazione Orario Scolastico')
    expect(timetableGenService.getPreferences).toHaveBeenCalled()
    expect(wrapper.vm.grid[1][1]).toBe('preferred')
    expect(wrapper.vm.grid[5][8]).toBe('unavailable')
  })

  it('cycles cell preference: neutral -> preferred -> unavailable -> neutral', async () => {
    await flushPromises()
    // Tuesday (2), hour 3 is initially neutral
    expect(wrapper.vm.grid[2][3]).toBe('neutral')

    wrapper.vm.cyclePreference(2, 3)
    expect(wrapper.vm.grid[2][3]).toBe('preferred')

    wrapper.vm.cyclePreference(2, 3)
    expect(wrapper.vm.grid[2][3]).toBe('unavailable')

    wrapper.vm.cyclePreference(2, 3)
    expect(wrapper.vm.grid[2][3]).toBe('neutral')
  })

  it('saves preferences dispatching to timetableGenService.savePreferences', async () => {
    await flushPromises()
    timetableGenService.savePreferences.mockResolvedValue({ data: { success: true } })

    // Change slot Wednesday (3), hour 2 to preferred
    wrapper.vm.grid[3][2] = 'preferred'

    await wrapper.vm.savePreferences()
    await flushPromises()

    expect(timetableGenService.savePreferences).toHaveBeenCalledWith({
      preferences: expect.arrayContaining([
        expect.objectContaining({ day_of_week: 3, hour_index: 2, preference_type: 'preferred' })
      ])
    })
  })
})
