import api from './api'

export const securityService = {
  async downloadCadPackage(academicYear = '2025/2026') {
    const res = await api.get(`/api/v1/signatures/cad-preservation/download?academic_year=${academicYear}`, { responseType: 'blob' })
    return res.data
  },

  async getImmutabilityChain() {
    const res = await api.get('/api/v1/audit-log/immutability-chain')
    return res.data
  },

  async signSubstitutionRegister(subId, notes = '') {
    const res = await api.post(`/api/v1/substitutions/${subId}/sign-register`, { notes })
    return res.data
  },

  async getRecommendedSubstitutes(classId, date, hour, subjectId = '') {
    const res = await api.get(`/api/v1/substitutions/recommend-substitutes?class_id=${classId}&date=${date}&hour=${hour}&subject_id=${subjectId}`)
    return res.data
  }
}

export default securityService
