import { defineStore } from 'pinia'
import pctoService from '@/services/pctoService'
import api from '@/services/api'

export const usePCTOStore = defineStore('pcto', {
    state: () => ({
        projects: [],
        loading: false,
        error: null,
        totalHours: 0
    }),

    getters: {
        completedProjects: (state) => state.projects.filter(p => (p.status || '').toLowerCase() === 'completed'),
        inProgressProjects: (state) => state.projects.filter(p => (p.status || '').toLowerCase() !== 'completed'),
    },

    actions: {
        async fetchProjects(studentId = null) {
            this.loading = true
            this.error = null
            try {
                let items = []
                if (studentId) {
                    const res = await pctoService.getStudentProjects(studentId).catch(() => null)
                    items = res?.data || []
                } else {
                    const res = await pctoService.getProjects().catch(() => null)
                    items = res?.data || []
                }

                if (!items || items.length === 0) {
                    // Fallback to direct api if needed
                    const apiRes = await api.get('/pcto/projects').catch(() => null)
                    items = apiRes?.data || []
                }

                this.projects = items
                this.totalHours = items.reduce((acc, p) => acc + (Number(p.hours) || 0), 0)
                return this.projects
            } catch (err) {
                this.error = err.response?.data?.error || err.userMessage || 'Failed to fetch PCTO projects'
                console.error('Error fetching PCTO projects:', err)
                return []
            } finally {
                this.loading = false
            }
        },

        async createProject(payload) {
            this.loading = true
            this.error = null
            try {
                const res = await pctoService.createProject(payload)
                const created = res.data || res
                if (created) {
                    this.projects.push(created)
                    this.totalHours += (Number(created.hours) || 0)
                }
                return created
            } catch (err) {
                this.error = err.response?.data?.error || err.userMessage || 'Failed to create PCTO project'
                console.error('Error creating PCTO project:', err)
                throw err
            } finally {
                this.loading = false
            }
        }
    }
})
