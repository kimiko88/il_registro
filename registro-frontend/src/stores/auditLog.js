import { defineStore } from 'pinia'
import api from '@/services/api'

export const useAuditLogStore = defineStore('auditLog', {
    state: () => ({
        logs: [],
        total: 0,
        page: 1,
        limit: 20,
        totalPages: 1,
        loading: false,
        error: null
    }),

    actions: {
        async fetchLogs(filters = {}, page = 1, limit = 20) {
            this.loading = true
            this.error = null
            try {
                const params = {
                    page,
                    limit
                }
                if (filters.actor_id) params.actor_id = filters.actor_id
                if (filters.action && filters.action !== 'all') params.action = filters.action
                if (filters.entity_type) params.entity_type = filters.entity_type
                if (filters.from) params.from = filters.from
                if (filters.to) params.to = filters.to

                const response = await api.get('/audit-log', { params })
                const data = response.data || {}
                this.logs = data.data || []
                this.total = data.total || 0
                this.page = data.page || page
                this.limit = data.limit || limit
                this.totalPages = data.total_pages || 1
                return data
            } catch (err) {
                this.error = err.response?.data?.error || err.userMessage || 'Errore caricamento audit log'
                throw err
            } finally {
                this.loading = false
            }
        },

        async exportLogs(filters = {}) {
            let url = null
            let link = null
            try {
                const params = {}
                if (filters.actor_id) params.actor_id = filters.actor_id
                if (filters.action && filters.action !== 'all') params.action = filters.action
                if (filters.entity_type) params.entity_type = filters.entity_type
                if (filters.from) params.from = filters.from
                if (filters.to) params.to = filters.to

                const response = await api.get('/audit-log/export', {
                    params,
                    responseType: 'blob',
                    timeout: 60000
                })
                const blob = new Blob([response.data], { type: 'text/csv' })
                url = window.URL.createObjectURL(blob)
                link = document.createElement('a')
                link.href = url
                link.setAttribute('download', `audit_logs_${new Date().toISOString().slice(0, 10)}.csv`)
                document.body.appendChild(link)
                link.click()
            } catch (err) {
                console.error('Error exporting audit logs:', err)
                throw err
            } finally {
                if (link) {
                    try {
                        if (link.parentNode) {
                            link.parentNode.removeChild(link)
                        } else if (typeof link.remove === 'function') {
                            link.remove()
                        }
                    } catch (e) {
                        console.debug('Failed to remove anchor element:', e)
                    }
                }
                if (url) {
                    window.URL.revokeObjectURL(url)
                }
            }
        }
    }
})
