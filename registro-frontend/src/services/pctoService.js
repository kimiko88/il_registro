import api from './api'

export const pctoService = {
    async getMyProjects() {
        return api.get('/pcto/my-projects')
    },
    async getProjectDetails(id) {
        return api.get(`/pcto/my-projects/${id}`)
    },
    async logHours(data) {
        return api.post('/pcto/hours', data)
    }
}
