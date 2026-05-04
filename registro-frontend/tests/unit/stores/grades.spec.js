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
        expect(store.grades).toEqual([])
        expect(store.loading).toBe(false)
        expect(store.subjects).toContain('Mathematics')
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
        const promise = store.addGrade(gradeData)
        const newGrade = await promise
        expect(newGrade.id).toBeDefined()
        expect(store.grades).toContainEqual(newGrade)
    })

    it('updates a grade', async () => {
        store.grades = [{ id: 'g1', value: 5 }]
        gradeService.updateGrade.mockResolvedValue({ data: { id: 'g1', value: 6 } })
        await store.updateGrade('g1', { value: 6 })
        expect(store.grades[0].value).toBe(6)
    })

    it('deletes a grade', async () => {
        store.grades = [{ id: 'g1' }, { id: 'g2' }]
        gradeService.deleteGrade.mockResolvedValue({})
        await store.deleteGrade('g1')
        expect(store.grades.length).toBe(1)
        expect(store.grades[0].id).toBe('g2')
    })

    it('calculates class average', () => {
        expect(store.classAverage).toBe(0)
        store.grades = [{ grade_value: 6 }, { grade_value: 8 }]
        expect(store.classAverage).toBe('7.0')
    })

    it('gets grades by student', () => {
        store.grades = [{ studentId: 's1' }, { studentId: 's2' }]
        expect(store.getGradesByStudent('s1').length).toBe(1)
    })
})
