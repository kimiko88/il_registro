import { api } from 'src/boot/axios'

export const attendanceService = {
    async markAttendance(data) {
        return api.post('/teacher/attendance/mark', data)
    },
    async getByClass(classId, date) {
        return api.get(`/teacher/attendance/class/${classId}`, { params: { date } })
    },
    async justify(absenceId, justification) {
        return api.post(`/parent/attendance/${absenceId}/justify`, justification)
    }
}
