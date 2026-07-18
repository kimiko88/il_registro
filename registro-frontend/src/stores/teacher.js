import { defineStore } from 'pinia';
import authService from 'src/services/authService';

// Normalize the user profile from backend (snake_case) to a consistent shape
function normalizeProfile(data) {
    if (!data) return null
    return {
        ...data,
        // Ensure both naming conventions work, preferring snake_case source
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
        notifications: []
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
                console.error("Error fetching teacher profile:", err);
            } finally {
                this.loading = false;
            }
        },

        async fetchNotifications() {
            // Notifications are currently empty or fetched via Communications
            this.notifications = [];
        }
    }
});
