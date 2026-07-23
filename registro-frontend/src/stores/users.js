import { defineStore } from 'pinia'
import api from '../services/api'

export const useUsersStore = defineStore('users', {
  state: () => ({
    users: [],
    loading: false,
    error: null
  }),

  actions: {
    async bulkImportUsers(formData) {
      this.loading = true
      this.error = null
      try {
        const response = await api.post('/users/bulk-import', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        return response.data
      } catch (err) {
        this.error = err.response?.data?.error || 'Errore durante l\'importazione massiva'
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    }
  }
})
