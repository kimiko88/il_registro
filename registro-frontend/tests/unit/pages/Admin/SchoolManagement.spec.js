import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import SchoolManagement from '@/pages/admin/SchoolManagement.vue'
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
        
        wrapper = mount(SchoolManagement, {
            global: {
                // We rely on global mocks and stubs from setup.js
                // But we add $router which is not global
                mocks: {
                    $router: { push: vi.fn() }
                },
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-table': { template: '<div><slot name="body" /></div>' },
                    'q-btn': { template: '<button @click="$emit(\'click\', $event)">{{ label }}<slot /></button>', props: ['label'] }
                },
                provide: {
                    'router': { push: vi.fn(), replace: vi.fn() }
                }
            }
        })
    })

    it('renders correctly and has export button', async () => {
        await wrapper.vm.$nextTick()

        // The title might be inside a stubbed component, check if it exists in the DOM
        expect(wrapper.html()).toContain('Gestione Scuole')
        
        const buttons = wrapper.findAll('button')
        const exportBtn = buttons.find(b => b.text().includes('Esporta'))
        expect(exportBtn).toBeDefined()
    })

    it('handles bulk delete', async () => {
        await wrapper.vm.$nextTick()

        // Select item
        wrapper.vm.selected = [{ id: '1' }]

        // Call delete
        await wrapper.vm.deleteSelected()

        expect(adminService.deleteSchool).toHaveBeenCalledWith('1')
    })

    it('opens create dialog', async () => {
        await wrapper.vm.$nextTick()
        wrapper.vm.openCreate()
        await wrapper.vm.$nextTick()
        expect(wrapper.vm.showCreateDialog).toBe(true)
        expect(wrapper.vm.editingSchool).toBeNull()
        expect(wrapper.vm.schoolForm.name).toBe('')
    })

    it('opens edit dialog', async () => {
        await wrapper.vm.$nextTick()
        const school = { id: '1', name: 'Original', address: 'Addr' }
        wrapper.vm.editSchool(school)
        await wrapper.vm.$nextTick()
        expect(wrapper.vm.showCreateDialog).toBe(true)
        expect(wrapper.vm.editingSchool).not.toBeNull()
        expect(wrapper.vm.schoolForm.name).toBe('Original')
    })

    it('creates school successfully', async () => {
        wrapper.vm.openCreate()
        wrapper.vm.schoolForm.name = 'New School'
        wrapper.vm.schoolForm.code = 'NS'

        await wrapper.vm.saveSchool()

        expect(adminService.createSchool).toHaveBeenCalled()
        expect(wrapper.vm.showCreateDialog).toBe(false)
    })

    it('updates school successfully', async () => {
        const school = { id: '1', name: 'Old' }
        wrapper.vm.editSchool(school)
        wrapper.vm.schoolForm.name = 'Updated'

        await wrapper.vm.saveSchool()

        expect(adminService.updateSchool).toHaveBeenCalledWith('1', expect.objectContaining({ name: 'Updated' }))
    })

    it('confirms delete school', async () => {
        const school = { id: '1' }
        wrapper.vm.confirmDelete(school)
        expect(adminService.deleteSchool).toHaveBeenCalledWith('1')
    })
})
