import api from './api'

export const verbaliService = {
  listMeetings(params) {
    return api.get('/verbali/meetings', { params })
  },
  createMeeting(data) {
    return api.post('/verbali/meetings', data)
  },
  updateMeeting(id, data) {
    return api.put(`/verbali/meetings/${id}`, data)
  },
  deleteMeeting(id) {
    return api.delete(`/verbali/meetings/${id}`)
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
  updateVerbale(id, data) {
    return api.put(`/verbali/${id}`, data)
  },
  deleteVerbale(id) {
    return api.delete(`/verbali/${id}`)
  },
  signVerbale(id) {
    return api.post(`/verbali/${id}/sign`)
  },
  getSignatures(id) {
    return api.get(`/verbali/${id}/signatures`)
  }
}

export default verbaliService

