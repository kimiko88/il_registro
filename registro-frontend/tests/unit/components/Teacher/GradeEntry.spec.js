import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar, Notify } from 'quasar'
import GradeEntry from '@/components/Teacher/GradeEntry.vue'
import { useGradesStore } from '@/stores/grades'
import { gradeService } from '@/services/gradeService'

// Mock Services
vi.mock('@/services/gradeService', () => ({
    gradeService: {
        saveGrade: vi.fn()
    }
}))

// Mock services if needed
const mockNotify = vi.hoisted(() => vi.fn())
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            notify: mockNotify,
            dark: { isActive: false },
            lang: { rtl: false }
        })
    }
})

describe('GradeEntry.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        wrapper = mount(GradeEntry, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            grades: {
                                grades: {
                                    students: [
                                        { student_id: 'S1', full_name: 'Mario Rossi', absences: 2, grades: [], average: '0.0' },
                                        { student_id: 'S2', full_name: 'Luigi Verdi', absences: 0, grades: [], average: '0.0' },
                                        { student_id: 'S3', full_name: 'Anna Neri', absences: 1, grades: [], average: '0.0' }
                                    ]
                                }
                            }
                        }
                    })
                ],
                // Stub complex Quasar components
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-table': {
                        template: `
                          <div class="q-table-stub">
                            <slot name="header" :props="{}" />
                            <template v-for="(row, index) in rows" :key="index">
                              <slot name="body" :row="row" :rowIndex="index" :props="{ row, rowIndex: index }" />
                            </template>
                          </div>
                        `,
                        props: ['rows']
                    },
                    'q-tr': { template: '<tr><slot /></tr>' },
                    'q-th': { template: '<th><slot /></th>' },
                    'q-td': { template: '<td><slot /></td>' },
                    'q-input': {
                        template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
                        props: ['modelValue']
                    },
                    'q-select': {
                        template: '<select :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)"><option v-for="opt in options" :key="opt" :value="opt">{{ opt }}</option></select>',
                        props: ['modelValue', 'options']
                    },
                    'q-btn': { template: '<button @click="$emit(\'click\')"></button>' },
                    'q-badge': { template: '<span><slot /></span>' },
                    'q-tooltip': { template: '<span><slot /></span>' },
                    'q-icon': true
                }
            },
            props: {
                classId: 'C1',
                subject: 'Math',
                date: '2025-01-20',
                type: 'Oral'
            }
        })
    })

    it('initializes with mock data', () => {
        // Check if initData was called (it runs on mounted via watch immediate)
        // We can check if internal state is set.
        // Since we stubbed q-table, we can't easily check rendered rows unless we inspect component state or stub props.
        // But testing-library or standard mounting renders stubs.
        expect(wrapper.vm.studentsWithGrades.length).toBe(3)
    })

    it('detects dirty state', async () => {
        const studentId = 'S1'
        // Modify value
        wrapper.vm.entryData[studentId].value = 9

        expect(wrapper.vm.isDirty(studentId)).toBe(true)
        expect(wrapper.vm.hasChanges).toBe(true)
    })

    it('saves a single line', async () => {
        const studentId = 'S1'
        wrapper.vm.entryData[studentId].value = 8

        await wrapper.vm.saveLine(studentId)

        expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({
            message: expect.stringContaining('Voto 8 salvato')
        }))
        expect(wrapper.vm.isDirty(studentId)).toBe(false)
    })

    it('validation: prevents saving empty', async () => {
        const studentId = 'S1'
        wrapper.vm.entryData[studentId].value = null

        await wrapper.vm.saveLine(studentId)

        expect(mockNotify).not.toHaveBeenCalled()
    })

    it('simulates bulk save', async () => {
        wrapper.vm.entryData['S1'].value = 8
        const savePromise = wrapper.vm.saveAll()
        expect(wrapper.vm.loading).toBe(true)

        await savePromise

        expect(wrapper.vm.loading).toBe(false)
        expect(mockNotify).toHaveBeenCalled()
    })

    it('recognizes religion subject and exempts non-avvalente students', async () => {
        const religionWrapper = mount(GradeEntry, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            grades: {
                                grades: {
                                    is_religion_subject: true,
                                    students: [
                                        { student_id: 'S1', full_name: 'Mario Rossi', religion_choice: 'avvalente', absences: 0, grades: [] },
                                        { student_id: 'S2', full_name: 'Luigi Verdi', religion_choice: 'non_avvalente', absences: 0, grades: [] },
                                        { student_id: 'S3', full_name: 'Sara Neri', religion_choice: 'attivita_alternativa', absences: 0, grades: [] }
                                    ]
                                }
                            }
                        }
                    })
                ],
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-table': { template: '<div class="q-table-stub"><slot /></div>' },
                    'q-tr': { template: '<tr><slot /></tr>' },
                    'q-th': { template: '<th><slot /></th>' },
                    'q-td': { template: '<td><slot /></td>' },
                    'q-input': true,
                    'q-select': true,
                    'q-btn': true,
                    'q-badge': true,
                    'q-tooltip': true,
                    'q-icon': true
                }
            },
            props: {
                classId: 'C1',
                subject: 'Religione Cattolica',
                isReligion: true,
                date: '2025-01-20',
                type: 'Oral'
            }
        })

        expect(religionWrapper.vm.isReligionSubject).toBe(true)

        // S1 is avvalente -> not exempt
        expect(religionWrapper.vm.isStudentExemptFromReligion({ religion_choice: 'avvalente' })).toBe(false)

        // S2 is non_avvalente -> exempt
        expect(religionWrapper.vm.isStudentExemptFromReligion({ religion_choice: 'non_avvalente' })).toBe(true)

        // S3 is attivita_alternativa -> exempt
        expect(religionWrapper.vm.isStudentExemptFromReligion({ religion_choice: 'attivita_alternativa' })).toBe(true)
    })

    it('saves religion judgment as grade_type: judgment with mapped numerical value', async () => {
        const religionWrapper = mount(GradeEntry, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            grades: {
                                grades: {
                                    is_religion_subject: true,
                                    students: [
                                        { student_id: 'S1', full_name: 'Mario Rossi', religion_choice: 'avvalente', absences: 0, grades: [] }
                                    ]
                                }
                            }
                        }
                    })
                ],
                stubs: {
                    'q-card': true,
                    'q-table': true,
                    'q-tr': true,
                    'q-th': true,
                    'q-td': true,
                    'q-input': true,
                    'q-select': true,
                    'q-btn': true,
                    'q-badge': true,
                    'q-tooltip': true,
                    'q-icon': true
                }
            },
            props: {
                classId: 'C1',
                subject: 'Religione',
                isReligion: true,
                date: '2025-01-20',
                type: 'Oral'
            }
        })

        religionWrapper.vm.entryData['S1'].value = 'Ottimo'
        await religionWrapper.vm.saveLine('S1')

        expect(gradeService.saveGrade).toHaveBeenCalledWith(expect.objectContaining({
            grade_value: 10,
            grade_type: 'judgment',
            description: expect.stringContaining('Ottimo')
        }))
    })
})
