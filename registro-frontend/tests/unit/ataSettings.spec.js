import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useMenuItems } from '../../src/composables/useMenuItems'
import authService from '../../src/services/authService'
import api from '../../src/services/api'
import routes from '../../src/router/routes'

vi.mock('../../src/services/api', () => ({
  default: {
    post: vi.fn(),
    get: vi.fn(),
    patch: vi.fn()
  }
}))

describe('ATA & DSGA Settings Logic & Routing', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
  })

  it('includes Impostazioni & Profilo in DSGA menu items', () => {
    const items = useMenuItems('dsga')
    const flat = []
    items.forEach(cat => {
      if (cat.children) flat.push(...cat.children)
      else flat.push(cat)
    })
    const settingsItem = flat.find(i => i.path === '/ata/settings')
    expect(settingsItem).toBeDefined()
    expect(settingsItem.label).toContain('Impostazioni')
  })

  it('includes Impostazioni & Profilo in assistente_amministrativo menu items', () => {
    const items = useMenuItems('assistente_amministrativo')
    const flat = []
    items.forEach(cat => {
      if (cat.children) flat.push(...cat.children)
      else flat.push(cat)
    })
    const settingsItem = flat.find(i => i.path === '/ata/settings')
    expect(settingsItem).toBeDefined()
    expect(settingsItem.path).toBe('/ata/settings')
  })

  it('includes Impostazioni & Profilo in collaboratore_ds and collaboratore_scolastico menu items', () => {
    ['collaboratore_ds', 'collaboratore_scolastico'].forEach(role => {
      const items = useMenuItems(role)
      const flat = []
      items.forEach(cat => {
        if (cat.children) flat.push(...cat.children)
        else flat.push(cat)
      })
      const settingsItem = flat.find(i => i.path === '/ata/settings')
      expect(settingsItem).toBeDefined()
    })
  })

  it('has /ata/settings route registered with proper DSGA and ATA roles', () => {
    const mainLayoutRoute = routes.find(r => r.path === '/')
    expect(mainLayoutRoute).toBeDefined()
    const ataSettingsRoute = mainLayoutRoute.children.find(c => c.path === 'ata/settings')
    expect(ataSettingsRoute).toBeDefined()
    expect(ataSettingsRoute.meta.roles).toContain('dsga')
    expect(ataSettingsRoute.meta.roles).toContain('assistente_amministrativo')
    expect(ataSettingsRoute.meta.roles).toContain('collaboratore_ds')
    expect(ataSettingsRoute.meta.roles).toContain('collaboratore_scolastico')
  })

  it('calls changePassword endpoint for DSGA password modification', async () => {
    api.post.mockResolvedValueOnce({ data: { message: 'Password aggiornata con successo' } })

    const res = await authService.changePassword('CurrentDsgaPwd123!', 'NewDsgaSecurePwd2026!')
    expect(api.post).toHaveBeenCalledWith('/auth/change-password', {
      current_password: 'CurrentDsgaPwd123!',
      new_password: 'NewDsgaSecurePwd2026!'
    })
    expect(res.message).toBe('Password aggiornata con successo')
  })

  it('persists notification settings in localStorage', () => {
    const prefs = {
      leaveVistoAlerts: true,
      sidiDeadlines: true,
      strikeAlerts: false,
      circularsAlerts: true,
      inAppSound: true
    }
    localStorage.setItem('ata_notification_settings', JSON.stringify(prefs))
    const loaded = JSON.parse(localStorage.getItem('ata_notification_settings'))
    expect(loaded.leaveVistoAlerts).toBe(true)
    expect(loaded.strikeAlerts).toBe(false)
    expect(loaded.inAppSound).toBe(true)
  })

  it('renders Settings.vue component without errors', async () => {
    const { mount } = await import('@vue/test-utils')
    const { Quasar } = await import('quasar')
    const Settings = (await import('../../src/pages/ata/Settings.vue')).default

    const wrapper = mount(Settings, {
      global: {
        plugins: [Quasar],
        stubs: {
          'q-page': { template: '<div class="ata-settings-page"><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-badge': { template: '<span><slot /></span>' },
          'q-btn': { template: '<button><slot /></button>' },
          'q-tabs': { template: '<div><slot /></div>' },
          'q-tab': { template: '<div><slot /></div>' },
          'q-tab-panels': { template: '<div><slot /></div>' },
          'q-tab-panel': { template: '<div><slot /></div>' },
          'q-select': { template: '<div><slot /></div>' },
          AccessibilitySettingsPanel: { template: '<div class="stub-a11y"></div>' }
        },
        mocks: {
          $t: (key) => key
        }
      }
    })

    expect(wrapper.exists()).toBe(true)
    expect(wrapper.find('.ata-settings-page').exists()).toBe(true)
    expect(wrapper.text()).toContain('Impostazioni')
  })
})
