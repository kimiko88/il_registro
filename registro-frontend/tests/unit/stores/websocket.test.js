import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useWebSocketStore } from '@/stores/websocket'
import { useAuthStore } from '@/stores/auth'

vi.mock('quasar', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    Notify: {
      create: vi.fn()
    }
  }
})

describe('useWebSocketStore — Connection, Ticket & Reconnect Lifecycle', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  afterEach(() => {
    const wsStore = useWebSocketStore()
    wsStore.disconnect(true)
  })

  it('initializes with default disconnected state', () => {
    const wsStore = useWebSocketStore()
    expect(wsStore.isConnected).toBe(false)
    expect(wsStore.reconnectAttempts).toBe(0)
    expect(wsStore.hasFailedPermanently).toBe(false)
  })

  it('skips connection if user is not authenticated', async () => {
    const authStore = useAuthStore()
    authStore.token = null

    const wsStore = useWebSocketStore()
    await wsStore.connect()

    expect(wsStore.isConnected).toBe(false)
  })

  it('disconnects and resets connection state cleanly', () => {
    const wsStore = useWebSocketStore()
    wsStore.isConnected = true
    wsStore.reconnectAttempts = 3

    wsStore.disconnect(true)

    expect(wsStore.isConnected).toBe(false)
    expect(wsStore.reconnectAttempts).toBe(0)
    expect(wsStore.hasFailedPermanently).toBe(false)
  })

  it('triggers sendDesktopNotification when Notification permission is granted', () => {
    const mockNotification = vi.fn()
    mockNotification.permission = 'granted'
    globalThis.Notification = mockNotification

    const wsStore = useWebSocketStore()
    const notif = wsStore.sendDesktopNotification({
      title: 'Nuovo Voto',
      body: 'Matematica: 8'
    })

    expect(mockNotification).toHaveBeenCalledWith('Nuovo Voto', expect.objectContaining({
      body: 'Matematica: 8'
    }))
  })

  it('dispatches desktop notification on incoming GRADE_ADDED message', () => {
    const mockNotification = vi.fn()
    mockNotification.permission = 'granted'
    globalThis.Notification = mockNotification

    const wsStore = useWebSocketStore()
    wsStore.handleMessage({
      type: 'GRADE_ADDED',
      payload: {
        grade_value: '9',
        subject_name: 'Storia'
      }
    })

    expect(mockNotification).toHaveBeenCalledWith('Registro Elettronico - Voti', expect.objectContaining({
      body: expect.stringContaining('Storia')
    }))
  })
})
