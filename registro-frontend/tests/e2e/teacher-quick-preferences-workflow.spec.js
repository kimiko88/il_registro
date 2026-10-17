import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import TimetableConstraints from '@/pages/secretary/TimetableConstraints.vue'
import timetableGenService from '@/services/timetableGenService'
import api from '@/services/api'

vi.mock('@/services/timetableGenService', () => ({
  default: {
    getRoomRequirements: vi.fn(),
    saveRoomRequirement: vi.fn(),
    deleteRoomRequirement: vi.fn(),
    getConstraints: vi.fn(),
    saveConstraint: vi.fn(),
    getDesiderataWindow: vi.fn(),
    setDesiderataWindow: vi.fn(),
    getPreferences: vi.fn(),
    savePreferences: vi.fn(),
    getAcademicYears: vi.fn(),
    getClassesCurriculumPlans: vi.fn(),
    getClassCurriculumPlan: vi.fn(),
    saveClassCurriculumPlan: vi.fn(),
    inheritClassCurriculumPlan: vi.fn(),
    inheritAllClassesCurriculumPlans: vi.fn(),
    getTeachersQuickPreferences: vi.fn(),
    saveTeachersQuickPreferences: vi.fn(),
    saveTeacherQuickPreference: vi.fn()
  }
}))

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn()
  }
}))

const mockNotify = vi.fn()
vi.mock('quasar', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useQuasar: () => ({
      dark: { isActive: false },
      notify: mockNotify,
      lang: { current: 'it' }
    })
  }
})

describe('E2E Workflow: Teacher Quick Preferences & Day-Off Balancing (Desiderata Docenti)', () => {
  let wrapper
  let piniaSecretary

  const mockTeachersQuickPrefs = [
    {
      teacher_id: 't-rossi',
      teacher_name: 'Mario Rossi',
      subject_name: 'Matematica',
      day_off: 1, // Lunedì
      time_slot_pref: 'early_hours',
      max_hours_per_day: 5,
      preferred_hours_count: 4,
      unavailable_hours_count: 6
    },
    {
      teacher_id: 't-bianchi',
      teacher_name: 'Giulia Bianchi',
      subject_name: 'Italiano',
      day_off: 1, // Lunedì
      time_slot_pref: 'late_hours',
      max_hours_per_day: 4,
      preferred_hours_count: 2,
      unavailable_hours_count: 8
    },
    {
      teacher_id: 't-verdi',
      teacher_name: 'Luca Verdi',
      subject_name: 'Scienze',
      day_off: 3, // Mercoledì
      time_slot_pref: 'none',
      max_hours_per_day: 6,
      preferred_hours_count: 0,
      unavailable_hours_count: 0
    }
  ]

  const mockDayOffCounts = {
    1: 2, // 2 docenti il Lunedì
    3: 1  // 1 docente il Mercoledì
  }

  beforeEach(async () => {
    vi.clearAllMocks()

    piniaSecretary = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: {
          user: { role: 'secretary', id: 'sec-1', school_id: 'school-1' },
          isAuthenticated: true
        }
      }
    })

    timetableGenService.getRoomRequirements.mockResolvedValue({ data: [] })
    timetableGenService.getConstraints.mockResolvedValue({ data: [] })
    timetableGenService.getDesiderataWindow.mockResolvedValue({ data: { is_open: true } })
    timetableGenService.getPreferences.mockResolvedValue({ data: [] })
    timetableGenService.getAcademicYears.mockResolvedValue({ data: ['2024/2025'] })
    timetableGenService.getClassesCurriculumPlans.mockResolvedValue({ data: [] })
    timetableGenService.getTeachersQuickPreferences.mockResolvedValue({
      data: {
        teachers: JSON.parse(JSON.stringify(mockTeachersQuickPrefs)),
        day_off_counts: mockDayOffCounts
      }
    })
    timetableGenService.saveTeachersQuickPreferences.mockResolvedValue({ data: { message: 'ok' } })
    timetableGenService.saveTeacherQuickPreference.mockResolvedValue({ data: { message: 'ok' } })

    api.get.mockResolvedValue({ data: [] })

    wrapper = mount(TimetableConstraints, {
      global: {
        plugins: [piniaSecretary],
        mocks: {
          t: (key) => key
        },
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-banner': { template: '<div class="q-banner"><slot /></div>' },
          'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
          'q-tab': { template: '<button class="q-tab"><slot /></button>' },
          'q-btn': {
            template: '<button class="q-btn" :disabled="loading" @click="$emit(\'click\')"><slot />{{ label }}</button>',
            props: ['label', 'loading']
          },
          'q-icon': true,
          'q-badge': { template: '<span class="q-badge"><slot /></span>' },
          'q-chip': {
            template: '<span class="q-chip" @click="$emit(\'click\')"><slot /></span>'
          },
          'q-table': {
            template: '<div class="q-table-stub"><slot name="body-cell-actions" :row="rows[0]" v-if="rows && rows.length" /></div>',
            props: ['rows']
          },
          'q-btn-dropdown': { template: '<div class="q-btn-dropdown-stub"><slot /></div>' },
          'q-list': { template: '<div class="q-list"><slot /></div>' },
          'q-item': { template: '<div class="q-item"><slot /></div>' },
          'q-item-section': { template: '<div class="q-item-section"><slot /></div>' },
          'q-select': { template: '<div class="q-select-stub"><slot /></div>' },
          'q-input': { template: '<div class="q-input-stub"><slot /></div>' },
          'q-toggle': true,
          'q-separator': true,
          'q-dialog': true,
          'q-tooltip': true
        }
      }
    })

    await flushPromises()
  })

  it('1. Loads Desiderata Docenti tab and fetches teachers quick preferences with day-off distribution', async () => {
    wrapper.vm.currentTab = 'desiderata'
    wrapper.vm.desiderataViewMode = 'quick_table'
    await flushPromises()

    expect(wrapper.vm.desiderataViewMode).toBe('quick_table')
    expect(timetableGenService.getTeachersQuickPreferences).toHaveBeenCalled()
    expect(wrapper.vm.quickPrefs).toHaveLength(3)
    expect(wrapper.vm.quickPrefsDayOffCounts[1]).toBe(2)
    expect(wrapper.vm.quickPrefsDayOffCounts[3]).toBe(1)
  })

  it('2. Filters teachers correctly using the day-off KPI distribution chip', async () => {
    wrapper.vm.currentTab = 'desiderata'
    wrapper.vm.desiderataViewMode = 'quick_table'
    await flushPromises()

    // No filter active: all 3 teachers returned
    expect(wrapper.vm.filteredQuickPrefs).toHaveLength(3)

    // Filter by Lunedì (1)
    wrapper.vm.filterDayOff = 1
    expect(wrapper.vm.filteredQuickPrefs).toHaveLength(2)
    expect(wrapper.vm.filteredQuickPrefs.map(t => t.teacher_id)).toEqual(['t-rossi', 't-bianchi'])

    // Filter by Mercoledì (3)
    wrapper.vm.filterDayOff = 3
    expect(wrapper.vm.filteredQuickPrefs).toHaveLength(1)
    expect(wrapper.vm.filteredQuickPrefs[0].teacher_id).toBe('t-verdi')

    // Reset filter
    wrapper.vm.filterDayOff = null
    expect(wrapper.vm.filteredQuickPrefs).toHaveLength(3)
  })

  it('3. Filters teachers using the text search query (name and subject)', async () => {
    wrapper.vm.currentTab = 'desiderata'
    wrapper.vm.desiderataViewMode = 'quick_table'
    await flushPromises()

    // Search by teacher name
    wrapper.vm.quickPrefsSearch = 'Bianchi'
    expect(wrapper.vm.filteredQuickPrefs).toHaveLength(1)
    expect(wrapper.vm.filteredQuickPrefs[0].teacher_name).toBe('Giulia Bianchi')

    // Search by subject name
    wrapper.vm.quickPrefsSearch = 'Scienze'
    expect(wrapper.vm.filteredQuickPrefs).toHaveLength(1)
    expect(wrapper.vm.filteredQuickPrefs[0].subject_name).toBe('Scienze')

    // Empty search restores all
    wrapper.vm.quickPrefsSearch = ''
    expect(wrapper.vm.filteredQuickPrefs).toHaveLength(3)
  })

  it('4. Saves single teacher preference and triggers notification and refresh', async () => {
    wrapper.vm.currentTab = 'desiderata'
    wrapper.vm.desiderataViewMode = 'quick_table'
    await flushPromises()

    const targetTeacher = wrapper.vm.quickPrefs[0] // Mario Rossi
    targetTeacher.day_off = 5 // Spostato a Venerdì
    targetTeacher.time_slot_pref = 'late_hours'

    await wrapper.vm.saveSingleRowQuickPreference(targetTeacher)

    expect(timetableGenService.saveTeacherQuickPreference).toHaveBeenCalledWith(
      't-rossi',
      expect.objectContaining({
        teacher_id: 't-rossi',
        day_off: 5,
        time_slot_pref: 'late_hours'
      })
    )
    expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({ type: 'positive' }))
  })

  it('5. Saves all teachers quick preferences in bulk', async () => {
    wrapper.vm.currentTab = 'desiderata'
    wrapper.vm.desiderataViewMode = 'quick_table'
    await flushPromises()

    await wrapper.vm.saveAllQuickPreferences()

    expect(timetableGenService.saveTeachersQuickPreferences).toHaveBeenCalledWith({
      preferences: expect.arrayContaining([
        expect.objectContaining({ teacher_id: 't-rossi' }),
        expect.objectContaining({ teacher_id: 't-bianchi' }),
        expect.objectContaining({ teacher_id: 't-verdi' })
      ])
    })
    expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({ type: 'positive' }))
  })

  it('6. Seamlessly switches view mode between quick table and single teacher matrix', async () => {
    wrapper.vm.currentTab = 'desiderata'
    wrapper.vm.desiderataViewMode = 'quick_table'
    await flushPromises()

    expect(wrapper.vm.desiderataViewMode).toBe('quick_table')

    // Switch to single teacher matrix mode
    wrapper.vm.desiderataViewMode = 'single_matrix'
    expect(wrapper.vm.desiderataViewMode).toBe('single_matrix')

    // Open teacher matrix from quick table helper
    wrapper.vm.openMatrixForTeacher('t-bianchi')
    expect(wrapper.vm.selectedDesiderataTeacherId).toBe('t-bianchi')
    expect(wrapper.vm.desiderataViewMode).toBe('single_matrix')
  })
})
