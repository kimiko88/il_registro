import api from './api'

export const gradeService = {
    async getMyGrades(params) {
        return params ? api.get('/grades/my-grades', { params }) : api.get('/grades/my-grades')
    },
    async getByClass(classId, subjectId, params = {}) {
        const queryParams = { ...params }
        if (subjectId) queryParams.subject_id = subjectId
        return api.get(`/grades/class/${classId}`, { params: queryParams })
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
    },
    async createClassTest(testData) {
        return api.post('/grades/tests', testData)
    },
    async createTestWithGrades(testData) {
        return this.createClassTest(testData)
    },
    async getClassTests(classId, subjectId) {
        return api.get('/grades/tests', { params: { class_id: classId, subject_id: subjectId } })
    },
    async getTestsByClass(classId, subjectId) {
        return this.getClassTests(classId, subjectId)
    },
    async deleteClassTest(id) {
        return api.delete(`/grades/tests/${id}`)
    },
    async deleteTest(id) {
        return this.deleteClassTest(id)
    },
    async updateClassTest(id, testData) {
        return api.patch(`/grades/tests/${id}`, testData)
    },
    async updateTestWithGrades(id, testData) {
        return this.updateClassTest(id, testData)
    },
    async getUpcomingTestsForClass(classId) {
        return api.get(`/grades/tests/class/${classId}`)
    },
    async bulkImport(formData) {
        return api.post('/grades/bulk-import', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
            timeout: 60000
        })
    },
    async getSemesterReport(semester = 1) {
        return api.get(`/grades/my-grades/semester/${semester}`)
    },
    async downloadReportCardPDF(semester = 1) {
        return api.get(`/grades/my-grades/semester/${semester}/pdf`, { responseType: 'blob', timeout: 60000 })
    }
}

export default gradeService

