import api from './api'

export const textbookService = {
    getAll() {
        return api.get('/textbooks')
    },
    create(data) {
        return api.post('/textbooks', data)
    },
    delete(id) {
        return api.delete(`/textbooks/${id}`)
    },
    listByClass(classId) {
        return api.get(`/textbooks/class/${classId}`)
    },
    assignToClass(classId, data) {
        return api.post(`/textbooks/class/${classId}`, data)
    },
    removeFromClass(assignmentId) {
        return api.delete(`/textbooks/class/assignment/${assignmentId}`)
    }
}
