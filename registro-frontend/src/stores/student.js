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

export const STUDENT_CACHE_TTL = 2 * 60 * 1000 // 2 minutes

export const useStudentStore = defineStore('student', {
    state: () => ({
        profile: null,
        loading: false,
        error: null,
        notifications: [],
        _lastFetchProfile: 0,
        _lastFetchNotifications: 0
    }),

    getters: {
        fullName: (state) => state.profile
            ? `${state.profile.first_name} ${state.profile.last_name}`.trim()
            : '',
        className: (state) => state.profile?.class_name || '',
        isAuthenticated: (state) => !!state.profile
    },

    actions: {
        invalidateCache() {
            this._lastFetchProfile = 0
            this._lastFetchNotifications = 0
        },

        async fetchProfile(options = {}) {
            const isFresh = !options.force && this.profile && (Date.now() - this._lastFetchProfile < STUDENT_CACHE_TTL)
            if (isFresh) {
                return this.profile
            }

            this.loading = true
            this.error = null
            try {
                const userData = await authService.getCurrentUser()
                this.profile = normalizeProfile(userData)
                this._lastFetchProfile = Date.now()
                return this.profile
            } catch (err) {
                this.error = err.response?.data?.error || err.userMessage || err.message || 'Error fetching student profile'
                console.error('Error fetching student profile:', err)
            } finally {
                this.loading = false
            }
        },

        async fetchNotifications(options = {}) {
            const isFresh = !options.force && this.notifications.length > 0 && (Date.now() - this._lastFetchNotifications < STUDENT_CACHE_TTL)
            if (isFresh) {
                return this.notifications
            }

            try {
                const response = await api.get('/notifications')
                this.notifications = response.data || []
                this._lastFetchNotifications = Date.now()
                return this.notifications
            } catch (err) {
                console.error('Error fetching notifications:', err)
                this.notifications = []
                return []
            }
        }
    }
})
