import { defineStore } from 'pinia';

export const usePCTOStore = defineStore('pcto', {
    state: () => ({
        projects: [],
        loading: false,
        totalHours: 0
    }),

    actions: {
        async fetchProjects() {
            this.loading = true;
            try {
                await new Promise(resolve => setTimeout(resolve, 500));
                this.projects = [
                    { id: 1, title: 'Company Internship', company: 'Tech Corp', hours: 40, status: 'Completed' },
                    { id: 2, title: 'Library Assistant', company: 'City Library', hours: 20, status: 'In Progress' }
                ];
                this.totalHours = 60;
            } finally {
                this.loading = false;
            }
        }
    }
});
