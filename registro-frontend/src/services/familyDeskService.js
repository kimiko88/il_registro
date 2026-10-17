import api from '@/services/api'

export const familyDeskService = {
    async submitRequest(payload) {
        const response = await api.post('/family-desk/requests', payload)
        return response.data
    },

    async getRequests(params = {}) {
        const response = await api.get('/family-desk/requests', { params })
        return response.data
    },

    async getRequestById(id) {
        const response = await api.get(`/family-desk/requests/${id}`)
        return response.data
    },

    async reviewRequest(id, payload) {
        const response = await api.put(`/family-desk/requests/${id}/review`, payload)
        return response.data
    },

    async getDelegates(studentId) {
        const response = await api.get('/family-desk/delegates', {
            params: { student_id: studentId }
        })
        return response.data
    }
}
