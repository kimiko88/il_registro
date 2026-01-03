import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import UserManagement from '@/components/Admin/UserManagement.vue'
import ConfirmDialog from '@/components/Common/ConfirmDialog.vue'

// Mock Admin Service
const mockAdminService = vi.hoisted(() => ({
    deleteUser: vi.fn()
}))
vi.mock('@/services/adminService', () => ({ default: mockAdminService }))

// Mock useApi
const mockFetch = vi.fn()
const mockUsers = ref([])
const mockLoading = ref(false)

import { ref } from 'vue'

vi.mock('@/composables/useApi', () => ({
    useApi: () => ({
        data: mockUsers,
        loading: mockLoading,
        fetch: mockFetch
    })
}))

// Mock Quasar
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            notify: vi.fn()
        })
    }
})

describe('Admin/UserManagement.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        mockUsers.value = [
            { id: 1, first_name: 'John', last_name: 'Doe', email: 'john@test.com', role: 'student' }
        ]
        wrapper = mount(UserManagement, {
            global: {
                stubs: {
                    'q-table': {
                        template: '<div><slot name="body-cell-actions" :props="{row: rows && rows[0] ? rows[0] : {} }" /></div>',
                        props: ['rows']
                    },
                    'q-btn': true,
                    'q-td': { template: '<div><slot /></div>' },
                    'confirm-dialog': true
                }
            }
        })
    })

    it('fetches users on mount', () => {
        expect(mockFetch).toHaveBeenCalled()
    })

    it('opens add user dialog', () => {
        const btn = wrapper.findAllComponents({ name: 'q-btn' }).find(b => b.attributes('icon') === 'add')
        btn.trigger('click')
        // Logic pending implementation
    })

    it('confirms delete', async () => {
        const user = mockUsers.value[0]
        wrapper.vm.confirmDelete(user)
        expect(wrapper.vm.selectedUser).toEqual(user)
        expect(wrapper.vm.showConfirm).toBe(true)
    })

    it('deletes user on confirmation', async () => {
        mockAdminService.deleteUser.mockResolvedValue({})

        wrapper.vm.selectedUser = mockUsers.value[0]
        await wrapper.vm.deleteUser()

        expect(mockAdminService.deleteUser).toHaveBeenCalledWith(1)
        expect(mockFetch).toHaveBeenCalledTimes(2) // 1 on mount, 1 after delete
    })
})
