import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import AuditLog from '@/pages/admin/AuditLog.vue'
import adminService from '@/services/adminService'

vi.mock('@/services/adminService', () => ({
  default: {
    getAuditLogs: vi.fn()
  }
}))

vi.mock('@/composables/useTableExport', () => ({
  useTableExport: () => ({
    exportTableCsv: vi.fn()
  })
}))

describe('AuditLog.vue - Advanced Filters & Compliance Export', () => {
  let wrapper

  beforeEach(() => {
    vi.clearAllMocks()

    adminService.getAuditLogs.mockResolvedValue({
      data: {
        total: 1,
        items: [
          {
            id: 'log-1',
            admin_id: 'adm-001',
            admin_name: 'DPO Rossi',
            action_type: 'update',
            target: 'user',
            target_id: 'u-123',
            school_name: 'Liceo Fermi',
            details: 'Modifica permessi utente',
            created_at: '2026-09-06T12:00:00Z'
          }
        ]
      }
    })

    wrapper = mount(AuditLog, {
      global: {
        plugins: [[Quasar, {}]],
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-select': true,
          'q-input': true,
          'q-btn': true,
          'q-table': true,
          'q-icon': true,
          'q-chip': true
        }
      }
    })
  })

  it('fetches logs on mount with default pagination', async () => {
    expect(wrapper.exists()).toBe(true)
    expect(adminService.getAuditLogs).toHaveBeenCalledWith(
      expect.objectContaining({ page: 1, page_size: 20 })
    )
  })

  it('sends actor, from_date and to_date when filters are applied', async () => {
    wrapper.vm.filters.actor = 'adm-001'
    wrapper.vm.filters.action = 'update'
    wrapper.vm.filters.fromDate = '2026-09-01'
    wrapper.vm.filters.toDate = '2026-09-06'

    await wrapper.vm.fetchLogs()

    expect(adminService.getAuditLogs).toHaveBeenCalledWith(
      expect.objectContaining({
        admin_id: 'adm-001',
        action: 'update',
        from_date: '2026-09-01',
        to_date: '2026-09-06'
      })
    )
  })

  it('resets filters correctly', async () => {
    wrapper.vm.filters.actor = 'adm-001'
    wrapper.vm.filters.action = 'update'
    wrapper.vm.filters.fromDate = '2026-09-01'
    wrapper.vm.filters.toDate = '2026-09-06'

    wrapper.vm.resetFilters()

    expect(wrapper.vm.filters.actor).toBe('')
    expect(wrapper.vm.filters.action).toBeNull()
    expect(wrapper.vm.filters.fromDate).toBe('')
    expect(wrapper.vm.filters.toDate).toBe('')
  })
})
