import api from './api'

export const teacherService = {
  getPersonalRegisterPDF(params = {}) {
    return api.get('/teachers/registro-personale/pdf', {
      params,
      responseType: 'blob'
    })
  }
}

export default teacherService
