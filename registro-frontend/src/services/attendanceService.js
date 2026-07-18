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
    async getMyAttendance() {
        return api.get('/attendance/my-attendance')
    }
}
