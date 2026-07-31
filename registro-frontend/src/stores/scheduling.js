import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export const useSchedulingStore = defineStore('scheduling', () => {
    const generatedSchedule = ref([])
    const conflicts = ref([])
    const isLoading = ref(false)
    const error = ref(null)

    async function generateSchedule(constraints) {
        isLoading.value = true
        error.value = null
        conflicts.value = []
        try {
            const response = await api.post('/scheduling/generate', constraints)
            generatedSchedule.value = response.data || []
            return generatedSchedule.value
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
            const response = await api.post('/scheduling/validate', { slots })
            conflicts.value = response.data?.conflicts || []
            return conflicts.value
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
