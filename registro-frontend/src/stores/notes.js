import { defineStore } from 'pinia'
import api from '../services/api'

export const useNotesStore = defineStore('notes', {
  state: () => ({
    notes: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchNotes(params = {}) {
      this.loading = true
      this.error = null
      try {
        let url = '/notes'
        if (typeof params === 'string') {
          url = `/notes/student/${params}`
          params = {}
        }
        const response = await api.get(url, { params })
        this.notes = response.data || []
        return this.notes
      } catch (err) {
        this.error = err.response?.data?.error || 'Errore durante il caricamento delle note'
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async addNote(payload) {
      this.loading = true
      this.error = null
      try {
        const response = await api.post('/notes', payload)
        if (response.data) {
          this.notes.unshift(response.data)
        }
        return response.data
      } catch (err) {
        this.error = err.response?.data?.error || 'Errore durante la creazione della nota'
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async deleteNote(id) {
      this.loading = true
      this.error = null
      try {
        await api.delete(`/notes/${id}`)
        this.notes = this.notes.filter(n => n.id !== id)
      } catch (err) {
        this.error = err.response?.data?.error || 'Errore durante l\'eliminazione della nota'
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    }
  }
})
