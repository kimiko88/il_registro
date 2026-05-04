import api from './api'

export const gradeService = {
    async getMyGrades() {
        return api.get('/grades/my-grades')
    },
    async getByClass(classId, subjectId) {
        return api.get(`/grades/class/${classId}`, { params: { subject_id: subjectId } })
    },
    async saveGrade(gradeData) {
        return api.post('/grades', gradeData)
    },
    async updateGrade(id, gradeData) {
        return api.patch(`/grades/${id}`, gradeData)
    },
    async deleteGrade(id) {
        return api.delete(`/grades/${id}`)
    },
    async getChildGrades(studentId, params) {
        return api.get(`/grades/child-grades/${studentId}`, { params })
    },
    async getStudentGrades(studentId, params) {
        return api.get(`/grades/student/${studentId}`, { params })
    }
}
