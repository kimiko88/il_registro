import api from './api'

export const accessibilityService = {
  /**
   * Invia una segnalazione di accessibilità (AgID / WCAG 2.2)
   * @param {Object} data { name, email, barrierType, description }
   */
  async submitFeedback(data) {
    const payload = {
      name: data.name,
      email: data.email,
      barrier_type: data.barrierType || data.barrier_type,
      description: data.description
    }
    const response = await api.post('/public/accessibility-feedback', payload)
    return response.data
  },

  /**
   * Recupera le segnalazioni di accessibilità per gli amministratori
   * @param {Object} params { status, school_id, limit, offset }
   */
  async listFeedbacks(params = {}) {
    const response = await api.get('/admin/accessibility-feedbacks', { params })
    return response.data
  },

  /**
   * Aggiorna lo stato di una segnalazione
   * @param {string} id 
   * @param {string} status 'open' | 'in_progress' | 'resolved'
   * @param {string} responseNotes 
   */
  async updateFeedbackStatus(id, status, responseNotes = '') {
    const response = await api.patch(`/admin/accessibility-feedbacks/${id}/status`, {
      status,
      response_notes: responseNotes
    })
    return response.data
  }
}

export default accessibilityService
