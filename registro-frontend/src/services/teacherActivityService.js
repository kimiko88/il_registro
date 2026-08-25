import api from './api'

/**
 * teacherActivityService — CRUD per le attività libere del docente
 * (ore a disposizione, gita, riunione, formazione, etc.)
 * Non collegate ad una classe specifica.
 *
 * Endpoints:
 *   GET    /teacher/free-activities?from=YYYY-MM-DD&to=YYYY-MM-DD
 *   POST   /teacher/free-activities
 *   GET    /teacher/free-activities/:id
 *   PUT    /teacher/free-activities/:id
 *   DELETE /teacher/free-activities/:id
 */
export const teacherActivityService = {
  /**
   * Restituisce le attività libere del docente autenticato nell'intervallo indicato.
   * Se from/to sono omessi, il backend usa gli ultimi 30 e i prossimi 30 giorni.
   * @param {string} [from] YYYY-MM-DD
   * @param {string} [to]   YYYY-MM-DD
   */
  async list (from, to) {
    return api.get('/teacher/free-activities', { params: { from, to } })
  },

  /**
   * Recupera una singola attività per ID.
   * @param {string} id UUID
   */
  async getById (id) {
    return api.get(`/teacher/free-activities/${id}`)
  },

  /**
   * Crea una nuova attività libera.
   * @param {Object} data
   * @param {string} data.date          YYYY-MM-DD
   * @param {number} data.start_hour    1–8
   * @param {number} data.duration      numero di ore >= 1
   * @param {string} data.activity_type disponibilita | riunione | formazione | ptof | gita | altro
   * @param {string} data.description   descrizione obbligatoria
   * @param {string} [data.notes]       note opzionali
   */
  async create (data) {
    return api.post('/teacher/free-activities', data)
  },

  /**
   * Aggiorna un'attività libera esistente (solo il proprietario).
   * @param {string} id   UUID
   * @param {Object} data campi da aggiornare (tutti opzionali)
   */
  async update (id, data) {
    return api.put(`/teacher/free-activities/${id}`, data)
  },

  /**
   * Elimina un'attività libera (solo il proprietario).
   * @param {string} id UUID
   */
  async delete (id) {
    return api.delete(`/teacher/free-activities/${id}`)
  }
}

export default teacherActivityService

