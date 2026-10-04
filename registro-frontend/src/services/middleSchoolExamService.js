import api from './api'

export const middleSchoolExamService = {
    getOrCreateExam(classId, presidentName = 'Presidente Commissione') {
        return api.post(`/middle-school-exam/classes/${classId}`, { president_name: presidentName })
    },
    listCandidates(examId) {
        return api.get(`/middle-school-exam/${examId}/candidates`)
    },
    saveAdmission(examId, studentId, data) {
        return api.post(`/middle-school-exam/${examId}/candidates/${studentId}/admission`, data)
    },
    evaluateCandidate(examId, studentId, data) {
        return api.post(`/middle-school-exam/${examId}/candidates/${studentId}/evaluate`, data)
    },
    generateDiploma(candidateId) {
        return api.get(`/middle-school-exam/candidates/${candidateId}/diploma`, {
            responseType: 'blob'
        })
    },
    updateStatus(examId, status) {
        return api.put(`/middle-school-exam/${examId}/status`, { status })
    }
}

export default middleSchoolExamService
