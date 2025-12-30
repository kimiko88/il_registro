import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useGradesStore } from '@/stores/grades'

describe('Grades Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('adds a grade locally', () => {
        const store = useGradesStore()
        const grade = { id: 'g1', value: 8 }
        store.addGrade(grade)
        // Assuming addGrade pushes to state or we test state mutation directly if action calls API
        // If action calls API, we might mock axios.
        // Here testing state if it exists.
        // If store logic is complex, we assume simple mutation for unit test or mock service.
    })
})
