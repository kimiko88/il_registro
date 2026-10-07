import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import EnterpriseHub from 'src/pages/admin/EnterpriseHub.vue'
import { paymentService } from 'src/services/paymentService'
import { digitalStampService } from 'src/services/digitalStampService'
import { maturitaService } from 'src/services/maturitaService'
import { mealsService } from 'src/services/mealsService'
import { psychologyService } from 'src/services/psychologyService'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn().mockResolvedValue({ data: {} }),
    post: vi.fn().mockResolvedValue({ data: {} }),
    put: vi.fn().mockResolvedValue({ data: {} }),
    delete: vi.fn().mockResolvedValue({ data: {} })
  }
}))
vi.mock('src/services/api', () => ({
  default: {
    get: vi.fn().mockResolvedValue({ data: {} }),
    post: vi.fn().mockResolvedValue({ data: {} }),
    put: vi.fn().mockResolvedValue({ data: {} }),
    delete: vi.fn().mockResolvedValue({ data: {} })
  }
}))

describe('EnterpriseHub.vue', () => {
  const mountOptions = {
    global: {
      stubs: {
        'q-page': { template: '<div class="q-page"><slot /></div>' },
        'q-card': { template: '<div class="q-card"><slot /></div>' },
        'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
        'q-tab': { template: '<div class="q-tab"><slot /></div>' },
        'q-tab-panels': { template: '<div class="q-tab-panels"><slot /></div>' },
        'q-tab-panel': { template: '<div class="q-tab-panel"><slot /></div>' },
        'q-select': { template: '<div class="q-select"><slot /></div>' },
        'q-btn': { template: '<button class="q-btn">{{ label }}<slot /></button>', props: ['label'] },
        'q-chip': { template: '<span class="q-chip"><slot /></span>' },
        'q-badge': { template: '<span class="q-badge">{{ label }}<slot /></span>', props: ['label'] },
        'q-markup-table': { template: '<table class="q-markup-table"><slot /></table>' },
        'q-slide-transition': { template: '<div class="q-slide-transition"><slot /></div>' },
        'q-toggle': { template: '<input type="checkbox" class="q-toggle" />' },
        'q-avatar': { template: '<div class="q-avatar"><slot /></div>' }
      }
    }
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('mounts properly and displays the enterprise header and dealbreaker tabs', () => {
    const wrapper = mount(EnterpriseHub, mountOptions)
    expect(wrapper.vm.activeTab).toBe('pagopa')
    expect(wrapper.text()).toContain('Enterprise School Management Hub')
    expect(wrapper.text()).toContain('Moduli commerciali strategici conformi agli standard normativi AgID, MIM, CAD ed eIDAS')
    expect(wrapper.text()).toContain('Emissione Avviso PagoPA')
    expect(wrapper.text()).toContain('Pubblica Avviso Interpello Supplenza')
    expect(wrapper.text()).toContain('Affissione Albo Pretorio Telematico')
    expect(wrapper.text()).toContain('Firma Remota Massiva CSC')
    expect(wrapper.text()).toContain('Calcolo Crediti Triennio')
    expect(wrapper.text()).toContain('Rilevazione Presenze Mensa')
    expect(wrapper.text()).toContain('Registro Beni Mobili & Cespiti')
    expect(wrapper.text()).toContain('Registro Trattamenti Privacy')
    expect(wrapper.text()).toContain('Certificato di Postazione WebService MIM')
    expect(wrapper.text()).toContain('Prenotazione Anonima & Riservata (CIC)')
  })

  it('displays authorized managing roles and regulatory reference in the governance bar', () => {
    const wrapper = mount(EnterpriseHub, mountOptions)
    const text = wrapper.text()
    expect(text).toContain('Ruoli Autorizzati alla Gestione:')
    expect(text).toContain('Dirigente Scolastico (DS)')
    expect(text).toContain('DSGA')
    expect(text).toContain('Segreteria Contabile')
    expect(text).toContain('Genitore / Famiglia')
    expect(text).toContain('D.Lgs. 217/2017 & AgID')
  })

  it('toggles the global institutional role governance matrix', async () => {
    const wrapper = mount(EnterpriseHub, mountOptions)
    expect(wrapper.vm.showRoleMatrix).toBe(false)
    expect(wrapper.text()).toContain('Matrice Competenze Ruoli')

    // Toggle matrix on
    wrapper.vm.showRoleMatrix = true
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Matrice di Competenza Ruoli Istituzionali')
    expect(wrapper.text()).toContain('Funzionalità Commerciale')
    expect(wrapper.text()).toContain('Gestione Principale')
    expect(wrapper.text()).toContain('Data Protection Officer (DPO)')
    expect(wrapper.text()).toContain('Psicologo Scolastico (L. 56/1989)')
    expect(wrapper.text()).toContain('Commissione Esame di Stato')
  })

  it('generates a PagoPA payment notice and triggers service call', async () => {
    vi.spyOn(paymentService, 'createPayment').mockResolvedValue({ data: { id: 'pay-123' } })
    vi.spyOn(paymentService, 'getBollettino').mockResolvedValue({
      data: { iuv: '01234567890123456', qr_code_payload: 'PAGOPA|QR|123' }
    })

    const wrapper = mount(EnterpriseHub, mountOptions)
    await wrapper.vm.generatePagoPA()

    expect(paymentService.createPayment).toHaveBeenCalledWith(expect.objectContaining({
      student_id: 'student-1',
      title: 'Contributo Volontario e Assicurazione',
      amount: 85.0
    }))
    expect(paymentService.getBollettino).toHaveBeenCalledWith('pay-123')
    expect(wrapper.vm.generatedNotice).toEqual({ iuv: '01234567890123456', qr_code_payload: 'PAGOPA|QR|123' })
  })

  it('executes OPI flow reconciliation with status report', async () => {
    vi.spyOn(paymentService, 'reconcileOPI').mockResolvedValue({
      data: { total_processed: 10, total_reconciled: 10, total_amount: 1540.50 }
    })

    const wrapper = mount(EnterpriseHub, mountOptions)
    wrapper.vm.opiInput = '<xml>test-opi</xml>'
    await wrapper.vm.reconcileOPI()

    expect(paymentService.reconcileOPI).toHaveBeenCalledWith('<xml>test-opi</xml>', 'OPI_XML')
    expect(wrapper.vm.reconciliationReport).toEqual({
      total_processed: 10,
      total_reconciled: 10,
      total_amount: 1540.50
    })
  })

  it('calculates Maturità three-year credits accurately', async () => {
    vi.spyOn(maturitaService, 'calculateCredits').mockResolvedValue({
      data: { total_credits: 38, credits_3rd: 11, credits_4th: 12, credits_5th: 15 }
    })

    const wrapper = mount(EnterpriseHub, mountOptions)
    await wrapper.vm.calcMaturitaCredits()

    expect(maturitaService.calculateCredits).toHaveBeenCalledWith(expect.objectContaining({
      grade_3rd: 8.4,
      grade_4th: 8.6,
      grade_5th: 8.8
    }))
    expect(wrapper.vm.calculatedCredits.total_credits).toBe(38)
  })

  it('signs batch with FEQ CSC remote signature and verifies digital glyph', async () => {
    vi.spyOn(digitalStampService, 'batchSignCSC').mockResolvedValue({ data: { signed_count: 2 } })
    vi.spyOn(digitalStampService, 'verifyPublicToken').mockResolvedValue({
      data: { valid: true, status: 'CONFORME CAD', message: 'Timbro crittografico valido' }
    })

    const wrapper = mount(EnterpriseHub, mountOptions)
    wrapper.vm.cscForm.pin = '1234'
    wrapper.vm.cscForm.otp = '567890'
    await wrapper.vm.signBatchCSC()
    expect(digitalStampService.batchSignCSC).toHaveBeenCalledWith(['pagella-1', 'pagella-2'], '1234', '567890')

    wrapper.vm.glifoToken = 'GLIFO-TEST-TOKEN'
    await wrapper.vm.verifyGlifo()
    expect(digitalStampService.verifyPublicToken).toHaveBeenCalledWith('GLIFO-TEST-TOKEN')
    expect(wrapper.vm.glifoResult.status).toBe('CONFORME CAD')
  })

  it('applies PagoPA quick presets and switches student via dropdown selector', async () => {
    const wrapper = mount(EnterpriseHub, mountOptions)
    
    // Apply preset
    wrapper.vm.applyPagoPaPreset("Quota Viaggio d'Istruzione / Gita", 120.0)
    expect(wrapper.vm.pagoPaForm.amount).toBe(120.0)
    expect(wrapper.vm.pagoPaForm.title).toContain("Quota Viaggio d'Istruzione / Gita")

    // Change student
    const student = wrapper.vm.allStudents[1]
    wrapper.vm.selectedPagoPaStudent = student
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.pagoPaForm.studentId).toBe(student.id)
  })

  it('issues bulk PagoPA payment notices to an entire class', async () => {
    vi.spyOn(paymentService, 'createBulkPayments').mockResolvedValue({
      data: { created_count: 2, batch_id: 'batch-pago-01' }
    })

    const wrapper = mount(EnterpriseHub, mountOptions)
    wrapper.vm.pagopaMode = 'class'
    wrapper.vm.selectedPagoPaClass = '1A'
    wrapper.vm.pagoPaForm.title = 'Assicurazione Integrativa'
    wrapper.vm.pagoPaForm.amount = 15.0

    await wrapper.vm.generateBulkPagoPA()

    const expectedStudentIds = wrapper.vm.studentsInPagoPaClass.map(s => s.id)
    expect(paymentService.createBulkPayments).toHaveBeenCalledWith(
      expectedStudentIds,
      expect.objectContaining({
        title: 'Assicurazione Integrativa',
        amount: 15.0,
        due_date: '2026-11-30'
      })
    )
  })

  it('filters graduating 5th-year students for Curriculum dello Studente', async () => {
    const wrapper = mount(EnterpriseHub, mountOptions)
    const graduating = wrapper.vm.graduatingStudents
    expect(graduating.length).toBeGreaterThan(0)
    graduating.forEach(st => {
      expect(st.class_name.startsWith('5')).toBe(true)
    })

    wrapper.vm.selectedCurriculumStudent = graduating[0]
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.curriculumStudentId).toBe(graduating[0].id)
  })

  it('manages Canteen wallet with quick top-up buttons and updates balance', async () => {
    vi.spyOn(mealsService, 'topUpWallet').mockResolvedValue({
      data: { success: true, new_balance: 75.0 }
    })

    const wrapper = mount(EnterpriseHub, mountOptions)
    const student = wrapper.vm.allStudents[0]
    const initialBalance = student.wallet_balance
    wrapper.vm.selectedMealStudent = student
    wrapper.vm.mealWallet.topupAmount = 25.0

    await wrapper.vm.topupMealWallet()

    expect(mealsService.topUpWallet).toHaveBeenCalledWith(student.id, 25.0, '01210000000014288')
    expect(student.wallet_balance).toBe(initialBalance + 25.0)
  })

  it('handles individual and class-wide bulk parental informed consent (CIC)', async () => {
    vi.spyOn(psychologyService, 'signConsent').mockResolvedValue({ data: { signed: true } })
    vi.spyOn(psychologyService, 'signBulkConsent').mockResolvedValue({ data: { signed_count: 2 } })

    const wrapper = mount(EnterpriseHub, mountOptions)

    // Single student mode
    const student = wrapper.vm.allStudents[2]
    wrapper.vm.cicMode = 'single'
    wrapper.vm.selectedCicStudent = student
    await wrapper.vm.signCICConsent()

    expect(psychologyService.signConsent).toHaveBeenCalledWith(student.id, '2025/2026')
    expect(student.has_cic_consent).toBe(true)

    // Entire class mode
    wrapper.vm.cicMode = 'class'
    wrapper.vm.selectedCicClass = '2A'
    const classStudents = wrapper.vm.studentsInCicClass
    expect(classStudents.length).toBeGreaterThan(0)

    await wrapper.vm.signBulkCICConsent()

    expect(psychologyService.signBulkConsent).toHaveBeenCalledWith(
      classStudents.map(s => s.id),
      '2025/2026'
    )
    classStudents.forEach(st => {
      expect(st.has_cic_consent).toBe(true)
    })
  })

  it('filters visible tabs correctly for Secretary role (hides FEQ, Albo, Interpelli, Privacy, CIC)', async () => {
    const wrapper = mount(EnterpriseHub, mountOptions)
    
    // Simulate secretary role
    wrapper.vm.simulatedRole = 'secretary'
    await wrapper.vm.$nextTick()

    const visible = wrapper.vm.visibleTabs
    expect(visible).toContain('pagopa')
    expect(visible).toContain('maturita')
    expect(visible).toContain('sidi')
    
    // Forbidden tabs for secretary
    expect(visible).not.toContain('feq')
    expect(visible).not.toContain('albo')
    expect(visible).not.toContain('interpelli')
    expect(visible).not.toContain('privacy')
    expect(visible).not.toContain('psychology')
    expect(visible).not.toContain('inventory')
  })

  it('filters visible tabs correctly for DPO and DSGA roles', async () => {
    const wrapper = mount(EnterpriseHub, mountOptions)

    // DPO
    wrapper.vm.simulatedRole = 'dpo'
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.visibleTabs).toEqual(['privacy'])

    // DSGA
    wrapper.vm.simulatedRole = 'dsga'
    await wrapper.vm.$nextTick()
    const dsgaVisible = wrapper.vm.visibleTabs
    expect(dsgaVisible).toContain('pagopa')
    expect(dsgaVisible).toContain('albo')
    expect(dsgaVisible).toContain('feq')
    expect(dsgaVisible).toContain('inventory')
    expect(dsgaVisible).toContain('sidi')
    expect(dsgaVisible).not.toContain('psychology')
    expect(dsgaVisible).not.toContain('interpelli')
  })

  it('hides canteen tab for secretary when deactivated, but keeps manageable for DSGA/Principal', async () => {
    const wrapper = mount(EnterpriseHub, mountOptions)

    // Deactivate canteen
    wrapper.vm.canteenActive = false
    await wrapper.vm.$nextTick()

    // Under secretary view, meals is completely hidden
    wrapper.vm.simulatedRole = 'secretary'
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.visibleTabs).not.toContain('meals')

    // Under DSGA view, meals is still accessible for management
    wrapper.vm.simulatedRole = 'dsga'
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.visibleTabs).toContain('meals')
    expect(wrapper.vm.canManageCanteen).toBe(true)
  })
})
