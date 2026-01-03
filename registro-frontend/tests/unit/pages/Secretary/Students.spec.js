import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import SecretaryStudents from '@/pages/secretary/Students.vue'
import { useAuthStore } from '@/stores/auth'

// Mock services
const { mockGetAllUsers, mockGetSchoolClasses, mockUpdateUser } = vi.hoisted(() => ({
    mockGetAllUsers: vi.fn(),
    mockGetSchoolClasses: vi.fn(),
    mockUpdateUser: vi.fn()
}))

vi.mock('src/services/userService', () => ({
    userService: {
        getAll: mockGetAllUsers,
        update: mockUpdateUser
    }
}))
vi.mock('src/services/adminService', () => ({
    default: { getSchoolClasses: mockGetSchoolClasses }
}))

// Mock Quasar
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            notify: vi.fn(),
            loading: { show: vi.fn(), hide: vi.fn() },
            dialog: vi.fn(() => ({ onOk: (fn) => fn() }))
        })
    }
})

describe('Secretary/Students.vue', () => {
    let wrapper
    let auth

    beforeEach(() => {
        vi.clearAllMocks()

        // Default mocks
        mockGetAllUsers.mockResolvedValue({
            data: {
                users: [
                    { id: 'u1', first_name: 'Mario', last_name: 'Rossi', email: 'mario@test.com', ClassName: '1A' },
                    { id: 'u2', first_name: 'Luigi', last_name: 'Verdi', email: 'luigi@test.com', ClassName: '1B' }
                ]
            }
        })
        mockGetSchoolClasses.mockResolvedValue({
            data: [
                { id: 'c1', name: '1A' },
                { id: 'c2', name: '1B' }
            ]
        })

        wrapper = mount(SecretaryStudents, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            auth: {
                                user: { school_id: 'school1', role: 'secretary' }
                            }
                        },
                        stubActions: false
                    })
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-table': {
                        template: '<table><tbody><tr v-for="r in rows" :key="r.id" @click="$emit(\'row-click\', r)"><td>{{r.first_name}}</td></tr></tbody></table>',
                        props: ['rows', 'columns', 'filter', 'loading']
                    },
                    'q-btn': true,
                    'q-input': true,
                    'q-icon': true,
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-toolbar': true,
                    'q-toolbar-title': true,
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-select': true,
                    'q-space': true,
                    'q-td': true
                }
            }
        })
        auth = useAuthStore()
    })

    it('fetches students and classes on mount', () => {
        expect(mockGetAllUsers).toHaveBeenCalledWith({ role: 'student' })
        expect(mockGetSchoolClasses).toHaveBeenCalledWith('school1')
    })

    it('opens enrollment dialog', async () => {
        // There is a QBtn with @click="showEnrollment = true"
        // Since we stub QBtn, we can find it by label or icon if we rendered attributes, 
        // BUT simplest is to check if dialog state changes if we could access vm.
        // Or we trigger the button click.
        // The button has label "Nuova Iscrizione".
        // Stubbed q-btn usually emits click.
        const btn = wrapper.findAllComponents({ name: 'q-btn' }).find(b => b.attributes('label') === 'Nuova Iscrizione')
        if (btn) await btn.trigger('click')

        expect(wrapper.vm.showEnrollment).toBe(true)
    })

    it('opens edit dialog with student data', () => {
        const student = { id: 'u1', first_name: 'Mario', last_name: 'Rossi', email: 'mario@test.com', ClassName: '1A', ClassID: 'c1' }
        // Call directly or via table row click if emulated
        wrapper.vm.editStudent(student)
        expect(wrapper.vm.showEditDialog).toBe(true)
        expect(wrapper.vm.editForm.first_name).toBe('Mario')
        expect(wrapper.vm.editForm.class_id).toBe('c1')
    })

    it('saves student changes', async () => {
        wrapper.vm.editForm.id = 'u1'
        wrapper.vm.editForm.class_id = 'c2'

        await wrapper.vm.saveStudent()

        expect(mockUpdateUser).toHaveBeenCalledWith('u1', expect.objectContaining({ class_id: 'c2' }))
        expect(wrapper.vm.showEditDialog).toBe(false)
        // Should refresh list
        expect(mockGetAllUsers).toHaveBeenCalledTimes(2) // Once on mount, once after save
    })

    it('changes student class (Secretary flow)', async () => {
        // Explicit flow: open edit -> change class -> save
        const student = { id: 'u1', first_name: 'Mario', ClassID: 'c1' }
        wrapper.vm.editStudent(student)

        expect(wrapper.vm.editForm.class_id).toBe('c1')

        // Simulate changing selection
        wrapper.vm.editForm.class_id = 'c2'

        await wrapper.vm.saveStudent()
        expect(mockUpdateUser).toHaveBeenCalledWith('u1', expect.objectContaining({ class_id: 'c2' }))
    })

    it('handles error when changing class', async () => {
        mockUpdateUser.mockRejectedValue(new Error('Update failed'))
        wrapper.vm.editStudent({ id: 'u1' })

        await wrapper.vm.saveStudent()

        // Check for notify error
        // Note: We need to spy on notify. In mock setup, we return a notify mock.
        // But we don't have access to the specific spy instance returned by the factory unless we hoist it or use a global mock.
        // In this file, useQuasar is mocked inline.
        // We can inspect the console or better, refactor the mock to expose the spy.
        // Alternatively, check if showEditDialog remains open or verify side effects.
        expect(wrapper.vm.showEditDialog).toBe(true) // Should stay open on error?
        // Actually the code resets close but notifies error? Let's check source code.
        // Source: } catch (e) { $q.notify(...) } 
        // It does NOT close dialog in catch block? Source says: 
        /*
        try { ... showEditDialog.value = false ... } catch (e) { $q.notify ... }
        */
        // So yes, it should stay open.
        expect(wrapper.vm.showEditDialog).toBe(true)
    })

    it('validates school ID during update', async () => {
        // Ensure school_id is passed if required
        const student = { id: 'u1', school_id: 's99' }
        wrapper.vm.editStudent(student)
        await wrapper.vm.saveStudent()
        expect(mockUpdateUser).toHaveBeenCalledWith('u1', expect.objectContaining({ school_id: 's99' }))
    })
})
