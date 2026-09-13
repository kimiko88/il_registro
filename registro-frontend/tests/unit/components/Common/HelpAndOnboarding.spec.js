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

    it('resolves principal role and loads 6 dedicated tour steps with valid titles', () => {
      const principalPinia = createTestingPinia({
        createSpy: vi.fn,
        initialState: {
          auth: {
            user: { role: 'principal', first_name: 'Giulia', last_name: 'Verdi' },
            userRole: 'principal',
            token: 'jwt-token'
          }
        }
      })

      const wrapper = mount(OnboardingTour, {
        global: {
          plugins: [principalPinia, i18n],
          stubs: {
            'q-dialog': { template: '<div class="q-dialog-stub" v-if="modelValue"><slot /></div>', props: ['modelValue'] },
            'q-icon': true,
            'q-chip': true,
            'q-btn': true
          }
        }
      })

      expect(wrapper.vm.userRole).toBe('principal')
      expect(wrapper.vm.tourSteps).toHaveLength(6)
      expect(wrapper.vm.tourSteps[0].title).toBe('Dashboard Direzione & Quadro Generale')
      expect(wrapper.vm.tourSteps[1].title).toBe('Personale, Nomine & Incarichi')
      expect(wrapper.vm.tourSteps[2].title).toBe('Decreti & Visti Personale')
      expect(wrapper.vm.tourSteps[3].title).toBe('Sostituzioni Docenti & Gestione Emergenze')
      expect(wrapper.vm.tourSteps[4].title).toBe('Atti, Verbali & Delibere Collegiali')
      expect(wrapper.vm.tourSteps[5].title).toBe('Monitoraggio Didattico, Scrutini & Dispersione')
    })

    it('correctly maps various institutional roles to canonical roles', () => {
      const rolesToTest = [
        { role: 'vice_principal', expected: 'principal' },
        { role: 'coordinator', expected: 'teacher' },
        { role: 'responsabile_servizio', expected: 'collaboratore_ds' },
        { role: 'assistente_alunni', expected: 'assistente_amministrativo' },
        { role: 'system_auditor', expected: 'admin' },
        { role: 'responsabile_conservazione', expected: 'secretary' }
      ]

      for (const { role, expected } of rolesToTest) {
        const testPinia = createTestingPinia({
          createSpy: vi.fn,
          initialState: {
            auth: {
              user: { role, first_name: 'Test', last_name: 'User' },
              userRole: role,
              token: 'jwt-token'
            }
          }
        })

        const wrapper = mount(OnboardingTour, {
          global: {
            plugins: [testPinia, i18n],
            stubs: {
              'q-dialog': true,
              'q-icon': true,
              'q-chip': true,
              'q-btn': true
            }
          }
        })

        expect(wrapper.vm.userRole).toBe(expected)
      }
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

    it('loads 5 categories and dedicated articles for principal', () => {
      const principalPinia = createTestingPinia({
        createSpy: vi.fn,
        initialState: {
          auth: {
            user: { role: 'principal', first_name: 'Giulia', last_name: 'Verdi' },
            userRole: 'principal',
            token: 'jwt-token'
          }
        }
      })

      const wrapper = mount(HelpDrawer, {
        global: {
          plugins: [principalPinia, i18n],
          stubs: {
            'q-drawer': { template: '<div class="q-drawer-stub" v-if="modelValue"><slot /></div>', props: ['modelValue'] },
            'q-icon': true,
            'q-btn': true,
            'q-input': true,
            'q-expansion-item': true
          }
        }
      })

      expect(wrapper.vm.userRole).toBe('principal')
      expect(wrapper.vm.categories).toHaveLength(5)
      expect(wrapper.vm.categories.map(c => c.key)).toEqual([
        'cat_direction',
        'cat_personnel',
        'cat_substitutions',
        'cat_verbali',
        'cat_strike'
      ])
      expect(wrapper.vm.articles.length).toBeGreaterThanOrEqual(6)
      expect(wrapper.vm.articles[0].question).toContain('quadro generale delle assenze')
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

    it('loads 5 categories and guides for principal', () => {
      const principalPinia = createTestingPinia({
        createSpy: vi.fn,
        initialState: {
          auth: {
            user: { role: 'principal', first_name: 'Giulia', last_name: 'Verdi' },
            userRole: 'principal',
            token: 'jwt-token'
          }
        }
      })

      const wrapper = mount(HelpCenterPanel, {
        global: {
          plugins: [principalPinia, i18n],
          stubs: {
            'q-dialog': { template: '<div class="q-dialog-stub" v-if="modelValue"><slot /></div>', props: ['modelValue'] },
            'q-icon': true,
            'q-btn': true,
            'q-input': true,
            'q-expansion-item': true
          }
        }
      })

      expect(wrapper.vm.userRole).toBe('principal')
      expect(wrapper.vm.categories).toHaveLength(5)
      expect(wrapper.vm.categories.map(c => c.key)).toEqual([
        'dashboard',
        'personnel',
        'substitutions',
        'verbali',
        'strike'
      ])
    })
  })
})
