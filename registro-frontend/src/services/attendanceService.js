import api from './api'

export const attendanceService = {
    async markAttendance(data) {
        return api.post('/attendance/mark', data)
    },
    async getByClass(classId, date) {
        return api.get(`/attendance/class/${classId}`, { params: { date } })
    },
    async justify(absenceId, justification) {
        return api.post('/attendance/justify', { absence_id: absenceId, ...justification })
    },
    async getChildAttendance(studentId) {
        return api.get(`/attendance/child-attendance/${studentId}`)
    },
    async getChildAttendanceSummary(studentId) {
        return api.get(`/attendance/child-attendance/${studentId}/summary`)
    },
    async getChildAttendanceTrends(studentId) {
        return api.get(`/attendance/child-attendance/${studentId}/trends`)
    },
    async getMyAttendance() {
        return api.get('/attendance/my-attendance')
    },
    async updateAttendance(id, data) {
        return api.put(`/attendance/${id}`, data)
    },
    async deleteHour(classId, date, hour) {
        return api.delete(`/attendance/class/${classId}/hour/${hour}`, { params: { date } })
    },
    async exportAttendance(classId, date) {
        return api.get('/attendance/export', { params: { class_id: classId, date }, responseType: 'blob' })
    },
    async getStudentSummary(studentId) {
        const res = await api.get(`/attendance/students/${studentId}/summary`)
        return res.data
    }
}
