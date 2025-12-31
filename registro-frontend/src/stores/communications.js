import { defineStore } from 'pinia';
import { api } from 'boot/axios';

export const useCommunicationsStore = defineStore('communications', {
    state: () => ({
        messages: [],
        loading: false,
        error: null,
    }),

    actions: {
        async fetchMessages() {
            this.loading = true;
            this.error = null;
            try {
                const response = await api.get('/communications');
                this.messages = response.data;
            } catch (err) {
                this.error = err.response?.data?.error || 'Failed to fetch messages';
                console.error('Error fetching messages:', err);
            } finally {
                this.loading = false;
            }
        },

        async sendMessage(payload) {
            this.loading = true;
            try {
                const response = await api.post('/communications', payload);
                this.messages.unshift(response.data);
                return response.data;
            } catch (err) {
                console.error('Error sending message:', err);
                throw err;
            } finally {
                this.loading = false;
            }
        }
    }
});
