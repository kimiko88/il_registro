import api from '@/services/api'

export const protocolService = {
    async protocolDocument(payload) {
        const response = await api.post('/protocol', payload)
        return response.data
    },

    async getEntries(params = {}) {
        const response = await api.get('/protocol', { params })
        return response.data
    },

    async getEntryById(id) {
        const response = await api.get(`/protocol/${id}`)
        return response.data
    },

    async getProtocolByEntity(entityType, entityId) {
        const response = await api.get(`/protocol/by-entity/${entityType}/${entityId}`)
        return response.data
    }
}
