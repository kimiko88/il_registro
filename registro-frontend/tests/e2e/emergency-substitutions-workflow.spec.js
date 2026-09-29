import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import EmergencySubstitutions from '@/pages/ata/EmergencySubstitutions.vue'
import substitutionService from '@/services/substitutionService'
import userService from '@/services/userService'
import schoolService from '@/services/schoolService'

vi.mock('@/services/substitutionService', () => ({
  default: {
    listBySchool: vi.fn(),
    todaySummary: vi.fn(),
    create: vi.fn(),
    assign: vi.fn(),
    recommendSubstitutes: vi.fn(),
    signRegister: vi.fn()
  }
}))

vi.mock('@/services/userService', () => ({
  default: {
    getUsers: vi.fn().mockResolvedValue({ data: { users: [] } })
  }
}))

vi.mock('@/services/schoolService', () => ({
  default: {
    getClasses: vi.fn().mockResolvedValue({ data: [] })
  }
}))

describe('Emergency Substitutions Workflow E2E', () => {
  let pinia

  const commonStubs = {
    'q-page': { template: '<div class="q-page"><slot /></div>' },
    'q-card': { template: '<div class="q-card"><slot /></div>' },
    'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
    'q-dialog': {
      props: ['modelValue'],
      template: '<div v-if="modelValue" class="q-dialog"><slot /></div>'
    },
    'q-table': {
      props: ['rows', 'columns'],
      template: '<div class="q-table" :data-count="rows?.length"><slot name="body" v-for="row in rows" :row="row" /><slot /></div>'
    },
    'q-btn': {
      props: ['label'],
      template: '<button class="q-btn">{{ label }}<slot /></button>'
    },
    'q-btn-group': { template: '<div class="q-btn-group"><slot /></div>' },
    'q-icon': true,
    'q-badge': {
      props: ['label'],
      template: '<span class="q-badge">{{ label }}<slot /></span>'
    },
    'q-chip': true,
    'q-input': {
      props: ['modelValue'],
      template: '<input class="q-input" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />'
    },
    'q-select': {
      template: '<div class="q-select"><slot /></div>'
    },
    'q-tooltip': true,
    'q-popup-proxy': true,
    'q-date': true
  }

  const mockSubs = [
    {
      id: 'sub-emergency-1',
      class_id: 'class-3A',
      class_name: '3A',
      absent_teacher_id: 't-absent-1',
      absent_teacher_name: 'Prof. Mario Rossi',
      date: '2026-10-15',
      slot: 1,
      hour: 1,
      subject: 'Matematica',
      status: 'pending',
      substitute_teacher_id: null,
      substitute_teacher_name: null
    },
    {
      id: 'sub-emergency-2',
      class_id: 'class-4B',
      class_name: '4B',
      absent_teacher_id: 't-absent-2',
      absent_teacher_name: 'Prof.ssa Anna Verdi',
      date: '2026-10-15',
      slot: 2,
      hour: 2,
      subject: 'Fisica',
      status: 'assigned',
      substitute_teacher_id: 't-sub-1',
      substitute_teacher_name: 'Prof. Luigi Neri'
    }
  ]

  const mockSummary = {
    total: 2,
    pending: 1,
    assigned: 1,
    confirmed: 0,
    cancelled: 0,
    substitutions: mockSubs
  }

  beforeEach(() => {
    vi.clearAllMocks()
    pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: { token: null, user: { role: 'vice_principal', permissions: ['manage_substitutions'] } }
      }
    })

    substitutionService.todaySummary.mockResolvedValue({ data: mockSummary })
    substitutionService.listBySchool.mockResolvedValue({ data: mockSubs })
  })

  it('loads today summary and displays statistics pills', async () => {
    const wrapper = mount(EmergencySubstitutions, {
      global: {
        plugins: [pinia, Quasar],
        stubs: commonStubs
      }
    })

    await flushPromises()

    expect(substitutionService.todaySummary).toHaveBeenCalled()
    expect(wrapper.text()).toContain('Emergenza Sostituzioni')

    // Check summary KPIs
    expect(wrapper.vm.summary.total).toBe(2)
    expect(wrapper.vm.summary.pending).toBe(1)
    expect(wrapper.vm.summary.assigned).toBe(1)
  })

  it('fetches AI/smart candidate recommendations and quick assigns substitute', async () => {
    const mockRecs = [
      {
        teacher_id: 't-best',
        teacher_name: 'Prof. Carlo Bianchi',
        score: 95,
        reason: 'Docente della stessa classe',
        teaches_class: true,
        teaches_subject: true
      },
      {
        teacher_id: 't-alt',
        teacher_name: 'Prof.ssa Elena Gialli',
        score: 70,
        reason: 'Insegna la stessa materia',
        teaches_class: false,
        teaches_subject: true
      }
    ]

    substitutionService.recommendSubstitutes.mockResolvedValue({ data: mockRecs })
    substitutionService.assign.mockResolvedValue({ data: { message: 'substitute assigned' } })

    const wrapper = mount(EmergencySubstitutions, {
      global: {
        plugins: [pinia, Quasar],
        stubs: commonStubs
      }
    })

    await flushPromises()

    // Trigger selectSubstitution on the pending substitution
    await wrapper.vm.selectSubstitution(mockSubs[0])
    await flushPromises()

    expect(substitutionService.recommendSubstitutes).toHaveBeenCalled()
    expect(wrapper.vm.recommendations.length).toBe(2)
    expect(wrapper.vm.recommendations[0].teacher_name).toBe('Prof. Carlo Bianchi')
    expect(wrapper.vm.recommendations[0].score).toBe(95)

    // Quick assign the recommended candidate
    await wrapper.vm.quickAssign(mockRecs[0])
    await flushPromises()

    expect(substitutionService.assign).toHaveBeenCalledWith('sub-emergency-1', expect.objectContaining({
      substitute_teacher_id: 't-best'
    }))
  })
})
