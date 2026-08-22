import api from './api'

export default {
  getTrips(params) {
    return api.get('/trips', { params })
  },
  createTrip(data) {
    return api.post('/trips', data)
  },
  signConsent(tripId, data = {}) {
    return api.post(`/trips/${tripId}/consent`, { trip_id: tripId, ...data })
  }
}
