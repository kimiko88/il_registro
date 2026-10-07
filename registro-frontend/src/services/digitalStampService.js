import api from './api'

export const digitalStampService = {
  // CSC Remote Batch Signing (Dirigente / DSGA)
  batchSignCSC(documentIds, pin, otp) {
    return api.post('/signatures/csc/batch-sign', {
      document_ids: documentIds,
      pin,
      otp
    })
  },

  // Timbro Digitale di Sicurezza / Glifo (CAD art. 23)
  createDigitalStamp(data) {
    return api.post('/signatures/digital-stamp', data)
  },
  verifyDigitalStamp(payload, originalSHA256) {
    return api.post('/signatures/digital-stamp/verify', {
      payload,
      original_sha256: originalSHA256
    })
  },
  verifyPublicToken(token) {
    return api.get(`/public/verifica-glifo/${token}`)
  }
}

export default digitalStampService
