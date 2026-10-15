import api from './api'

export const mealsService = {
  // Rilevazione presenze rapida mattutina per la ditta di ristorazione
  recordRollCall(classId, date, presentIds = []) {
    return api.post('/meals/roll-call', {
      class_id: classId,
      date,
      present_ids: presentIds
    })
  },
  getCateringReport(date) {
    return api.get('/meals/catering-report', { params: { date } })
  },

  // Anagrafica Diete Speciali Certificate (sanitarie ed etico-religiose)
  registerDiet(data) {
    return api.post('/meals/special-diet', data)
  },
  getDiet(studentId) {
    return api.get(`/meals/special-diet/${studentId}`)
  },
  listDiets() {
    return api.get('/meals/special-diets')
  },

  // Borsellino Elettronico Mensa
  getWallet(studentId) {
    return api.get(`/meals/wallet/${studentId}`)
  },
  topUpWallet(studentId, amount, pagopaIUV = '') {
    return api.post(`/meals/wallet/${studentId}/topup`, {
      amount,
      pagopa_iuv: pagopaIUV
    })
  }
}

export default mealsService
