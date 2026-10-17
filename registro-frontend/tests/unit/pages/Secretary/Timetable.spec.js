import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import Timetable from '@/pages/secretary/Timetable.vue'
import timetableGenService from '@/services/timetableGenService'

vi.mock('@/services/timetableGenService', () => ({
  default: {
    startGeneration: vi.fn(),
    getJobStatus: vi.fn(),
    publishSchedule: vi.fn(),
    adjustSchedule: vi.fn(),
    listConstraints: vi.fn().mockResolvedValue([]),
    listRoomRequirements: vi.fn().mockResolvedValue([]),
  }
}))

vi.mock('@/services/timetablesService', () => ({
  default: {
    getClassSchedule: vi.fn().mockResolvedValue([]),
    getTeacherSchedule: vi.fn().mockResolvedValue([]),
    updateClassSchedule: vi.fn().mockResolvedValue({ success: true }),
    updateTeacherSchedule: vi.fn().mockResolvedValue({ success: true }),
  }
}))

vi.mock('@/services/classesService', () => ({
  default: {
    getClasses: vi.fn().mockResolvedValue([
      { id: 'class-1', name: '1A', year: 1, section: 'A' },
      { id: 'class-2', name: '2A', year: 2, section: 'A' }
    ])
  }
}))

vi.mock('@/services/teachersService', () => ({
  default: {
    getTeachers: vi.fn().mockResolvedValue([
      { id: 'teacher-1', first_name: 'Mario', last_name: 'Rossi' }
    ])
  }
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key) => key,
    locale: { value: 'it' }
  }),
  createI18n: () => ({
    global: {
      t: (key) => key,
      locale: { value: 'it' }
    }
  })
}))

describe('Secretary Timetable Page — Multi-Alternative Generation & Management', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  const mountComponent = () => {
    return mount(Timetable, {
      global: {
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-btn-toggle': { template: '<div><slot /></div>' },
          'q-select': { template: '<div class="q-select"><slot /></div>' },
          'q-icon': { template: '<i></i>' },
          'q-badge': { template: '<span><slot /></span>' },
          'q-banner': { template: '<div><slot /></div>' },
          'q-dialog': { template: '<div v-if="modelValue"><slot /></div>', props: ['modelValue'] },
          'q-toolbar': { template: '<div><slot /></div>' },
          'q-toolbar-title': { template: '<div><slot /></div>' },
          'q-linear-progress': { template: '<div></div>' }
        },
        mocks: {
          $q: {
            dark: { isActive: false },
            notify: vi.fn(),
            dialog: vi.fn(() => ({
              onOk: (cb) => { cb(); return { onCancel: vi.fn(), onDismiss: vi.fn() } },
              onCancel: vi.fn(),
              onDismiss: vi.fn()
            }))
          }
        }
      }
    })
  }

  it('renders correctly with default class view mode', () => {
    const wrapper = mountComponent()
    expect(wrapper.exists()).toBe(true)
    expect(wrapper.vm.viewMode).toBe('class')
  })

  it('allows toggling between class schedule and teacher schedule view modes', async () => {
    const wrapper = mountComponent()
    wrapper.vm.viewMode = 'teacher'
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.viewMode).toBe('teacher')
  })

  it('opens automatic generation dialog and initializes options', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.showGenerateDialog).toBe(false)
    wrapper.vm.openGenerateDialog()
    expect(wrapper.vm.showGenerateDialog).toBe(true)
  })

  it('handles 3 alternative schedules returned from generator', async () => {
    const wrapper = mountComponent()

    const mockAlternatives = [
      {
        id: 1,
        label: 'Proposta 1: Bilanciata',
        strategy: 'balanced',
        description: 'Equilibrio ottimale',
        coverage_pct: 100,
        assigned_slots: 30,
        total_slots: 30,
        slots: [{ id: 's1', day_of_week: 1, hour_index: 1, subject_name: 'Italiano' }],
        hard_conflicts: [],
        score: 850
      },
      {
        id: 2,
        label: 'Proposta 2: Didattica & Prime Ore',
        strategy: 'didactic_first',
        description: 'Priorità apprendimento',
        coverage_pct: 100,
        assigned_slots: 30,
        total_slots: 30,
        slots: [{ id: 's2', day_of_week: 1, hour_index: 1, subject_name: 'Matematica' }],
        hard_conflicts: [],
        score: 920
      },
      {
        id: 3,
        label: 'Proposta 3: Compatta',
        strategy: 'compact_teacher',
        description: 'Minimizza buchi docenti',
        coverage_pct: 96.6,
        assigned_slots: 29,
        total_slots: 30,
        slots: [{ id: 's3', day_of_week: 1, hour_index: 2, subject_name: 'Fisica' }],
        hard_conflicts: [],
        score: 890
      }
    ]

    wrapper.vm.generationResult = {
      job_id: 'job-123',
      total_slots: 30,
      assigned_slots: 30,
      coverage_pct: 100,
      duration_ms: 120,
      alternatives: mockAlternatives
    }
    wrapper.vm.selectedAlternativeId = 1
    await wrapper.vm.$nextTick()

    // 1. Initial selection is Alternative 1
    expect(wrapper.vm.selectedAlternativeId).toBe(1)
    expect(wrapper.vm.activeAlternativeStats.label).toBe('Proposta 1: Bilanciata')

    // 2. Select Alternative 2 (Didattica)
    wrapper.vm.selectAlternative(mockAlternatives[1])
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.selectedAlternativeId).toBe(2)
    expect(wrapper.vm.activeAlternativeStats.strategy).toBe('didactic_first')

    // 3. Select Alternative 3 (Compatta)
    wrapper.vm.selectAlternative(mockAlternatives[2])
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.selectedAlternativeId).toBe(3)
    expect(wrapper.vm.activeAlternativeStats.strategy).toBe('compact_teacher')
  })

  it('publishes chosen alternative passing alternative_id to timetableGenService', async () => {
    timetableGenService.publishSchedule.mockResolvedValueOnce({ success: true })
    const wrapper = mountComponent()

    wrapper.vm.generationJobId = 'job-xyz-456'
    wrapper.vm.selectedAlternativeId = 2
    wrapper.vm.generationResult = {
      job_id: 'job-xyz-456',
      alternatives: [
        { id: 1, slots: [] },
        { id: 2, slots: [{ day_of_week: 1, hour_index: 1 }] }
      ]
    }

    await wrapper.vm.publishGeneratedSchedule()
    expect(timetableGenService.publishSchedule).toHaveBeenCalledWith(
      'job-xyz-456',
      expect.objectContaining({ alternative_id: 2 })
    )
  })

  it('loads working slots into adjust modal and tracks manual slot adjustments', async () => {
    const wrapper = mountComponent()

    const sampleSlots = [
      {
        id: 'slot-1',
        class_id: 'class-1',
        class_name: '1A',
        day_of_week: 1,
        hour_index: 1,
        subject_id: 'sub-mat',
        subject_name: 'Matematica',
        teacher_id: 'teacher-1',
        teacher_name: 'Prof Rossi'
      }
    ]

    wrapper.vm.generationResult = {
      job_id: 'job-adjust-1',
      slots: sampleSlots,
      alternatives: [{ id: 1, slots: sampleSlots }]
    }
    wrapper.vm.selectedAlternativeId = 1

    wrapper.vm.openAdjustModal()
    expect(wrapper.vm.showAdjustModal).toBe(true)
    expect(wrapper.vm.workingSlots.length).toBe(1)
    expect(wrapper.vm.workingSlots[0].subject_name).toBe('Matematica')
  })
})
