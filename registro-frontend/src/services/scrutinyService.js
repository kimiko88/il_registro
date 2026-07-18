import api from './api'

export const scrutinyService = {
    getMatrix(classId, semester = 1) {
        return api.get(`/scrutiny/matrix/${classId}`, { params: { semester } })
    },
    save(data) {
        return api.post('/scrutiny/save', data)
    }
}
