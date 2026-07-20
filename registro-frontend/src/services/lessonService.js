import api from './api'

export const lessonService = {
    // Lessons (Registro di Classe)
    async getLessons(classId, subjectId, date) {
        return api.get(`/lessons/class/${classId}`, { params: { subject_id: subjectId, date } })
    },
    async getLessonsByGroup(groupId, date) {
        return api.get(`/lessons/group/${groupId}`, { params: { date } })
    },
    async createLesson(data) {
        return api.post('/lessons', data)
    },

    // Homework (Compiti)
    async getHomeworks(classId) {
        return api.get(`/homeworks/class/${classId}`)
    },
    async createHomework(data) {
        return api.post('/homeworks', data)
    },

    // Student & Parent access to homeworks
    async getMyHomeworks(classId) {
        return api.get(`/homeworks/class/${classId}`)
    }
}
