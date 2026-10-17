import api from './api'

export const strikeService = {
  /**
   * Elenco avvisi di sciopero visibili per la scuola corrente
   */
  async getStrikeNotices() {
    const res = await api.get('/strike-notices')
    return res.data
  },

  /**
   * Dettaglio singolo avviso con stato dichiarazione dell'utente
   */
  async getStrikeNotice(id) {
    const res = await api.get(`/strike-notices/${id}`)
    return res.data
  },

  /**
   * Crea una nuova comunicazione di sciopero con scadenza (DSGA / Preside / Admin)
   * @param {Object} data - { title, proclaimed_by, strike_date, declaration_deadline, notes, publish_to_bacheca }
   */
  async createStrikeNotice(data) {
    const res = await api.post('/strike-notices', data)
    return res.data
  },

  /**
   * Elimina un avviso di sciopero
   */
  async deleteStrikeNotice(id) {
    const res = await api.delete(`/strike-notices/${id}`)
    return res.data
  },

  /**
   * Invia la dichiarazione preventiva di adesione (dipendente)
   * @param {string} noticeId - ID dell'avviso di sciopero
   * @param {string} intention - 'participates' | 'not_participates' | 'undecided'
   */
  async submitDeclaration(noticeId, intention) {
    const res = await api.post(`/strike-notices/${noticeId}/declare`, { intention })
    return res.data
  },

  /**
   * Ottiene il quadro preventivo con conteggi, percentuali e lista nominativa (DSGA / Preside / Admin)
   * @param {string} noticeId - ID dell'avviso di sciopero
   */
  async getNoticeSummary(noticeId) {
    const res = await api.get(`/strike-notices/${noticeId}/summary`)
    return res.data
  }
}

export default strikeService
