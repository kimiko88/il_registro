import api from './api'

export const alboPretorioService = {
  // Public Albo Pretorio (L. 69/2009 & D.Lgs. 33/2013)
  listPublic(params = {}) {
    return api.get('/public/albo-pretorio', { params })
  },
  listTrasparenza(params = {}) {
    return api.get('/public/amministrazione-trasparente', { params })
  },
  downloadANACXML(params = {}) {
    return api.get('/public/amministrazione-trasparente/anac.xml', {
      params,
      responseType: 'blob'
    })
  },

  // Protected Administration
  publishAct(data) {
    return api.post('/albo-pretorio', data)
  },
  defiggiAtto(id, dsName = 'Dirigente Scolastico', force = false) {
    return api.post(`/albo-pretorio/${id}/defissione`, null, {
      params: { ds_name: dsName, force }
    })
  },
  getCertificato(id, dsName = 'Dirigente Scolastico') {
    return api.get(`/albo-pretorio/${id}/certificato`, {
      params: { ds_name: dsName }
    })
  }
}

export default alboPretorioService
