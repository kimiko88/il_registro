import api from './api'

export const tenantsService = {
  getTenants() {
    return api.get('/tenants')
  },
  getTenant(id) {
    return api.get(`/tenants/${id}`)
  },
  createTenant(data) {
    return api.post('/tenants', data)
  },
  updateTenant(id, data) {
    return api.put(`/tenants/${id}`, data)
  },
  deleteTenant(id) {
    return api.delete(`/tenants/${id}`)
  }
}

export default tenantsService

