import { defineStore } from 'pinia';
// import { api } from 'src/boot/axios'; // Removed invalid import

export const useTeacherStore = defineStore('teacher', {
    state: () => ({
        profile: null,
        loading: false,
        error: null,
        notifications: []
    }),

    getters: {
        isAuthenticated: (state) => !!state.profile,
        isCoordinator: (state) => !!state.profile?.isCoordinator,
        fullName: (state) => state.profile ? `${state.profile.firstName} ${state.profile.lastName}` : ''
    },

    actions: {
        async fetchProfile() {
            this.loading = true;
            try {
                // Mock API call
                // const response = await api.get('/teacher/profile');
                // this.profile = response.data;

                // Mock Data
                await new Promise(resolve => setTimeout(resolve, 500));
                this.profile = {
                    id: 'te-1',
                    firstName: 'Mario',
                    lastName: 'Rossi',
                    email: 'mario.rossi@school.it',
                    subjects: ['Mathematics', 'Physics'],
                    isCoordinator: true, // For testing Coordinator view
                    avatar: 'https://cdn.quasar.dev/img/avatar.png'
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
                { id: 1, title: 'Meeting Reminder', message: 'Colloquio with Parent A at 10:00', read: false, type: 'info' },
                { id: 2, title: 'Document Signed', message: 'Director signed PDP for Student B', read: false, type: 'positive' }
            ];
        }
    }
});
