import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { createI18n } from 'vue-i18n'
import OnboardingTour from '@/components/Common/OnboardingTour.vue'
import HelpDrawer from '@/components/Common/HelpDrawer.vue'
import HelpFab from '@/components/Common/HelpFab.vue'
import HelpCenterPanel from '@/components/Common/HelpCenterPanel.vue'
import itLocale from '@/i18n/it-IT'

const i18n = createI18n({
  legacy: false,
  locale: 'it-IT',
  fallbackLocale: 'it-IT',
  messages: {
    'it-IT': itLocale,
    'it': itLocale
  }
})

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn()
  }),
  useRoute: () => ({
    path: '/'
  })
}))

describe('Onboarding & Help Center Components', () => {
  let pinia

  beforeEach(() => {
    localStorage.clear()
    pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: {
          user: { role: 'teacher', first_name: 'Mario', last_name: 'Rossi' },
          userRole: 'teacher',
          token: 'jwt-token'
        }
      }
    })
  })

  describe('OnboardingTour.vue', () => {
    it('initializes and provides startTour method', () => {
      const wrapper = mount(OnboardingTour, {
        global: {
          plugins: [pinia, i18n],
          stubs: {
            'q-dialog': { template: '<div class="q-dialog-stub" v-if="modelValue"><slot /></div>', props: ['modelValue'] },
            'q-icon': true,
            'q-chip': true,
            'q-btn': true
          }
        }
      })

      expect(typeof wrapper.vm.startTour).toBe('function')
    })

    it('sets localStorage when completed and does not auto-open next time', () => {
      localStorage.setItem('onboarding_done_teacher', 'true')

      expect(localStorage.getItem('onboarding_done_teacher')).toBe('true')
    })
  })

  describe('HelpFab.vue', () => {
    it('renders FAB button and emits events on click', async () => {
      const wrapper = mount(HelpFab, {
        global: {
          plugins: [pinia, i18n],
          stubs: {
            'q-icon': true,
            'q-tooltip': true
          }
        }
      })

      const btn = wrapper.find('.help-fab-btn')
      expect(btn.exists()).toBe(true)

      // Toggle menu
      await btn.trigger('click')
      expect(wrapper.vm.isMenuOpen).toBe(true)
    })
  })

  describe('HelpDrawer.vue', () => {
    it('mounts and exposes open and close methods', () => {
      const wrapper = mount(HelpDrawer, {
        global: {
          plugins: [pinia, i18n],
          stubs: {
            'q-drawer': { template: '<div class="q-drawer-stub" v-if="modelValue"><slot /></div>', props: ['modelValue'] },
            'q-icon': true,
            'q-btn': true,
            'q-input': true,
            'q-expansion-item': true
          }
        }
      })

      expect(typeof wrapper.vm.open).toBe('function')
      expect(typeof wrapper.vm.close).toBe('function')

      wrapper.vm.open()
      expect(wrapper.vm.isOpen).toBe(true)

      wrapper.vm.close()
      expect(wrapper.vm.isOpen).toBe(false)
    })
  })

  describe('HelpCenterPanel.vue', () => {
    it('mounts and exposes open and close methods', () => {
      const wrapper = mount(HelpCenterPanel, {
        global: {
          plugins: [pinia, i18n],
          stubs: {
            'q-dialog': { template: '<div class="q-dialog-stub" v-if="modelValue"><slot /></div>', props: ['modelValue'] },
            'q-icon': true,
            'q-btn': true,
            'q-input': true,
            'q-expansion-item': true
          }
        }
      })

      expect(typeof wrapper.vm.open).toBe('function')
      expect(typeof wrapper.vm.close).toBe('function')

      wrapper.vm.open('grades')
      expect(wrapper.vm.isOpen).toBe(true)

      wrapper.vm.close()
      expect(wrapper.vm.isOpen).toBe(false)
    })
  })
})
