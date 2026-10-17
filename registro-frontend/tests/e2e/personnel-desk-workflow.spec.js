import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import PersonnelDesk from '@/pages/ata/PersonnelDesk.vue'
import personnelDeskService from '@/services/personnelDeskService'

vi.mock('@/services/personnelDeskService', () => ({
  default: {
    listRequests: vi.fn(),
    createRequest: vi.fn(),
    submitRequest: vi.fn(),
    aaReview: vi.fn(),
    dsgaSign: vi.fn(),
    dsApprove: vi.fn(),
    deleteRequest: vi.fn()
  }
}))

describe('Personnel Desk Workflow E2E', () => {
  let pinia

  const commonStubs = {
    'q-page': { template: '<div class="q-page"><slot /></div>' },
    'q-card': { template: '<div class="q-card"><slot /></div>' },
    'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
    'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
    'q-table': {
      props: ['rows', 'columns'],
      template: '<div class="q-table" :data-count="rows?.length"><slot name="body" v-for="row in rows" :row="row" /><slot /></div>'
    },
    'q-dialog': {
      props: ['modelValue'],
      template: '<div v-if="modelValue" class="q-dialog"><slot /></div>'
    },
    'q-btn': {
      props: ['label'],
      template: '<button class="q-btn">{{ label }}<slot /></button>'
    },
    'q-icon': true,
    'q-badge': true,
    'q-chip': true,
    'q-input': {
      props: ['modelValue'],
      template: '<input class="q-input" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />'
    },
    'q-select': { template: '<div class="q-select"><slot /></div>' },
    'q-btn-toggle': true,
    'q-tooltip': true,
    'q-separator': true,
    'q-space': true
  }

  const mockRequestsList = [
    {
      id: 'req-001',
      applicant_id: 'emp-1',
      applicant_name: 'Giulia Bianchi',
      applicant_role: 'docente',
      category: 'permesso_breve',
      start_date: '2026-10-05',
      end_date: '2026-10-05',
      hours: 2,
      days: 0,
      description: 'Permesso breve 2 ore',
      status: 'submitted'
    }
  ]

  beforeEach(() => {
    vi.clearAllMocks()
    personnelDeskService.listRequests.mockResolvedValue({ data: mockRequestsList })

    pinia = createTestingPinia({
      initialState: {
        auth: {
          user: { id: 'admin-1', role: 'admin', name: 'Amministratore' },
          token: null
        }
      },
      stubActions: false
    })
  })

  it('1. loads and displays submitted personnel requests on mount', async () => {
    const wrapper = mount(PersonnelDesk, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k }
      }
    })

    await flushPromises()
    expect(personnelDeskService.listRequests).toHaveBeenCalled()
    expect(wrapper.vm.requests.length).toBe(1)
    expect(wrapper.vm.requests[0].description).toBe('Permesso breve 2 ore')
  })

  it('2. employee creates a new leave/permit digital request', async () => {
    const wrapper = mount(PersonnelDesk, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k }
      }
    })

    await flushPromises()
    wrapper.vm.openNewRequestDialog()
    expect(wrapper.vm.createDialog).toBe(true)

    wrapper.vm.createForm.category = 'ferie'
    wrapper.vm.createForm.start_date = '2026-10-15'
    wrapper.vm.createForm.end_date = '2026-10-18'
    wrapper.vm.createForm.days = 3
    wrapper.vm.createForm.description = 'Ferie autunnali concordate'
    wrapper.vm.createForm.submit_now = true

    personnelDeskService.createRequest.mockResolvedValue({ data: { id: 'req-002' } })

    await wrapper.vm.submitCreate()
    expect(personnelDeskService.createRequest).toHaveBeenCalledWith(expect.objectContaining({
      category: 'ferie',
      days: 3,
      description: 'Ferie autunnali concordate'
    }))
    expect(wrapper.vm.createDialog).toBe(false)
  })

  it('3. administrative assistant (AA) completes istruttoria review', async () => {
    personnelDeskService.aaReview.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(PersonnelDesk, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k }
      }
    })

    await flushPromises()
    const targetReq = wrapper.vm.requests[0]
    wrapper.vm.openAAReviewDialog(targetReq)

    expect(wrapper.vm.aaDialog).toBe(true)
    expect(wrapper.vm.activeReq).toEqual(targetReq)

    wrapper.vm.aaForm.note = 'Documentazione conforme e monte ore verificato'
    wrapper.vm.aaForm.approve = true

    await wrapper.vm.submitAAReview()
    expect(personnelDeskService.aaReview).toHaveBeenCalledWith('req-001', {
      note: 'Documentazione conforme e monte ore verificato',
      approve: true
    })
    expect(wrapper.vm.aaDialog).toBe(false)
  })

  it('4. DSGA affixes visto contabile and financial endorsement', async () => {
    personnelDeskService.dsgaSign.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(PersonnelDesk, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k }
      }
    })

    await flushPromises()
    const targetReq = wrapper.vm.requests[0]
    wrapper.vm.openDSGASignDialog(targetReq)

    expect(wrapper.vm.dsgaDialog).toBe(true)

    wrapper.vm.dsgaForm.note = 'Capitolo di spesa capiente'
    wrapper.vm.dsgaForm.approve = true

    await wrapper.vm.submitDSGASign()
    expect(personnelDeskService.dsgaSign).toHaveBeenCalledWith('req-001', {
      note: 'Capitolo di spesa capiente',
      approve: true
    })
    expect(wrapper.vm.dsgaDialog).toBe(false)
  })

  it('5. school principal (DS) issues final decree approval', async () => {
    personnelDeskService.dsApprove.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(PersonnelDesk, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k }
      }
    })

    await flushPromises()
    const targetReq = wrapper.vm.requests[0]
    wrapper.vm.openDSApproveDialog(targetReq)

    expect(wrapper.vm.dsDialog).toBe(true)
    wrapper.vm.dsForm.decree_num = 'DECR-2026/105'
    wrapper.vm.dsForm.note = 'Si decreta l\'autorizzazione richiesta.'
    wrapper.vm.dsForm.approve = true

    await wrapper.vm.submitDSApprove()
    expect(personnelDeskService.dsApprove).toHaveBeenCalledWith('req-001', {
      decree_num: 'DECR-2026/105',
      note: 'Si decreta l\'autorizzazione richiesta.',
      approve: true
    })
    expect(wrapper.vm.dsDialog).toBe(false)
  })
})
