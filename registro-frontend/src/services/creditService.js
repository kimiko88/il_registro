import api from './api'

export const creditService = {
  calculateSuggested(params) {
    return api.get('/credits/calculate', { params })
  },
  assignCredit(payload) {
    return api.post('/credits/assign', payload)
  },
  listClassCredits(classId, academicYear = '') {
    return api.get(`/credits/class/${classId}`, { params: { academic_year: academicYear } })
  },
  getStudentSummary(studentId) {
    return api.get(`/credits/student/${studentId}/summary`)
  }
}

export default creditService
