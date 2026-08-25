import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import App from '@/App.vue'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { useWebSocketStore } from '@/stores/websocket'

// Mock i18n
vi.mock('vue-i18n', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useI18n: () => ({
      t: (key) => key
    })
  }
})

describe('App.vue & Root Lifecycle', () => {
  let wrapper

  beforeEach(() => {
    setActivePinia(createPinia())
  })

  afterEach(() => {
    if (wrapper) {
      wrapper.unmount()
    }
  })

  it('mounts App.vue and initializes theme', () => {
    const themeStore = useThemeStore()
    const initSpy = vi.spyOn(themeStore, 'initTheme')

    wrapper = mount(App, {
      global: {
        stubs: {
          'router-view': true,
          'q-dialog': true,
          'q-card': true,
          'q-card-section': true,
          'q-card-actions': true,
          'q-btn': true,
          'q-icon': true
        }
      }
    })

    expect(wrapper.exists()).toBe(true)
    expect(initSpy).toHaveBeenCalled()
  })

  it('attaches window event listeners on mount and removes them on unmount', () => {
    const addSpy = vi.spyOn(window, 'addEventListener')
    const removeSpy = vi.spyOn(window, 'removeEventListener')

    wrapper = mount(App, {
      global: {
        stubs: {
          'router-view': true,
          'q-dialog': true,
          'q-card': true,
          'q-card-section': true,
          'q-card-actions': true,
          'q-btn': true,
          'q-icon': true
        }
      }
    })

    expect(addSpy).toHaveBeenCalledWith('beforeunload', expect.any(Function))
    expect(addSpy).toHaveBeenCalledWith('keydown', expect.any(Function))

    wrapper.unmount()
    expect(removeSpy).toHaveBeenCalledWith('beforeunload', expect.any(Function))
    expect(removeSpy).toHaveBeenCalledWith('keydown', expect.any(Function))
    wrapper = null
  })

  it('connects websocket when user becomes authenticated and disconnects on logout', async () => {
    const authStore = useAuthStore()
    const wsStore = useWebSocketStore()
    const connectSpy = vi.spyOn(wsStore, 'connect').mockImplementation(() => {})
    const disconnectSpy = vi.spyOn(wsStore, 'disconnect').mockImplementation(() => {})

    wrapper = mount(App, {
      global: {
        stubs: {
          'router-view': true,
          'q-dialog': true,
          'q-card': true,
          'q-card-section': true,
          'q-card-actions': true,
          'q-btn': true,
          'q-icon': true
        }
      }
    })

    const header = btoa(JSON.stringify({ alg: 'HS256', typ: 'JWT' }))
    const payload = btoa(JSON.stringify({
      sub: 'u1',
      role: 'teacher',
      exp: Math.floor(Date.now() / 1000) + 3600
    }))
    const validToken = `${header}.${payload}.signature`

    // Simulate login
    authStore.token = validToken
    await wrapper.vm.$nextTick()
    expect(connectSpy).toHaveBeenCalled()

    // Simulate logout
    authStore.token = null
    await wrapper.vm.$nextTick()
    expect(disconnectSpy).toHaveBeenCalled()
  })
})
