import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import EmergencySubstitutions from '@/pages/ata/EmergencySubstitutions.vue'
import VisitorRegistry from '@/pages/ata/VisitorRegistry.vue'
import Timecard from '@/pages/ata/Timecard.vue'
import PersonnelDesk from '@/pages/ata/PersonnelDesk.vue'

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
