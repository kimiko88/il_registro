import api from './api'

export const securityService = {
  async downloadCadPackage(academicYear = '2025/2026') {
    const res = await api.get('/signatures/cad-preservation/download', {
      params: { academic_year: academicYear },
      responseType: 'blob',
      timeout: 60000
    })
    return res.data
  },

  async getImmutabilityChain() {
    const res = await api.get('/audit-log/immutability-chain')
    return res.data
  },

  async signSubstitutionRegister(subId, notes = '') {
    const res = await api.post(`/substitutions/${subId}/sign-register`, { notes })
    return res.data
  },

  async getRecommendedSubstitutes(classId, date, hour, subjectId = '') {
    const params = { class_id: classId, date, hour }
    if (subjectId) params.subject_id = subjectId
    const res = await api.get('/substitutions/recommend-substitutes', { params })
    return res.data
  }
}

export default securityService

