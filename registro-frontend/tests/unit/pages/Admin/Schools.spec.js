import { describe, it, expect, vi } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import Schools from '@/pages/admin/Schools.vue'
import SchoolTable from '@/components/Admin/SchoolTable.vue'
import SchoolForm from '@/components/Admin/SchoolForm.vue'

// Mock composable
const mockUseSchoolManagement = {
    showDialog: { value: false },
    editingSchool: { value: null },
    openCreate: vi.fn(),
    openEdit: vi.fn(),
    submitSchool: vi.fn(),
    confirmDelete: vi.fn(),
    onRequest: vi.fn()
}

vi.mock('src/composables/useSchoolManagement', () => ({
    useSchoolManagement: () => mockUseSchoolManagement
}))

describe('Schools Page', () => {
    it('renders and wires events to composable', async () => {
        const wrapper = shallowMount(Schools, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-dialog': { template: '<div><slot /></div>' }
                }
            }
        })

        const table = wrapper.findComponent(SchoolTable)
        expect(table.exists()).toBe(true)

        await table.vm.$emit('create')
        expect(mockUseSchoolManagement.openCreate).toHaveBeenCalled()

        await table.vm.$emit('edit', { id: 1 })
        expect(mockUseSchoolManagement.openEdit).toHaveBeenCalledWith({ id: 1 })

        await table.vm.$emit('delete', { id: 1 })
        expect(mockUseSchoolManagement.confirmDelete).toHaveBeenCalledWith({ id: 1 })

        await table.vm.$emit('request', { page: 2 })
        expect(mockUseSchoolManagement.onRequest).toHaveBeenCalledWith({ page: 2 })
    })

    it('submits form via composable', async () => {
        const wrapper = shallowMount(Schools, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-dialog': { template: '<div><slot /></div>' }
                }
            }
        })

        const form = wrapper.findComponent(SchoolForm)
        expect(form.exists()).toBe(true)

        await form.vm.$emit('submit', { name: 'New School' })

        expect(mockUseSchoolManagement.submitSchool).toHaveBeenCalledWith({ name: 'New School' })
    })
})
