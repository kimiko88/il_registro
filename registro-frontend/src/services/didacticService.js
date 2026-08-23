import api from './api'

export const didacticService = {
    getMaterials(classId) {
        return api.get(`/didactic-materials/class/${classId}`)
    },
    createMaterial(data) {
        if (data instanceof FormData) {
            return api.post('/didactic-materials', data, {
                headers: { 'Content-Type': 'multipart/form-data' },
                timeout: 60000
            })
        }
        return api.post('/didactic-materials', data)
    },
    deleteMaterial(id) {
        return api.delete(`/didactic-materials/${id}`)
    },
    downloadMaterial(id) {
        return api.get(`/didactic-materials/${id}/download`, {
            responseType: 'blob',
            timeout: 60000
        })
    }
}

export default didacticService

