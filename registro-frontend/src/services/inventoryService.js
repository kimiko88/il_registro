import api from './api'

export const inventoryService = {
  // Registro Cespiti e Beni Mobili (D.I. 129/2018)
  createAsset(data) {
    return api.post('/inventory/assets', data)
  },
  listAssets(params = {}) {
    return api.get('/inventory/assets', { params })
  },
  getLabel(assetId) {
    return api.get(`/inventory/assets/${assetId}/label`)
  },

  // Contratti Comodato d'Uso Gratuito Dispositivi
  createLoan(data) {
    return api.post('/inventory/loans', data)
  },
  returnLoan(loanId, condition = 'buona') {
    return api.post(`/inventory/loans/${loanId}/return`, { condition })
  }
}

export default inventoryService
