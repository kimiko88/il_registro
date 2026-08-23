import api from './api'

export const studentGoalService = {
  listByStudent(studentId) {
    return api.get(`/student-goals/student/${studentId}`)
  },
  create(data) {
    return api.post('/student-goals', data)
  },
  update(id, data) {
    return api.put(`/student-goals/${id}`, data)
  },
  updateStatus(id, status) {
    return api.patch(`/student-goals/${id}/status`, { status })
  },
  delete(id) {
    return api.delete(`/student-goals/${id}`)
  }
}

export default studentGoalService

