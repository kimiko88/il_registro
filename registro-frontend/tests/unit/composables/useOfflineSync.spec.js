import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useOfflineSync } from '@/composables/useOfflineSync'
import { useOutboxStore } from '@/stores/outbox'
import { Notify } from 'quasar'
import api from '@/services/api'

vi.mock('quasar', () => ({
  Notify: {
    create: vi.fn()
  }
}))

vi.mock('@/services/api', () => ({
  default: {
    request: vi.fn()
  }
}))

describe('useOfflineSync Composable', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
    Object.defineProperty(navigator, 'onLine', { value: true, writable: true })
  })

  it('executes API request directly when online and successful', async () => {
    const { executeWithOfflineQueue } = useOfflineSync()
    api.request.mockResolvedValue({ status: 200, data: { ok: true } })

    const res = await executeWithOfflineQueue({ url: '/test', method: 'post', data: { val: 1 } })
    expect(res.data.ok).toBe(true)
    expect(api.request).toHaveBeenCalledTimes(1)
    expect(Notify.create).not.toHaveBeenCalled()
  })

  it('enqueues to outbox and notifies when device is offline', async () => {
    Object.defineProperty(navigator, 'onLine', { value: false, writable: true })
    const { executeWithOfflineQueue, outboxStore } = useOfflineSync()

    const res = await executeWithOfflineQueue({ url: '/test-offline', method: 'post', data: { grade: 8 } }, { title: 'Nuovo Voto' })

    expect(res.offline).toBe(true)
    expect(res.enqueued).toBe(true)
    expect(outboxStore.pendingCount).toBe(1)
    expect(outboxStore.queue[0].title).toBe('Nuovo Voto')
    expect(api.request).not.toHaveBeenCalled()
    expect(Notify.create).toHaveBeenCalledWith(expect.objectContaining({
      type: 'warning',
      icon: 'wifi_off'
    }))
  })

  it('catches network errors when online and gracefully enqueues to outbox', async () => {
    const { executeWithOfflineQueue, outboxStore } = useOfflineSync()
    api.request.mockRejectedValue({
      code: 'ERR_NETWORK',
      message: 'Network Error'
    })

    const res = await executeWithOfflineQueue({ url: '/lessons/sign', method: 'post' }, { title: 'Firma Lezione' })

    expect(res.offline).toBe(true)
    expect(res.enqueued).toBe(true)
    expect(outboxStore.pendingCount).toBe(1)
    expect(outboxStore.queue[0].title).toBe('Firma Lezione')
    expect(Notify.create).toHaveBeenCalledWith(expect.objectContaining({
      type: 'warning',
      icon: 'cloud_queue'
    }))
  })

  it('re-throws non-network errors (e.g. 422 Unprocessable) without enqueuing', async () => {
    const { executeWithOfflineQueue, outboxStore } = useOfflineSync()
    const error422 = {
      response: { status: 422, data: { error: 'Validation failed' } }
    }
    api.request.mockRejectedValue(error422)

    await expect(executeWithOfflineQueue({ url: '/invalid', method: 'post' })).rejects.toEqual(error422)
    expect(outboxStore.pendingCount).toBe(0)
  })

  it('triggers sync manually and emits success notification upon synced items', async () => {
    const outboxStore = useOutboxStore()
    outboxStore.enqueue({ url: '/synced-1', method: 'post' })

    api.request.mockResolvedValue({ status: 200, data: { ok: true } })

    const { triggerSync } = useOfflineSync()
    await triggerSync()

    expect(outboxStore.pendingCount).toBe(0)
    expect(Notify.create).toHaveBeenCalledWith(expect.objectContaining({
      type: 'positive',
      icon: 'cloud_done'
    }))
  })
})
