import { defineStore } from 'pinia'
import api from '@/services/api'

export const useOrientamentoStore = defineStore('orientamento', {
    state: () => ({
        events: [],
        capolavori: [],
        curriculum: null,
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
        },

        async fetchCapolavori() {
            try {
                const response = await api.get('/orientamento/capolavori').catch(() => null)
                this.capolavori = response?.data || []
                return this.capolavori
            } catch (err) {
                console.error('Error fetching capolavori:', err)
                return []
            }
        },

        async saveCapolavoro(capolavoro) {
            const response = await api.post('/orientamento/capolavoro', capolavoro)
            await this.fetchCapolavori()
            return response.data
        },

        async fetchCurriculumStudente() {
            try {
                const response = await api.get('/orientamento/curriculum-studente').catch(() => null)
                this.curriculum = response?.data || null
                return this.curriculum
            } catch (err) {
                console.error('Error fetching curriculum studente:', err)
                return null
            }
        }
    }
})
