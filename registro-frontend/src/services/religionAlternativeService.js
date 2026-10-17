import api from '@/services/api'

export const religionAlternativeService = {
    async getOptions(academicYear = '2026/2027') {
        const response = await api.get('/religion-alternative/options', {
            params: { academic_year: academicYear }
        })
        return response.data
    },

    async saveOption(payload) {
        const response = await api.post('/religion-alternative/options', payload)
        return response.data
    },

    async getOptionByStudent(studentId, academicYear = '2026/2027') {
        const response = await api.get(`/religion-alternative/options/${studentId}`, {
            params: { academic_year: academicYear }
        })
        return response.data
    },

    async getEvaluations(params = {}) {
        const response = await api.get('/religion-alternative/evaluations', { params })
        return response.data
    },

    async saveEvaluation(payload) {
        const response = await api.post('/religion-alternative/evaluations', payload)
        return response.data
    }
}
