import api from './api'

export const extracurricularService = {
  getCourses(params) {
    return api.get('/extracurricular/courses', { params })
  },
  getCourse(id) {
    return api.get(`/extracurricular/courses/${id}`)
  },
  createCourse(data) {
    return api.post('/extracurricular/courses', data)
  },
  updateCourse(id, data) {
    return api.put(`/extracurricular/courses/${id}`, data)
  },
  deleteCourse(id) {
    return api.delete(`/extracurricular/courses/${id}`)
  },
  enrollStudent(courseId, data) {
    return api.post(`/extracurricular/courses/${courseId}/enroll`, data)
  },
  unenrollStudent(courseId, studentId) {
    return api.delete(`/extracurricular/courses/${courseId}/students/${studentId}`)
  },
  recordAttendance(courseId, data) {
    return api.post(`/extracurricular/courses/${courseId}/attendance`, data)
  }
}

export default extracurricularService

