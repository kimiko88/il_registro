import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import SecretaryStudents from '@/pages/secretary/Students.vue'
import { useAuthStore } from '@/stores/auth'

const { mockGetAllUsers, mockGetSchoolClasses, mockUpdateUser, mockListChoices, mockSetStudentChoice } = vi.hoisted(() => ({
    mockGetAllUsers: vi.fn(),
    mockGetSchoolClasses: vi.fn(),
    mockUpdateUser: vi.fn(),
    mockListChoices: vi.fn(),
    mockSetStudentChoice: vi.fn()
}))

vi.mock('quasar', async () => {
    const actual = await vi.importActual('quasar')
    return {
        ...actual,
        useQuasar: () => ({
            notify: vi.fn(),
            dialog: vi.fn(() => ({
                onOk: (cb) => {
                    cb()
                    return { onCancel: (cb2) => cb2() }
                }
            }))
        })
    }
})

vi.mock('@/services/userService', () => ({
    userService: {
        getAll: mockGetAllUsers,
        update: mockUpdateUser
    }
}))

vi.mock('@/services/adminService', () => ({
    default: { getSchoolClasses: mockGetSchoolClasses }
}))

vi.mock('@/services/religionService', () => ({
    religionService: {
        listChoices: mockListChoices,
        setStudentChoice: mockSetStudentChoice
    }
}))

describe('Secretary/Students.vue - Religion & Avvalimento Workflow', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()

        mockGetAllUsers.mockResolvedValue({
            data: {
                users: [
                    { id: 'u1', student_id: 's1', first_name: 'Mario', last_name: 'Rossi', email: 'mario@test.com', ClassName: '1A' },
                    { id: 'u2', student_id: 's2', first_name: 'Luigi', last_name: 'Verdi', email: 'luigi@test.com', ClassName: '1A' },
                    { id: 'u3', student_id: 's3', first_name: 'Chiara', last_name: 'Bianchi', email: 'chiara@test.com', ClassName: '1B' }
                ]
            }
        })

        mockGetSchoolClasses.mockResolvedValue({
            data: [
                { id: 'c1', name: '1A' },
                { id: 'c2', name: '1B' }
            ]
        })

        mockListChoices.mockResolvedValue({
            data: [
                { student_id: 's1', choice: 'avvalente' },
                { student_id: 's2', choice: 'non_avvalente' },
                { student_id: 's3', choice: 'attivita_alternativa' }
            ]
        })

        mockSetStudentChoice.mockResolvedValue({
            data: { message: 'Scelta aggiornata' }
        })

        wrapper = mount(SecretaryStudents, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            auth: {
                                user: { school_id: 'school-1', role: 'secretary' }
                            }
                        },
                        stubActions: false
                    })
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-table': {
                        template: '<table><tbody><tr v-for="r in rows" :key="r.id"><td>{{r.first_name}}</td><td>{{r.religion_choice}}</td></tr></tbody></table>',
                        props: ['rows', 'columns', 'filter', 'loading']
                    },
                    'q-btn': true,
                    'q-chip': true,
                    'q-input': true,
                    'q-icon': true,
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-toolbar': true,
                    'q-toolbar-title': true,
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-select': true,
                    'q-radio': true,
                    'q-space': true,
                    'q-td': true
                }
            }
        })
    })

    it('fetches religion choices and maps them onto student rows', async () => {
        await wrapper.vm.$nextTick()
        // Wait for Promise.allSettled
        await new Promise(resolve => setTimeout(resolve, 10))

        expect(mockListChoices).toHaveBeenCalledWith({ school_id: 'school-1' })
        const students = wrapper.vm.students
        expect(students).toHaveLength(3)

        const s1 = students.find(s => s.id === 'u1')
        expect(s1.religion_choice).toBe('avvalente')

        const s2 = students.find(s => s.id === 'u2')
        expect(s2.religion_choice).toBe('non_avvalente')

        const s3 = students.find(s => s.id === 'u3')
        expect(s3.religion_choice).toBe('attivita_alternativa')
    })

    it('defaults students with no explicit choice to avvalente', async () => {
        mockListChoices.mockResolvedValueOnce({
            data: [
                { student_id: 's1', choice: 'non_avvalente' }
            ]
        })

        await wrapper.vm.fetchStudents()
        await new Promise(resolve => setTimeout(resolve, 10))

        const students = wrapper.vm.students
        const s2 = students.find(s => s.id === 'u2')
        // No explicit row in religion list -> defaults to avvalente
        expect(s2.religion_choice).toBe('avvalente')
    })

    it('returns correct chip color for each religion choice', () => {
        expect(wrapper.vm.getReligionChipColor('avvalente')).toBe('teal-7')
        expect(wrapper.vm.getReligionChipColor('non_avvalente')).toBe('orange-8')
        expect(wrapper.vm.getReligionChipColor('attivita_alternativa')).toBe('purple-7')
        expect(wrapper.vm.getReligionChipColor(undefined)).toBe('teal-7')
    })

    it('returns correct chip icon for each religion choice', () => {
        expect(wrapper.vm.getReligionChipIcon('avvalente')).toBe('check_circle')
        expect(wrapper.vm.getReligionChipIcon('non_avvalente')).toBe('block')
        expect(wrapper.vm.getReligionChipIcon('attivita_alternativa')).toBe('swap_horiz')
        expect(wrapper.vm.getReligionChipIcon('unknown')).toBe('check_circle')
    })

    it('returns correct label for each religion choice', () => {
        expect(wrapper.vm.getReligionLabel('avvalente')).toBe('Avvalente IRC')
        expect(wrapper.vm.getReligionLabel('non_avvalente')).toBe('Non avvalente')
        expect(wrapper.vm.getReligionLabel('attivita_alternativa')).toBe('Attività Alternativa')
        expect(wrapper.vm.getReligionLabel('unknown')).toBe('Avvalente IRC')
    })

    it('opens religion choice dialog with preselected student choice', () => {
        const student = { id: 'u2', student_id: 's2', first_name: 'Luigi', last_name: 'Verdi', religion_choice: 'non_avvalente' }
        wrapper.vm.openReligionDialog(student)

        expect(wrapper.vm.showReligionDialog).toBe(true)
        expect(wrapper.vm.selectedStudentForReligion).toEqual(student)
        expect(wrapper.vm.selectedReligionChoice).toBe('non_avvalente')
    })

    it('saves religion choice and updates the student record locally', async () => {
        const student = { id: 'u1', student_id: 's1', first_name: 'Mario', last_name: 'Rossi', religion_choice: 'avvalente' }
        wrapper.vm.openReligionDialog(student)
        wrapper.vm.selectedReligionChoice = 'attivita_alternativa'

        await wrapper.vm.saveReligionChoice()

        expect(mockSetStudentChoice).toHaveBeenCalledWith('s1', 'attivita_alternativa')
        expect(student.religion_choice).toBe('attivita_alternativa')
        expect(wrapper.vm.showReligionDialog).toBe(false)
    })
})
