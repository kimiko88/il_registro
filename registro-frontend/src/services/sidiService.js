import api from './api'

export const sidiService = {
  async downloadStudentsXml(classId) {
    const res = await api.get(`/reports/sidi/students?class_id=${classId}`, { responseType: 'blob' })
    return res.data
  },

  async downloadScrutiniXml(classId, semester = 2) {
    const res = await api.get(`/reports/sidi/scrutini?class_id=${classId}&semester=${semester}`, { responseType: 'blob' })
    return res.data
  },

  async downloadAttendanceCsv(classId) {
    const res = await api.get(`/reports/sidi/attendance?class_id=${classId}`, { responseType: 'blob' })
    return res.data
  }
}

export default sidiService
