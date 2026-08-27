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
  deleteSlot(id) {
    return api.delete(`/colloqui/slots/${id}`)
  },
  bookSlot(data) {
    return api.post('/colloqui/bookings', data)
  },
  cancelBooking(id) {
    return api.delete(`/colloqui/bookings/${id}`)
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
  },
  createGeneralMeeting(data) {
    return api.post('/colloqui/general-meetings', data)
  },
  listGeneralMeetings() {
    return api.get('/colloqui/general-meetings')
  },
  getGeneralMeeting(id) {
    return api.get(`/colloqui/general-meetings/${id}`)
  },
  bookQueueTicket(data) {
    return api.post('/colloqui/general-meetings/book-ticket', data)
  },
  listQueueTickets(meetingId, params = {}) {
    return api.get(`/colloqui/general-meetings/${meetingId}/tickets`, { params })
  },
  updateTicketStatus(ticketId, data) {
    return api.patch(`/colloqui/general-meetings/tickets/${ticketId}/status`, data)
  }
}

export default colloquiService
