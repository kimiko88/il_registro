import { defineStore } from 'pinia';
import api from '../services/api';

export const useSettingsStore = defineStore('settings', {
    state: () => ({
        config: {
            email: {},
            backup: {},
            auth: {},
        },
        loading: false,
    }),
    actions: {
        async fetchSettings() {
            this.loading = true;
            try {
                const response = await api.get('/settings');
                this.config = response.data;
            } finally {
                this.loading = false;
            }
        },
        async updateSettings(section, data) {
            await api.patch('/settings', { [section]: data });
            this.config[section] = { ...this.config[section], ...data };
        }
    }
});
