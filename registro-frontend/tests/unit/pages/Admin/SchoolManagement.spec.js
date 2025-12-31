import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { exportFile } from 'quasar'
import SchoolManagement from '@/pages/Admin/SchoolManagement.vue'
import adminService from '@/services/adminService'

// Mock adminService
vi.mock('@/services/adminService', () => ({
    default: {
        getSchools: vi.fn(),
        deleteSchool: vi.fn(),
        createSchool: vi.fn(),
        updateSchool: vi.fn()
    }
}))

// Mock permissions
vi.mock('@/composables/usePermissions', () => ({
    usePermissions: () => ({
        isSuperAdmin: true,
        canCreateSchools: true,
        canDeleteSchools: true,
        canEditSchool: () => true
    })
}))

// Mock exportFile from quasar
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        exportFile: vi.fn().mockReturnValue(true),
        useQuasar: () => ({
            notify: vi.fn(),
            dialog: vi.fn().mockImplementation(({ title, message }) => ({
                onOk: (fn) => fn() // Auto confirm
            }))
        })
    }
})

describe('SchoolManagement.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        adminService.getSchools.mockResolvedValue({
            data: {
                items: [
                    { id: '1', name: 'School A', code: 'A', city: 'City', student_count: 10, teacher_count: 2, is_active: true }
                ],
                total: 1
            }
        })
    })

    it('renders correctly and has export button', async () => {
        wrapper = mount(SchoolManagement, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' }
                }
            }
        })
        await wrapper.vm.$nextTick()

        expect(wrapper.text()).toContain('Gestione Scuole')
        expect(wrapper.text()).toContain('Gestione Scuole')

        // Check for Export button stub
        const exportBtn = wrapper.findAll('q-btn-stub').find(w => w.attributes('label') === 'Export CSV')
        expect(exportBtn).toBeDefined()
        expect(exportBtn.exists()).toBe(true)
    })

    it('calls exportFile when export button clicked', async () => {
        wrapper = mount(SchoolManagement, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' }
                }
            }
        })
        await wrapper.vm.$nextTick()

        // Trigger export
        await wrapper.vm.exportTable() // Calling method directly for simplicity

        // Check if exportFile was called (it is mocked above)
        // Note: Since we mocked the whole module, we need to import it to check
        const { exportFile } = await import('quasar')
        expect(exportFile).toHaveBeenCalled()
    })

    it('handles bulk delete', async () => {
        wrapper = mount(SchoolManagement, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' }
                }
            }
        })
        await wrapper.vm.$nextTick()

        // Select item
        wrapper.vm.selected = [{ id: '1' }]

        // Call delete
        await wrapper.vm.deleteSelected()

        // Since dialog auto-confirms in mock, it should call service
        expect(adminService.deleteSchool).toHaveBeenCalledWith('1')
    })
})
