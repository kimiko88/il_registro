import api from './api'

export default {
    getMaterials(classId) {
        return api.get(`/didactic-materials/class/${classId}`)
    },
    createMaterial(data) {
        return api.post('/didactic-materials', data)
    },
    deleteMaterial(id) {
        return api.delete(`/didactic-materials/${id}`)
    }
}
