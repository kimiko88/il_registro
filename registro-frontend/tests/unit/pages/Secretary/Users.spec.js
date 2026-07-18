
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { exportFile } from 'quasar'
import Users from '@/pages/secretary/Users.vue'
import { userService } from 'src/services/userService'
import adminService from 'src/services/adminService'

// Mock Quasar
const mockDialog = vi.fn()
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        exportFile: vi.fn().mockReturnValue(true),
        useQuasar: () => ({
            loading: { show: vi.fn(), hide: vi.fn() },
            notify: vi.fn(),
            dialog: mockDialog.mockImplementation(() => ({ onOk: (fn) => fn() }))
        })
    }
})

// Mock Services
vi.mock('src/services/userService', () => ({
    userService: {
        getAll: vi.fn().mockResolvedValue({ data: { users: [{ id: 1, first_name: 'Test', role: 'student' }] } }),
        create: vi.fn().mockResolvedValue({ data: { id: 2 } }),
        update: vi.fn().mockResolvedValue({}),
        delete: vi.fn().mockResolvedValue({}),
        forceResetPassword: vi.fn().mockResolvedValue({})
    }
}))

vi.mock('src/services/adminService', () => ({
    default: {
        getSchoolClasses: vi.fn().mockResolvedValue({ data: [{ id: 'c1', name: '1A' }] }),
        getSubjects: vi.fn().mockResolvedValue({ data: [{ id: 's1', name: 'Math' }] }),
        getTeachersList: vi.fn().mockResolvedValue({ data: [{ id: 't1', user_id: 10 }] }),
        createClass: vi.fn().mockResolvedValue({}),
        getTeacherSubjects: vi.fn().mockResolvedValue({ data: [] }),
        assignSubjectToTeacher: vi.fn().mockResolvedValue({}),
        removeTeacherSubject: vi.fn().mockResolvedValue({})
    }
}))

describe('Secretary Users Page (Users.vue)', () => {
    let wrapper

    beforeEach(() => {
        vi.useFakeTimers()
        vi.clearAllMocks()
        wrapper = mount(Users, {
            global: {
                plugins: [createTestingPinia({
                    initialState: {
                        auth: { 
                            user: { school_id: '1' },
                            userRole: 'secretary'
                        }
                    },
                    createSpy: vi.fn
                })],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-dialog': { template: '<div><slot /></div>' },
                    'UserTable': {
                        template: '<div class="user-table"></div>',
                        emits: ['create', 'edit', 'delete', 'reset-pwd', 'filter-role', 'export', 'bulk-delete', 'manage-subjects']
                    },
                    'q-btn': true,
                    'q-input': true,
                    'q-select': true,
                    'q-file': true,
                    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-list': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-space': true,
                    'q-icon': true,
                    'q-toolbar': true,
                    'q-toolbar-title': true
                }
            }
        })
    })

    afterEach(() => {
        vi.useRealTimers()
    })

    it('fetches initial data', async () => {
        await vi.runAllTimersAsync()
        expect(userService.getAll).toHaveBeenCalled()
        expect(adminService.getSchoolClasses).toHaveBeenCalled()
    })

    it('opens class dialog and creates class', async () => {
        const btn = wrapper.find('q-btn[label="Nuova Classe"]')
        // Direct VM access for testing
        wrapper.vm.openClassDialog()
        expect(wrapper.vm.showClassDialog).toBe(true)

        wrapper.vm.classForm.name = '1A'
        await wrapper.vm.saveClass()

        expect(adminService.createClass).toHaveBeenCalledWith(expect.objectContaining({ name: '1A' }))
        expect(wrapper.vm.showClassDialog).toBe(false)
    })

    it('manages teacher subjects', async () => {
        const user = { id: 10, last_name: 'Teach', first_name: 'Er' }
        await wrapper.vm.openManageSubjects(user)

        expect(wrapper.vm.currentTeacherId).toBe('t1')
        expect(wrapper.vm.showSubjectsDialog).toBe(true)
        expect(adminService.getTeacherSubjects).toHaveBeenCalledWith('t1')

        // Add subject
        wrapper.vm.selectedSubjectToAdd = 's1'
        await wrapper.vm.addTeacherSubject()
        expect(adminService.assignSubjectToTeacher).toHaveBeenCalledWith('t1', 's1')

        // Remove subject
        await wrapper.vm.removeTeacherSubject('s1')
        expect(adminService.removeTeacherSubject).toHaveBeenCalledWith('t1', 's1')
    })

    it('handles import mock', async () => {
        wrapper.vm.importFile = new File([''], 'test.csv')
        await wrapper.vm.handleImport()
        expect(wrapper.vm.showImport).toBe(false)
    })

    it('handles bulk delete mock', () => {
        wrapper.vm.bulkDelete([])
    })

    it('confirms and deletes user via dialog', async () => {
        const user = { id: 7, first_name: 'Delete', last_name: 'Me' }

        await wrapper.vm.confirmDelete(user)

        expect(mockDialog).toHaveBeenCalled()
        expect(userService.delete).toHaveBeenCalledWith(7)
    })

    it('exports users to CSV', () => {
        wrapper.vm.users = [
            { id: 1, first_name: 'Test', last_name: 'User', email: 'test@test.com', role: 'student' }
        ]
        wrapper.vm.exportUsers()
        expect(exportFile).toHaveBeenCalled()
    })

    it('opens reset password dialog and handles reset', async () => {
        const user = { id: 8, email: 'reset@test.com', first_name: 'Test', last_name: 'User' }
 
        wrapper.vm.openResetPwd(user)
        expect(wrapper.vm.showResetPwdDialog).toBe(true)
        expect(wrapper.vm.resetTargetId).toBe(8)

        wrapper.vm.newPassword = 'newPassword123'
        await wrapper.vm.handleResetPwd()

        expect(userService.forceResetPassword).toHaveBeenCalledWith(8, 'newPassword123')
        expect(wrapper.vm.showResetPwdDialog).toBe(false)
    })
})
