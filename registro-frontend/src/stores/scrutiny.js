import { defineStore } from 'pinia'
import api from '../services/api'

export const useScrutinyStore = defineStore('scrutiny', {
  state: () => ({
    overview: [],
    loading: false,
    error: null,
    currentReport: null
  }),

  getters: {
    classesWithRisk: (state) => state.overview.filter(c => c.risk_count > 0 || c.has_risk),
    promotedCount: (state) => state.currentReport?.students?.filter(s => s.promoted === 'SÌ' || s.promoted === true).length || 0,
    pendingCount: (state) => state.overview.filter(c => c.status === 'pending').length || 0
  },

  actions: {
    async fetchOverview(schoolID = null) {
      this.loading = true
      this.error = null
      try {
        const params = schoolID ? { school_id: schoolID } : {}
        const response = await api.get('/scrutiny/overview', { params })
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
      this.error = null
      try {
        const response = await api.get(`/scrutiny/class/${classID}/report`)
        this.currentReport = response.data
        return response.data
      } catch (err) {
        this.error = err.response?.data?.error || 'Errore caricamento report classe'
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async finalizeScrutiny(classID, schoolID = null) {
      this.loading = true
      this.error = null
      try {
        const response = await api.post(`/scrutiny/class/${classID}/finalize`)
        await this.fetchOverview(schoolID)
        return response.data
      } catch (err) {
        this.error = err.response?.data?.error || 'Errore finalizzazione dello scrutinio'
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async exportAll() {
      let url = null
      let link = null
      try {
        const response = await api.get('/scrutiny/export', { responseType: 'blob' })
        const blob = new Blob([response.data], { type: 'text/csv' })
        url = window.URL.createObjectURL(blob)
        link = document.createElement('a')
        link.href = url
        link.setAttribute('download', 'scrutini_supervisione.csv')
        document.body.appendChild(link)
        link.click()
      } catch (err) {
        console.error(err)
        throw err
      } finally {
        if (link && link.parentNode) {
          link.remove()
        }
        if (url) {
          setTimeout(() => {
            window.URL.revokeObjectURL(url)
          }, 100)
        }
      }
    }
  }
})
