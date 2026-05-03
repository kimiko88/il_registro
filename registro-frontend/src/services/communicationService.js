import api from './api'

export const communicationService = {
    async getMessages() {
        const response = await api.get('/communications')
        return response
    },
    async sendMessage(data) {
        const response = await api.post('/communications', data)
        return response
    }
}
