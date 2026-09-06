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
    }
}

export const TEACHER_CACHE_TTL = 2 * 60 * 1000 // 2 minutes

export const useTeacherStore = defineStore('teacher', {
    state: () => ({
        profile: null,
        loading: false,
        error: null,
        notifications: [],
        pendingJustifications: 0,
        upcomingColloqui: 0,
        _lastFetchProfile: 0,
        _lastFetchNotifications: 0,
        _lastFetchJustifications: 0,
        _lastFetchColloqui: 0
    }),

    getters: {
        isAuthenticated: (state) => !!state.profile,
        isCoordinator: (state) => !!state.profile?.is_coordinator,
        fullName: (state) => state.profile
            ? `${state.profile.first_name} ${state.profile.last_name}`.trim()
            : ''
    },

    actions: {
        invalidateCache() {
            this._lastFetchProfile = 0
            this._lastFetchNotifications = 0
            this._lastFetchJustifications = 0
            this._lastFetchColloqui = 0
        },

        async fetchProfile(options = {}) {
            const isFresh = !options.force && this.profile && (Date.now() - this._lastFetchProfile < TEACHER_CACHE_TTL)
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
                this.error = err.response?.data?.error || err.userMessage || err.message || 'Error fetching teacher profile'
                console.error('Error fetching teacher profile:', err)
            } finally {
                this.loading = false
            }
        },

        async fetchNotifications(options = {}) {
            const isFresh = !options.force && this.notifications.length > 0 && (Date.now() - this._lastFetchNotifications < TEACHER_CACHE_TTL)
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
        },

        async fetchPendingJustifications(options = {}) {
            const isFresh = !options.force && this._lastFetchJustifications > 0 && (Date.now() - this._lastFetchJustifications < TEACHER_CACHE_TTL)
            if (isFresh) {
                return this.pendingJustifications
            }

            try {
                const response = await api.get('/attendance/pending-justifications')
                const items = Array.isArray(response.data)
                    ? response.data
                    : (response.data?.items || [])
                this.pendingJustifications = items.length
                this._lastFetchJustifications = Date.now()
                return this.pendingJustifications
            } catch (err) {
                console.error('Error fetching pending justifications:', err)
                this.pendingJustifications = 0
                return 0
            }
        },

        async fetchUpcomingColloqui(options = {}) {
            const isFresh = !options.force && this._lastFetchColloqui > 0 && (Date.now() - this._lastFetchColloqui < TEACHER_CACHE_TTL)
            if (isFresh) {
                return this.upcomingColloqui
            }

            try {
                const response = await api.get('/colloqui/slots/my?upcoming=true')
                const items = Array.isArray(response.data)
                    ? response.data
                    : (response.data?.items || [])
                this.upcomingColloqui = items.filter(c => c.booked).length
                this._lastFetchColloqui = Date.now()
                return this.upcomingColloqui
            } catch (err) {
                console.error('Error fetching upcoming colloqui:', err)
                this.upcomingColloqui = 0
                return 0
            }
        }
    }
})
