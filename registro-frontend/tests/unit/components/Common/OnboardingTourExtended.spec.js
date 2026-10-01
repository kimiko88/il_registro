import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { createI18n } from 'vue-i18n'
import OnboardingTour from '@/components/Common/OnboardingTour.vue'
import itLocale from '@/i18n/it-IT/index.js'

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

describe('OnboardingTourExtended.spec.js — Interactive Tour Comprehensive Suite', () => {
  let i18n

  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()

    i18n = createI18n({
      legacy: false,
      locale: 'it-IT',
      messages: {
        'it-IT': itLocale
      }
    })
  })

  function createWrapper(userRole = 'student') {
    const pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: {
          user: { role: userRole, first_name: 'Test', last_name: 'User' },
          userRole
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

  describe('Route Integrity Across All 10 Roles', () => {
    const canonicalRoles = [
      'principal',
      'teacher',
      'student',
      'parent',
      'secretary',
      'admin',
      'assistente_amministrativo',
      'collaboratore_ds',
      'collaboratore_scolastico',
      'dsga'
    ]

    canonicalRoles.forEach(role => {
      it(`role "${role}" has valid routes defined for every single step`, () => {
        const wrapper = createWrapper(role)
        const steps = wrapper.vm.tourSteps

        expect(steps.length).toBeGreaterThanOrEqual(4)
        steps.forEach((step, index) => {
          expect(step.route, `Step ${index + 1} of ${role} missing route`).toBeDefined()
          expect(step.route.startsWith('/'), `Step ${index + 1} route must start with '/'`).toBe(true)
          expect(step.icon).toBeTruthy()
          expect(step.color).toBeTruthy()
          expect(step.title).toBeTruthy()
        })
      })
    })
  })

  describe('Direct Section Navigation (navigateTo)', () => {
    it('navigateTo sets completion localStorage, closes modal, and pushes route to router', () => {
      const wrapper = createWrapper('collaboratore_scolastico')
      wrapper.vm.startTour()
      expect(wrapper.vm.isVisible).toBe(true)

      wrapper.vm.navigateTo('/ata/visitor-registry')

      expect(localStorage.getItem('onboarding_done_collaboratore_scolastico')).toBe('true')
      expect(wrapper.vm.isVisible).toBe(false)
      expect(mockPush).toHaveBeenCalledWith('/ata/visitor-registry')
      expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({ type: 'positive' }))
    })

    it('navigateTo with empty route only completes tour without navigating', () => {
      const wrapper = createWrapper('student')
      wrapper.vm.navigateTo('')

      expect(localStorage.getItem('onboarding_done_student')).toBe('true')
      expect(mockPush).not.toHaveBeenCalled()
    })
  })

  describe('Step Navigation, Progress Calculation, and Boundary Handling', () => {
    it('advances steps sequentially with nextStep and caps at COMPLETION_STEP', () => {
      const wrapper = createWrapper('teacher')
      const total = wrapper.vm.tourSteps.length

      wrapper.vm.goToStep(1)
      expect(wrapper.vm.currentStep).toBe(1)
      expect(wrapper.vm.progressPercent).toBeCloseTo((1 / total) * 100, 1)

      wrapper.vm.nextStep()
      expect(wrapper.vm.currentStep).toBe(2)

      // Jump to last step
      wrapper.vm.goToStep(total)
      wrapper.vm.nextStep()
      expect(wrapper.vm.currentStep).toBe(999) // COMPLETION_STEP
    })

    it('decrements steps with prevStep and does not go below step 1', () => {
      const wrapper = createWrapper('teacher')
      wrapper.vm.goToStep(2)
      wrapper.vm.prevStep()
      expect(wrapper.vm.currentStep).toBe(1)

      // Should not go below step 1
      wrapper.vm.prevStep()
      expect(wrapper.vm.currentStep).toBe(1)
    })
  })

  describe('Keyboard Navigation', () => {
    it('handles ArrowRight by invoking nextStep', () => {
      const wrapper = createWrapper('student')
      wrapper.vm.goToStep(1)
      const event = new KeyboardEvent('keydown', { key: 'ArrowRight' })
      wrapper.vm.handleKeydown(event)
      expect(wrapper.vm.currentStep).toBe(2)
    })

    it('handles ArrowLeft by invoking prevStep', () => {
      const wrapper = createWrapper('student')
      wrapper.vm.goToStep(3)
      const event = new KeyboardEvent('keydown', { key: 'ArrowLeft' })
      wrapper.vm.handleKeydown(event)
      expect(wrapper.vm.currentStep).toBe(2)
    })

    it('handles Escape by calling completeTour', () => {
      const wrapper = createWrapper('parent')
      wrapper.vm.startTour()
      const event = new KeyboardEvent('keydown', { key: 'Escape' })
      wrapper.vm.handleKeydown(event)
      expect(wrapper.vm.isVisible).toBe(false)
      expect(localStorage.getItem('onboarding_done_parent')).toBe('true')
    })
  })

  describe('Role Canonical Resolution Fallbacks', () => {
    it('falls back unknown roles to "student"', () => {
      const wrapper = createWrapper('unknown_external_guest')
      expect(wrapper.vm.userRole).toBe('student')
    })

    it('resolves case-insensitive trimmed roles', () => {
      const wrapper = createWrapper('  DIRIGENTE_SCOLASTICO  ')
      expect(wrapper.vm.userRole).toBe('principal')
    })
  })
})
