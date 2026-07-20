import api from './api'

export default {
  getVerbali(params) {
    return api.get('/verbali', { params })
  },
  getVerbale(id) {
    return api.get(`/verbali/${id}`)
  },
  createVerbale(data) {
    return api.post('/verbali', data)
  },
  signVerbale(id) {
    return api.post(`/verbali/${id}/sign`)
  }
}
