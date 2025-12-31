
import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import AuditLog from '@/pages/Admin/AuditLog.vue'
import adminService from '@/services/adminService'

// Mock adminService
vi.mock('@/services/adminService', () => ({
    default: {
        getAuditLogs: vi.fn()
    }
}))

describe('AuditLog.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        adminService.getAuditLogs.mockResolvedValue({
            data: {
                items: [
                    {
                        id: '1',
                        admin_name: 'Test Admin',
                        action_type: 'create',
                        target: 'school',
                        details: 'Created school Test',
                        created_at: '2023-01-01T12:00:00Z'
                    }
                ],
                total: 1
            }
        })
    })

    it('renders correctly and fetches data', async () => {
        wrapper = mount(AuditLog, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' }
                }
            }
        })

        expect(wrapper.text()).toContain('Audit Logs')
        expect(adminService.getAuditLogs).toHaveBeenCalled()
    })

    it('filters fetch logs on action change', async () => {
        wrapper = mount(AuditLog, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' }
                }
            }
        })

        // Simulate filter change (if accessible via UI or vm)
        // Accessing reactive data directly for unit test simplicity
        wrapper.vm.filters.action = 'create'
        await wrapper.vm.fetchLogs()

        expect(adminService.getAuditLogs).toHaveBeenLastCalledWith(expect.objectContaining({
            action: 'create'
        }))
    })
})
