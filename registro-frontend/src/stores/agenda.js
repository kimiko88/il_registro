import { defineStore } from 'pinia'
import api from '../services/api'

export const useAgendaStore = defineStore('agenda', {
  state: () => ({
    events: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchAgenda(params = {}) {
      this.loading = true
      this.error = null
      try {
        const response = await api.get('/agenda', { params })
        this.events = response.data || []
        return this.events
      } catch (err) {
        this.error = err.response?.data?.error || 'Errore durante il recupero dell\'agenda'
        console.error('Errore fetchAgenda:', err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async createEvent(payload) {
      this.loading = true
      this.error = null
      try {
        const response = await api.post('/agenda', payload)
        if (response.data) {
          this.events.push(response.data)
        }
        return response.data
      } catch (err) {
        this.error = err.response?.data?.error || 'Errore durante la creazione dell\'evento'
        console.error('Errore createEvent:', err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async updateEvent(id, payload) {
      this.loading = true
      this.error = null
      try {
        const response = await api.patch(`/agenda/${id}`, payload)
        const updated = response.data
        const index = this.events.findIndex(e => e.id === id)
        if (index !== -1 && updated) {
          this.events[index] = { ...this.events[index], ...updated }
        }
        return updated
      } catch (err) {
        this.error = err.response?.data?.error || 'Errore durante la modifica dell\'evento'
        console.error('Errore updateEvent:', err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async deleteEvent(id) {
      this.loading = true
      this.error = null
      try {
        await api.delete(`/agenda/${id}`)
        this.events = this.events.filter(e => e.id !== id)
      } catch (err) {
        this.error = err.response?.data?.error || 'Errore durante l\'eliminazione dell\'evento'
        console.error('Errore deleteEvent:', err)
        throw err
      } finally {
        this.loading = false
      }
    }
  }
})
