import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import AuditLog from '@/pages/admin/AuditLog.vue'
import adminService from '@/services/adminService'
import { createTestingPinia } from '@pinia/testing'

// Mock Quasar
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            notify: vi.fn(),
            dialog: vi.fn()
        })
    }
})

// Mock Service
vi.mock('@/services/adminService', () => ({
    default: {
        getAuditLogs: vi.fn()
    }
}))

describe('AuditLog', () => {
    let wrapper

    beforeEach(() => {
        wrapper = mount(AuditLog, {
            global: {
                plugins: [createTestingPinia()],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-select': true,
                    'q-btn': true,
                    'q-table': true, // Simplified
                    'q-tr': true,
                    'q-td': true,
                    'q-chip': true
                }
            }
        })
    })

    it('fetches logs on mount', () => {
        expect(adminService.getAuditLogs).toHaveBeenCalled()
    })

    it('loads data correctly', async () => {
        adminService.getAuditLogs.mockResolvedValue({
            data: { items: [{ id: 1, action_type: 'create' }], total: 1 }
        })
        await wrapper.vm.fetchLogs()
        expect(wrapper.vm.logs.length).toBe(1)
        expect(wrapper.vm.loading).toBe(false)
    })

    it('handles pagination request', () => {
        // We can't spy on internal fetchLogs easily, but we can check if service is called again with new params
        adminService.getAuditLogs.mockClear()
        wrapper.vm.onRequest({ pagination: { page: 2, rowsPerPage: 20 } })
        expect(wrapper.vm.pagination.page).toBe(2)
        expect(adminService.getAuditLogs).toHaveBeenCalledWith(expect.objectContaining({ page: 2 }))
    })

    it('gets correct action color', () => {
        expect(wrapper.vm.getActionColor('create')).toBe('positive')
        expect(wrapper.vm.getActionColor('delete')).toBe('negative')
    })
})
