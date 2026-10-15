import api from './api'

export const interpelliService = {
  // Public Interpelli (O.M. 88/2024)
  listPublic(params = {}) {
    return api.get('/public/interpelli', { params })
  },
  getPublic(id) {
    return api.get(`/public/interpelli/${id}`)
  },
  submitCandidatura(id, data) {
    return api.post(`/public/interpelli/${id}/candidatura`, data)
  },

  // Protected Administration
  createNotice(data) {
    return api.post('/interpelli', data)
  },
  getGraduatoria(noticeId) {
    return api.get(`/interpelli/${noticeId}/graduatoria`)
  },
  convoca(candidateId, hoursToRespond = 24) {
    return api.post(`/interpelli/candidature/${candidateId}/convoca`, {
      hours_to_respond: hoursToRespond
    })
  },
  rispondi(candidateId, risposta, notes = '') {
    return api.post(`/interpelli/candidature/${candidateId}/risposta`, {
      risposta,
      notes
    })
  }
}

export default interpelliService
