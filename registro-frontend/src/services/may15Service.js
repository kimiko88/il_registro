import api from '@/services/api'

export const may15Service = {
  getDocument(classId, academicYear = '2025/2026') {
    return api.get(`/may15/class/${classId}`, { params: { academic_year: academicYear } })
  },

  saveDocument(classId, payload) {
    return api.put(`/may15/class/${classId}`, payload)
  },

  publishDocument(classId, academicYear = '2025/2026') {
    return api.post(`/may15/class/${classId}/publish`, null, { params: { academic_year: academicYear } })
  }
}

export default may15Service
