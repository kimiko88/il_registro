import api from '@/services/api'

export const helpDeskService = {
    async getSlots(params = {}) {
        const response = await api.get('/help-desk/slots', { params })
        return response.data
    },

    async createSlot(payload) {
        const response = await api.post('/help-desk/slots', payload)
        return response.data
    },

    async bookSlot(slotId, payload) {
        const response = await api.post(`/help-desk/slots/${slotId}/book`, payload)
        return response.data
    },

    async getBookings(slotId) {
        const response = await api.get(`/help-desk/slots/${slotId}/bookings`)
        return response.data
    },

    async markAttendance(bookingId, status) {
        const response = await api.put(`/help-desk/bookings/${bookingId}/attendance`, { status })
        return response.data
    },

    async completeSlot(slotId) {
        const response = await api.put(`/help-desk/slots/${slotId}/complete`)
        return response.data
    },

    async getFISReport() {
        const response = await api.get('/help-desk/fis-report')
        return response.data
    }
}
