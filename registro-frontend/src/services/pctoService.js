import api from './api'

export const pctoService = {
    async getMyProjects() {
        return api.get('/pcto/my-projects')
    },
    async getProjectDetails(id) {
        return api.get(`/pcto/my-projects/${id}`)
    },
    async getStudentProjects(studentId) {
        return api.get(`/pcto/student/${studentId}`)
    },
    async createProject(data) {
        return api.post('/pcto/projects', data)
    },
    async updateProject(id, data) {
        return api.put(`/pcto/projects/${id}`, data)
    },
    async logHours(data) {
        return api.post('/pcto/hours', data)
    },
    async approveHours(id, approved) {
        return api.post(`/pcto/hours/${id}/approve`, { approved })
    }
}

export default pctoService

