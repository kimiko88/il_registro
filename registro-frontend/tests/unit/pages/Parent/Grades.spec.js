import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import ParentGrades from '@/pages/parent/Grades.vue'

// Mock services
const { mockGetChildGrades } = vi.hoisted(() => ({
    mockGetChildGrades: vi.fn(),
}))

vi.mock('src/services/gradeService', () => ({
    gradeService: {
        getChildGrades: mockGetChildGrades
    }
}))

describe('Parent/Grades.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        mockGetChildGrades.mockResolvedValue({
            data: {
                semesters: [
                    {
                        semester: 1,
                        grades: [
                            { id: 1, subject_id: 'Math', grade_value: 8, date: '2023-01-01', grade_type: 'oral', description: 'Good' }
                        ]
                    }
                ]
            }
        })

        wrapper = mount(ParentGrades, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            parent: {
                                children: [{ id: 'child1', firstName: 'Mario' }],
                                selectedChildId: 'child1'
                            }
                        },
                        stubActions: false
                    })
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-select': true,
                    'q-table': {
                        template: '<table><tbody><tr v-for="r in rows" :key="r.id"><td>{{r.value}}</td></tr></tbody></table>',
                        props: ['rows', 'columns']
                    },
                    'q-badge': true,
                    'q-td': true,
                    'q-btn': true
                }
            }
        })
    })

    it('fetches grades when child is selected', () => {
        expect(mockGetChildGrades).toHaveBeenCalledWith('child1')
    })

    it('renders grades in table', async () => {
        await wrapper.vm.$nextTick()
        await new Promise(resolve => setTimeout(resolve, 10))
        await wrapper.vm.$nextTick()

        expect(wrapper.vm.currentGrades).toHaveLength(1)
        expect(wrapper.vm.currentGrades[0].value).toBe(8)
    })

    it('filters by period', async () => {
        await wrapper.vm.$nextTick()
        await new Promise(resolve => setTimeout(resolve, 10))
        await wrapper.vm.$nextTick() // Default 'Primo Quadrimestre' -> 1 grade

        expect(wrapper.vm.currentGrades).toHaveLength(1)

        // Switch to Second Semester
        wrapper.vm.period = 'Secondo Quadrimestre'
        await wrapper.vm.$nextTick()

        // Mock data only has semester 1.
        expect(wrapper.vm.currentGrades).toHaveLength(0)
    })
})
