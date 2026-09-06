import api from './api'

function cleanSem(semester) {
    if (typeof semester === 'number') return semester
    const s = String(semester || '').toLowerCase()
    if (s.includes('2') || s.includes('second')) return 2
    return 1
}

export const scrutinyService = {
    getMatrix(classId, semester = 1) {
        return api.get(`/scrutiny/matrix/${classId}`, { params: { semester: cleanSem(semester) } })
    },
    save(data) {
        return api.post('/scrutiny/save', data)
    },
    startScrutiny(classId, semester = 1) {
        return api.post(`/scrutiny/class/${classId}/start`, null, { params: { semester: cleanSem(semester) } })
    },
    validateScrutiny(classId, semester = 1) {
        return api.post(`/scrutiny/class/${classId}/validate`, null, { params: { semester: cleanSem(semester) } })
    },
    closeScrutiny(classId, semester = 1) {
        return api.post(`/scrutiny/class/${classId}/close`, null, { params: { semester: cleanSem(semester) } })
    },
    exportPagellaPDF(studentId, classId, semester = 1) {
        return api.get(`/scrutiny/export/${studentId}/pdf`, {
            params: { class_id: classId, semester: cleanSem(semester) },
            responseType: 'blob',
            timeout: 60000
        })
    },
    exportClassScrutinyZip(classId, semester = 1) {
        return api.get(`/scrutiny/class/${classId}/export-zip`, {
            params: { semester: cleanSem(semester) },
            responseType: 'blob',
            timeout: 90000
        })
    },

    // Deficiency & Deferred Scrutiny API helpers
    saveDeficiency(data) {
        return api.post('/scrutiny/deficiencies', data)
    },
    getStudentDeficiencies(studentId) {
        return api.get(`/scrutiny/deficiencies/student/${studentId}`)
    },
    getClassDeficiencies(classId, semester = 0) {
        return api.get(`/scrutiny/deficiencies/class/${classId}`, { params: { semester } })
    },
    saveDeferredScrutiny(data) {
        return api.post('/scrutiny/deferred', data)
    }
}

export default scrutinyService

