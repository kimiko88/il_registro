import { defineStore } from 'pinia';
import authService from 'src/services/authService';
import api from '../services/api';

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

export const useTeacherStore = defineStore('teacher', {
    state: () => ({
        profile: null,
        loading: false,
        error: null,
        notifications: [],
        pendingJustifications: 0,
        upcomingColloqui: 0
    }),

    getters: {
        isAuthenticated: (state) => !!state.profile,
        isCoordinator: (state) => !!state.profile?.is_coordinator,
        fullName: (state) => state.profile
            ? `${state.profile.first_name} ${state.profile.last_name}`.trim()
            : ''
    },

    actions: {
        async fetchProfile() {
            this.loading = true;
            try {
                const userData = await authService.getCurrentUser();
                this.profile = normalizeProfile(userData);
            } catch (err) {
                this.error = err.message;
                console.error('Error fetching teacher profile:', err);
            } finally {
                this.loading = false;
            }
        },

        async fetchNotifications() {
            try {
                const response = await api.get('/users/me/notifications')
                this.notifications = response.data || []
            } catch (err) {
                console.error('Error fetching notifications:', err)
                this.notifications = []
            }
        },

        async fetchPendingJustifications() {
            try {
                const response = await api.get('/attendance/justifications/pending')
                const items = Array.isArray(response.data)
                    ? response.data
                    : (response.data?.items || [])
                this.pendingJustifications = items.length
            } catch (err) {
                console.error('Error fetching pending justifications:', err)
                this.pendingJustifications = 0
            }
        },

        async fetchUpcomingColloqui() {
            try {
                const response = await api.get('/colloqui/my-slots?upcoming=true')
                const items = Array.isArray(response.data)
                    ? response.data
                    : (response.data?.items || [])
                this.upcomingColloqui = items.filter(c => c.booked).length
            } catch (err) {
                console.error('Error fetching upcoming colloqui:', err)
                this.upcomingColloqui = 0
            }
        }
    }
});
