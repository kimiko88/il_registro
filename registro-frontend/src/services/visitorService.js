import api from './api'

export const visitorService = {
  // Visitatori esterni
  listTodayVisitors(date) {
    return api.get('/visitors/today', { params: { date } })
  },
  registerVisitor(data) {
    return api.post('/visitors', data)
  },
  recordVisitorExit(id, notes = '') {
    return api.patch(`/visitors/${id}/exit`, { notes })
  },

  // Uscite anticipate studenti
  listTodayEarlyExits(date) {
    return api.get('/visitors/early-exits', { params: { date } })
  },
  recordEarlyExit(data) {
    return api.post('/visitors/early-exits', data)
  },
  recordStudentReturn(id, notes = '') {
    return api.patch(`/visitors/early-exits/${id}/return`, { notes })
  },

  // Segnalazioni guasti / anomalie
  listMaintenanceReports(status = '') {
    return api.get('/visitors/maintenance', { params: { status } })
  },
  createMaintenanceReport(data) {
    return api.post('/visitors/maintenance', data)
  },
  updateMaintenanceStatus(id, data) {
    return api.patch(`/visitors/maintenance/${id}/status`, data)
  }
}

export default visitorService
