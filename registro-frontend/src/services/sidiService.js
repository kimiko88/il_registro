import api from './api'

export const sidiService = {
  // Legacy Export Flows
  async downloadStudentsXml(classId) {
    const res = await api.get('/reports/sidi/students', {
      params: { class_id: classId },
      responseType: 'blob',
      timeout: 60000
    })
    return res.data
  },

  async downloadScrutiniXml(classId, semester = 2) {
    const res = await api.get('/reports/sidi/scrutini', {
      params: { class_id: classId, semester },
      responseType: 'blob',
      timeout: 60000
    })
    return res.data
  },

  async downloadAttendanceCsv(classId) {
    const res = await api.get('/reports/sidi/attendance', {
      params: { class_id: classId },
      responseType: 'blob',
      timeout: 60000
    })
    return res.data
  },

  // Cooperazione Applicativa Diretta con i WebService SIDI (MIM)
  getCooperationConfig() {
    return api.get('/sidi/cooperation-config')
  },
  syncStudentCodes(students = []) {
    return api.post('/sidi/sync-student-codes', { students })
  },
  pushScrutinyResults(classId, schoolYear, results = []) {
    return api.post('/sidi/push-scrutiny-results', {
      class_id: classId,
      school_year: schoolYear,
      results
    })
  },

  // Tracciati Ufficiali Posizionali MIM / ANS / Nuova PassWeb INPS
  getPositionalTemplates() {
    return api.get('/sidi/positional/templates')
  },
  generatePositional(type, schoolYear = '2025/2026', data = {}) {
    return api.post('/sidi/positional/generate', {
      type,
      school_year: schoolYear,
      ...data
    })
  },
  downloadPositional(id) {
    return api.get(`/sidi/positional/download/${id}`, { responseType: 'blob' })
  }
}

export default sidiService
