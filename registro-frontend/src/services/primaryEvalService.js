import api from './api'

export const primaryEvalService = {
    async getObjectives(params = {}) {
        return api.get('/primary/objectives', { params })
    },
    async createObjective(data) {
        return api.post('/primary/objectives', data)
    },
    async deleteObjective(id) {
        return api.delete(`/primary/objectives/${id}`)
    },
    async getEvaluations(params = {}) {
        return api.get('/primary/evaluations', { params })
    },
    async saveEvaluationsBatch(data) {
        return api.post('/primary/evaluations/batch', data)
    },
    async getMatrix(params = {}) {
        return api.get('/primary/matrix', { params })
    },
    async getStudentEvaluations(studentId, params = {}) {
        return api.get(`/primary/student/${studentId}`, { params })
    }
}

export default primaryEvalService
