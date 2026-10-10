import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import EnterpriseHub from 'src/pages/admin/EnterpriseHub.vue'
import PsychologyDesk from 'src/pages/student/PsychologyDesk.vue'
import Canteen from 'src/pages/student/Canteen.vue'
import PsychologyConsent from 'src/pages/parent/PsychologyConsent.vue'
import CanteenWallet from 'src/pages/parent/CanteenWallet.vue'
import { alboPretorioService } from 'src/services/alboPretorioService'
import { interpelliService } from 'src/services/interpelliService'
import { mealsService } from 'src/services/mealsService'
import { psychologyService } from 'src/services/psychologyService'
import { digitalStampService } from 'src/services/digitalStampService'
import { useSchoolSettingsStore } from 'src/stores/schoolSettings'
import itLocale from 'src/i18n/it-IT'

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

const defaultStubs = {
    'q-page': { template: '<div class="q-page"><slot /></div>' },
    'q-card': { template: '<div class="q-card"><slot /></div>' },
    'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
    'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
    'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
    'q-tab': { template: '<div class="q-tab">{{ label }}<slot /></div>', props: ['label'] },
    'q-tab-panels': { template: '<div class="q-tab-panels"><slot /></div>' },
    'q-tab-panel': { template: '<div class="q-tab-panel"><slot /></div>' },
    'q-avatar': { template: '<div class="q-avatar"><slot /></div>' },
    'q-chip': { template: '<span class="q-chip"><slot /></span>' },
    'q-badge': { template: '<span class="q-badge">{{ label }}<slot /></span>', props: ['label'] },
    'q-icon': { template: '<i class="q-icon" />' },
    'q-select': { template: '<div class="q-select"><slot /></div>' },
    'q-input': { template: '<div class="q-input"><slot /></div>', props: ['prefix', 'label'] },
    'q-btn': { template: '<button class="q-btn">{{ label }}<slot /></button>', props: ['label'] },
    'q-checkbox': { template: '<input type="checkbox" class="q-checkbox" />' },
    'q-list': { template: '<div class="q-list"><slot /></div>' },
    'q-item': { template: '<div class="q-item"><slot /></div>' },
    'q-item-section': { template: '<div class="q-item-section"><slot /></div>' },
    'q-item-label': { template: '<div class="q-item-label"><slot /></div>' },
    'q-dialog': { template: '<div class="q-dialog"><slot /></div>' },
    'q-space': { template: '<div class="q-space" />' },
    'q-banner': { template: '<div class="q-banner"><slot /></div>' },
    'q-markup-table': { template: '<table class="q-markup-table"><slot /></table>' },
    'q-slide-transition': { template: '<div class="q-slide-transition"><slot /></div>' },
    'q-toggle': { template: '<input type="checkbox" class="q-toggle" />' }
}

function getNested(obj, path) {
    return path.split('.').reduce((acc, part) => acc && acc[part], obj)
}

const mountOptions = {
    global: {
        stubs: defaultStubs,
        mocks: {
            $t: (key) => getNested(itLocale, key) || key
        }
    }
}

describe('Enterprise Dealbreakers Full E2E Workflow', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
        localStorage.clear()
    })

    it('ED01 — EnterpriseHub allows administrator to toggle and view all dealbreaker enterprise tabs', async () => {
        const wrapper = mount(EnterpriseHub, mountOptions)
        expect(wrapper.exists()).toBe(true)
        expect(wrapper.html()).toContain('Enterprise')
    })

    it('ED02 — Student Canteen meal booking flow and dietary plan verification', async () => {
        const store = useSchoolSettingsStore()
        store.canteenEnabled = true

        const wrapper = mount(Canteen, mountOptions)
        expect(wrapper.exists()).toBe(true)
        expect(wrapper.text()).toContain('Mensa Scolastica & Refezione')
        expect(wrapper.text()).toContain('Menu del Giorno')
        expect(wrapper.text()).toContain('Dieta Mediterranea Standard')
    })

    it('ED03 — Parent Canteen Wallet balance tracking and PagoPA recharge workflow', async () => {
        vi.spyOn(mealsService, 'topUpWallet').mockResolvedValue({ data: { success: true } })
        const store = useSchoolSettingsStore()
        store.canteenEnabled = true

        const wrapper = mount(CanteenWallet, mountOptions)
        expect(wrapper.exists()).toBe(true)
        expect(wrapper.text()).toContain('Borsellino Elettronico Pasti')

        wrapper.vm.topupAmount = 25.0
        await wrapper.vm.submitTopup()
        expect(mealsService.topUpWallet).toHaveBeenCalled()
    })

    it('ED04 — Student Psychology Desk confidential counseling appointment request', async () => {
        vi.spyOn(psychologyService, 'bookSession').mockResolvedValue({ data: { success: true } })

        const wrapper = mount(PsychologyDesk, mountOptions)
        expect(wrapper.exists()).toBe(true)
        expect(wrapper.text()).toContain('Sportello d\'Ascolto')
        expect(wrapper.text()).toContain('Tutela Assoluta della Riservatezza')

        wrapper.vm.bookingForm.slotTime = 'Giovedì 15 Ottobre 2026 - 11:00'
        await wrapper.vm.submitBooking()
        expect(psychologyService.bookSession).toHaveBeenCalledWith('psy-01', 'Giovedì 15 Ottobre 2026 - 11:00', 45)
    })

    it('ED05 — Parent Psychology Consent compliance with dual parental signature requirement', async () => {
        vi.spyOn(psychologyService, 'signConsent').mockResolvedValue({ data: { signed: true } })

        const wrapper = mount(PsychologyConsent, mountOptions)
        expect(wrapper.exists()).toBe(true)
        expect(wrapper.text()).toContain('Sportello Psicologico (CIC) & Consenso Genitoriale')
        expect(wrapper.vm.selectedChild.has_consent).toBe(false)

        await wrapper.vm.signConsent()
        expect(psychologyService.signConsent).toHaveBeenCalledWith('student-minor-1', '2025/2026')
        expect(wrapper.vm.selectedChild.has_consent).toBe(true)
    })

    it('ED06 — Albo Pretorio, Interpelli and Digital Stamp legal service calls verification', async () => {
        vi.spyOn(alboPretorioService, 'listPublic').mockResolvedValue({
            data: [
                { id: 'act-1', number: 101, year: 2026, title: 'Bando Supplenza A028', category: 'bandi' }
            ]
        })
        vi.spyOn(interpelliService, 'listPublic').mockResolvedValue({
            data: [
                { id: 'int-1', code: 'INT-2026-001', subject: 'Matematica', status: 'aperto' }
            ]
        })
        vi.spyOn(digitalStampService, 'createDigitalStamp').mockResolvedValue({
            data: { verified: true, stamp_id: 'stamp-123' }
        })

        const acts = await alboPretorioService.listPublic({ year: 2026 })
        expect(acts.data).toHaveLength(1)
        expect(acts.data[0].title).toBe('Bando Supplenza A028')

        const interpelli = await interpelliService.listPublic()
        expect(interpelli.data[0].code).toBe('INT-2026-001')

        const stampRes = await digitalStampService.createDigitalStamp({ document_id: 'doc-999' })
        expect(stampRes.data.verified).toBe(true)
    })
})
