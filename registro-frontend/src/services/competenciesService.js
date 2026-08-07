import api from './api'

export const competenciesService = {
  async getStudentEvaluations(studentId, semester = 1) {
    const res = await api.get(`/api/v1/competencies/student/${studentId}?semester=${semester}`)
    return res.data
  },

  async saveEvaluation(data) {
    const res = await api.post('/api/v1/competencies/evaluations', data)
    return res.data
  }
}

export default competenciesService
