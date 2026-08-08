import { defineStore } from 'pinia';
import api from '../services/api';

export const useCommunicationsStore = defineStore('communications', {
    state: () => ({
        communications: [],
        loading: false,
        error: null,
    }),

    actions: {
        async fetchCommunications(schoolID = null) {
            this.loading = true;
            this.error = null;
            try {
                const response = schoolID
                    ? await api.get('/communications', { params: { school_id: schoolID } })
                    : await api.get('/communications');
                this.communications = response.data || [];
            } catch (err) {
                this.error = err.response?.data?.error || 'Failed to fetch communications';
                console.error('Error fetching communications:', err);
            } finally {
                this.loading = false;
            }
        },

        async sendMessage(payload) {
            this.loading = true;
            this.error = null;
            try {
                const response = await api.post('/communications', payload);
                this.communications.unshift(response.data);
                return response.data;
            } catch (err) {
                this.error = err.response?.data?.error || 'Error sending message';
                console.error('Error sending message:', err);
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async deleteCommunication(id) {
            this.loading = true;
            this.error = null;
            try {
                await api.delete(`/communications/${id}`);
                this.communications = this.communications.filter(c => c.id !== id);
            } catch (err) {
                this.error = err.response?.data?.error || 'Error deleting communication';
                console.error('Error deleting communication:', err);
                throw err;
            } finally {
                this.loading = false;
            }
        }
    }
});
