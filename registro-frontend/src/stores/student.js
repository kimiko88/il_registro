import { defineStore } from 'pinia';

export const useStudentStore = defineStore('student', {
    state: () => ({
        profile: null,
        loading: false,
        error: null,
        notifications: []
    }),

    getters: {
        fullName: (state) => state.profile ? `${state.profile.firstName} ${state.profile.lastName}` : '',
        className: (state) => state.profile ? state.profile.className : '',
        isAuthenticated: (state) => !!state.profile
    },

    actions: {
        async fetchProfile() {
            this.loading = true;
            try {
                await new Promise(resolve => setTimeout(resolve, 500));
                // Mock Data
                this.profile = {
                    id: 's1',
                    firstName: 'Marco',
                    lastName: 'Rossi',
                    className: '5A Scientifico',
                    email: 'marco.rossi@studenti.school.it',
                    avatar: 'https://cdn.quasar.dev/img/boy-avatar.png'
                };
            } catch (err) {
                this.error = err.message;
            } finally {
                this.loading = false;
            }
        },

        async fetchNotifications() {
            // Mock notifications
            this.notifications = [
                { id: 1, title: 'New Grade', message: 'Math grade posted', date: '2025-01-20', type: 'info', read: false },
                { id: 2, title: 'Attendance Alert', message: 'You were marked absent yesterday', date: '2025-01-19', type: 'warning', read: false }
            ];
        }
    }
});
