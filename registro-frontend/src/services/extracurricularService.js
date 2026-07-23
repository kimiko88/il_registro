import api from './api'

export default {
  getCourses(params) {
    return api.get('/extracurricular/courses', { params })
  },
  createCourse(data) {
    return api.post('/extracurricular/courses', data)
  },
  enrollStudent(courseId, data) {
    return api.post(`/extracurricular/courses/${courseId}/enroll`, data)
  },
  recordAttendance(courseId, data) {
    return api.post(`/extracurricular/courses/${courseId}/attendance`, data)
  }
}
