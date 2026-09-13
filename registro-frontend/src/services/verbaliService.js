import api from './api'

export const verbaliService = {
  // Riunioni
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

  // Verbali
  getAllVerbali(params) {
    return api.get('/verbali', { params })
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
  },
  exportPdf(id) {
    return api.get(`/verbali/${id}/pdf`, { responseType: 'blob' })
  },
  enqueueAsyncPdf(id) {
    return api.post(`/verbali/${id}/async-pdf`)
  },
  getPdfJobStatus(jobId) {
    return api.get(`/verbali/pdf-jobs/${jobId}`)
  },

  // Modelli / Templates per Verbali & Ordini del Giorno (ODG)
  getTemplates(params) {
    return api.get('/verbali/templates', { params })
  },
  getTemplate(id) {
    return api.get(`/verbali/templates/${id}`)
  },
  createTemplate(data) {
    return api.post('/verbali/templates', data)
  },
  updateTemplate(id, data) {
    return api.put(`/verbali/templates/${id}`, data)
  },
  deleteTemplate(id) {
    return api.delete(`/verbali/templates/${id}`)
  }
}

export default verbaliService
