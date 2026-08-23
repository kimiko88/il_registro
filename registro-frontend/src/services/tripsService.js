import api from './api'

export const tripsService = {
  getTrips(params) {
    return api.get('/trips', { params })
  },
  getTrip(id) {
    return api.get(`/trips/${id}`)
  },
  createTrip(data) {
    return api.post('/trips', data)
  },
  updateTrip(id, data) {
    return api.put(`/trips/${id}`, data)
  },
  deleteTrip(id) {
    return api.delete(`/trips/${id}`)
  },
  signConsent(tripId, data = {}) {
    return api.post(`/trips/${tripId}/consent`, { trip_id: tripId, ...data })
  },
  getConsents(tripId) {
    return api.get(`/trips/${tripId}/consents`)
  }
}

export default tripsService

