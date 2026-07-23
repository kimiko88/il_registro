import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useGradesStore } from '@/stores/grades'
import { gradeService } from '@/services/gradeService'

// Mock gradeService
vi.mock('@/services/gradeService', () => ({
    gradeService: {
        getByClass: vi.fn(),
        getMyGrades: vi.fn(),
        saveGrade: vi.fn(),
        updateGrade: vi.fn(),
        deleteGrade: vi.fn()
    }
}))

describe('Grades Store', () => {
    let store

    beforeEach(() => {
        setActivePinia(createPinia())
        store = useGradesStore()
        vi.useFakeTimers()
    })

    afterEach(() => {
        vi.useRealTimers()
    })

    it('initializes with default state', () => {
        expect(store.grades).toBeNull()
        expect(store.loading).toBe(false)
        expect(store.subjects).toEqual([])
    })

    it('fetches grades for teacher', async () => {
        gradeService.getByClass.mockResolvedValue({ data: [{ id: 'g1', subject: 'Math' }] })
        const promise = store.fetchGrades('c1', 'Math')
        expect(store.loading).toBe(true)
        await promise
        expect(store.loading).toBe(false)
        expect(store.grades.length).toBeGreaterThan(0)
    })

    it('fetches grades for student', async () => {
        gradeService.getMyGrades.mockResolvedValue({ data: [{ id: 'g1', subject_id: 'Mathematics', grade_value: 8 }] })
        const promise = store.fetchMyGrades()
        expect(store.loading).toBe(true)
        await promise
        expect(store.loading).toBe(false)
        expect(store.grades).toHaveLength(1)
    })

    it('adds a grade', async () => {
        const gradeData = { studentId: 's5', value: 9 }
        gradeService.saveGrade.mockResolvedValue({ data: { id: 'g_new', ...gradeData } })
        gradeService.getByClass.mockResolvedValue({ data: { students: [] } })
        
        // Mock a previous fetch so it refetches
        store._lastClassId = 'c1'
        store._lastSubjectId = 'Math'
        
        const promise = store.addGrade(gradeData)
        const newGrade = await promise
        expect(newGrade.id).toBeDefined()
        expect(gradeService.getByClass).toHaveBeenCalledWith('c1', 'Math')
    })

    it('updates a grade', async () => {
        store.grades = { semesters: [{ grades: [{ id: 'g1', value: 5 }] }] }
        gradeService.updateGrade.mockResolvedValue({ data: { id: 'g1', value: 6 } })
        gradeService.getByClass.mockResolvedValue({ data: { students: [] } })
        
        store._lastClassId = 'c1'
        store._lastSubjectId = 'Math'

        await store.updateGrade('g1', { value: 6 })
        expect(gradeService.updateGrade).toHaveBeenCalledWith('g1', { value: 6 })
        expect(gradeService.getByClass).toHaveBeenCalledWith('c1', 'Math')
    })

    it('deletes a grade', async () => {
        store.grades = { students: [{ grades: [{ id: 'g1' }, { id: 'g2' }] }] }
        gradeService.deleteGrade.mockResolvedValue({})
        gradeService.getByClass.mockResolvedValue({ data: { students: [] } })
        
        // Mock a previous fetch so it refetches
        store._lastClassId = 'c1'
        store._lastSubjectId = 'Math'

        await store.deleteGrade('g1')
        expect(gradeService.deleteGrade).toHaveBeenCalledWith('g1')
        expect(gradeService.getByClass).toHaveBeenCalledWith('c1', 'Math')
    })

    it('calculates class average', () => {
        expect(store.classAverage).toBe(0)
        store.grades = { students: [{ grades: [{ grade_value: 6 }, { grade_value: 8 }] }] }
        expect(store.classAverage).toBe(7)
    })

    it('gets grades by student', () => {
        store.grades = { students: [{ student_id: 's1', grades: [{ id: '1' }] }, { student_id: 's2', grades: [{ id: '2' }] }] }
        expect(store.getGradesByStudent('s1').length).toBe(1)
        expect(store.getGradesByStudent('s3')).toEqual([])
    })
})
