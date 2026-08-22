import api from './api'

export const udaService = {
  async getByClass(classId) {
    const res = await api.get(`/uda/class/${classId}`)
    return res.data
  },

  async create(data) {
    const res = await api.post('/uda', data)
    return res.data
  },

  async update(id, data) {
    const res = await api.put(`/uda/${id}`, data)
    return res.data
  },

  async delete(id) {
    const res = await api.delete(`/uda/${id}`)
    return res.data
  }
}

export default udaService
