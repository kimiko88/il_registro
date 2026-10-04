import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import StaffAttendance from '@/pages/ata/StaffAttendance.vue'
import staffAttendanceService from '@/services/staffAttendanceService'

vi.mock('@/services/staffAttendanceService', () => ({
  default: {
    getDailySummary: vi.fn(),
    getList: vi.fn(),
    setStrikeMode: vi.fn(),
    recordAttendance: vi.fn(),
    registerBadgeSwipe: vi.fn(),
    processBadgeSwipes: vi.fn()
  }
}))

describe('Staff Attendance & Strike Mode Workflow E2E', () => {
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
    'q-btn-group': { template: '<div class="q-btn-group"><slot /></div>' },
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
    'q-select': {
      props: ['modelValue'],
      template: '<div class="q-select"><slot /></div>'
    },
    'q-avatar': true,
    'q-tooltip': true,
    'q-separator': true,
    'q-popup-proxy': true,
    'q-date': true
  }

  const mockSummary = {
    date: '2026-10-01',
    is_strike_day: false,
    teachers: { present: 18, total: 20, on_strike: 0 },
    ata: { present: 8, total: 10, on_strike: 0 },
    total_strike: 0,
    by_role: [
      { role: 'teacher', count: 20, present: 18 },
      { role: 'collaboratore_scolastico', count: 6, present: 5 },
      { role: 'assistente_amministrativo', count: 4, present: 3 }
    ]
  }

  const mockStaffList = [
    {
      id: 'rec-1',
      user_id: 'usr-1',
      first_name: 'Giovanni',
      last_name: 'Bianchi',
      email: 'giovanni.bianchi@scuola.it',
      role: 'teacher',
      status: 'present',
      is_strike_recorded: false,
      badge_entry_time: '2026-10-01T08:00:00Z',
      badge_exit_time: null,
      notes: ''
    },
    {
      id: 'rec-2',
      user_id: 'usr-2',
      first_name: 'Maria',
      last_name: 'Verdi',
      email: 'maria.verdi@scuola.it',
      role: 'collaboratore_scolastico',
      status: 'absent',
      is_strike_recorded: false,
      badge_entry_time: null,
      badge_exit_time: null,
      notes: 'Assenza ingiustificata'
    },
    {
      id: 'rec-3',
      user_id: 'usr-3',
      first_name: 'Carla',
      last_name: 'Neri',
      email: 'carla.neri@scuola.it',
      role: 'assistente_amministrativo',
      status: 'sick_leave',
      is_strike_recorded: false,
      badge_entry_time: null,
      badge_exit_time: null,
      notes: 'Certificato medico'
    }
  ]

  beforeEach(() => {
    vi.clearAllMocks()
    staffAttendanceService.getDailySummary.mockResolvedValue(mockSummary)
    staffAttendanceService.getList.mockResolvedValue({ data: mockStaffList })
    staffAttendanceService.setStrikeMode.mockResolvedValue({ success: true })
    staffAttendanceService.recordAttendance.mockResolvedValue({ success: true })
    staffAttendanceService.registerBadgeSwipe.mockResolvedValue({ success: true })
    staffAttendanceService.processBadgeSwipes.mockResolvedValue({ success: true })

    pinia = createTestingPinia({
      initialState: {
        auth: {
          user: { id: 'dsga-1', role: 'dsga', first_name: 'Lucia', last_name: 'Fontana', email: 'lucia.fontana@scuola.it' },
          token: null
        }
      },
      stubActions: false
    })
  })

  it('1. loads daily summary KPIs and staff attendance list on mount for DSGA', async () => {
    const wrapper = mount(StaffAttendance, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: {
          t: (k, params) => {
            if (k === 'staffAttendance.presenceRate') return `${params?.rate}%`
            return k
          },
          te: () => false
        }
      }
    })

    await flushPromises()

    expect(staffAttendanceService.getDailySummary).toHaveBeenCalledWith(expect.any(String))
    expect(staffAttendanceService.getList).toHaveBeenCalledWith(expect.any(String))
    expect(wrapper.vm.summaryData).toBeDefined()
    expect(wrapper.vm.staffList.length).toBe(3)
    expect(wrapper.vm.canWriteAttendance).toBe(true)
    expect(wrapper.vm.canActivateStrike).toBe(true)
    expect(wrapper.vm.canManageStrike).toBe(true)
    expect(wrapper.vm.isStrikeMode).toBe(false)
  })

  it('2. filters staff by category role and attendance status correctly', async () => {
    const wrapper = mount(StaffAttendance, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k, te: () => false }
      }
    })

    await flushPromises()

    // Filter by role teacher
    wrapper.vm.filterRole = 'teacher'
    expect(wrapper.vm.filteredStaffList.length).toBe(1)
    expect(wrapper.vm.filteredStaffList[0].last_name).toBe('Bianchi')

    // Filter by role collaboratore_scolastico
    wrapper.vm.filterRole = 'collaboratore_scolastico'
    expect(wrapper.vm.filteredStaffList.length).toBe(1)
    expect(wrapper.vm.filteredStaffList[0].last_name).toBe('Verdi')

    // Reset role filter, filter by status 'leave' (sick_leave, permit, mission)
    wrapper.vm.filterRole = 'all'
    wrapper.vm.filterStatus = 'leave'
    expect(wrapper.vm.filteredStaffList.length).toBe(1)
    expect(wrapper.vm.filteredStaffList[0].status).toBe('sick_leave')

    // Text search filter
    wrapper.vm.filterStatus = 'all'
    wrapper.vm.filterSearch = 'Carla'
    expect(wrapper.vm.filteredStaffList.length).toBe(1)
    expect(wrapper.vm.filteredStaffList[0].first_name).toBe('Carla')
  })

  it('3. navigates between days with changeDate and returns with goToToday', async () => {
    const wrapper = mount(StaffAttendance, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k, te: () => false }
      }
    })

    await flushPromises()

    const initialDate = wrapper.vm.selectedDate

    // Go back 1 day
    wrapper.vm.changeDate(-1)
    await flushPromises()
    expect(wrapper.vm.selectedDate).not.toBe(initialDate)
    expect(staffAttendanceService.getList).toHaveBeenCalledTimes(2)

    // Go to today
    wrapper.vm.goToToday()
    await flushPromises()
    const todayStr = new Date().toISOString().substring(0, 10)
    expect(wrapper.vm.selectedDate).toBe(todayStr)
  })

  it('4. activates and deactivates strike mode lifecycle seamlessly', async () => {
    const wrapper = mount(StaffAttendance, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k, te: () => false }
      }
    })

    await flushPromises()
    expect(wrapper.vm.isStrikeMode).toBe(false)

    // Activate strike mode: loadData will fetch updated summary
    staffAttendanceService.getDailySummary.mockResolvedValueOnce({
      ...mockSummary,
      is_strike_day: true
    })
    await wrapper.vm.toggleStrikeMode()
    expect(staffAttendanceService.setStrikeMode).toHaveBeenCalledWith({
      date: wrapper.vm.selectedDate,
      is_strike_day: true
    })
    expect(wrapper.vm.isStrikeMode).toBe(true)

    // Deactivate strike mode: loadData will fetch updated summary
    staffAttendanceService.getDailySummary.mockResolvedValueOnce({
      ...mockSummary,
      is_strike_day: false
    })
    await wrapper.vm.toggleStrikeMode()
    expect(staffAttendanceService.setStrikeMode).toHaveBeenCalledWith({
      date: wrapper.vm.selectedDate,
      is_strike_day: false
    })
    expect(wrapper.vm.isStrikeMode).toBe(false)
  })

  it('5. opens edit attendance dialog and saves attendance status with strike details', async () => {
    const wrapper = mount(StaffAttendance, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k, te: () => false }
      }
    })

    await flushPromises()

    // Open edit dialog for second row (Maria Verdi)
    const targetRow = wrapper.vm.staffList[1]
    wrapper.vm.openEditDialog(targetRow)
    expect(wrapper.vm.showEditDialog).toBe(true)
    expect(wrapper.vm.editingRecord.user_id).toBe('usr-2')

    // Change status to on_strike
    wrapper.vm.editForm.status = 'on_strike'
    wrapper.vm.editForm.strike_code = 'FLC-CGIL-2026'
    wrapper.vm.editForm.notes = 'Adesione sciopero nazionale'

    await wrapper.vm.saveAttendance()
    expect(staffAttendanceService.recordAttendance).toHaveBeenCalledWith(expect.objectContaining({
      user_id: 'usr-2',
      date: wrapper.vm.selectedDate,
      status: 'on_strike',
      strike_code: 'FLC-CGIL-2026',
      notes: 'Adesione sciopero nazionale'
    }))
    expect(wrapper.vm.showEditDialog).toBe(false)
  })

  it('6. simulates badge swipe recording from hardware terminal', async () => {
    const wrapper = mount(StaffAttendance, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k, te: () => false }
      }
    })

    await flushPromises()

    // Open swipe dialog
    wrapper.vm.openBadgeSwipeDialog()
    expect(wrapper.vm.showBadgeSwipeDialog).toBe(true)

    wrapper.vm.badgeForm.badge_code = 'BDG-ATA-007'
    wrapper.vm.badgeForm.swipe_type = 'in'
    wrapper.vm.badgeForm.swipe_time = '07:45'
    wrapper.vm.badgeForm.device_id = 'PORTINERIA-NORD'

    await wrapper.vm.submitBadgeSwipe()
    expect(staffAttendanceService.registerBadgeSwipe).toHaveBeenCalledWith(expect.objectContaining({
      badge_code: 'BDG-ATA-007',
      swipe_type: 'in',
      device_id: 'PORTINERIA-NORD'
    }))
    expect(staffAttendanceService.processBadgeSwipes).toHaveBeenCalled()
    expect(wrapper.vm.showBadgeSwipeDialog).toBe(false)
  })

  it('7. restricts editing and strike management for non-authorized roles (teacher & student)', async () => {
    const teacherPinia = createTestingPinia({
      initialState: {
        auth: {
          user: { id: 't-1', role: 'teacher', first_name: 'Paolo', last_name: 'Docente' },
          token: null
        }
      },
      stubActions: false
    })

    const wrapper = mount(StaffAttendance, {
      global: {
        plugins: [Quasar, teacherPinia],
        stubs: commonStubs,
        mocks: { t: (k) => k, te: () => false }
      }
    })

    await flushPromises()

    // Teachers cannot write attendance nor activate strike
    expect(wrapper.vm.canWriteAttendance).toBe(false)
    expect(wrapper.vm.canActivateStrike).toBe(false)
    expect(wrapper.vm.canManageStrike).toBe(false)

    // Columns should not include 'actions' column
    const columnNames = wrapper.vm.columns.map(c => c.name)
    expect(columnNames).not.toContain('actions')

    // Calling openEditDialog should be a no-op
    wrapper.vm.openEditDialog(mockStaffList[0])
    expect(wrapper.vm.showEditDialog).toBe(false)

    // Calling toggleStrikeMode should be a no-op
    await wrapper.vm.toggleStrikeMode()
    expect(staffAttendanceService.setStrikeMode).not.toHaveBeenCalled()
  })
})
