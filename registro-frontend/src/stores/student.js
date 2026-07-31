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
            this.loading = true;
            try {
                const userData = await authService.getCurrentUser();
                this.profile = normalizeProfile(userData);
            } catch (err) {
                this.error = err.message;
                console.error('Error fetching profile:', err);
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
        }
    }
});
