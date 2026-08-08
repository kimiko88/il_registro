import { defineStore } from 'pinia'
import api from '../services/api'
import { ref, computed } from 'vue'

export const useParentStore = defineStore('parent', () => {
    const children = ref([])
    const selectedChildId = ref(localStorage.getItem('selectedChildId') || null)
    const loading = ref(false)
    const error = ref(null)

    const selectedChild = computed(() => {
        if (!selectedChildId.value) return children.value[0] || null
        return children.value.find(c => c.id === selectedChildId.value) || children.value[0] || null
    })

    async function fetchChildren() {
        loading.value = true
        error.value = null
        try {
            const response = await api.get('/users/me/children')
            // Normalize API snake_case fields to camelCase only when the raw
            // API fields are present; preserve already-mapped values otherwise.
            children.value = (response.data || []).map(c => ({
                ...c,
                firstName: c.first_name ?? c.firstName,
                lastName: c.last_name ?? c.lastName,
                schoolName: c.school_name ?? c.schoolName,
                className: c['class'] ?? c.className
            }))

            // Validate the stored selectedChildId still belongs to this user's children
            if (selectedChildId.value) {
                const stillValid = children.value.find(c => c.id === selectedChildId.value)
                if (!stillValid) {
                    selectedChildId.value = children.value[0]?.id || null
                    if (selectedChildId.value) {
                        localStorage.setItem('selectedChildId', selectedChildId.value)
                    } else {
                        localStorage.removeItem('selectedChildId')
                    }
                }
            } else if (children.value.length > 0) {
                selectedChildId.value = children.value[0].id
                localStorage.setItem('selectedChildId', selectedChildId.value)
            }
        } catch (err) {
            console.error('Failed to fetch children:', err)
            children.value = []
            error.value = 'Failed to load children'
        } finally {
            loading.value = false
        }
    }

    function selectChild(id) {
        if (!id) return
        if (children.value.length > 0) {
            const exists = children.value.find(c => c.id === id)
            if (!exists) return
        }
        selectedChildId.value = id
        localStorage.setItem('selectedChildId', id)
    }

    function reset() {
        children.value = []
        selectedChildId.value = null
        localStorage.removeItem('selectedChildId')
    }

    return {
        children,
        selectedChildId,
        selectedChild,
        loading,
        error,
        fetchChildren,
        selectChild,
        reset
    }
})
