import { mount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, beforeEach, vi } from 'vitest'
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

vi.mock('quasar', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useQuasar: () => ({
      dark: { isActive: false },
      notify: vi.fn(),
      lang: { current: 'it' }
    })
  }
})

describe('Secretary TimetableConstraints.vue — Timetable Constraints & Room Requirements', () => {
  let wrapper

  const mockReqs = [
    {
      id: 'req-1',
      subject_id: 'sub-info',
      subject_name: 'Informatica',
      required_room_type: 'lab_informatica',
      is_mandatory: true
    }
  ]

  const mockSubjects = [
    { id: 'sub-info', name: 'Informatica' },
    { id: 'sub-chim', name: 'Chimica' }
  ]

  const mockClassPlans = [
    {
      class_id: 'c-1',
      class_name: '1A',
      section: 'A',
      academic_year: '2024/2025',
      min_hours_per_day: 4,
      max_hours_per_day: 6,
      total_hours_week: 5,
      subjects: [
        { subject_id: 'sub-info', subject_name: 'Informatica', hours_per_week: 5 }
      ]
    }
  ]

  beforeEach(() => {
    vi.clearAllMocks()
    timetableGenService.getRoomRequirements.mockResolvedValue({ data: mockReqs })
    timetableGenService.getConstraints.mockResolvedValue({ data: [] })
    timetableGenService.getDesiderataWindow.mockResolvedValue({ data: { is_open: true } })
    timetableGenService.getPreferences.mockResolvedValue({ data: [] })
    timetableGenService.getAcademicYears.mockResolvedValue({ data: ['2024/2025', '2023/2024'] })
    timetableGenService.getClassesCurriculumPlans.mockResolvedValue({ data: mockClassPlans })
    timetableGenService.getClassCurriculumPlan.mockResolvedValue({ data: mockClassPlans[0] })
    timetableGenService.saveClassCurriculumPlan.mockResolvedValue({ data: mockClassPlans[0] })
    timetableGenService.inheritClassCurriculumPlan.mockResolvedValue({ data: mockClassPlans[0] })
    timetableGenService.inheritAllClassesCurriculumPlans.mockResolvedValue({ data: { classes_updated: 1, subjects_copied: 2 } })
    timetableGenService.getTeachersQuickPreferences.mockResolvedValue({
      data: {
        teachers: [
          {
            teacher_id: 't-1',
            teacher_name: 'Prof Rossi',
            subject_name: 'Matematica',
            day_off: 1,
            time_slot_pref: 'early_hours',
            max_hours_per_day: 5,
            preferred_hours_count: 6,
            unavailable_hours_count: 8
          }
        ],
        day_off_counts: { 1: 1 }
      }
    })
    timetableGenService.saveTeachersQuickPreferences.mockResolvedValue({ data: { message: 'ok' } })
    timetableGenService.saveTeacherQuickPreference.mockResolvedValue({ data: { message: 'ok' } })

    api.get.mockImplementation((url) => {
      if (url.includes('/subjects')) {
        return Promise.resolve({ data: mockSubjects })
      }
      return Promise.resolve({ data: [] })
    })

    wrapper = mount(TimetableConstraints, {
      global: {
        stubs: {
          'q-page': { template: '<div class="q-page"><slot /></div>' },
          'q-banner': { template: '<div class="q-banner"><slot name="avatar" /><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
          'q-btn': { template: '<button @click="$emit(\'click\')">{{ label }}<slot /></button>', props: ['label'] },
          'q-icon': true,
          'q-chip': { template: '<span><slot /></span>' },
          'q-badge': { template: '<span><slot /></span>' },
          'q-select': { template: '<div class="q-select"><slot /></div>' },
          'q-input': { template: '<div class="q-input"><slot /></div>' },
          'q-toggle': { template: '<div class="q-toggle"><slot /></div>' },
          'q-table': { template: '<div class="q-table"><slot /></div>' },
          'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
          'q-tab': { template: '<div class="q-tab"><slot /></div>' },
          'q-dialog': { template: '<div><slot /></div>' },
          'q-separator': true,
          'q-spinner-dots': true
        }
      }
    })
  })

  it('renders header and loads initial room requirements', async () => {
    await flushPromises()
    expect(wrapper.text()).toContain('Vincoli e Requisiti Orario Scolastico')
    expect(wrapper.text()).toContain('Assegnazione Aule Speciali / Laboratori per Materia')
    expect(timetableGenService.getRoomRequirements).toHaveBeenCalled()
    expect(wrapper.vm.roomReqs.length).toBe(1)
  })

  it('saves new room requirement calling timetableGenService.saveRoomRequirement', async () => {
    await flushPromises()
    timetableGenService.saveRoomRequirement.mockResolvedValue({ data: { id: 'req-2' } })

    wrapper.vm.openRoomReqModal()
    wrapper.vm.reqForm.subject_id = 'sub-chim'
    wrapper.vm.reqForm.required_room_type = 'lab_chimica'

    await wrapper.vm.saveReq()
    await flushPromises()

    expect(timetableGenService.saveRoomRequirement).toHaveBeenCalledWith(
      expect.objectContaining({
        subject_id: 'sub-chim',
        required_room_type: 'lab_chimica'
      })
    )
  })

  it('renders curriculum plans, daily limits, and allows editing hours', async () => {
    await flushPromises()
    wrapper.vm.currentTab = 'curriculum'
    await flushPromises()

    expect(wrapper.text()).toContain('Piano Orario Materie')
    expect(wrapper.vm.currentClassPlan).toBeDefined()
    expect(wrapper.vm.currentClassPlan.class_name).toBe('1A')
    expect(wrapper.vm.currentClassPlan.min_hours_per_day).toBe(4)
    expect(wrapper.vm.currentClassPlan.max_hours_per_day).toBe(6)
  })

  it('saves class curriculum plan and daily limits', async () => {
    await flushPromises()
    wrapper.vm.currentTab = 'curriculum'
    wrapper.vm.currentClassPlan = {
      class_id: 'c-1',
      class_name: '1A',
      min_hours_per_day: 5,
      max_hours_per_day: 6,
      subjects: [{ subject_id: 'sub-info', subject_name: 'Informatica', hours_per_week: 4 }]
    }

    timetableGenService.saveClassCurriculumPlan.mockResolvedValue({
      data: {
        class_id: 'c-1',
        class_name: '1A',
        min_hours_per_day: 5,
        max_hours_per_day: 6,
        total_hours_week: 4,
        subjects: [{ subject_id: 'sub-info', subject_name: 'Informatica', hours_per_week: 4 }]
      }
    })

    await wrapper.vm.saveCurrentClassPlan()
    await flushPromises()

    expect(timetableGenService.saveClassCurriculumPlan).toHaveBeenCalledWith(
      'c-1',
      expect.objectContaining({
        min_hours_per_day: 5,
        max_hours_per_day: 6
      })
    )
  })

  it('inherits class curriculum plan from previous academic year', async () => {
    await flushPromises()
    wrapper.vm.currentTab = 'curriculum'
    wrapper.vm.currentClassPlan = { class_id: 'c-1', class_name: '1A' }
    wrapper.vm.inheritSourceYear = '2023/2024'
    wrapper.vm.inheritScope = 'single'

    timetableGenService.inheritClassCurriculumPlan.mockResolvedValue({
      data: {
        class_id: 'c-1',
        class_name: '1A',
        min_hours_per_day: 4,
        max_hours_per_day: 6,
        total_hours_week: 6,
        subjects: [
          { subject_id: 'sub-info', subject_name: 'Informatica', hours_per_week: 4 },
          { subject_id: 'sub-chim', subject_name: 'Chimica', hours_per_week: 2 }
        ]
      }
    })

    await wrapper.vm.executeInherit()
    await flushPromises()

    expect(timetableGenService.inheritClassCurriculumPlan).toHaveBeenCalledWith(
      'c-1',
      expect.objectContaining({
        source_academic_year: '2023/2024'
      })
    )
  })

  it('renders tabular representation of teachers quick preferences and day-off distribution', async () => {
    await flushPromises()
    wrapper.vm.currentTab = 'desiderata'
    wrapper.vm.desiderataViewMode = 'quick_table'
    await flushPromises()

    expect(wrapper.text()).toContain('Rappresentazione Tabellare (Tutti i Docenti)')
    expect(wrapper.text()).toContain('Bilanciamento Giorni Liberi Richiesti')
    expect(timetableGenService.getTeachersQuickPreferences).toHaveBeenCalled()
    expect(wrapper.vm.quickPrefs.length).toBe(1)
    expect(wrapper.vm.quickPrefs[0].teacher_name).toBe('Prof Rossi')
    expect(wrapper.vm.quickPrefs[0].day_off).toBe(1)
    expect(wrapper.vm.quickPrefs[0].time_slot_pref).toBe('early_hours')
  })

  it('saves all quick preferences for teachers in batch', async () => {
    await flushPromises()
    wrapper.vm.currentTab = 'desiderata'
    wrapper.vm.desiderataViewMode = 'quick_table'
    await flushPromises()

    wrapper.vm.quickPrefs[0].day_off = 3
    wrapper.vm.quickPrefs[0].time_slot_pref = 'late_hours'

    await wrapper.vm.saveAllQuickPreferences()
    await flushPromises()

    expect(timetableGenService.saveTeachersQuickPreferences).toHaveBeenCalledWith(
      expect.objectContaining({
        preferences: expect.arrayContaining([
          expect.objectContaining({
            teacher_id: 't-1',
            day_off: 3,
            time_slot_pref: 'late_hours'
          })
        ])
      })
    )
  })

  it('saves single teacher quick preference', async () => {
    await flushPromises()
    wrapper.vm.currentTab = 'desiderata'
    wrapper.vm.desiderataViewMode = 'quick_table'
    await flushPromises()

    const row = wrapper.vm.quickPrefs[0]
    row.day_off = 2
    row.time_slot_pref = 'early_hours'

    await wrapper.vm.saveSingleRowQuickPreference(row)
    await flushPromises()

    expect(timetableGenService.saveTeacherQuickPreference).toHaveBeenCalledWith(
      't-1',
      expect.objectContaining({
        teacher_id: 't-1',
        day_off: 2,
        time_slot_pref: 'early_hours'
      })
    )
  })
})


