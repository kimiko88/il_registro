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
  }
}

export default supportService
