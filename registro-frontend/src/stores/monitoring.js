import { defineStore } from 'pinia';
import monitoringService from '../services/monitoringService';

export const useMonitoringStore = defineStore('monitoring', {
    state: () => ({
        health: null,
        metrics: [],
        activeUsers: 0,
        loading: false,
    }),
    actions: {
        async fetchHealth() {
            try {
                const response = await monitoringService.getSystemHealth();
                this.health = response.data;
            } catch (error) {
                console.error('Health check failed', error);
            }
        },
        async fetchMetrics(range = '24h') {
            try {
                const response = await monitoringService.getMetrics(range);
                this.metrics = response.data;
            } catch (error) {
                console.error('Metrics fetch failed', error);
            }
        },
        async fetchActiveUsers() {
            try {
                const response = await monitoringService.getActiveUsers();
                this.activeUsers = response.data.count;
            } catch (error) {
                console.error('Active users fetch failed', error);
            }
        }
    }
});
