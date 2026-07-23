import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useSchedulingStore = defineStore('scheduling', () => {
    const generatedSchedule = ref([])
    const conflicts = ref([])
    const isLoading = ref(false)
    const error = ref(null)

    async function generateSchedule(_constraints) {
        isLoading.value = true
        error.value = null
        conflicts.value = []
        try {
            // Mock API delay
            await new Promise(r => setTimeout(r, 500))

            const mockResult = [
                { id: '1', teacher_id: 'T1', class_id: '1A', subject_id: 'Math', day: 1, hour: 8 },
                { id: '2', teacher_id: 'T2', class_id: '1B', subject_id: 'Hist', day: 1, hour: 8 },
            ]
            generatedSchedule.value = mockResult
            return mockResult
        } catch (err) {
            error.value = err.message || 'Failed to generate schedule'
            console.error('Failed to generate schedule', err)
            throw err
        } finally {
            isLoading.value = false
        }
    }

    async function validateSchedule(slots) {
        isLoading.value = true
        error.value = null
        try {
            await new Promise(r => setTimeout(r, 300))

            const localConflicts = []
            const viewed = new Set()
            slots.forEach(s => {
                const key = `${s.day}-${s.hour}-${s.teacher_id}`
                if (viewed.has(key)) {
                    localConflicts.push({ type: 'TEACHER_OVERLAP', description: `Teacher ${s.teacher_id} double booked` })
                }
                viewed.add(key)
            })

            conflicts.value = localConflicts
            return localConflicts
        } catch (err) {
            error.value = err.message || 'Validation failed'
            console.error('Validation failed', err)
            throw err
        } finally {
            isLoading.value = false
        }
    }

    return {
        generatedSchedule,
        conflicts,
        isLoading,
        error,
        generateSchedule,
        validateSchedule
    }
})
