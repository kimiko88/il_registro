import api from './api'

export const maturitaService = {
  // Commissione d'Esame
  saveCommission(data) {
    return api.post('/maturita/commission', data)
  },
  getCommission(params = {}) {
    return api.get('/maturita/commission', { params })
  },

  // Conversione Crediti e Calcolo Punteggi Esame di Stato (D.Lgs. 62/2017)
  calculateCredits(data) {
    return api.post('/maturita/calculate-credits', data)
  },
  calculateScores(data) {
    return api.post('/maturita/calculate-scores', data)
  },

  // Fascicolo Candidato e Tabellone Finale
  saveStudentRecord(data) {
    return api.post('/maturita/student-record', data)
  },
  getStudentRecord(params = {}) {
    return api.get('/maturita/student-record', { params })
  },
  getTabellone(params = {}) {
    return api.get('/maturita/tabellone', { params })
  },

  // Curriculum dello Studente XML (D.M. 88/2020)
  downloadCurriculumXML(params = {}) {
    return api.get('/maturita/curriculum-xml', {
      params,
      responseType: 'blob'
    })
  }
}

export default maturitaService
