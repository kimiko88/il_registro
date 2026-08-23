import api from './api'

export const communicationService = {
    getMessages(params) {
        return params ? api.get('/communications', { params }) : api.get('/communications')
    },
    sendMessage(data) {
        return api.post('/communications', data)
    },
    signMessage(id) {
        return api.post(`/communications/${id}/sign`)
    },
    getSignatures(id) {
        return api.get(`/communications/${id}/signatures`)
    },
    markAsRead(id) {
        return api.post(`/communications/${id}/read`)
    },
    getUnreadUsers(id) {
        return api.get(`/communications/${id}/unread-users`)
    },
    getUnreadCount() {
        return api.get('/communications/unread-count')
    },
    uploadAttachment(file) {
        const formData = new FormData()
        formData.append('file', file)
        return api.post('/communications/upload', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
            timeout: 60000
        })
    }
}

export default communicationService

