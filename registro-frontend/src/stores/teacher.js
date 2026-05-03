import { defineStore } from 'pinia';
import authService from 'src/services/authService';

export const useTeacherStore = defineStore('teacher', {
    state: () => ({
        profile: null,
        loading: false,
        error: null,
        notifications: []
    }),

    getters: {
        isAuthenticated: (state) => !!state.profile,
        isCoordinator: (state) => !!state.profile?.is_coordinator || !!state.profile?.isCoordinator,
        fullName: (state) => state.profile ? `${state.profile.firstName || state.profile.first_name} ${state.profile.lastName || state.profile.last_name}` : ''
    },

    actions: {
        async fetchProfile() {
            this.loading = true;
            try {
                const userData = await authService.getCurrentUser();
                this.profile = userData;
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
