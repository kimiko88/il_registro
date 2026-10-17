import api from '@/services/api'

export const electionsService = {
    async getElections() {
        const response = await api.get('/elections')
        return response.data
    },

    async getElectionDetails(id) {
        const response = await api.get(`/elections/${id}`)
        return response.data
    },

    async castVote(id, payload) {
        const response = await api.post(`/elections/${id}/vote`, payload)
        return response.data
    },

    async getScrutiny(id, seats = 4) {
        const response = await api.get(`/elections/${id}/scrutiny`, {
            params: { seats }
        })
        return response.data
    },

    async closeElection(id) {
        const response = await api.put(`/elections/${id}/close`)
        return response.data
    },

    async createElection(payload) {
        const response = await api.post('/elections', payload)
        return response.data
    }
}
