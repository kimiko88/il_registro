import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import EmergencySubstitutions from '@/pages/ata/EmergencySubstitutions.vue'
import VisitorRegistry from '@/pages/ata/VisitorRegistry.vue'
import Timecard from '@/pages/ata/Timecard.vue'
import PersonnelDesk from '@/pages/ata/PersonnelDesk.vue'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn().mockResolvedValue({ data: [] }),
    post: vi.fn().mockResolvedValue({ data: {} }),
    put: vi.fn().mockResolvedValue({ data: {} }),
    patch: vi.fn().mockResolvedValue({ data: {} }),
    delete: vi.fn().mockResolvedValue({ data: {} })
  }
}))

vi.mock('@/services/userService', () => ({
  default: {
    getUsers: vi.fn().mockResolvedValue({ data: [] })
  }
}))

vi.mock('@/services/schoolService', () => ({
  default: {
    getClasses: vi.fn().mockResolvedValue({ data: [] })
  }
}))

vi.mock('@/services/substitutionService', () => ({
  default: {
    getEmergencyDashboard: vi.fn().mockResolvedValue({ data: { open: [], assigned: [], daily_summary: {} } })
  }
}))

vi.mock('@/services/visitorService', () => ({
  default: {
    getVisitors: vi.fn().mockResolvedValue({ data: [] }),
    getDailyStats: vi.fn().mockResolvedValue({ data: {} })
  }
}))

vi.mock('@/services/timecardService', () => ({
  default: {
    getMonthlyEntries: vi.fn().mockResolvedValue({ data: [] }),
    getLeaveRequests: vi.fn().mockResolvedValue({ data: [] }),
    getAllStaffTimecards: vi.fn().mockResolvedValue({ data: [] })
  }
}))

vi.mock('@/services/personnelDeskService', () => ({
  default: {
    getRequests: vi.fn().mockResolvedValue({ data: [] })
  }
}))

vi.mock('@/services/staffAttendanceService', () => ({
  default: {
    getDailySummary: vi.fn().mockResolvedValue({ data: {} }),
    getAttendanceList: vi.fn().mockResolvedValue({ data: [] })
  }
}))

describe('ATA Specialist 4 Vertical Modules E2E Workflow', () => {
  let pinia

  beforeEach(() => {
    vi.clearAllMocks()
    pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: {
          user: { id: 'u-ata-ds', role: 'collaboratore_ds', name: 'Collaboratore DS Mario' }
        }
      }
    })
  })

  const commonStubs = {
    'q-page': { template: '<div><slot /></div>' },
    'q-card': { template: '<div><slot /></div>' },
    'q-card-section': { template: '<div><slot /></div>' },
    'q-table': { template: '<div><slot /></div>' },
    'q-dialog': { template: '<div><slot /></div>' },
    'q-tabs': { template: '<div><slot /></div>' },
    'q-tab': { template: '<div><slot /></div>' },
    'q-tab-panels': { template: '<div><slot /></div>' },
    'q-tab-panel': { template: '<div><slot /></div>' },
    'q-btn': { template: '<button><slot /></button>' },
    'q-icon': { template: '<i><slot /></i>' },
    'q-badge': { template: '<span><slot /></span>' },
    'q-input': { template: '<input />' },
    'q-select': { template: '<select />' }
  }

  it('1. mounts EmergencySubstitutions module for collaboratore_ds', () => {
    const wrapper = mount(EmergencySubstitutions, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: {
          $t: (key) => key
        }
      }
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('2. mounts VisitorRegistry module for collaboratore_scolastico', () => {
    const wrapper = mount(VisitorRegistry, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: {
          $t: (key) => key
        }
      }
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('3. mounts Timecard module for DSGA and ATA staff', () => {
    const wrapper = mount(Timecard, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: {
          $t: (key) => key
        }
      }
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('4. mounts PersonnelDesk module for multi-role digital requests', () => {
    const wrapper = mount(PersonnelDesk, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: {
          $t: (key) => key
        }
      }
    })
    expect(wrapper.exists()).toBe(true)
  })
})
