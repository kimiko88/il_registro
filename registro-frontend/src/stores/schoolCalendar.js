import { defineStore } from 'pinia'
import api from '../services/api'

export const useSchoolCalendarStore = defineStore('schoolCalendar', {
  state: () => ({
    events: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchCalendar(year, month) {
      this.loading = true
      this.error = null
      try {
        const params = {}
        if (year) params.year = year
        if (month) params.month = month

        const response = await api.get('/school-calendar/student', { params })
        this.events = response.data || []
        return this.events
      } catch (err) {
        this.error = err.response?.data?.error || 'Errore caricamento calendario'
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    }
  }
})
