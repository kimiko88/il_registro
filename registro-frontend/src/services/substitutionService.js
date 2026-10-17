import api from './api'

export const substitutionService = {
  listBySchool(params) {
    return api.get('/substitutions', { params })
  },
  listMy() {
    return api.get('/substitutions/my')
  },
  create(data) {
    return api.post('/substitutions', data)
  },
  update(id, data) {
    return api.put(`/substitutions/${id}`, data)
  },
  delete(id) {
    return api.delete(`/substitutions/${id}`)
  },
  assign(id, data) {
    return api.put(`/substitutions/${id}/assign`, data)
  },
  recommendSubstitutes(params) {
    return api.get('/substitutions/recommend-substitutes', { params })
  },
  signRegister(subId, notes = '') {
    return api.post(`/substitutions/${subId}/sign-register`, { notes })
  },
  todaySummary(date) {
    return api.get('/substitutions/today-summary', { params: { date } })
  }
}

export default substitutionService


