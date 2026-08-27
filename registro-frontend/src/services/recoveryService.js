import api from './api'

export const recoveryService = {
  createCourse(payload) {
    return api.post('/recovery/courses', payload)
  },
  listCourses(params = {}) {
    return api.get('/recovery/courses', { params })
  },
  getCourse(id) {
    return api.get(`/recovery/courses/${id}`)
  },
  updateStatus(id, status) {
    return api.patch(`/recovery/courses/${id}/status`, { status })
  },
  updateAttendance(courseId, studentId, payload) {
    return api.patch(`/recovery/courses/${courseId}/students/${studentId}/attendance`, payload)
  },
  recordTest(payload) {
    return api.post('/recovery/tests', payload)
  },
  listTests(params = {}) {
    return api.get('/recovery/tests', { params })
  }
}

export default recoveryService
