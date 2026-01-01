import { api } from 'src/boot/axios'

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

    // Delete a note
    async deleteNote(id) {
        return api.delete(`/notes/${id}`)
    }
}
