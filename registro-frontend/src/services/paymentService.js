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
  },
  createBulkPayments(studentIds, paymentData) {
    return Promise.all(studentIds.map(id => this.createPayment({ ...paymentData, student_id: id })))
  },
  // PagoPA & OPI Enterprise Dealbreakers
  getBollettino(id) {
    return api.get(`/payments/${id}/bollettino`)
  },
  createCheckoutSession(id, returnUrl) {
    return api.post(`/payments/${id}/checkout`, { return_url: returnUrl })
  },
  reconcileOPI(streamData, format = 'OPI_XML') {
    return api.post(`/payments/reconcile-opi?format=${format}`, streamData, {
      headers: {
        'Content-Type': format === 'OPI_XML' ? 'application/xml' : 'text/csv'
      }
    })
  },
  getSolleciti() {
    return api.get('/payments/solleciti')
  }
}

export default paymentService
