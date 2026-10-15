import api from './api'

export const privacyService = {
  // Consensi Annuali Famiglia e Semaforo Privacy Docenti
  saveConsent(data) {
    return api.post('/privacy/consent', data)
  },
  getConsent(studentId, schoolYear = '2025/2026') {
    return api.get(`/privacy/consent/${studentId}`, { params: { school_year: schoolYear } })
  },
  getClassBadges(studentIds = [], schoolYear = '2025/2026') {
    return api.post('/privacy/class-badges', {
      student_ids: studentIds,
      school_year: schoolYear
    })
  },

  // Registro Trattamenti Art. 30 GDPR
  createTreatment(data) {
    return api.post('/privacy/treatments', data)
  },
  listTreatments() {
    return api.get('/privacy/treatments')
  }
}

export default privacyService
