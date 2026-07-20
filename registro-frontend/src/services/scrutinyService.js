import api from './api'

export const scrutinyService = {
    getMatrix(classId, semester = 1) {
        return api.get(`/scrutiny/matrix/${classId}`, { params: { semester } })
    },
    save(data) {
        return api.post('/scrutiny/save', data)
    },
    startScrutiny(classId, semester = 1) {
        return api.post(`/scrutiny/class/${classId}/start`, null, { params: { semester } })
    },
    validateScrutiny(classId, semester = 1) {
        return api.post(`/scrutiny/class/${classId}/validate`, null, { params: { semester } })
    },
    closeScrutiny(classId, semester = 1) {
        return api.post(`/scrutiny/class/${classId}/close`, null, { params: { semester } })
    },
    exportPagellaPDF(studentId, classId, semester = 1) {
        return api.get(`/scrutiny/export/${studentId}/pdf`, { params: { class_id: classId, semester }, responseType: 'blob' })
    }
}
