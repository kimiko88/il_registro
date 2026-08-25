import { defineStore } from 'pinia'
import api from '@/services/api'

export const useOrientamentoStore = defineStore('orientamento', {
    state: () => ({
        events: [],
        loading: false,
        error: null
    }),

    actions: {
        async fetchEvents() {
            this.loading = true
            this.error = null
            try {
                const response = await api.get('/orientamento/events').catch(() => null)
                this.events = response?.data || []
                return this.events
            } catch (err) {
                this.error = err.response?.data?.error || err.userMessage || 'Failed to fetch orientamento events'
                console.error('Error fetching orientamento events:', err)
                this.events = []
                return []
            } finally {
                this.loading = false
            }
        }
    }
})
