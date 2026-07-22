import { defineStore } from 'pinia';
import api from 'src/services/api';

export const useCertificatesStore = defineStore('certificates', {
    state: () => ({
        certificates: [],
        loading: false,
        error: null
    }),

    actions: {
        async fetchCertificates(filters = {}) {
            this.loading = true;
            this.error = null;
            try {
                const params = new URLSearchParams();
                if (filters.student_id) params.append('student_id', filters.student_id);
                if (filters.type && filters.type !== 'all') params.append('type', filters.type);
                if (filters.academic_year) params.append('academic_year', filters.academic_year);

                const response = await api.get(`/certificates?${params.toString()}`);
                this.certificates = response.data || [];
                return this.certificates;
            } catch (err) {
                this.error = err.response?.data?.error || 'Errore durante il caricamento dei certificati';
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async generateCertificate(payload) {
            this.loading = true;
            try {
                const response = await api.post('/certificates/generate', payload);
                const newCert = response.data?.cert || response.data;
                if (newCert) {
                    this.certificates.unshift(newCert);
                }
                return response.data;
            } catch (err) {
                this.error = err.response?.data?.error || 'Errore durante la generazione del certificato';
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async downloadPDF(id) {
            try {
                const response = await api.get(`/certificates/${id}/pdf`, {
                    responseType: 'blob'
                });
                const blob = new Blob([response.data], { type: 'application/pdf' });
                const url = window.URL.createObjectURL(blob);
                const link = document.createElement('a');
                link.href = url;
                link.setAttribute('download', `certificato_${id}.pdf`);
                document.body.appendChild(link);
                link.click();
                link.remove();
                window.URL.revokeObjectURL(url);
            } catch (err) {
                throw err;
            }
        },

        async deleteCertificate(id) {
            try {
                await api.delete(`/certificates/${id}`);
                this.certificates = this.certificates.filter(c => c.id !== id);
            } catch (err) {
                throw err;
            }
        }
    }
});
