import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import AdminUsers from '@/pages/admin/AdminUsers.vue'

// Mock Services
const { mockGetAdmins, mockCreateAdmin, mockUpdateAdmin, mockDeleteAdmin, mockGetSchools } = vi.hoisted(() => ({
    mockGetAdmins: vi.fn(),
    mockCreateAdmin: vi.fn(),
    mockUpdateAdmin: vi.fn(),
    mockDeleteAdmin: vi.fn(),
    mockGetSchools: vi.fn()
}))

vi.mock('@/services/adminService', () => ({
    default: {
        getAdmins: mockGetAdmins,
        createAdmin: mockCreateAdmin,
        updateAdmin: mockUpdateAdmin,
        deleteAdmin: mockDeleteAdmin,
        getSchools: mockGetSchools,
        resetAdminPassword: vi.fn(),
        getAdminActivity: vi.fn(() => Promise.resolve({ data: [] }))
    }
}))

// Mock Quasar
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            notify: vi.fn(),
            dialog: vi.fn(() => ({ onOk: (fn) => fn() }))
        }),
        debounce: (fn) => fn
    }
})

describe('Admin/AdminUsers.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        mockGetAdmins.mockResolvedValue({
            data: {
                items: [
                    { id: 1, first_name: 'John', last_name: 'Doe', email: 'john@test.com', is_active: true, school_name: 'School A' }
                ],
                total: 1
            }
        })
        mockGetSchools.mockResolvedValue({
            data: { items: [{ id: 101, name: 'School A' }] }
        })
    })

    const mountComponent = () => {
        return mount(AdminUsers, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-table': {
                        template: `
                            <table>
                                <tr v-for="row in rows" :key="row.id">
                                    <td><slot name="body-cell-user" :row="row" :props="{row}" /></td>
                                    <td><slot name="body-cell-school" :row="row" :props="{row}" /></td>
                                    <td><slot name="body-cell-status" :row="row" :props="{row}" /></td>
                                    <td><slot name="body-cell-last_login" :row="row" :props="{row}" /></td>
                                    <td><slot name="body-cell-actions" :row="row" :props="{row}" /></td>
                                </tr>
                            </table>
                        `,
                        props: ['rows', 'columns']
                    },
                    'q-dialog': { template: '<div><slot /></div>', props: ['modelValue'] },
                    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-input': { template: '<input class="q-input" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />', props: ['modelValue'] },
                    'q-select': { template: '<select class="q-select" :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)" ></select>', props: ['modelValue'] },
                    'q-btn': true,
                    'q-icon': true,
                    'q-td': { template: '<div><slot /></div>' },
                    'q-badge': { template: '<span><slot /></span>' },
                    'q-btn-dropdown': { template: '<div><slot /></div>' },
                    'q-list': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div @click="$emit(\'click\')"><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-toggle': { template: '<input type="checkbox" />' },
                    'q-separator': true,
                    'q-space': true
                }
            }
        })
    }

    it('fetches admins on mount', async () => {
        wrapper = mountComponent()
        await wrapper.vm.$nextTick() // wait for mount
        // Mock resolves immediately but in real fetch logic there might be another tick.
        await new Promise(resolve => setTimeout(resolve, 0)) // wait async

        expect(mockGetAdmins).toHaveBeenCalledTimes(1)
        expect(wrapper.vm.admins).toHaveLength(1)
        expect(wrapper.vm.admins[0].email).toBe('john@test.com')
    })

    it('opens create dialog', async () => {
        wrapper = mountComponent()
        await wrapper.vm.$nextTick()

        wrapper.vm.openCreateDialog()
        await wrapper.vm.$nextTick()

        expect(wrapper.vm.showDialog).toBe(true)
        expect(wrapper.vm.editingAdmin).toBeNull()
    })

    it('creates a new admin', async () => {
        wrapper = mountComponent()
        mockCreateAdmin.mockResolvedValue({})

        wrapper.vm.openCreateDialog()
        await wrapper.vm.$nextTick()

        wrapper.vm.adminForm.first_name = 'Jane'
        wrapper.vm.adminForm.last_name = 'Doe'
        wrapper.vm.adminForm.email = 'jane@test.com'
        wrapper.vm.adminForm.school_id = 101

        await wrapper.vm.saveAdmin()

        expect(mockCreateAdmin).toHaveBeenCalledWith(expect.objectContaining({
            first_name: 'Jane',
            email: 'jane@test.com'
        }))
        expect(wrapper.vm.showDialog).toBe(false)
        expect(mockGetAdmins).toHaveBeenCalledTimes(2)
    })

    it('filters admins', async () => {
        wrapper = mountComponent()
        await wrapper.vm.$nextTick()

        wrapper.vm.filters.search = 'Jane'
        wrapper.vm.debouncedFetch()

        expect(mockGetAdmins).toHaveBeenCalledWith(expect.objectContaining({
            search: 'Jane'
        }))
    })
})
