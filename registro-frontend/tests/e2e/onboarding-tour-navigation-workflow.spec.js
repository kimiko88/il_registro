import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
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

describe('E2E Workflow: Onboarding Tour Navigation & Direct Action Across Roles', () => {
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

  function mountTour(role) {
    const pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: {
          user: { role, first_name: 'Workflow', last_name: 'Tester' },
          userRole: role
        }
      }
    })

    return mount(OnboardingTour, {
      props: { modelValue: true },
      global: {
        plugins: [pinia, i18n],
        stubs: {
          'q-dialog': { template: '<div class="q-dialog-stub" v-if="modelValue"><slot /></div>', props: ['modelValue'] },
          'q-icon': true,
          'q-chip': { template: '<span class="q-chip" @click="$emit(\'click\')"><slot /></span>' },
          'q-btn': {
            template: '<button class="q-btn" @click="$emit(\'click\')"><slot />{{ label }}</button>',
            props: ['label', 'color', 'icon', 'disable']
          }
        }
      }
    })
  }

  it('1. Collaboratore Scolastico onboarding workflow: specialized steps and direct section navigation', async () => {
    const wrapper = mountTour('collaboratore_scolastico')
    wrapper.vm.startTour()
    await flushPromises()

    expect(wrapper.vm.isVisible).toBe(true)
    expect(wrapper.vm.userRole).toBe('collaboratore_scolastico')

    const steps = wrapper.vm.tourSteps
    expect(steps.length).toBeGreaterThanOrEqual(4)

    // Navigate to step 1 (Portineria)
    wrapper.vm.goToStep(1)
    await flushPromises()
    expect(wrapper.vm.currentStep).toBe(1)
    expect(wrapper.vm.activeStep.route).toBe('/ata')

    // Navigate to step 2 (Registro Visitatori)
    wrapper.vm.goToStep(2)
    await flushPromises()
    expect(wrapper.vm.activeStep.route).toBe('/ata/visitor-registry')

    // Click "Vai alla sezione"
    wrapper.vm.navigateTo(wrapper.vm.activeStep.route)
    await flushPromises()

    // Assert tour marked completed and router directed to /ata/visitor-registry
    expect(localStorage.getItem('onboarding_done_collaboratore_scolastico')).toBe('true')
    expect(wrapper.vm.isVisible).toBe(false)
    expect(mockPush).toHaveBeenCalledWith('/ata/visitor-registry')
  })

  it('2. Secretary onboarding workflow: step progression, jumping to step, and completing tour', async () => {
    const wrapper = mountTour('secretary')
    wrapper.vm.startTour()
    await flushPromises()

    expect(wrapper.vm.userRole).toBe('secretary')
    const totalSteps = wrapper.vm.tourSteps.length

    // Jump to step 3
    wrapper.vm.jumpToStep(3)
    await flushPromises()
    expect(wrapper.vm.currentStep).toBe(3)

    // Advance to end
    wrapper.vm.goToStep(totalSteps)
    await flushPromises()
    expect(wrapper.vm.currentStep).toBe(totalSteps)

    // Next step triggers completion screen
    wrapper.vm.nextStep()
    await flushPromises()
    expect(wrapper.vm.currentStep).toBe(999) // COMPLETION_STEP

    // Finish tour
    wrapper.vm.completeTour()
    await flushPromises()

    expect(localStorage.getItem('onboarding_done_secretary')).toBe('true')
    expect(wrapper.vm.isVisible).toBe(false)
    expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({ type: 'positive' }))
  })

  it('3. Teacher onboarding workflow: keyboard shortcuts navigation', async () => {
    const wrapper = mountTour('teacher')
    wrapper.vm.startTour()
    await flushPromises()

    expect(wrapper.vm.userRole).toBe('teacher')

    // Go to step 1
    wrapper.vm.goToStep(1)
    await flushPromises()

    // Press ArrowRight to move to step 2
    wrapper.vm.handleKeydown(new KeyboardEvent('keydown', { key: 'ArrowRight' }))
    expect(wrapper.vm.currentStep).toBe(2)

    // Press ArrowDown to move to step 3
    wrapper.vm.handleKeydown(new KeyboardEvent('keydown', { key: 'ArrowDown' }))
    expect(wrapper.vm.currentStep).toBe(3)

    // Press ArrowLeft to move back to step 2
    wrapper.vm.handleKeydown(new KeyboardEvent('keydown', { key: 'ArrowLeft' }))
    expect(wrapper.vm.currentStep).toBe(2)

    // Press Escape to dismiss tour
    wrapper.vm.handleKeydown(new KeyboardEvent('keydown', { key: 'Escape' }))
    expect(wrapper.vm.isVisible).toBe(false)
    expect(localStorage.getItem('onboarding_done_teacher')).toBe('true')
  })
})
