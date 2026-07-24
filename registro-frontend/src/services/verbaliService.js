import api from './api'

export default {
  listMeetings(params) {
    return api.get('/verbali/meetings', { params })
  },
  createMeeting(data) {
    return api.post('/verbali/meetings', data)
  },
  getVerbali(meetingId) {
    if (!meetingId || meetingId === 'undefined' || meetingId === 'null') {
      return Promise.resolve({ data: [] })
    }
    return api.get(`/verbali/meeting/${meetingId}`)
  },
  getVerbale(id) {
    return api.get(`/verbali/${id}`)
  },
  createVerbale(data) {
    return api.post('/verbali', data)
  },
  signVerbale(id) {
    return api.post(`/verbali/${id}/sign`)
  },
  getSignatures(id) {
    return api.get(`/verbali/${id}/signatures`)
  }
}
