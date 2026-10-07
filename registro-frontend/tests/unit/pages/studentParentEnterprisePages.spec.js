import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import PsychologyDesk from 'src/pages/student/PsychologyDesk.vue'
import Canteen from 'src/pages/student/Canteen.vue'
import PsychologyConsent from 'src/pages/parent/PsychologyConsent.vue'
import CanteenWallet from 'src/pages/parent/CanteenWallet.vue'
import { psychologyService } from 'src/services/psychologyService'
import { mealsService } from 'src/services/mealsService'
import { useSchoolSettingsStore } from 'src/stores/schoolSettings'
import { createPinia, setActivePinia } from 'pinia'
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
  'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
  'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
  'q-space': { template: '<div class="q-space" />' }
}

function getNested(obj, path) {
  return path.split('.').reduce((acc, part) => acc && acc[part], obj)
}

const defaultOptions = {
  global: {
    stubs: defaultStubs,
    mocks: {
      $t: (key) => getNested(itLocale, key) || key
    }
  }
}

describe('Student & Parent Psychology and Canteen Pages', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorage.clear()
  })

  describe('Student: PsychologyDesk.vue', () => {
    it('mounts and displays confidentiality information and anonymous booking form', async () => {
      vi.spyOn(psychologyService, 'bookSession').mockResolvedValue({ data: { success: true } })

      const wrapper = mount(PsychologyDesk, defaultOptions)

      expect(wrapper.text()).toContain('Sportello d\'Ascolto')
      expect(wrapper.text()).toContain('Tutela Assoluta della Riservatezza')
      expect(wrapper.vm.hasConsent).toBe(true)

      // Submit booking
      wrapper.vm.bookingForm.slotTime = 'Venerdì 16 Ottobre 2026 - 10:30'
      await wrapper.vm.submitBooking()

      expect(psychologyService.bookSession).toHaveBeenCalledWith('psy-01', 'Venerdì 16 Ottobre 2026 - 10:30', 45)
      expect(wrapper.vm.myBookings.length).toBeGreaterThan(1)
    })
  })

  describe('Student: Canteen.vue', () => {
    it('displays inactive warning when canteen is disabled in school settings', () => {
      const store = useSchoolSettingsStore()
      store.canteenEnabled = false

      const wrapper = mount(Canteen, defaultOptions)

      expect(wrapper.text()).toContain('Il Servizio Mensa non è attivo per questo istituto')
    })

    it('displays today menu and dietary plan when canteen is enabled', () => {
      const store = useSchoolSettingsStore()
      store.canteenEnabled = true

      const wrapper = mount(Canteen, defaultOptions)

      expect(wrapper.text()).toContain('Mensa Scolastica & Refezione')
      expect(wrapper.text()).toContain('Menu del Giorno')
      expect(wrapper.text()).toContain('Dieta Mediterranea Standard')
    })
  })

  describe('Parent: PsychologyConsent.vue', () => {
    it('allows parent to sign mandatory informed consent for child', async () => {
      vi.spyOn(psychologyService, 'signConsent').mockResolvedValue({ data: { signed: true } })

      const wrapper = mount(PsychologyConsent, defaultOptions)

      expect(wrapper.text()).toContain('Sportello Psicologico (CIC) & Consenso Genitoriale')
      expect(wrapper.vm.selectedChild.has_consent).toBe(false)

      await wrapper.vm.signConsent()

      expect(psychologyService.signConsent).toHaveBeenCalledWith('student-minor-1', '2025/2026')
      expect(wrapper.vm.selectedChild.has_consent).toBe(true)
    })
  })

  describe('Parent: CanteenWallet.vue', () => {
    it('displays current wallet balance and processes PagoPA topup', async () => {
      vi.spyOn(mealsService, 'topUpWallet').mockResolvedValue({ data: { success: true } })

      const store = useSchoolSettingsStore()
      store.canteenEnabled = true

      const wrapper = mount(CanteenWallet, defaultOptions)

      expect(wrapper.text()).toContain('Borsellino Elettronico Pasti')
      const initialBal = wrapper.vm.selectedChild.balance

      wrapper.vm.topupAmount = 25.0
      await wrapper.vm.submitTopup()

      expect(mealsService.topUpWallet).toHaveBeenCalledWith(
        'student-minor-1',
        25.0,
        expect.stringMatching(/^0121/)
      )
      expect(wrapper.vm.selectedChild.balance).toBe(initialBal + 25.0)
    })
  })
})
