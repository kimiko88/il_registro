import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import Timecard from '@/pages/ata/Timecard.vue'
import staffAttendanceService from '@/services/staffAttendanceService'
import userService from '@/services/userService'

vi.mock('@/services/staffAttendanceService', () => ({
  default: {
    getTimecard: vi.fn(),
    listLeaves: vi.fn(),
    createLeave: vi.fn(),
    approveLeave: vi.fn(),
    rejectLeave: vi.fn(),
    deleteLeave: vi.fn(),
    exportTimecardCsv: vi.fn()
  }
}))

vi.mock('@/services/userService', () => ({
  default: {
    getUsers: vi.fn().mockResolvedValue({ users: [] }),
    getAll: vi.fn().mockResolvedValue([])
  }
}))

describe('Timecard and Leave Workflow E2E', () => {
  let pinia

  const commonStubs = {
    'q-page': { template: '<div class="q-page"><slot /></div>' },
    'q-card': { template: '<div class="q-card"><slot /></div>' },
    'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
    'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
    'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
    'q-tab': { template: '<div class="q-tab"><slot /></div>' },
    'q-tab-panels': { template: '<div class="q-tab-panels"><slot /></div>' },
    'q-tab-panel': { template: '<div class="q-tab-panel"><slot /></div>' },
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

  const mockTimecard = {
    user_id: 'user-ata-1',
    user_name: 'Mario Rossi',
    role: 'collaboratore_scolastico',
    badge_code: 'BADGE-7788',
    month: '2026-10',
    contract_hours: 156.0,
    worked_hours: 164.5,
    overtime_hours: 8.5,
    absence_days: 0,
    leave_days: 2,
    sick_days: 0,
    permit_hours: 0.0,
    daily_entries: [
      {
        date: '2026-10-01',
        entry_time: '2026-10-01T07:30:00Z',
        exit_time: '2026-10-01T14:42:00Z',
        status: 'present',
        notes: ''
      }
    ]
  }

  const mockLeaves = [
    {
      id: 'leave-101',
      user_id: 'user-ata-1',
      user_name: 'Mario Rossi',
      type: 'ferie',
      start_date: '2026-10-20',
      end_date: '2026-10-22',
      days: 3,
      hours: 0,
      status: 'pending',
      notes: 'Richiesta ferie autunnali'
    }
  ]

  beforeEach(() => {
    vi.clearAllMocks()
    pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: {
          token: null,
          user: {
            id: 'user-ata-1',
            role: 'collaboratore_scolastico',
            badge_code: 'BADGE-7788'
          }
        }
      }
    })

    staffAttendanceService.getTimecard.mockResolvedValue(mockTimecard)
    staffAttendanceService.listLeaves.mockResolvedValue(mockLeaves)
  })

  it('renders timecard stats and monthly worked hours', async () => {
    const wrapper = mount(Timecard, {
      global: {
        plugins: [pinia, Quasar],
        stubs: commonStubs
      }
    })

    await flushPromises()

    expect(staffAttendanceService.getTimecard).toHaveBeenCalled()
    expect(wrapper.text()).toContain('Cartellino & Piano Ferie')

    // Badge code and timecard data
    expect(wrapper.vm.badgeCode).toBe('BADGE-7788')
    expect(wrapper.vm.timecardData.worked_hours).toBe(164.5)
    expect(wrapper.vm.overtimeBalance).toBe(8.5)
  })

  it('submits a new leave request dialog and refreshes the leaves list', async () => {
    staffAttendanceService.createLeave.mockResolvedValue({ id: 'leave-102', status: 'pending' })

    const wrapper = mount(Timecard, {
      global: {
        plugins: [pinia, Quasar],
        stubs: commonStubs
      }
    })

    await flushPromises()

    // Open leave request dialog
    wrapper.vm.openNewLeaveDialog()
    expect(wrapper.vm.leaveDialog).toBe(true)

    // Fill form
    wrapper.vm.leaveForm.type = 'ferie'
    wrapper.vm.leaveForm.start_date = '2026-11-02'
    wrapper.vm.leaveForm.end_date = '2026-11-04'
    wrapper.vm.leaveForm.days = 3
    wrapper.vm.leaveForm.notes = 'Ponte Ognissanti'

    // Submit leave request
    await wrapper.vm.submitLeave()
    await flushPromises()

    expect(staffAttendanceService.createLeave).toHaveBeenCalledWith(expect.objectContaining({
      type: 'ferie',
      start_date: '2026-11-02',
      end_date: '2026-11-04',
      days: 3
    }))

    expect(wrapper.vm.leaveDialog).toBe(false)
    expect(staffAttendanceService.listLeaves).toHaveBeenCalled()
  })

  it('allows DSGA/Admin to approve a pending leave request', async () => {
    const adminPinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: {
          token: null,
          user: {
            id: 'dsga-1',
            role: 'dsga'
          }
        }
      }
    })

    staffAttendanceService.approveLeave.mockResolvedValue({ message: 'richiesta approvata' })

    const wrapper = mount(Timecard, {
      global: {
        plugins: [adminPinia, Quasar],
        stubs: commonStubs
      }
    })

    await flushPromises()

    expect(wrapper.vm.isDSGAOrAdmin).toBe(true)

    // Approve the pending leave
    await wrapper.vm.approveLeave(mockLeaves[0])
    await flushPromises()

    expect(staffAttendanceService.approveLeave).toHaveBeenCalledWith('leave-101', expect.objectContaining({
      notes: expect.any(String)
    }))
  })
})
