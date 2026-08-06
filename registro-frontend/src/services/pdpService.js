import api from './api'

export const pdpService = {
  getStudentPlans(studentId, year = '') {
    return api.get(`/pdp/student/${studentId}`, { params: { year } })
  },
  getClassPlans(classId, year = '') {
    return api.get(`/pdp/class/${classId}`, { params: { year } })
  },
  getPlan(id) {
    return api.get(`/pdp/${id}`)
  },
  createPlan(data) {
    return api.post('/pdp', data)
  },
  updatePlan(id, data) {
    return api.put(`/pdp/${id}`, data)
  },
  deletePlan(id) {
    return api.delete(`/pdp/${id}`)
  },
  shareWithFamily(id, share = true) {
    return api.post(`/pdp/${id}/share`, { share })
  },
  approveByFamily(id, comment = '') {
    return api.post(`/pdp/${id}/approve`, { comment })
  },
  getCompensativeMeasures() {
    return api.get('/pdp/measures/compensative')
  },
  getDispensativeMeasures() {
    return api.get('/pdp/measures/dispensative')
  }
}
