import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import PdpPlans from '@/pages/teacher/PdpPlans.vue'
import { pdpService } from '@/services/pdpService'
import api from '@/services/api'
import { useClassesStore } from '@/stores/classes'

vi.mock('@/services/pdpService', () => ({
  pdpService: {
    getClassPlans: vi.fn(),
    createPlan: vi.fn(),
    updatePlan: vi.fn(),
    deletePlan: vi.fn(),
    shareWithFamily: vi.fn(),
    exportPdpPdf: vi.fn()
  }
}))

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn()
  }
}))

describe('PDP / PEI Plans Workflow E2E', () => {
  let pinia

  const commonStubs = {
    'q-page': { template: '<div class="q-page"><slot /></div>' },
    'q-card': { template: '<div class="q-card"><slot /></div>' },
    'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
    'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
    'q-dialog': {
      props: ['modelValue'],
      template: '<div v-if="modelValue" class="q-dialog"><slot /></div>'
    },
    'q-form': {
      template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>'
    },
    'q-btn': {
      props: ['label'],
      template: '<button class="q-btn">{{ label }}<slot /></button>'
    },
    'q-icon': true,
    'q-avatar': true,
    'q-separator': true,
    'q-badge': true,
    'q-chip': {
      props: ['color'],
      template: '<span class="q-chip" :data-color="color"><slot /></span>'
    },
    'q-input': {
      props: ['modelValue'],
      template: '<input class="q-input" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />'
    },
    'q-select': {
      template: '<div class="q-select"><slot /></div>'
    },
    'q-file': {
      template: '<div class="q-file"><slot /></div>'
    },
    'CompensativeMeasuresSelector': {
      props: ['modelValue'],
      template: '<div class="compensative-measures-selector" />'
    },
    'SkeletonCard': true
  }

  const mockPlans = [
    {
      id: 'pdp-101',
      student_id: 'stu-1',
      student_name: 'Giovanni Alunni',
      class_id: 'class-1A',
      academic_year: '2025/2026',
      plan_type: 'pdp',
      diagnosis: 'Dislessia Evolutiva (F81.0)',
      shared_with_family: true,
      family_approved_at: '2026-09-20T10:00:00Z',
      content: {
        compensative: ['calc', 'extra_time'],
        dispensative: ['no_read_aloud'],
        notes: 'Uso mappe concettuali'
      }
    },
    {
      id: 'pei-102',
      student_id: 'stu-2',
      student_name: 'Marco Inclusione',
      class_id: 'class-1A',
      academic_year: '2025/2026',
      plan_type: 'pei',
      diagnosis: 'Disabilità certificata L. 104/92',
      shared_with_family: false,
      family_approved_at: null,
      content: {
        compensative: ['tutor'],
        dispensative: [],
        notes: 'Sostegno 18h'
      }
    }
  ]

  beforeEach(() => {
    vi.clearAllMocks()
    pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: { token: null, user: { role: 'teacher' } }
      }
    })

    const classesStore = useClassesStore()
    classesStore.classes = [
      { id: 'class-1A', name: 'Classe 1A' },
      { id: 'class-2B', name: 'Classe 2B' }
    ]

    pdpService.getClassPlans.mockResolvedValue({
      data: { plans: mockPlans }
    })

    api.get.mockResolvedValue({
      data: {
        users: [
          { id: 'stu-1', first_name: 'Giovanni', last_name: 'Alunni' },
          { id: 'stu-2', first_name: 'Marco', last_name: 'Inclusione' }
        ]
      }
    })
  })

  it('renders PDP/PEI plans for the selected class', async () => {
    const wrapper = mount(PdpPlans, {
      global: {
        plugins: [pinia, Quasar],
        stubs: commonStubs
      }
    })

    await flushPromises()

    expect(pdpService.getClassPlans).toHaveBeenCalledWith('class-1A', '2025/2026')
    expect(wrapper.text()).toContain('Giovanni Alunni')
    expect(wrapper.text()).toContain('Marco Inclusione')
    expect(wrapper.text()).toContain('Dislessia Evolutiva (F81.0)')

    // Verify plans array
    expect(wrapper.vm.plans.length).toBe(2)
  })

  it('opens dialog to create a new PDP plan and saves successfully', async () => {
    pdpService.createPlan.mockResolvedValue({
      data: { plan: { id: 'pdp-new', student_id: 'stu-1' } }
    })

    const wrapper = mount(PdpPlans, {
      global: {
        plugins: [pinia, Quasar],
        stubs: commonStubs
      }
    })

    await flushPromises()

    // Trigger openCreateDialog
    wrapper.vm.openCreateDialog()
    expect(wrapper.vm.showDialog).toBe(true)
    expect(wrapper.vm.isEditing).toBe(false)

    // Fill form
    wrapper.vm.form.student_id = 'stu-1'
    wrapper.vm.form.diagnosis = 'Disortografia (F81.1)'
    wrapper.vm.form.content.compensative = ['calc', 'pc']

    // Save
    await wrapper.vm.savePlan()
    await flushPromises()

    expect(pdpService.createPlan).toHaveBeenCalledWith(expect.objectContaining({
      student_id: 'stu-1',
      diagnosis: 'Disortografia (F81.1)',
      class_id: 'class-1A'
    }))

    expect(wrapper.vm.showDialog).toBe(false)
    expect(pdpService.getClassPlans).toHaveBeenCalledTimes(2) // Initial + refresh
  })

  it('shares plan with family', async () => {
    pdpService.shareWithFamily.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(PdpPlans, {
      global: {
        plugins: [pinia, Quasar],
        stubs: commonStubs
      }
    })

    await flushPromises()

    // Share unshared plan (mockPlans[1])
    await wrapper.vm.toggleShare(mockPlans[1])
    await flushPromises()

    expect(pdpService.shareWithFamily).toHaveBeenCalledWith('pei-102', true)
  })
})
