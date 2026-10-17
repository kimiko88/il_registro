import api from './api'

export const textbookService = {
    getAll() {
        return api.get('/textbooks')
    },
    create(data) {
        return api.post('/textbooks', data)
    },
    update(id, data) {
        return api.put(`/textbooks/${id}`, data)
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
    },

    // AIE & Spending Limits endpoints
    importAIE(formData) {
        return api.post('/textbooks/aie/import', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        })
    },
    searchAIECatalog(params) {
        return api.get('/textbooks/aie/catalog', { params })
    },
    getSpendingReport(classId) {
        return api.get(`/textbooks/classes/${classId}/spending-report`)
    },
    listClassAdoptions(classId) {
        return api.get(`/textbooks/classes/${classId}/adoptions`)
    },
    adoptBook(classId, data) {
        return api.post(`/textbooks/classes/${classId}/adoptions`, data)
    },
    deleteAdoption(id) {
        return api.delete(`/textbooks/adoptions/${id}`)
    },
    upsertSpendingLimit(data) {
        return api.post('/textbooks/spending-limits', data)
    },
    exportClassAIE(classId) {
        return api.get(`/textbooks/classes/${classId}/aie-export`, {
            responseType: 'blob'
        })
    }
}

export default textbookService
