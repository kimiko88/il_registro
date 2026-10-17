import { describe, it, expect, vi, beforeEach } from 'vitest'
import verbaliService from '@/services/verbaliService'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn()
  }
}))

describe('verbaliService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Riunioni', () => {
    it('listMeetings calls /verbali/meetings with params', () => {
      const params = { class_id: 'c-1' }
      verbaliService.listMeetings(params)
      expect(api.get).toHaveBeenCalledWith('/verbali/meetings', { params })
    })

    it('createMeeting calls /verbali/meetings with payload', () => {
      const payload = { title: 'Consiglio di Classe', meeting_type: 'consiglio_classe' }
      verbaliService.createMeeting(payload)
      expect(api.post).toHaveBeenCalledWith('/verbali/meetings', payload)
    })

    it('deleteMeeting calls DELETE on /verbali/meetings/:id', () => {
      verbaliService.deleteMeeting('m-1')
      expect(api.delete).toHaveBeenCalledWith('/verbali/meetings/m-1')
    })
  })

  describe('Verbali & Lifecycle', () => {
    it('getAllVerbali calls /verbali with params', () => {
      verbaliService.getAllVerbali({ school_id: 's-1' })
      expect(api.get).toHaveBeenCalledWith('/verbali', { params: { school_id: 's-1' } })
    })

    it('getVerbale calls /verbali/:id', () => {
      verbaliService.getVerbale('v-1')
      expect(api.get).toHaveBeenCalledWith('/verbali/v-1')
    })

    it('updateVerbale calls PUT on /verbali/:id', () => {
      const payload = { title: 'Nuovo Titolo', content: 'Nuovo Contenuto' }
      verbaliService.updateVerbale('v-1', payload)
      expect(api.put).toHaveBeenCalledWith('/verbali/v-1', payload)
    })

    it('signVerbale calls POST on /verbali/:id/sign', () => {
      verbaliService.signVerbale('v-1')
      expect(api.post).toHaveBeenCalledWith('/verbali/v-1/sign')
    })

    it('getSignatures calls GET on /verbali/:id/signatures', () => {
      verbaliService.getSignatures('v-1')
      expect(api.get).toHaveBeenCalledWith('/verbali/v-1/signatures')
    })

    it('exportPdf calls GET on /verbali/:id/pdf with blob responseType', () => {
      verbaliService.exportPdf('v-1')
      expect(api.get).toHaveBeenCalledWith('/verbali/v-1/pdf', { responseType: 'blob' })
    })
  })

  describe('Modelli (Templates) & Ordini del Giorno (ODG)', () => {
    it('getTemplates calls /verbali/templates', () => {
      verbaliService.getTemplates({ meeting_type: 'consiglio_classe' })
      expect(api.get).toHaveBeenCalledWith('/verbali/templates', { params: { meeting_type: 'consiglio_classe' } })
    })

    it('createTemplate calls POST on /verbali/templates', () => {
      const payload = { title: 'Modello ODG', meeting_type: 'collegio_docenti', default_agenda: '1. Approvazione' }
      verbaliService.createTemplate(payload)
      expect(api.post).toHaveBeenCalledWith('/verbali/templates', payload)
    })

    it('updateTemplate calls PUT on /verbali/templates/:id', () => {
      const payload = { title: 'Modello Aggiornato' }
      verbaliService.updateTemplate('tpl-1', payload)
      expect(api.put).toHaveBeenCalledWith('/verbali/templates/tpl-1', payload)
    })

    it('deleteTemplate calls DELETE on /verbali/templates/:id', () => {
      verbaliService.deleteTemplate('tpl-1')
      expect(api.delete).toHaveBeenCalledWith('/verbali/templates/tpl-1')
    })
  })
})
