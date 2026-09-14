import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import Users from '@/pages/secretary/Users.vue'

const { mockUserService } = vi.hoisted(() => {
    return {
        mockUserService: {
            getAll: vi.fn(),
            getUsers: vi.fn(),
            getClasses: vi.fn(),
            createUser: vi.fn(),
            updateUser: vi.fn(),
            deleteUser: vi.fn(),
            getAssignments: vi.fn(),
            addAssignment: vi.fn(),
            deleteAssignment: vi.fn(),
            setCoordinatedClasses: vi.fn()
        }
    }
})

vi.mock('@/services/userService', () => ({
    userService: mockUserService,
    default: mockUserService
}))

vi.mock('@/services/schoolService', () => ({
    default: {
        getClasses: vi.fn().mockResolvedValue([
            { id: 'class-2a', name: '2A', section: 'A' },
            { id: 'class-2b', name: '2B', section: 'B' }
        ])
    }
}))

describe('Users.vue Assignments & Governance Integration Tests', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        mockUserService.getAll.mockResolvedValue({
            data: {
                users: [
                    {
                        id: 'teacher-1',
                        first_name: 'Marco',
                        last_name: 'Rossi',
                        email: 'docente1@scuola.it',
                        role: 'teacher',
                        is_active: true,
                        assignments: [
                            { id: 'asgn-1', assignment_type: 'coordinatore_classe', scope_id: 'class-2a', title: 'Coordinatore 2A', is_active: true }
                        ]
                    },
                    {
                        id: 'tec-1',
                        first_name: 'Alessandro',
                        last_name: 'Tecnico',
                        email: 'tecnico.prova@scuola.it',
                        role: 'assistente_tecnico',
                        is_active: true,
                        assignments: []
                    }
                ],
                total: 2
            }
        })
        mockUserService.getClasses.mockResolvedValue({
            data: [
                { id: 'class-2a', name: '2A', section: 'A' },
                { id: 'class-2b', name: '2B', section: 'B' }
            ]
        })
        mockUserService.setCoordinatedClasses.mockResolvedValue({ success: true })
        mockUserService.addAssignment.mockResolvedValue({ id: 'asgn-new', success: true })
        mockUserService.deleteAssignment.mockResolvedValue({ success: true })
        mockUserService.getAssignments.mockResolvedValue({
            data: [
                { id: 'asgn-1', assignment_type: 'coordinatore_classe', scope_id: 'class-2a', title: 'Coordinatore 2A', is_active: true }
            ]
        })

        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'admin-1', role: 'admin', email: 'admin@scuola.it' }
                }
            }
        })
    })

    function createWrapper() {
        return mount(Users, {
            global: {
                plugins: [Quasar, pinia],
                stubs: {
                    UserTable: {
                        template: '<div class="user-table-stub"><button class="btn-manage-assignments" @click="$emit(\'manage-assignments\', users[0])">Incarichi</button></div>',
                        props: ['users', 'loading', 'pagination', 'roleOptions']
                    },
                    CsvUserImportDialog: true,
                    'q-page': { template: '<div><slot /></div>' },
                    'q-dialog': { template: '<div v-if="modelValue" class="q-dialog-stub"><slot /></div>', props: ['modelValue'] },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-select': true,
                    'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
                    'q-chip': { template: '<span><slot /></span>' },
                    'q-icon': true,
                    'q-tooltip': true
                }
            }
        })
    }

    it('opens assignments dialog when manage-assignments event is emitted', async () => {
        const wrapper = createWrapper()
        await wrapper.vm.$nextTick()

        const sampleTeacher = {
            id: 'teacher-1',
            first_name: 'Marco',
            last_name: 'Rossi',
            email: 'docente1@scuola.it',
            role: 'teacher',
            assignments: [
                { id: 'asgn-1', assignment_type: 'coordinatore_classe', scope_id: 'class-2a', title: 'Coordinatore 2A', is_active: true }
            ]
        }

        await wrapper.vm.openManageAssignments(sampleTeacher)
        await wrapper.vm.$nextTick()

        expect(wrapper.vm.showAssignmentsDialog).toBe(true)
        expect(wrapper.vm.targetUserForAssignments.id).toBe('teacher-1')
        expect(wrapper.vm.selectedCoordinatedClasses).toEqual(['class-2a'])
    })

    it('calls userService.setCoordinatedClasses when updating coordination', async () => {
        const wrapper = createWrapper()
        const sampleTeacher = {
            id: 'teacher-1',
            first_name: 'Marco',
            last_name: 'Rossi',
            email: 'docente1@scuola.it',
            role: 'teacher',
            assignments: []
        }

        await wrapper.vm.openManageAssignments(sampleTeacher)
        wrapper.vm.selectedCoordinatedClasses = ['class-2a', 'class-2b']

        await wrapper.vm.saveCoordinatedClasses()

        expect(mockUserService.setCoordinatedClasses).toHaveBeenCalledWith('teacher-1', ['class-2a', 'class-2b'])
        expect(wrapper.vm.savingCoordinatedClasses).toBe(false)
    })

    it('assigns quick pedagogical duties and reloads assignments', async () => {
        const wrapper = createWrapper()
        const sampleTeacher = {
            id: 'teacher-1',
            first_name: 'Marco',
            last_name: 'Rossi',
            email: 'docente1@scuola.it',
            role: 'teacher',
            assignments: []
        }

        await wrapper.vm.openManageAssignments(sampleTeacher)

        const inclusionDuty = {
            type: 'referente_inclusione',
            label: 'Referente Inclusione (BES / DSA)',
            desc: 'Coordina i piani PDP/PEI e le misure compensative'
        }

        await wrapper.vm.assignDutyQuick(inclusionDuty)

        expect(mockUserService.addAssignment).toHaveBeenCalledWith('teacher-1', expect.objectContaining({
            assignment_type: 'referente_inclusione',
            scope_type: 'school',
            title: 'Referente Inclusione (BES / DSA)'
        }))
        expect(mockUserService.getAssignments).toHaveBeenCalledWith('teacher-1')
    })

    it('revokes an active duty via userService.deleteAssignment', async () => {
        const wrapper = createWrapper()
        const sampleTeacher = {
            id: 'teacher-1',
            first_name: 'Marco',
            last_name: 'Rossi',
            email: 'docente1@scuola.it',
            role: 'teacher',
            assignments: [
                { id: 'asgn-1', assignment_type: 'referente_inclusione', is_active: true }
            ]
        }

        mockUserService.getAssignments.mockResolvedValueOnce({
            data: [
                { id: 'asgn-1', assignment_type: 'referente_inclusione', is_active: true }
            ]
        })

        await wrapper.vm.openManageAssignments(sampleTeacher)
        await wrapper.vm.removeDutyByType('referente_inclusione')

        expect(mockUserService.deleteAssignment).toHaveBeenCalledWith('teacher-1', 'asgn-1')
    })
})
