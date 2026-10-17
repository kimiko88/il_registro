import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { createI18n } from 'vue-i18n'
import OnboardingTour from '@/components/Common/OnboardingTour.vue'
import itLocale from '@/i18n/it-IT/index.js'
import enLocale from '@/i18n/en-US/index.js'

const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush
  })
}))

const mockNotify = vi.fn()
vi.mock('quasar', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useQuasar: () => ({
      notify: mockNotify,
      screen: { lt: { sm: false } }
    })
  }
})

describe('EnterpriseRoleTours.spec.js — New Features Tour Coverage', () => {
  let i18n

  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()

    i18n = createI18n({
      legacy: false,
      locale: 'it-IT',
      fallbackLocale: 'it-IT',
      messages: {
        'it-IT': itLocale,
        'it': itLocale,
        'en-US': enLocale,
        'en': enLocale
      }
    })
  })

  function createTour(role) {
    const pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: {
          user: { role, first_name: 'Test', last_name: 'User' },
          userRole: role
        }
      }
    })

    return mount(OnboardingTour, {
      props: { modelValue: true },
      global: {
        plugins: [pinia, i18n],
        stubs: {
          'q-dialog': { template: '<div class="q-dialog-stub"><slot /></div>' },
          'q-icon': true,
          'q-chip': true,
          'q-btn': true
        }
      }
    })
  }

  it('Principal tour includes Enterprise Hub with remote FEQ signature and notice board', () => {
    const wrapper = createTour('principal')
    const steps = wrapper.vm.tourSteps
    expect(steps).toHaveLength(7)

    const enterpriseStep = steps[6]
    expect(enterpriseStep.route).toBe('/admin/enterprise')
    expect(enterpriseStep.title).toBe('Enterprise Hub & Firma Remota FEQ')
    expect(enterpriseStep.bullets.length).toBeGreaterThanOrEqual(3)
  })

  it('Student tour includes Psychology Desk and Canteen meals', () => {
    const wrapper = createTour('student')
    const steps = wrapper.vm.tourSteps
    expect(steps).toHaveLength(10)

    const psyStep = steps[8]
    expect(psyStep.route).toBe('/student/psychology')
    expect(psyStep.title).toContain('Sportello d\'Ascolto')
    expect(psyStep.bullets.length).toBeGreaterThanOrEqual(3)

    const canteenStep = steps[9]
    expect(canteenStep.route).toBe('/student/canteen')
    expect(canteenStep.title).toContain('Mensa Scolastica')
    expect(canteenStep.bullets.length).toBeGreaterThanOrEqual(3)
  })

  it('Parent tour includes Psychology Parental Consent and Canteen Electronic Wallet', () => {
    const wrapper = createTour('parent')
    const steps = wrapper.vm.tourSteps
    expect(steps).toHaveLength(10)

    const psyConsentStep = steps[8]
    expect(psyConsentStep.route).toBe('/parent/psychology')
    expect(psyConsentStep.title).toContain('Sportello Psicologico')
    expect(psyConsentStep.bullets.length).toBeGreaterThanOrEqual(3)

    const canteenWalletStep = steps[9]
    expect(canteenWalletStep.route).toBe('/parent/canteen')
    expect(canteenWalletStep.title).toContain('Mensa Scolastica & Borsellino')
    expect(canteenWalletStep.bullets.length).toBeGreaterThanOrEqual(3)
  })

  it('Secretary tour includes Enterprise Hub with PagoPA and Graduation/SIDI', () => {
    const wrapper = createTour('secretary')
    const steps = wrapper.vm.tourSteps
    expect(steps).toHaveLength(9)

    const enterpriseStep = steps[8]
    expect(enterpriseStep.route).toBe('/admin/enterprise')
    expect(enterpriseStep.title).toContain('Enterprise Hub')
    expect(enterpriseStep.bullets.length).toBeGreaterThanOrEqual(3)
  })

  it('DSGA tour includes Enterprise Hub with notice board, inventory, and OPI reconciliation', () => {
    const wrapper = createTour('dsga')
    const steps = wrapper.vm.tourSteps
    expect(steps).toHaveLength(6)

    const enterpriseStep = steps[5]
    expect(enterpriseStep.route).toBe('/admin/enterprise')
    expect(enterpriseStep.title).toContain('Enterprise Hub')
    expect(enterpriseStep.bullets.length).toBeGreaterThanOrEqual(3)
  })

  it('Assistente Amministrativo tour includes Administrative Enterprise Hub', () => {
    const wrapper = createTour('assistente_amministrativo')
    const steps = wrapper.vm.tourSteps
    expect(steps).toHaveLength(6)

    const enterpriseStep = steps[5]
    expect(enterpriseStep.route).toBe('/admin/enterprise')
    expect(enterpriseStep.title).toContain('Enterprise Hub Amministrativo')
    expect(enterpriseStep.bullets.length).toBeGreaterThanOrEqual(3)
  })

  it('Admin tour includes Unified Enterprise School Management Hub with role simulator', () => {
    const wrapper = createTour('admin')
    const steps = wrapper.vm.tourSteps
    expect(steps).toHaveLength(9)

    const enterpriseStep = steps[8]
    expect(enterpriseStep.route).toBe('/admin/enterprise')
    expect(enterpriseStep.title).toBe('Enterprise School Management Hub')
    expect(enterpriseStep.bullets.length).toBeGreaterThanOrEqual(3)
  })
})
