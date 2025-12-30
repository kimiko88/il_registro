import { defineStore } from 'pinia';

export const useParentStore = defineStore('parent', {
    state: () => ({
        profile: null,
        loading: false,
        error: null,
        notifications: []
    }),

    getters: {
        fullName: (state) => state.profile ? `${state.profile.firstName} ${state.profile.lastName}` : '',
        isAuthenticated: (state) => !!state.profile
    },

    actions: {
        async fetchProfile() {
            this.loading = true;
            try {
                await new Promise(resolve => setTimeout(resolve, 500));
                // Mock Data
                this.profile = {
                    id: 'p1',
                    firstName: 'Giulia',
                    lastName: 'Mancini',
                    email: 'giulia.mancini@email.com',
                    avatar: 'https://cdn.quasar.dev/img/avatar2.jpg'
                };
            } finally {
                this.loading = false;
            }
        },

        async fetchNotifications() {
            this.notifications = [
                { id: 1, title: 'New Grade', message: 'Mario received a grade', date: '2025-01-20', type: 'info' },
                { id: 2, title: 'Meeting', message: 'Colloquio tomorrow', date: '2025-01-21', type: 'warning' }
            ];
        }
    }
});
