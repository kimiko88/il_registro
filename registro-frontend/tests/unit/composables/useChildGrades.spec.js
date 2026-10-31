import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { nextTick } from 'vue'
import { useChildrenStore } from '@/stores/children'
import { useGradesStore } from '@/stores/grades'
import { useChildGrades } from '@/composables/useChildGrades'

describe('useChildGrades Composable', () => {
    let pinia

    beforeEach(() => {
        pinia = createPinia()
        setActivePinia(pinia)
        vi.clearAllMocks()
    })

    it('returns selectedChild, gradesBySubject, averages, getTrend, loading', () => {
        const composable = useChildGrades()

        expect(composable.selectedChild).toBeDefined()
        expect(composable.gradesBySubject).toBeDefined()
        expect(composable.averages).toBeDefined()
        expect(typeof composable.getTrend).toBe('function')
        expect(composable.loading).toBeDefined()
    })

    it('watches selectedChildId and calls fetchMyGrades when set', async () => {
        const childrenStore = useChildrenStore()
        const gradesStore = useGradesStore()
        gradesStore.fetchMyGrades = vi.fn().mockResolvedValue({})

        childrenStore.children = [
            { id: 'child-1', firstName: 'Lorenzo', lastName: 'Rossi' },
            { id: 'child-2', firstName: 'Chiara', lastName: 'Rossi' }
        ]

        useChildGrades()

        childrenStore.selectedChildId = 'child-1'
        await nextTick()

        expect(gradesStore.fetchMyGrades).toHaveBeenCalledWith('child-1')

        childrenStore.selectedChildId = 'child-2'
        await nextTick()

        expect(gradesStore.fetchMyGrades).toHaveBeenCalledWith('child-2')
    })

    it('exposes selectedChild matching childrenStore.selectedChild', async () => {
        const childrenStore = useChildrenStore()
        childrenStore.children = [
            { id: 'child-1', firstName: 'Lorenzo', lastName: 'Rossi' }
        ]
        childrenStore.selectedChildId = 'child-1'

        const { selectedChild } = useChildGrades()
        expect(selectedChild.value).toBeDefined()
        expect(selectedChild.value.firstName).toBe('Lorenzo')
    })

    it('groups grades by subject and calculates averages correctly', async () => {
        const gradesStore = useGradesStore()
        gradesStore.grades = {
            semesters: [
                {
                    grades: [
                        { subject_id: 'MAT', grade_value: 8.0 },
                        { subject_id: 'MAT', grade_value: 6.0 },
                        { subject_id: 'ITA', grade_value: 9.0 }
                    ]
                }
            ]
        }

        const { gradesBySubject, averages, getTrend } = useChildGrades()

        expect(gradesBySubject.value['MAT']).toHaveLength(2)
        expect(gradesBySubject.value['ITA']).toHaveLength(1)
        expect(averages.value['MAT']).toBe(7.0)
        expect(averages.value['ITA']).toBe(9.0)

        // MAT went from 8 to 6 -> down
        expect(getTrend('MAT')).toBe('down')
        // Single grade in ITA -> stable
        expect(getTrend('ITA')).toBe('stable')
    })

    it('returns stable trend when grade trajectory is unchanged or insufficient', () => {
        const gradesStore = useGradesStore()
        gradesStore.grades = {
            semesters: [
                {
                    grades: [
                        { subject_id: 'STORIA', grade_value: 7.0 },
                        { subject_id: 'STORIA', grade_value: 7.0 }
                    ]
                }
            ]
        }

        const { getTrend } = useChildGrades()
        expect(getTrend('STORIA')).toBe('stable')
        expect(getTrend('NON_EXISTENT')).toBe('stable')
    })
})

