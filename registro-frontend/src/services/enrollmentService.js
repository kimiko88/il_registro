import api from './api'

export const enrollmentService = {
    importSIDI(formData, academicYear = '2026/2027') {
        return api.post('/enrollment/import-sidi', formData, {
            params: { academic_year: academicYear },
            headers: { 'Content-Type': 'multipart/form-data' }
        })
    },
    listApplications(params) {
        return api.get('/enrollment/applications', { params })
    },
    generateFormationDraft(data) {
        return api.post('/enrollment/formation-drafts/generate', data)
    },
    listDrafts() {
        return api.get('/enrollment/formation-drafts')
    },
    getDraft(id) {
        return api.get(`/enrollment/formation-drafts/${id}`)
    },
    updateDraftAssignments(id, assignments) {
        return api.put(`/enrollment/formation-drafts/${id}`, assignments)
    },
    finalizeDraft(id) {
        return api.post(`/enrollment/formation-drafts/${id}/finalize`)
    }
}

export default enrollmentService
