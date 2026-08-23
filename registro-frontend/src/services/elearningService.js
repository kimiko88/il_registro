import api from './api'

export const elearningService = {
  getProviders() {
    return api.get('/elearning/providers')
  },
  connectGoogleClassroom(authCode) {
    return api.post('/elearning/google/connect', { code: authCode })
  },
  connectMicrosoftTeams(authCode) {
    return api.post('/elearning/microsoft/connect', { code: authCode })
  },
  syncCourses(provider) {
    return api.post(`/elearning/${provider}/sync-courses`)
  },
  syncAssignments(provider, classId) {
    return api.post(`/elearning/${provider}/sync-assignments`, { class_id: classId })
  },
  syncGrades(provider, classId) {
    return api.post(`/elearning/${provider}/sync-grades`, { class_id: classId })
  },
  disconnect(provider) {
    return api.post(`/elearning/${provider}/disconnect`)
  },
  getStatus(provider) {
    return api.get(`/elearning/${provider}/status`)
  }
}

export default elearningService

