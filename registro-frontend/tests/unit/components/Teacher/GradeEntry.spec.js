import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar, Notify } from 'quasar'
import GradeEntry from '@/components/Teacher/GradeEntry.vue'

// Mock Quasar's useQuasar
const mockNotify = vi.hoisted(() => vi.fn())
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            notify: mockNotify
        }),
        Notify: {
            create: mockNotify
        }
    }
})

describe('GradeEntry.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        wrapper = mount(GradeEntry, {
            global: {
                plugins: [
                    [Quasar, {}]
                ],
                // Stub complex Quasar components
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-table': {
                        template: '<div><slot name="header" /><slot name="body" :row="rows[0]" /><slot name="bottom" /></div>',
                        props: ['rows']
                    },
                    'q-tr': { template: '<tr><slot /></tr>' },
                    'q-th': { template: '<th><slot /></th>' },
                    'q-td': { template: '<td><slot /></td>' },
                    'q-input': {
                        template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
                        props: ['modelValue']
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
        vi.useFakeTimers()
        wrapper.vm.saveAll()
        expect(wrapper.vm.loading).toBe(true)

        vi.advanceTimersByTime(1000)

        expect(wrapper.vm.loading).toBe(false)
        expect(mockNotify).toHaveBeenCalled()
        vi.useRealTimers()
    })
})
