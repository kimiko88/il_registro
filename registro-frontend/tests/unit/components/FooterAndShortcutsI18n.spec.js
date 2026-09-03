import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

const { currentLocale } = vi.hoisted(() => {
  const { ref } = require('vue')
  return { currentLocale: ref('it-IT') }
})

vi.mock('vue-i18n', async (importOriginal) => {
  const actual = await importOriginal().catch(() => ({}))
  const itMessages = (await import('@/i18n/it-IT/index.js')).default
  const enMessages = (await import('@/i18n/en-US/index.js')).default
  function getNestedValue(obj, path) {
    if (!obj || !path) return null
    return path.split('.').reduce((prev, curr) => (prev && prev[curr] !== undefined ? prev[curr] : null), obj)
  }
  const translate = (key, params) => {
    const messages = currentLocale.value === 'en-US' ? enMessages : itMessages
    let val = getNestedValue(messages, key) || key
    if (typeof val === 'string' && params && typeof params === 'object') {
      Object.keys(params).forEach(k => {
        val = val.replace(new RegExp(`\\{${k}\\}`, 'g'), params[k])
      })
    }
    return val
  }
  return {
    ...actual,
    useI18n: () => ({
      locale: currentLocale,
      t: translate,
      te: (key) => !!getNestedValue(currentLocale.value === 'en-US' ? enMessages : itMessages, key)
    }),
    createI18n: () => ({
      install: () => {},
      global: {
        locale: currentLocale,
        t: translate
      }
    })
  }
})

vi.mock('vue-router', () => ({
  useRoute: () => ({ path: '/teacher/dashboard', meta: {} }),
  useRouter: () => ({ push: vi.fn() })
}))

import KeyboardShortcutsDialog from '@/components/Common/KeyboardShortcutsDialog.vue'
import AccessibilityStatement from '@/pages/AccessibilityStatement.vue'
import MainLayout from '@/layouts/MainLayout.vue'
import { useThemeStore } from '@/stores/theme'
import { useAuthStore } from '@/stores/auth'

describe('Footer, KeyboardShortcutsDialog, and AccessibilityStatement i18n & Reactivity', () => {
  let wrapper

  beforeEach(() => {
    setActivePinia(createPinia())
    currentLocale.value = 'it-IT'
  })

  afterEach(() => {
    if (wrapper) wrapper.unmount()
  })

  it('KeyboardShortcutsDialog renders all items translated in Italian and updates to English', async () => {
    const themeStore = useThemeStore()
    themeStore.keyboardShortcutsHelpOpen = true

    wrapper = mount(KeyboardShortcutsDialog, {
      global: {
        stubs: {
          'q-dialog': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-input': { template: '<input />' },
          'q-list': { template: '<div><slot /></div>' },
          'q-item': { template: '<div><slot /></div>' },
          'q-item-section': { template: '<div><slot /></div>' },
          'q-item-label': { template: '<div><slot /></div>' },
          'q-icon': true,
          'q-btn': true
        }
      }
    })

    // Verify initial Italian labels in shortcutGroups computed
    expect(wrapper.vm.shortcutGroups[0].category).toContain('Navigazione Globale')
    expect(wrapper.vm.shortcutGroups[0].items[0].description).toBe('Vai alla Dashboard')
    expect(wrapper.vm.shortcutGroups[1].category).toContain('Funzionalità di Accessibilità')
    expect(wrapper.vm.shortcutGroups[1].items[0].description).toBe('Mostra / Nascondi Guida Scorciatoie')

    // Switch locale to English
    currentLocale.value = 'en-US'
    await wrapper.vm.$nextTick()

    // Verify dynamically recomputed English labels
    expect(wrapper.vm.shortcutGroups[0].category).toContain('Global Navigation')
    expect(wrapper.vm.shortcutGroups[0].items[0].description).toBe('Go to Dashboard')
    expect(wrapper.vm.shortcutGroups[1].category).toContain('Accessibility Features')
    expect(wrapper.vm.shortcutGroups[1].items[0].description).toBe('Show / Hide Shortcuts Guide')
  })

  it('MainLayout footer contains translated copyright, a11y statement, and shortcuts in Italian and English', async () => {
    const authStore = useAuthStore()
    authStore.user = { id: 1, name: 'Mario Rossi', role: 'teacher' }

    wrapper = mount(MainLayout, {
      global: {
        stubs: {
          'q-layout': { template: '<div class="q-layout"><slot /></div>' },
          'q-header': { template: '<header><slot /></header>' },
          'q-toolbar': { template: '<nav><slot /></nav>' },
          'q-toolbar-title': { template: '<div><slot /></div>' },
          'q-btn': { template: '<button><slot /></button>' },
          'q-btn-dropdown': { template: '<div class="btn-dropdown"><slot /></div>' },
          'q-avatar': { template: '<div class="avatar"><slot /></div>' },
          'q-icon': { template: '<i class="icon" />' },
          'q-space': { template: '<div class="space" />' },
          'q-badge': { template: '<span class="badge"><slot /></span>' },
          'q-chip': { template: '<span class="chip"><slot /></span>' },
          'q-tooltip': { template: '<span class="tooltip"><slot /></span>' },
          'q-drawer': { template: '<aside><slot /></aside>' },
          'q-scroll-area': { template: '<div><slot /></div>' },
          'q-list': { template: '<ul><slot /></ul>' },
          'q-item': { template: '<li><slot /></li>' },
          'q-item-section': { template: '<div><slot /></div>' },
          'q-item-label': { template: '<div><slot /></div>' },
          'q-expansion-item': { template: '<div><slot /></div>' },
          'q-separator': { template: '<hr />' },
          'q-page-container': { template: '<main><slot /></main>' },
          'q-footer': { template: '<footer><slot /></footer>' },
          'q-select': true,
          'q-toggle': true,
          'q-spinner': true,
          'q-breadcrumbs': { template: '<nav class="breadcrumbs"><slot /></nav>' },
          'q-breadcrumbs-el': { template: '<span><slot /></span>' },
          'router-link': { template: '<a><slot /></a>' },
          'router-view': true,
          'GlobalSearch': true,
          'OnboardingTour': true,
          'HelpDrawer': true,
          'HelpCenterPanel': true,
          'SessionReauthDialog': true,
          'ReadingRuler': true,
          'KeyboardShortcutsDialog': true
        }
      }
    })

    const footer = wrapper.find('footer')
    expect(footer.exists()).toBe(true)
    expect(footer.text()).toContain('Registro Elettronico Scolastico')
    expect(footer.text()).toContain('Dichiarazione di Accessibilità (AgID)')
    expect(footer.text()).toContain('Scorciatoie ( ? )')

    // Switch locale to English
    currentLocale.value = 'en-US'
    await wrapper.vm.$nextTick()

    expect(footer.text()).toContain('Electronic School Register')
    expect(footer.text()).toContain('Accessibility Statement (AgID)')
    expect(footer.text()).toContain('Shortcuts ( ? )')
  })

  it('AccessibilityStatement renders in Italian and dynamically switches to English', async () => {
    wrapper = mount(AccessibilityStatement, {
      global: {
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-avatar': { template: '<div><slot /></div>' },
          'q-icon': true,
          'q-banner': { template: '<div><slot /><slot name="action" /></div>' },
          'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
          'q-input': true,
          'q-select': true,
          'q-btn': true
        }
      }
    })

    // Initial Italian check
    expect(wrapper.text()).toContain('Dichiarazione di Accessibilità')
    expect(wrapper.text()).toContain('Stato Conformità')
    expect(wrapper.text()).toContain('Pienamente Conforme WCAG 2.2 AA')
    expect(wrapper.text()).toContain("1. Impegno per l'Inclusione Scolastica Digitale")
    expect(wrapper.text()).toContain('4. Procedura di Attuazione (Difensore Civico per il Digitale)')

    // Switch locale to English
    currentLocale.value = 'en-US'
    await wrapper.vm.$nextTick()

    // English check
    expect(wrapper.text()).toContain('Accessibility Statement')
    expect(wrapper.text()).toContain('Compliance Status')
    expect(wrapper.text()).toContain('Fully Compliant WCAG 2.2 AA')
    expect(wrapper.text()).toContain('1. Commitment to Digital School Inclusion')
    expect(wrapper.text()).toContain('4. Implementation Procedure (Digital Ombudsman)')
  })
})

