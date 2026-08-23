import api from './api'

export const sidiService = {
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
  }
}

export default sidiService

