import api from './api'

export default {
    // Create a new note
    async createNote(data) {
        return api.post('/notes', data)
    },

    // List notes with filters
    async getNotes(params) {
        return api.get('/notes', { params })
    },

    // Update a note
    async updateNote(id, data) {
        return api.patch(`/notes/${id}`, data)
    },

    // Approve a note (Dirigenza)
    async approveNote(id) {
        return api.post(`/notes/${id}/approve`)
    },

    // Delete a note
    async deleteNote(id) {
        return api.delete(`/notes/${id}`)
    }
}
