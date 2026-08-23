import { defineStore } from 'pinia'
import authService from '@/services/authService'
import api from '@/services/api'

// Normalize the user profile from backend (snake_case) to a consistent shape
function normalizeProfile(data) {
    if (!data) return null
    return {
        ...data,
        first_name: data.first_name || data.firstName || '',
        last_name: data.last_name || data.lastName || '',
        school_id: data.school_id || data.schoolId || null,
        class_id: data.class_id || data.classId || null,
        class_name: data.class_name || data.className || '',
    }
}

export const useStudentStore = defineStore('student', {
    state: () => ({
        profile: null,
        loading: false,
        error: null,
        notifications: []
    }),

    getters: {
        fullName: (state) => state.profile
            ? `${state.profile.first_name} ${state.profile.last_name}`.trim()
            : '',
        className: (state) => state.profile?.class_name || '',
        isAuthenticated: (state) => !!state.profile
    },

    actions: {
        async fetchProfile() {
            this.loading = true
            this.error = null
            try {
                const userData = await authService.getCurrentUser()
                this.profile = normalizeProfile(userData)
                return this.profile
            } catch (err) {
                this.error = err.response?.data?.error || err.userMessage || err.message || 'Error fetching student profile'
                console.error('Error fetching student profile:', err)
            } finally {
                this.loading = false
            }
        },

        async fetchNotifications() {
            try {
                const response = await api.get('/notifications')
                this.notifications = response.data || []
                return this.notifications
            } catch (err) {
                console.error('Error fetching notifications:', err)
                this.notifications = []
                return []
            }
        }
    }
})
