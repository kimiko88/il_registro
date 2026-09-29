import { api } from '@/boot/axios'

export const religionService = {
  getStudentChoice(studentId) {
    return api.get(`/students/${studentId}/religion-choice`)
  },

  setStudentChoice(studentId, choice) {
    return api.put(`/students/${studentId}/religion-choice`, { choice })
  },

  listChoices(params = {}) {
    return api.get('/secretary/religion-choices', { params })
  },

  batchSetChoices(studentIds, choice) {
    return api.post('/secretary/religion-choices/batch', {
      student_ids: studentIds,
      choice
    })
  },

  getClassChoices(classId) {
    return api.get(`/classes/${classId}/religion-choices`)
  }
}

export default religionService
