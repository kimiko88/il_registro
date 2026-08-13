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
  }
}
