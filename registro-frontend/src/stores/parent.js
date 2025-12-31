import { defineStore } from 'pinia'
import { api } from '../boot/axios'
import { ref, computed } from 'vue'

export const useParentStore = defineStore('parent', () => {
    const children = ref([])
    const selectedChildId = ref(localStorage.getItem('selectedChildId') || null)
    const loading = ref(false)
    const error = ref(null)

    // Mock data for development if API fails or is empty
    const mockChildren = [
        {
            id: 'student-123',
            firstName: 'Luigi',
            lastName: 'Rossi',
            schoolName: 'Liceo Scientifico "A. Einstein"',
            class: '3A',
            avatar: null
        },
        {
            id: 'student-456',
            firstName: 'Maria',
            lastName: 'Rossi',
            schoolName: 'Liceo Scientifico "A. Einstein"',
            class: '1B',
            avatar: null
        }
    ]

    const selectedChild = computed(() => {
        if (!selectedChildId.value) return children.value[0] || null
        return children.value.find(c => c.id === selectedChildId.value) || children.value[0] || null
    })

    async function fetchChildren() {
        loading.value = true
        error.value = null
        try {
            // Try fetching from API
            // const response = await api.get('/users/me/children')
            // children.value = response.data

            // Fallback to mock for now
            await new Promise(r => setTimeout(r, 500))
            children.value = mockChildren

            // Select first if none selected
            if (!selectedChildId.value && children.value.length > 0) {
                selectedChildId.value = children.value[0].id
                localStorage.setItem('selectedChildId', selectedChildId.value)
            }
        } catch (err) {
            console.error('Failed to fetch children:', err)
            children.value = mockChildren // Fallback
            error.value = 'Failed to load children'
        } finally {
            loading.value = false
        }
    }

    function selectChild(id) {
        if (children.value.find(c => c.id === id)) {
            selectedChildId.value = id
            localStorage.setItem('selectedChildId', id)
        }
    }

    return {
        children,
        selectedChildId,
        selectedChild,
        loading,
        error,
        fetchChildren,
        selectChild
    }
})
