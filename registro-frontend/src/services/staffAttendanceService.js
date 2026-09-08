import api from './api'

export const staffAttendanceService = {
  /**
   * Ottiene il riepilogo giornaliero delle presenze di tutte le categorie di personale (Docenti + ATA)
   * @param {string} date - Formato YYYY-MM-DD
   */
  async getDailySummary(date) {
    const res = await api.get('/staff-attendance/summary', { params: { date } })
    return res.data
  },

  /**
   * Ottiene la lista nominativa del personale con stato presenza per una data
   * @param {string} date - Formato YYYY-MM-DD
   */
  async getList(date) {
    const res = await api.get('/staff-attendance', { params: { date } })
    return res.data
  },

  /**
   * Registra o aggiorna la presenza di un singolo membro del personale
   * @param {Object} data - { user_id, date, status, notes, strike_code, is_strike_recorded }
   */
  async recordAttendance(data) {
    const res = await api.post('/staff-attendance', data)
    return res.data
  },

  /**
   * Registrazione massiva presenze
   * @param {Object} data - { date, attendances: [{ user_id, status, notes, strike_code }] }
   */
  async recordBulk(data) {
    const res = await api.post('/staff-attendance/bulk', data)
    return res.data
  },

  /**
   * Elimina una registrazione di presenza
   * @param {string} id - ID della presenza
   */
  async deleteAttendance(id) {
    const res = await api.delete(`/staff-attendance/${id}`)
    return res.data
  },

  /**
   * Registra una timbratura badge da hardware
   * @param {Object} data - { badge_code, device_id, swipe_type, raw_data }
   */
  async registerBadgeSwipe(data) {
    const res = await api.post('/staff-attendance/badge-swipe', data)
    return res.data
  },

  /**
   * Elabora manualmente le timbrature pendenti
   */
  async processBadgeSwipes() {
    const res = await api.post('/staff-attendance/badge-swipe/process')
    return res.data
  },

  /**
   * Assegna o aggiorna un badge a un utente
   * @param {Object} data - { user_id, badge_code, notes }
   */
  async assignBadge(data) {
    const res = await api.post('/staff-attendance/badges', data)
    return res.data
  }
}

export default staffAttendanceService
