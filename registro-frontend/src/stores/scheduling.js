import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/services/api'

export const useSchedulingStore = defineStore('scheduling', () => {
    const generatedSchedule = ref([])
    const conflicts = ref([])
    const isLoading = ref(false)

    async function generateSchedule(constraints) {
        isLoading.value = true
        conflicts.value = []
        try {
            // Mock API call for now until backend endpoint is officially wired to a route
            // In real integration: const res = await api.post('/scheduling/generate', constraints)
            // returning stub for TDD

            // Simulating API delay
            await new Promise(r => setTimeout(r, 500))

            // Stub result
            const mockResult = [
                { id: '1', teacher_id: 'T1', class_id: '1A', subject_id: 'Math', day: 1, hour: 8 },
                { id: '2', teacher_id: 'T2', class_id: '1B', subject_id: 'Hist', day: 1, hour: 8 },
            ]
            generatedSchedule.value = mockResult
            return mockResult
        } catch (error) {
            console.error('Failed to generate schedule', error)
            throw error
        } finally {
            isLoading.value = false
        }
    }

    async function validateSchedule(slots) {
        isLoading.value = true
        try {
            // Mock API call
            // const res = await api.post('/scheduling/validate', slots)
            await new Promise(r => setTimeout(r, 300))

            // Stub simplified conflict check
            const localConflicts = []
            // naive check
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
        } catch (error) {
            console.error('Validation failed', error)
            throw error
        } finally {
            isLoading.value = false
        }
    }

    return {
        generatedSchedule,
        conflicts,
        isLoading,
        generateSchedule,
        validateSchedule
    }
})
