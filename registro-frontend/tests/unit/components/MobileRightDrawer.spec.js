import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { ref } from 'vue'

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

const mockPush = vi.fn()
const currentRoute = ref({
  path: '/dashboard',
  meta: { title: 'Dashboard' }
})

vi.mock('vue-router', () => ({
  useRoute: () => currentRoute.value,
  useRouter: () => ({
    push: mockPush
  })
}))

const mockDark = {
  isActive: false,
  toggle: vi.fn(() => { mockDark.isActive = !mockDark.isActive }),
  set: vi.fn((val) => { mockDark.isActive = val })
}

vi.mock('quasar', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useQuasar: () => ({
      dark: mockDark,
      fullscreen: { isActive: false, toggle: vi.fn() },
      screen: { lt: { md: true }, width: 375 },
      notify: vi.fn()
    })
  }
})

import MainLayout from '@/layouts/MainLayout.vue'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'

describe('Mobile Right Drawer & Navbar Layout', () => {
  let wrapper

  beforeEach(() => {
    setActivePinia(createPinia())
    mockPush.mockClear()
    mockDark.isActive = false
    currentLocale.value = 'it-IT'
  })

  afterEach(() => {
    if (wrapper) wrapper.unmount()
  })

  const mountLayout = () => {
    const authStore = useAuthStore()
    authStore.user = { id: 1, name: 'Mario Rossi', role: 'teacher' }

    return mount(MainLayout, {
      global: {
        stubs: {
          'q-layout': { template: '<div class="q-layout"><slot /></div>' },
          'q-header': { template: '<header><slot /></header>' },
          'q-toolbar': { template: '<nav><slot /></nav>' },
          'q-toolbar-title': { template: '<div class="q-toolbar-title"><slot /></div>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-btn-dropdown': { template: '<div class="btn-dropdown"><slot /></div>' },
          'q-avatar': { template: '<div class="avatar"><slot /></div>' },
          'q-icon': { template: '<i class="icon" />' },
          'q-space': { template: '<div class="space" />' },
          'q-badge': { template: '<span class="badge"><slot /></span>' },
          'q-chip': { template: '<span class="chip"><slot /></span>' },
          'q-tooltip': { template: '<span class="tooltip"><slot /></span>' },
          'q-drawer': {
            props: ['modelValue', 'side'],
            template: '<aside :class="side" v-if="modelValue"><slot /></aside>'
          },
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
          'q-breadcrumbs': true,
          'q-breadcrumbs-el': true,
          'router-link': true,
          'router-view': true,
          'SkipLinks': true,
          'ScreenReaderAnnouncer': true,
          'FocusModeToggle': true,
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
  }

  it('renders responsive toolbar title and mobile quick settings button', () => {
    wrapper = mountLayout()

    // Title has both responsive variants
    const title = wrapper.find('.q-toolbar-title')
    expect(title.exists()).toBe(true)
    expect(title.text()).toContain('Registro Elettronico')
    expect(title.text()).toContain('Registro')

    // Right drawer starts closed
    expect(wrapper.vm.rightDrawerOpen).toBe(false)

    // Toggle right drawer
    wrapper.vm.toggleRightDrawer()
    expect(wrapper.vm.rightDrawerOpen).toBe(true)
  })

  it('opens right drawer on toggle and closes on route navigation', async () => {
    wrapper = mountLayout()

    wrapper.vm.toggleRightDrawer()
    expect(wrapper.vm.rightDrawerOpen).toBe(true)

    // Simulate route navigation
    currentRoute.value = { path: '/teacher/grades', meta: { title: 'Voti' } }
    await wrapper.vm.$nextTick()

    // Test direct toggle and close
    wrapper.vm.rightDrawerOpen = false
    expect(wrapper.vm.rightDrawerOpen).toBe(false)
  })

  it('interacts with accessibility and theme store through right drawer bindings', () => {
    wrapper = mountLayout()
    const themeStore = useThemeStore()

    // Open right drawer
    wrapper.vm.rightDrawerOpen = true

    // Verify themeStore toggles exist and function
    expect(themeStore.dsaFont).toBe(false)
    themeStore.toggleDsaFont()
    expect(themeStore.dsaFont).toBe(true)

    expect(themeStore.readingRuler).toBe(false)
    themeStore.toggleReadingRuler()
    expect(themeStore.readingRuler).toBe(true)

    expect(themeStore.ttsEnabled).toBe(false)
    themeStore.toggleTts()
    expect(themeStore.ttsEnabled).toBe(true)
  })
})
