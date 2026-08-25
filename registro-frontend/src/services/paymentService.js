import api from './api'

export const paymentService = {
  getPayments(params) {
    return api.get('/payments', { params })
  },
  getPayment(id) {
    return api.get(`/payments/${id}`)
  },
  pay(id, data = {}) {
    return api.post(`/payments/${id}/pay`, data)
  },
  createPayment(data) {
    return api.post('/payments', data)
  }
}

export default paymentService
