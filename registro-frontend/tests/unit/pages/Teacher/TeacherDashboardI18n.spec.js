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

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn()
  })
}))

import { Quasar } from 'quasar'
import TeacherIndex from '@/pages/Teacher/Index.vue'
import { useTeacherStore } from '@/stores/teacher'
import { useClassesStore } from '@/stores/classes'

describe('Teacher Dashboard (Index.vue) i18n & Reactivity', () => {
  let wrapper

  beforeEach(() => {
    setActivePinia(createPinia())
    currentLocale.value = 'it-IT'
  })

  afterEach(() => {
    if (wrapper) wrapper.unmount()
  })

  const mountPage = () => {
    const teacherStore = useTeacherStore()
    const classesStore = useClassesStore()

    teacherStore.profile = { full_name: 'Prof. Rossi', next_lesson: null }
    teacherStore.upcomingColloqui = 3
    teacherStore.notifications = [
      { id: 1, title: 'Consiglio di Classe', message: 'Ore 15:00', read: false }
    ]
    teacherStore.fetchProfile = vi.fn().mockResolvedValue({})
    teacherStore.fetchNotifications = vi.fn().mockResolvedValue([])
    teacherStore.fetchPendingJustifications = vi.fn().mockResolvedValue(0)
    teacherStore.fetchUpcomingColloqui = vi.fn().mockResolvedValue(3)

    classesStore.classes = [
      { id: 1, name: '3', section: 'A', academic_year: '2025/2026', coordinator_id: 1 }
    ]
    classesStore.fetchAssignedClasses = vi.fn().mockResolvedValue([])

    return mount(TeacherIndex, {
      global: {
        plugins: [Quasar],
        directives: {
          Ripple: { mounted() {}, updated() {}, unmounted() {} },
          ripple: { mounted() {}, updated() {}, unmounted() {} }
        },
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-avatar': { template: '<div><slot /></div>' },
          'q-icon': true,
          'q-btn': {
            props: ['label'],
            template: '<button><slot />{{ label }}</button>'
          },
          'q-list': { template: '<ul><slot /></ul>' },
          'q-item': { template: '<li><slot /></li>' },
          'q-item-section': { template: '<div><slot /></div>' },
          'q-item-label': { template: '<div><slot /></div>' }
        }
      }
    })
  }

  it('renders all Teacher Dashboard sections in Italian and switches dynamically to English', async () => {
    wrapper = mountPage()

    // Italian tests
    expect(wrapper.text()).toContain('Pannello Docente')
    expect(wrapper.text()).toContain('Colloqui')
    expect(wrapper.text()).toContain('3 Prenotazioni')
    expect(wrapper.text()).toContain("Controlla l'agenda")
    expect(wrapper.text()).toContain('Messaggi')
    expect(wrapper.text()).toContain('1 Nuovi')
    expect(wrapper.text()).toContain('Comunicazioni interne')
    expect(wrapper.text()).toContain('Azioni Rapide')
    expect(wrapper.text()).toContain('Registra Voti')
    expect(wrapper.text()).toContain('Presenze')
    expect(wrapper.text()).toContain('Firma Doc')
    expect(wrapper.text()).toContain('Notifiche & Attività')
    expect(wrapper.text()).toContain('Le Mie Classi')

    // Verify month in todayDate is localized in Italian
    const dateTextIt = wrapper.vm.todayDate
    expect(dateTextIt).toBeTruthy()

    // Switch to English
    currentLocale.value = 'en-US'
    await wrapper.vm.$nextTick()

    // English assertions
    expect(wrapper.text()).toContain('Teacher Panel')
    expect(wrapper.text()).toContain('Parent Meetings')
    expect(wrapper.text()).toContain('3 Bookings')
    expect(wrapper.text()).toContain('Check agenda')
    expect(wrapper.text()).toContain('Messages')
    expect(wrapper.text()).toContain('1 New')
    expect(wrapper.text()).toContain('Internal communications')
    expect(wrapper.text()).toContain('Quick Actions')
    expect(wrapper.text()).toContain('Record Grades')
    expect(wrapper.text()).toContain('Attendance')
    expect(wrapper.text()).toContain('Sign Docs')
    expect(wrapper.text()).toContain('Notifications & Activities')
    expect(wrapper.text()).toContain('My Classes')

    // Verify month in todayDate switches to English
    const dateTextEn = wrapper.vm.todayDate
    expect(dateTextEn).not.toBe(dateTextIt)
  })
})
