import api from './api'

export const notesService = {
    // Create a new note
    async createNote(data) {
        return api.post('/notes', data)
    },

    // List notes with filters
    async getNotes(params) {
        return api.get('/notes', { params })
    },

    // Get notes for a specific student
    async getStudentNotes(studentId, params = {}) {
        return api.get(`/notes/student/${studentId}`, { params })
    },

    // Get notes for a specific class
    async getClassNotes(classId, params = {}) {
        return api.get(`/notes/class/${classId}`, { params })
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

export default notesService

