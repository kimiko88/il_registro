import api from '@/services/api'

export const pctoCompanyTutorService = {
    async getSession(token) {
        const response = await api.get('/pcto-tutor/session', {
            params: { token },
            headers: { 'X-Tutor-Token': token }
        })
        return response.data
    },

    async getAssignedStudents(token) {
        const response = await api.get('/pcto-tutor/students', {
            headers: { 'X-Tutor-Token': token }
        })
        return response.data
    },

    async verifyTimesheet(token, payload) {
        const response = await api.post('/pcto-tutor/timesheets/verify', payload, {
            headers: { 'X-Tutor-Token': token }
        })
        return response.data
    },

    async submitEvaluation(token, payload) {
        const response = await api.post('/pcto-tutor/evaluations', payload, {
            headers: { 'X-Tutor-Token': token }
        })
        return response.data
    },

    async registerTutor(payload) {
        const response = await api.post('/pcto-tutor/tutors', payload)
        return response.data
    }
}
