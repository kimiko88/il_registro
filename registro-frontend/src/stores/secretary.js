import { defineStore } from 'pinia';
import api from '../services/api';

export const useSecretaryStore = defineStore('secretary', {
    state: () => ({
        stats: {
            pendingDocuments: 0,
            pendingRequests: 0,
            totalStudents: 0
        },
        loading: false
    }),
    actions: {
        async fetchDashboardStats() {
            this.loading = true;
            try {
                const res = await api.get('/secretary/stats');
                this.stats = res.data;
            } catch (e) {
                console.error(e);
            } finally {
                this.loading = false;
            }
        }
    }
});
