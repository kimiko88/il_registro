import api from './api'

export const colloquiService = {
  getSlots(params) {
    return api.get('/colloqui/slots', { params })
  },
  createSlot(data) {
    return api.post('/colloqui/slots', data)
  },
  patchSlot(id, data) {
    return api.patch(`/colloqui/slots/${id}`, data)
  },
  cancelSlot(id) {
    return api.delete(`/colloqui/slots/${id}`)
  },
  bookSlot(data) {
    return api.post('/colloqui/bookings', data)
  },
  getMyBookings() {
    return api.get('/colloqui/my-bookings')
  },
  getBooking(id) {
    return api.get(`/colloqui/bookings/${id}`)
  },
  getSlotBookings(slotId) {
    return api.get(`/colloqui/slots/${slotId}/bookings`)
  },
  updateBookingStatus(id, data) {
    return api.put(`/colloqui/bookings/${id}/status`, data)
  },
  createAssembly(data) {
    return api.post('/colloqui/assemblies', data)
  }
}

export default colloquiService

