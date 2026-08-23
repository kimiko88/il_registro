import { defineStore } from 'pinia'
import api from '@/services/api'

export const useCertificatesStore = defineStore('certificates', {
    state: () => ({
        certificates: [],
        loading: false,
        error: null
    }),

    actions: {
        async fetchCertificates(filters = {}) {
            this.loading = true
            this.error = null
            try {
                const params = {}
                if (filters.student_id) params.student_id = filters.student_id
                if (filters.type && filters.type !== 'all') params.type = filters.type
                if (filters.academic_year) params.academic_year = filters.academic_year

                const response = await api.get('/certificates', { params })
                this.certificates = response.data || []
                return this.certificates
            } catch (err) {
                this.error = err.response?.data?.error || err.userMessage || 'Errore durante il caricamento dei certificati'
                throw err
            } finally {
                this.loading = false
            }
        },

        async generateCertificate(payload) {
            this.loading = true
            try {
                const response = await api.post('/certificates/generate', payload)
                const newCert = response.data?.cert || response.data
                if (newCert) {
                    this.certificates.unshift(newCert)
                }
                return response.data
            } catch (err) {
                this.error = err.response?.data?.error || err.userMessage || 'Errore durante la generazione del certificato'
                throw err
            } finally {
                this.loading = false
            }
        },

        async downloadPDF(id) {
            let url = null
            let link = null
            try {
                const response = await api.get(`/certificates/${id}/pdf`, {
                    responseType: 'blob',
                    timeout: 60000
                })
                const blob = new Blob([response.data], { type: 'application/pdf' })
                url = window.URL.createObjectURL(blob)
                link = document.createElement('a')
                link.href = url
                link.setAttribute('download', `certificato_${id}.pdf`)
                document.body.appendChild(link)
                link.click()
            } catch (err) {
                console.error('Error downloading certificate PDF:', err)
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
        },

        async deleteCertificate(id) {
            await api.delete(`/certificates/${id}`)
            this.certificates = this.certificates.filter(c => c.id !== id)
        }
    }
})
