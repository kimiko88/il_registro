import api from './api'

export default {
  getTrips(params) {
    return api.get('/trips', { params })
  },
  createTrip(data) {
    return api.post('/trips', data)
  },
  signConsent(tripId) {
    return api.post(`/trips/${tripId}/consent`)
  }
}
