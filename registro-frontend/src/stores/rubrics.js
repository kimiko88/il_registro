import { defineStore } from 'pinia'
import api from '../services/api'

export const useRubricsStore = defineStore('rubrics', {
  state: () => ({
    rubrics: [],
    assessments: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchRubrics() {
      this.loading = true
      this.error = null
      try {
        const response = await api.get('/rubrics')
        this.rubrics = response.data || []
        return this.rubrics
      } catch (err) {
        this.error = err.response?.data?.error || 'Errore durante il recupero delle rubriche'
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async getRubric(id) {
      this.loading = true
      try {
        const response = await api.get(`/rubrics/${id}`)
        return response.data
      } catch (err) {
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async createRubric(payload) {
      this.loading = true
      try {
        const response = await api.post('/rubrics', payload)
        if (response.data) {
          this.rubrics.unshift(response.data)
        }
        return response.data
      } catch (err) {
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async updateRubric(id, payload) {
      this.loading = true
      try {
        const response = await api.put(`/rubrics/${id}`, payload)
        const updated = response.data
        const index = this.rubrics.findIndex(r => r.id === id)
        if (index !== -1 && updated) {
          this.rubrics[index] = { ...this.rubrics[index], ...updated }
        }
        return updated
      } catch (err) {
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async deleteRubric(id) {
      this.loading = true
      try {
        await api.delete(`/rubrics/${id}`)
        this.rubrics = this.rubrics.filter(r => r.id !== id)
      } catch (err) {
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async assessStudent(rubricId, payload) {
      this.loading = true
      try {
        const response = await api.post(`/rubrics/${rubricId}/assessments`, payload)
        return response.data
      } catch (err) {
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async fetchStudentAssessments(studentId) {
      this.loading = true
      try {
        const response = await api.get(`/rubrics/assessments/student/${studentId}`)
        this.assessments = response.data || []
        return this.assessments
      } catch (err) {
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    },

    async fetchClassAssessments(classId) {
      this.loading = true
      try {
        const response = await api.get(`/rubrics/assessments/class/${classId}`)
        this.assessments = response.data || []
        return this.assessments
      } catch (err) {
        console.error(err)
        throw err
      } finally {
        this.loading = false
      }
    }
  }
})
