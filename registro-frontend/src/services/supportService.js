import api from './api'

export const supportService = {
  createDiaryEntry(payload) {
    return api.post('/support/diaries', payload)
  },
  listDiaryEntries(params = {}) {
    return api.get('/support/diaries', { params })
  },
  deleteDiaryEntry(id) {
    return api.delete(`/support/diaries/${id}`)
  },
  createPeiGoal(payload) {
    return api.post('/support/pei-goals', payload)
  },
  updateGoalProgress(id, progressStatus) {
    return api.patch(`/support/pei-goals/${id}/progress`, { progress_status: progressStatus })
  },
  listPeiGoals(params = {}) {
    return api.get('/support/pei-goals', { params })
  },
  deletePeiGoal(id) {
    return api.delete(`/support/pei-goals/${id}`)
  },

  // GLO Meetings and Joint Digital Signatures
  createGloMeeting(payload) {
    return api.post('/support/glo-meetings', payload)
  },
  listGloMeetings(params = {}) {
    return api.get('/support/glo-meetings', { params })
  },
  getGloMeeting(id) {
    return api.get(`/support/glo-meetings/${id}`)
  },
  signGloMeeting(payload) {
    return api.post('/support/glo-meetings/sign', payload)
  },

  // Support Hours Quantification (D.I. 182/2020 e D.M. 153/2023)
  saveSupportHours(payload) {
    return api.post('/support/hours', payload)
  },
  getSupportHours(params = {}) {
    return api.get('/support/hours', { params })
  }
}

export default supportService
