import api from './api'

export default {
  getTenants() {
    return api.get('/tenants')
  },
  getTenant(id) {
    return api.get(`/tenants/${id}`)
  },
  createTenant(data) {
    return api.post('/tenants', data)
  }
}
