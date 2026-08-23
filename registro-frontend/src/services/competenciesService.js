import api from './api'

export const competenciesService = {
  async getStudentEvaluations(studentId, semester = 1) {
    const res = await api.get(`/competencies/student/${studentId}`, { params: { semester } })
    return res.data
  },

  async saveEvaluation(data) {
    const res = await api.post('/competencies/evaluations', data)
    return res.data
  },

  async getCompetencies(subjectId, gradeLevel) {
    const params = {}
    if (subjectId) params.subject_id = subjectId
    if (gradeLevel) params.grade_level = gradeLevel
    const res = await api.get('/competencies', { params })
    return res.data
  }
}

export default competenciesService

