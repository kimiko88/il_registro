import { defineStore } from 'pinia';
import authService from 'src/services/authService';

export const useStudentStore = defineStore('student', {
    state: () => ({
        profile: null,
        loading: false,
        error: null,
        notifications: []
    }),

    getters: {
        fullName: (state) => state.profile ? `${state.profile.firstName || state.profile.first_name} ${state.profile.lastName || state.profile.last_name}` : '',
        className: (state) => state.profile ? state.profile.className || state.profile.class_name : '',
        isAuthenticated: (state) => !!state.profile
    },

    actions: {
        async fetchProfile() {
            this.loading = true;
            try {
                const userData = await authService.getCurrentUser();
                this.profile = userData;
            } catch (err) {
                this.error = err.message;
                console.error("Error fetching profile:", err);
            } finally {
                this.loading = false;
            }
        },

        async fetchNotifications() {
            // Notifications are currently handled via Communications or WebSocket
            // For now, keep it empty or fetch from communications
            this.notifications = [];
        }
    }
});
