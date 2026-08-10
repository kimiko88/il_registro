import api from './api'

export const lessonService = {
    // Lessons (Registro di Classe)
    async getLessons(classId, subjectId, date) {
        return api.get(`/lessons/class/${classId}`, { params: { subject_id: subjectId, date } })
    },
    async getLessonsByGroup(groupId, date) {
        return api.get(`/lessons/group/${groupId}`, { params: { date } })
    },
    async getTeacherDiary(from, to) {
        return api.get('/lessons/my-diary', { params: { from, to } })
    },
    async getLessonById(id) {
        return api.get(`/lessons/${id}`)
    },
    async createLesson(data) {
        return api.post('/lessons', data)
    },
    async updateLesson(id, data) {
        return api.put(`/lessons/${id}`, data)
    },
    async deleteLesson(id) {
        return api.delete(`/lessons/${id}`)
    },
    async getActivityHours(classId) {
        return api.get(`/lessons/class/${classId}/activity-hours`)
    },

    // Homework (Compiti)
    async getHomeworks(classId) {
        return api.get(`/homeworks/class/${classId}`)
    },
    async createHomework(data) {
        return api.post('/homeworks', data)
    },
    async updateHomework(id, data) {
        return api.put(`/homeworks/${id}`, data)
    },
    async deleteHomework(id) {
        return api.delete(`/homeworks/${id}`)
    },

    // Student & Parent access to homeworks
    async getMyHomeworks(classId) {
        return api.get(`/homeworks/class/${classId}`)
    }
}
