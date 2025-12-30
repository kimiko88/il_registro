import { api } from 'src/boot/axios'

export const gradeService = {
    async getMyGrades() {
        return api.get('/student/grades')
    },
    async getByClass(classId, subjectId) {
        return api.get(`/teacher/grades/class/${classId}`, { params: { subjectId } })
    },
    async saveGrade(gradeData) {
        return api.post('/teacher/grades', gradeData)
    },
    async updateGrade(id, gradeData) {
        return api.put(`/teacher/grades/${id}`, gradeData)
    },
    async deleteGrade(id) {
        return api.delete(`/teacher/grades/${id}`)
    }
}
