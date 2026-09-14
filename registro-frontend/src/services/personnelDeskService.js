import api from './api'

export const personnelDeskService = {
  listRequests(status = '') {
    return api.get('/personnel-desk/requests', { params: { status } })
  },
  getRequest(id) {
    return api.get(`/personnel-desk/requests/${id}`)
  },
  createRequest(data) {
    return api.post('/personnel-desk/requests', data)
  },
  submitRequest(id) {
    return api.patch(`/personnel-desk/requests/${id}/submit`)
  },
  aaReview(id, data) {
    return api.patch(`/personnel-desk/requests/${id}/aa-review`, data)
  },
  dsgaSign(id, data) {
    return api.patch(`/personnel-desk/requests/${id}/dsga-sign`, data)
  },
  dsApprove(id, data) {
    return api.patch(`/personnel-desk/requests/${id}/ds-approve`, data)
  },
  deleteRequest(id) {
    return api.delete(`/personnel-desk/requests/${id}`)
  }
}

export default personnelDeskService
