import { defineStore } from 'pinia'
import api from '../services/api'

export const useScrutinyStore = defineStore('scrutiny', {
  state: () => ({
    overview: [],
    loading: false,
    error: null,
    currentReport: null
  }),

  actions: {
    async fetchOverview() {
      this.loading = true
      this.error = null
      try {
        const response = await api.get('/scrutiny/overview')
        this.overview = response.data || []
        return this.overview
      } catch (err) {
        this.error = err.response?.data?.error || 'Errore caricamento supervisione scrutini'
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async fetchClassReport(classID) {
      this.loading = true
      try {
        const response = await api.get(`/scrutiny/class/${classID}/report`)
        this.currentReport = response.data
        return response.data
      } catch (err) {
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async finalizeScrutiny(classID) {
      try {
        const response = await api.post(`/scrutiny/class/${classID}/finalize`)
        await this.fetchOverview()
        return response.data
      } catch (err) {
        console.error(err)
        throw err
      }
    },

    async exportAll() {
      try {
        const response = await api.get('/scrutiny/export', { responseType: 'blob' })
        const blob = new Blob([response.data], { type: 'text/csv' })
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', 'scrutini_supervisione.csv')
        document.body.appendChild(link)
        link.click()
        link.remove()
      } catch (err) {
        console.error(err)
        throw err
      }
    }
  }
})
