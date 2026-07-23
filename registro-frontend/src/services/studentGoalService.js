import api from './api'

export const studentGoalService = {
  listByStudent(studentId) {
    return api.get(`/student-goals/student/${studentId}`)
  },
  create(data) {
    return api.post('/student-goals', data)
  },
  updateStatus(id, status) {
    return api.patch(`/student-goals/${id}/status`, { status })
  }
}
