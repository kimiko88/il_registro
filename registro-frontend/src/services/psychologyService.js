import api from './api'

export const psychologyService = {
  // Consenso Informato Genitoriale Preventivo Obbligatorio (entrambi i genitori per minori)
  signConsent(studentId, schoolYear = '2025/2026') {
    return api.post('/psychology/consent', {
      student_id: studentId,
      school_year: schoolYear
    })
  },

  signBulkConsent(studentIds, schoolYear = '2025/2026') {
    return Promise.all(studentIds.map(id => this.signConsent(id, schoolYear)))
  },

  // Prenotazione Anonima e Riservata Sportello d'Ascolto (CIC)
  bookSession(psychologistId, slotTime, durationMinutes = 45) {
    return api.post('/psychology/book', {
      psychologist_id: psychologistId,
      slot_time: slotTime,
      duration_minutes: durationMinutes
    })
  },

  listMySessions() {
    return api.get('/psychology/sessions')
  },

  getSession(id) {
    return api.get(`/psychology/sessions/${id}`)
  },

  // Note Cliniche Riservate Protette da Segreto Professionale (L. 56/1989)
  updateClinicalNotes(sessionId, notes) {
    return api.put(`/psychology/sessions/${sessionId}/notes`, { notes })
  }
}

export default psychologyService
