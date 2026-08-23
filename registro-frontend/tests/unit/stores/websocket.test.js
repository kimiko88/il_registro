import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useWebSocketStore } from '@/stores/websocket'
import { useAuthStore } from '@/stores/auth'

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
})
