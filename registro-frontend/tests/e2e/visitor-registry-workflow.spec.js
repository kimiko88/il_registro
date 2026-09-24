import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import VisitorRegistry from '@/pages/ata/VisitorRegistry.vue'
import visitorService from '@/services/visitorService'

vi.mock('@/services/visitorService', () => ({
  default: {
    listTodayVisitors: vi.fn(),
    registerVisitor: vi.fn(),
    recordVisitorExit: vi.fn(),
    listTodayEarlyExits: vi.fn(),
    recordEarlyExit: vi.fn(),
    recordStudentReturn: vi.fn(),
    listMaintenanceReports: vi.fn(),
    createMaintenanceReport: vi.fn(),
    updateMaintenanceStatus: vi.fn()
  }
}))

vi.mock('@/services/userService', () => ({
  default: {
    getUsers: vi.fn().mockResolvedValue({ data: [] })
  }
}))

describe('Visitor Registry Workflow E2E', () => {
  let pinia

  const commonStubs = {
    'q-page': { template: '<div class="q-page"><slot /></div>' },
    'q-card': { template: '<div class="q-card"><slot /></div>' },
    'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
    'q-table': {
      props: ['rows', 'columns'],
      template: '<div class="q-table" :data-count="rows?.length"><slot name="body" v-for="row in rows" :row="row" /><slot /></div>'
    },
    'q-dialog': {
      props: ['modelValue'],
      template: '<div v-if="modelValue" class="q-dialog"><slot /></div>'
    },
    'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
    'q-tab': { template: '<div class="q-tab"><slot /></div>' },
    'q-tab-panels': { template: '<div class="q-tab-panels"><slot /></div>' },
    'q-tab-panel': { template: '<div class="q-tab-panel"><slot /></div>' },
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
    'q-avatar': true,
    'q-tooltip': true,
    'q-separator': true,
    'q-space': true
  }

  const mockVisitors = [
    {
      id: 'v-101',
      name: 'Mario Rossi',
      document_id: 'CI-12345',
      purpose: 'parent',
      host_name: 'Prof. Verdi',
      badge_number: 'B-01',
      entry_time: '2026-09-24T09:00:00Z',
      exit_time: null
    }
  ]

  const mockEarlyExits = [
    {
      id: 'exit-101',
      student_id: 's-1',
      student_name: 'Luigi Rossi',
      class_name: '3A',
      delegatee_name: 'Mario Rossi',
      delegate_rel: 'Genitore',
      exit_time: '2026-09-24T11:30:00Z',
      return_time: null,
      reason_code: 'visita_medica'
    }
  ]

  const mockMaintenance = [
    {
      id: 'm-101',
      location: 'Aula Magna',
      category: 'elettrico',
      description: 'Presa guasta',
      priority: 'alta',
      status: 'aperto'
    }
  ]

  beforeEach(() => {
    vi.clearAllMocks()
    visitorService.listTodayVisitors.mockResolvedValue({ data: mockVisitors })
    visitorService.listTodayEarlyExits.mockResolvedValue({ data: mockEarlyExits })
    visitorService.listMaintenanceReports.mockResolvedValue({ data: mockMaintenance })

    pinia = createTestingPinia({
      initialState: {
        auth: {
          user: { id: 'staff-1', role: 'collaboratore_scolastico', name: 'Custode Mario' },
          token: null
        }
      },
      stubActions: false
    })
  })

  it('1. loads active visitors and displays badges on initial mount', async () => {
    const wrapper = mount(VisitorRegistry, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k }
      }
    })

    await flushPromises()
    expect(visitorService.listTodayVisitors).toHaveBeenCalled()
    expect(wrapper.vm.visitorsList.length).toBe(1)
    expect(wrapper.vm.currentVisitors.length).toBe(1)
    expect(wrapper.vm.visitorsList[0].name).toBe('Mario Rossi')
  })

  it('2. registers incoming visitor with badge assignment', async () => {
    const wrapper = mount(VisitorRegistry, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k }
      }
    })

    await flushPromises()
    wrapper.vm.openNewVisitorDialog()
    expect(wrapper.vm.visitorDialog).toBe(true)

    wrapper.vm.visitorForm.name = 'Giuseppe Verdi'
    wrapper.vm.visitorForm.document_id = 'PAT-9988'
    wrapper.vm.visitorForm.purpose = 'supplier'
    wrapper.vm.visitorForm.host_name = 'Segreteria'
    wrapper.vm.visitorForm.badge_number = 'B-02'

    visitorService.registerVisitor.mockResolvedValue({ data: { id: 'v-102' } })

    await wrapper.vm.submitVisitor()
    expect(visitorService.registerVisitor).toHaveBeenCalledWith(expect.objectContaining({
      name: 'Giuseppe Verdi',
      badge_number: 'B-02'
    }))
    expect(wrapper.vm.visitorDialog).toBe(false)
  })

  it('3. records visitor checkout and frees badge', async () => {
    visitorService.recordVisitorExit.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(VisitorRegistry, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k }
      }
    })

    await flushPromises()
    await wrapper.vm.recordExit(wrapper.vm.visitorsList[0])

    expect(visitorService.recordVisitorExit).toHaveBeenCalledWith('v-101', expect.any(String))
  })

  it('4. records student early exit and subsequent return', async () => {
    visitorService.recordStudentReturn.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(VisitorRegistry, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k }
      }
    })

    wrapper.vm.activeTab = 'early_exits'
    await wrapper.vm.loadEarlyExits()

    expect(visitorService.listTodayEarlyExits).toHaveBeenCalled()
    expect(wrapper.vm.earlyExitsList.length).toBe(1)

    await wrapper.vm.recordReturn(wrapper.vm.earlyExitsList[0])
    expect(visitorService.recordStudentReturn).toHaveBeenCalledWith('exit-101', expect.any(String))
  })

  it('5. updates maintenance report status', async () => {
    visitorService.updateMaintenanceStatus.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(VisitorRegistry, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k }
      }
    })

    wrapper.vm.activeTab = 'maintenance'
    await wrapper.vm.loadMaintenance()

    expect(visitorService.listMaintenanceReports).toHaveBeenCalled()
    expect(wrapper.vm.maintenanceList.length).toBe(1)

    await wrapper.vm.updateMaintStatus(wrapper.vm.maintenanceList[0], 'in_lavorazione')
    expect(visitorService.updateMaintenanceStatus).toHaveBeenCalledWith('m-101', {
      status: 'in_lavorazione'
    })
  })
})
