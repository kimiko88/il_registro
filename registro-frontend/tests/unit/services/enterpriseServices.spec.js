import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockApi = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn()
}))

vi.mock('@/services/api', () => ({
  default: mockApi
}))
vi.mock('./api', () => ({
  default: mockApi
}))

import { paymentService } from 'src/services/paymentService'
import { interpelliService } from 'src/services/interpelliService'
import { alboPretorioService } from 'src/services/alboPretorioService'
import { digitalStampService } from 'src/services/digitalStampService'
import { maturitaService } from 'src/services/maturitaService'
import { mealsService } from 'src/services/mealsService'
import { inventoryService } from 'src/services/inventoryService'
import { privacyService } from 'src/services/privacyService'
import { sidiService } from 'src/services/sidiService'
import { psychologyService } from 'src/services/psychologyService'

describe('Enterprise Dealbreakers Frontend Services', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('1. PagoPA & OPI (paymentService)', () => {
    it('fetches bollettino notice with IUV and QR Code', async () => {
      mockApi.get.mockResolvedValue({ data: { iuv: '01210000000014288' } })
      const res = await paymentService.getBollettino('pay-1')
      expect(mockApi.get).toHaveBeenCalledWith('/payments/pay-1/bollettino')
      expect(res.data.iuv).toBe('01210000000014288')
    })

    it('creates online direct checkout session', async () => {
      mockApi.post.mockResolvedValue({ data: { checkout_url: 'https://checkout.pagopa.it/sess' } })
      await paymentService.createCheckoutSession('pay-1', 'https://return.url')
      expect(mockApi.post).toHaveBeenCalledWith('/payments/pay-1/checkout', {
        return_url: 'https://return.url'
      })
    })

    it('submits OPI reconciliation file stream', async () => {
      mockApi.post.mockResolvedValue({ data: { total_reconciled: 5 } })
      await paymentService.reconcileOPI('<xml></xml>', 'OPI_XML')
      expect(mockApi.post).toHaveBeenCalledWith(
        '/payments/reconcile-opi?format=OPI_XML',
        '<xml></xml>',
        { headers: { 'Content-Type': 'application/xml' } }
      )
    })

    it('fetches overdue solleciti list', async () => {
      mockApi.get.mockResolvedValue({ data: { total: 2, solleciti: [] } })
      await paymentService.getSolleciti()
      expect(mockApi.get).toHaveBeenCalledWith('/payments/solleciti')
    })
  })

  describe('2. Interpelli Supplenze (interpelliService)', () => {
    it('fetches public notices', async () => {
      mockApi.get.mockResolvedValue({ data: [] })
      await interpelliService.listPublic({ concorso_class: 'A026' })
      expect(mockApi.get).toHaveBeenCalledWith('/public/interpelli', {
        params: { concorso_class: 'A026' }
      })
    })

    it('submits candidacy with DPR 445 declaration', async () => {
      mockApi.post.mockResolvedValue({ data: { id: 'cand-1' } })
      const candData = { candidate_name: 'Elena', dpr445_declared: true }
      await interpelliService.submitCandidatura('notice-1', candData)
      expect(mockApi.post).toHaveBeenCalledWith('/public/interpelli/notice-1/candidatura', candData)
    })

    it('creates administrative notice and gets graduatoria', async () => {
      mockApi.post.mockResolvedValue({ data: { id: 'notice-1' } })
      mockApi.get.mockResolvedValue({ data: { graduatoria: [] } })

      await interpelliService.createNotice({ title: 'A026' })
      expect(mockApi.post).toHaveBeenCalledWith('/interpelli', { title: 'A026' })

      await interpelliService.getGraduatoria('notice-1')
      expect(mockApi.get).toHaveBeenCalledWith('/interpelli/notice-1/graduatoria')
    })

    it('convokes candidate and records response', async () => {
      mockApi.post.mockResolvedValue({ data: { status: 'convocata' } })
      await interpelliService.convoca('cand-1', 24)
      expect(mockApi.post).toHaveBeenCalledWith('/interpelli/candidature/cand-1/convoca', {
        hours_to_respond: 24
      })

      await interpelliService.rispondi('cand-1', 'accettata', 'Disponibile da subito')
      expect(mockApi.post).toHaveBeenCalledWith('/interpelli/candidature/cand-1/risposta', {
        risposta: 'accettata',
        notes: 'Disponibile da subito'
      })
    })
  })

  describe('3. Albo Pretorio Online (alboPretorioService)', () => {
    it('lists public legal board acts and transparency items', async () => {
      mockApi.get.mockResolvedValue({ data: { items: [] } })
      await alboPretorioService.listPublic({ anno: 2026 })
      expect(mockApi.get).toHaveBeenCalledWith('/public/albo-pretorio', { params: { anno: 2026 } })

      await alboPretorioService.listTrasparenza()
      expect(mockApi.get).toHaveBeenCalledWith('/public/amministrazione-trasparente', { params: {} })
    })

    it('downloads ANAC L. 190/2012 XML file', async () => {
      mockApi.get.mockResolvedValue({ data: new Blob() })
      await alboPretorioService.downloadANACXML({ anno: 2026 })
      expect(mockApi.get).toHaveBeenCalledWith('/public/amministrazione-trasparente/anac.xml', {
        params: { anno: 2026 },
        responseType: 'blob'
      })
    })

    it('publishes and defigges act with legal relata', async () => {
      mockApi.post.mockResolvedValue({ data: { id: 'act-1' } })
      await alboPretorioService.publishAct({ subject: 'Bando PNRR' })
      expect(mockApi.post).toHaveBeenCalledWith('/albo-pretorio', { subject: 'Bando PNRR' })

      await alboPretorioService.defiggiAtto('act-1', 'DS Mario Rossi', true)
      expect(mockApi.post).toHaveBeenCalledWith('/albo-pretorio/act-1/defissione', null, {
        params: { ds_name: 'DS Mario Rossi', force: true }
      })
    })
  })

  describe('4. FEQ CSC Batch & Digital Stamp (digitalStampService)', () => {
    it('signs documents batch via Cloud Signature Consortium (CSC)', async () => {
      mockApi.post.mockResolvedValue({ data: { total_signed: 10 } })
      await digitalStampService.batchSignCSC(['doc-1', 'doc-2'], '123456', '987654')
      expect(mockApi.post).toHaveBeenCalledWith('/signatures/csc/batch-sign', {
        document_ids: ['doc-1', 'doc-2'],
        pin: '123456',
        otp: '987654'
      })
    })

    it('creates and verifies digital stamp (glifo CAD art. 23)', async () => {
      mockApi.post.mockResolvedValue({ data: { glifo_token: 'GLF-123' } })
      await digitalStampService.createDigitalStamp({ document_id: 'doc-1' })
      expect(mockApi.post).toHaveBeenCalledWith('/signatures/digital-stamp', { document_id: 'doc-1' })

      await digitalStampService.verifyDigitalStamp('CAD-GLIFO|...', 'hash-1')
      expect(mockApi.post).toHaveBeenCalledWith('/signatures/digital-stamp/verify', {
        payload: 'CAD-GLIFO|...',
        original_sha256: 'hash-1'
      })

      mockApi.get.mockResolvedValue({ data: { status: 'VALIDO_CONFORME_CAD_ART_23' } })
      await digitalStampService.verifyPublicToken('GLF-123')
      expect(mockApi.get).toHaveBeenCalledWith('/public/verifica-glifo/GLF-123')
    })
  })

  describe('5. Maturità & Curriculum dello Studente (maturitaService)', () => {
    it('saves commission and calculates credits and scores', async () => {
      mockApi.post.mockResolvedValue({ data: {} })
      await maturitaService.saveCommission({ commission_code: 'COMM-1' })
      expect(mockApi.post).toHaveBeenCalledWith('/maturita/commission', { commission_code: 'COMM-1' })

      await maturitaService.calculateCredits({ grade_3rd: 8.5 })
      expect(mockApi.post).toHaveBeenCalledWith('/maturita/calculate-credits', { grade_3rd: 8.5 })

      await maturitaService.calculateScores({ written1_score: 18.0 })
      expect(mockApi.post).toHaveBeenCalledWith('/maturita/calculate-scores', { written1_score: 18.0 })
    })

    it('manages candidate record and downloads curriculum XML', async () => {
      mockApi.get.mockResolvedValue({ data: {} })
      await maturitaService.getTabellone({ class_id: '5A' })
      expect(mockApi.get).toHaveBeenCalledWith('/maturita/tabellone', { params: { class_id: '5A' } })

      await maturitaService.downloadCurriculumXML({ student_id: 'st-1' })
      expect(mockApi.get).toHaveBeenCalledWith('/maturita/curriculum-xml', {
        params: { student_id: 'st-1' },
        responseType: 'blob'
      })
    })
  })

  describe('6. Refezione Scolastica & Diete Speciali (mealsService)', () => {
    it('records morning roll call and fetches catering report', async () => {
      mockApi.post.mockResolvedValue({ data: { total_meals: 20 } })
      await mealsService.recordRollCall('class-1', '2026-10-07', ['st-1', 'st-2'])
      expect(mockApi.post).toHaveBeenCalledWith('/meals/roll-call', {
        class_id: 'class-1',
        date: '2026-10-07',
        present_ids: ['st-1', 'st-2']
      })

      mockApi.get.mockResolvedValue({ data: {} })
      await mealsService.getCateringReport('2026-10-07')
      expect(mockApi.get).toHaveBeenCalledWith('/meals/catering-report', {
        params: { date: '2026-10-07' }
      })
    })

    it('registers certified special diet and handles wallet top-up', async () => {
      mockApi.post.mockResolvedValue({ data: { id: 'diet-1' } })
      await mealsService.registerDiet({ student_id: 'st-1', specific_diet: 'celiachia' })
      expect(mockApi.post).toHaveBeenCalledWith('/meals/special-diet', {
        student_id: 'st-1',
        specific_diet: 'celiachia'
      })

      await mealsService.topUpWallet('st-1', 50.0, '01210000000014288')
      expect(mockApi.post).toHaveBeenCalledWith('/meals/wallet/st-1/topup', {
        amount: 50.0,
        pagopa_iuv: '01210000000014288'
      })
    })
  })

  describe('7. Inventario Scolastico e Comodato (inventoryService)', () => {
    it('manages assets and barcode labels', async () => {
      mockApi.post.mockResolvedValue({ data: { id: 'ast-1' } })
      await inventoryService.createAsset({ description: 'Notebook Lenovo' })
      expect(mockApi.post).toHaveBeenCalledWith('/inventory/assets', { description: 'Notebook Lenovo' })

      mockApi.get.mockResolvedValue({ data: { barcode_data: '123' } })
      await inventoryService.getLabel('ast-1')
      expect(mockApi.get).toHaveBeenCalledWith('/inventory/assets/ast-1/label')
    })

    it('manages device loans and return tracking', async () => {
      mockApi.post.mockResolvedValue({ data: { id: 'loan-1' } })
      await inventoryService.createLoan({ asset_id: 'ast-1', student_id: 'st-1' })
      expect(mockApi.post).toHaveBeenCalledWith('/inventory/loans', {
        asset_id: 'ast-1',
        student_id: 'st-1'
      })

      await inventoryService.returnLoan('loan-1', 'ottima')
      expect(mockApi.post).toHaveBeenCalledWith('/inventory/loans/loan-1/return', {
        condition: 'ottima'
      })
    })
  })

  describe('8. Privacy GDPR & Semaforo Consensi (privacyService)', () => {
    it('saves family consent and fetches class badges', async () => {
      mockApi.post.mockResolvedValue({ data: { traffic_light_badge: 'VERDE' } })
      await privacyService.saveConsent({ student_id: 'st-1', photo_video_social_consent: true })
      expect(mockApi.post).toHaveBeenCalledWith('/privacy/consent', {
        student_id: 'st-1',
        photo_video_social_consent: true
      })

      await privacyService.getClassBadges(['st-1', 'st-2'], '2025/2026')
      expect(mockApi.post).toHaveBeenCalledWith('/privacy/class-badges', {
        student_ids: ['st-1', 'st-2'],
        school_year: '2025/2026'
      })
    })

    it('manages Art. 30 GDPR treatments register', async () => {
      mockApi.post.mockResolvedValue({ data: { id: 'tr-1' } })
      await privacyService.createTreatment({ activity_name: 'Foto' })
      expect(mockApi.post).toHaveBeenCalledWith('/privacy/treatments', { activity_name: 'Foto' })

      mockApi.get.mockResolvedValue({ data: [] })
      await privacyService.listTreatments()
      expect(mockApi.get).toHaveBeenCalledWith('/privacy/treatments')
    })
  })

  describe('9. WebService SIDI Cooperation (sidiService)', () => {
    it('gets configuration, syncs student codes and pushes scrutiny results', async () => {
      mockApi.get.mockResolvedValue({ data: { endpoint_url: 'https://cooperazione...' } })
      await sidiService.getCooperationConfig()
      expect(mockApi.get).toHaveBeenCalledWith('/sidi/cooperation-config')

      mockApi.post.mockResolvedValue({ data: { total_updated: 2 } })
      await sidiService.syncStudentCodes([{ codice_fiscale: 'RSSMRA...' }])
      expect(mockApi.post).toHaveBeenCalledWith('/sidi/sync-student-codes', {
        students: [{ codice_fiscale: 'RSSMRA...' }]
      })

      await sidiService.pushScrutinyResults('5A', '2024/2025', [{ codice_sidi: '123' }])
      expect(mockApi.post).toHaveBeenCalledWith('/sidi/push-scrutiny-results', {
        class_id: '5A',
        school_year: '2024/2025',
        results: [{ codice_sidi: '123' }]
      })
    })
  })

  describe('10. Sportello Psicologico CIC (psychologyService)', () => {
    it('signs parental consent and books anonymous appointment', async () => {
      mockApi.post.mockResolvedValue({ data: { status: 'approvato_entrambi' } })
      await psychologyService.signConsent('st-1', '2025/2026')
      expect(mockApi.post).toHaveBeenCalledWith('/psychology/consent', {
        student_id: 'st-1',
        school_year: '2025/2026'
      })

      mockApi.post.mockResolvedValue({ data: { anonymous_alias: 'ALUNNO-CIC-A1B2C3' } })
      await psychologyService.bookSession('psy-1', '2026-10-15 10:00', 45)
      expect(mockApi.post).toHaveBeenCalledWith('/psychology/book', {
        psychologist_id: 'psy-1',
        slot_time: '2026-10-15 10:00',
        duration_minutes: 45
      })
    })

    it('fetches sessions and updates private clinical notes', async () => {
      mockApi.get.mockResolvedValue({ data: { id: 'sess-1' } })
      await psychologyService.getSession('sess-1')
      expect(mockApi.get).toHaveBeenCalledWith('/psychology/sessions/sess-1')

      mockApi.put.mockResolvedValue({ data: { status: 'svolto' } })
      await psychologyService.updateClinicalNotes('sess-1', 'Note riservate')
      expect(mockApi.put).toHaveBeenCalledWith('/psychology/sessions/sess-1/notes', {
        notes: 'Note riservate'
      })
    })
  })
})
