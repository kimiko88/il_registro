import api from './api'

export const udaService = {
  async getByClass(classId) {
    const res = await api.get(`/api/v1/uda/class/${classId}`)
    return res.data
  },

  async create(data) {
    const res = await api.post('/api/v1/uda', data)
    return res.data
  },

  async update(id, data) {
    const res = await api.put(`/api/v1/uda/${id}`, data)
    return res.data
  },

  async delete(id) {
    const res = await api.delete(`/api/v1/uda/${id}`)
    return res.data
  }
}

export default udaService
