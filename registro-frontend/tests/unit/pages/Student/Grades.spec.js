import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import StudentGrades from '@/pages/student/Grades.vue'

// Mock services
const { mockGetMyGrades } = vi.hoisted(() => ({
    mockGetMyGrades: vi.fn(),
}))

vi.mock('src/services/gradeService', () => ({
    gradeService: {
        getMyGrades: mockGetMyGrades
    }
}))

describe('Student/Grades.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        mockGetMyGrades.mockResolvedValue({
            data: {
                semesters: [
                    {
                        semester: 1,
                        grades: [
                            { id: 1, subject_id: 'Math', grade_value: 8, date: '2023-01-01', grade_type: 'oral', description: 'Test 1', semester: 1 },
                            { id: 2, subject_id: 'History', grade_value: 5, date: '2023-01-02', grade_type: 'written', description: 'Test 2', semester: 1 }
                        ]
                    }
                ]
            }
        })

        wrapper = mount(StudentGrades, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({ createSpy: vi.fn })
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-btn': true,
                    'q-select': true,
                    'q-table': {
                        template: '<table><tbody><tr v-for="r in rows" :key="r.id"><td>{{r.value}}</td></tr></tbody></table>',
                        props: ['rows', 'columns']
                    },
                    'q-linear-progress': true,
                    'q-badge': true,
                    'q-td': true
                }
            }
        })
    })

    it('fetches grades on mount', () => {
        expect(mockGetMyGrades).toHaveBeenCalled()
    })

    it('calculates averages correctly', async () => {
        await wrapper.vm.$nextTick()
        await new Promise(resolve => setTimeout(resolve, 10))
        await wrapper.vm.$nextTick()

        // Check computed property via vm if possible, or render result
        // wrapper.vm.subjectAverages is computed
        // Math: 8, History: 5
        const avgs = wrapper.vm.subjectAverages
        expect(avgs).toHaveLength(2)
        expect(avgs.find(a => a.name === 'Math').avg).toBe('8.0')
    })

    it('filters by semester', async () => {
        await wrapper.vm.$nextTick() // wait load
        await new Promise(resolve => setTimeout(resolve, 10))
        await wrapper.vm.$nextTick()

        // Default semester 1 -> 2 items (Math, History)
        expect(wrapper.vm.filteredGrades).toHaveLength(2)

        // Switch to semester 2
        wrapper.vm.filters.semester = 2
        await wrapper.vm.$nextTick()
        expect(wrapper.vm.filteredGrades).toHaveLength(0)
    })

    it('filters by period', async () => {
        await wrapper.vm.$nextTick()
        await new Promise(resolve => setTimeout(resolve, 10))
        await wrapper.vm.$nextTick()

        // Data dates: 2023-01-01.
        // Mock system date to be 2023-01-05 (so last week includes it)
        // Or just filtering logic: 'Ultimo Mese'
        // If current date is 2026, 2023 is old.
        // 'Ultimo Mese' -> fails

        wrapper.vm.filters.period = 'Ultimo Mese'
        await wrapper.vm.$nextTick()
        expect(wrapper.vm.filteredGrades).toHaveLength(0)

        // Reset
        wrapper.vm.filters.period = 'Tutti'
        await wrapper.vm.$nextTick()
        expect(wrapper.vm.filteredGrades).toHaveLength(2)
    })
})
