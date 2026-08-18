import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import AttendanceMarking from '@/components/Teacher/AttendanceMarking.vue'
import { attendanceService } from '@/services/attendanceService'

// Mock Quasar
vi.mock('quasar', async () => {
    const actual = await vi.importActual('quasar')
    return {
        ...actual,
        useQuasar: () => ({
            notify: vi.fn(),
            dialog: vi.fn(() => ({
                onOk: vi.fn(cb => ({ onCancel: vi.fn(), ...cb() }))
            }))
        })
    }
})

// Mock attendance service
vi.mock('@/services/attendanceService', () => {
    const mockData = {
        data: [
            { student_id: 's1', first_name: 'Mario', last_name: 'Rossi', status: 'present' },
            { student_id: 's2', first_name: 'Luigi', last_name: 'Verdi', status: 'absent' }
        ]
    }
    return {
        attendanceService: {
            getByClass: vi.fn().mockResolvedValue(mockData),
            getAttendance: vi.fn().mockResolvedValue(mockData),
            markAttendance: vi.fn().mockResolvedValue({}),
            recordBulk: vi.fn().mockResolvedValue({}),
            justify: vi.fn().mockResolvedValue({})
        }
    }
})

describe('Teacher/AttendanceMarking.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        wrapper = mount(AttendanceMarking, {
            global: {
                plugins: [
                    Quasar,
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            classes: {
                                classes: [{ id: 'c1', name: '1A', section: 'A' }]
                            }
                        }
                    })
                ],
                stubs: {
                    'q-select': true,
                    'q-input': true,
                    'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
                    'q-table': {
                        template: '<div class="q-table-stub"><slot name="header" :props="{}" /><template v-for="(row, index) in rows" :key="index"><slot name="body" :row="row" :rowIndex="index" :props="{ row, rowIndex: index }" /></template></div>',
                        props: ['rows']
                    },
                    'q-tr': { template: '<tr><slot /></tr>' },
                    'q-td': { template: '<td><slot /></td>' },
                    'q-btn-toggle': {
                        template: '<div class="toggle-stub" @click="$emit(\'update:modelValue\', \'late\')"></div>',
                        props: ['modelValue']
                    },
                    'q-banner': true,
                    'q-chip': true,
                    'q-dialog': true
                }
            }
        })
    })

    it('renders correctly and toggles status', async () => {
        // Wait for onMounted and fetch
        await new Promise(resolve => setTimeout(resolve, 0))
        await wrapper.vm.$nextTick()

        // Check if class was auto-selected
        expect(wrapper.vm.selectedClassId).toBe('c1')
        
        // Verify table renders (via our stub class)
        expect(wrapper.find('.q-table-stub').exists()).toBe(true)
        expect(wrapper.vm.students.length).toBe(2)

        // Find toggle for first student
        const firstToggle = wrapper.find('.toggle-stub')
        await firstToggle.trigger('click')

        // Verify local state updated in attendanceMap
        expect(wrapper.vm.attendanceMap['s1'].status).toBe('late')
    })

    it('saves attendance changes', async () => {
        await new Promise(resolve => setTimeout(resolve, 0))
        
        // Change something to make it hasChanges = true
        wrapper.vm.attendanceMap['s1'].status = 'absent'
        
        await wrapper.vm.saveAll()
        
        expect(attendanceService.recordBulk).toHaveBeenCalled()
    })
})
