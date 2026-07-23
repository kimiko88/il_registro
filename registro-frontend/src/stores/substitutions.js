import { defineStore } from 'pinia'
import api from '../services/api'

export const useSubstitutionsStore = defineStore('substitutions', {
  state: () => ({
    mySubstitutions: [],
    todaySubstitutions: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchMySubstitutions(weekDate) {
      this.loading = true
      this.error = null
      try {
        const params = weekDate ? { date: weekDate } : {}
        const response = await api.get('/substitutions/my', { params })
        this.mySubstitutions = response.data || []
        return this.mySubstitutions
      } catch (err) {
        this.error = err.response?.data?.error || 'Errore durante il recupero delle sostituzioni'
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async fetchTodaySubstitutions() {
      this.loading = true
      try {
        const response = await api.get('/substitutions/my-today')
        this.todaySubstitutions = response.data || []
        return this.todaySubstitutions
      } catch (err) {
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async confirmSubstitution(id) {
      this.loading = true
      try {
        const response = await api.patch(`/substitutions/${id}/confirm`)
        const index = this.mySubstitutions.findIndex(s => s.id === id)
        if (index !== -1) {
          this.mySubstitutions[index].status = 'confirmed'
        }
        const todayIdx = this.todaySubstitutions.findIndex(s => s.id === id)
        if (todayIdx !== -1) {
          this.todaySubstitutions[todayIdx].status = 'confirmed'
        }
        return response.data
      } catch (err) {
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async markAttendanceForSubstitution(payload) {
      this.loading = true
      try {
        const response = await api.post('/attendance/batch', payload)
        return response.data
      } catch (err) {
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    }
  }
})
