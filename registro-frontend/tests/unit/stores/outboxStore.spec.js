import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useOutboxStore } from '@/stores/outbox'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    request: vi.fn()
  }
}))

describe('useOutboxStore — Offline Queue & FIFO Sync', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
  })

  it('initializes with empty queue', () => {
    const store = useOutboxStore()
    expect(store.queue).toEqual([])
    expect(store.pendingCount).toBe(0)
    expect(store.hasPending).toBe(false)
  })

  it('enqueues actions and persists them to localStorage', () => {
    const store = useOutboxStore()
    const id = store.enqueue({
      url: '/attendance',
      method: 'post',
      data: { lesson_id: 'l-1', status: 'present' },
      title: 'Firma Presenze'
    })

    expect(id).toMatch(/^outbox_/)
    expect(store.pendingCount).toBe(1)
    expect(store.hasPending).toBe(true)
    expect(store.queue[0].url).toBe('/attendance')
    expect(store.queue[0].title).toBe('Firma Presenze')

    // Verify localStorage persistence
    const raw = localStorage.getItem('registro_offline_outbox')
    expect(raw).toBeTruthy()
    const parsed = JSON.parse(raw)
    expect(parsed.length).toBe(1)
    expect(parsed[0].id).toBe(id)
  })

  it('dequeues specific item and updates storage', () => {
    const store = useOutboxStore()
    const id1 = store.enqueue({ url: '/act-1', method: 'post' })
    const id2 = store.enqueue({ url: '/act-2', method: 'post' })

    expect(store.pendingCount).toBe(2)

    store.dequeue(id1)
    expect(store.pendingCount).toBe(1)
    expect(store.queue[0].id).toBe(id2)

    const parsed = JSON.parse(localStorage.getItem('registro_offline_outbox'))
    expect(parsed.length).toBe(1)
    expect(parsed[0].id).toBe(id2)
  })

  it('clears all items in queue and storage', () => {
    const store = useOutboxStore()
    store.enqueue({ url: '/act-1' })
    store.enqueue({ url: '/act-2' })

    store.clear()
    expect(store.pendingCount).toBe(0)
    expect(localStorage.getItem('registro_offline_outbox')).toBeNull()
  })

  it('synchronizes queue items in sequential FIFO order on network connection', async () => {
    const store = useOutboxStore()
    store.enqueue({ url: '/act-1', method: 'post', data: { order: 1 } })
    store.enqueue({ url: '/act-2', method: 'put', data: { order: 2 } })

    api.request.mockResolvedValue({ status: 200, data: { success: true } })

    const res = await store.syncQueue()

    expect(res.synced).toBe(2)
    expect(res.remaining).toBe(0)
    expect(store.pendingCount).toBe(0)
    expect(api.request).toHaveBeenCalledTimes(2)
    expect(api.request.mock.calls[0][0].url).toBe('/act-1')
    expect(api.request.mock.calls[1][0].url).toBe('/act-2')
  })

  it('halts sync on transient network failure to preserve FIFO ordering', async () => {
    const store = useOutboxStore()
    store.enqueue({ url: '/act-1', method: 'post' })
    store.enqueue({ url: '/act-2', method: 'post' })

    // Simulate network error on first item
    api.request.mockRejectedValueOnce({
      code: 'ERR_NETWORK',
      message: 'Network Error'
    })

    const res = await store.syncQueue()

    expect(res.synced).toBe(0)
    expect(res.remaining).toBe(2)
    expect(store.pendingCount).toBe(2)
    // Only tried first item and stopped
    expect(api.request).toHaveBeenCalledTimes(1)
  })

  it('discards item after 3 persistent client errors to prevent head-of-line blocking', async () => {
    const store = useOutboxStore()
    store.enqueue({ url: '/act-invalid', method: 'post' })

    const error400 = {
      response: { status: 400, data: { error: 'Bad request' } }
    }
    api.request.mockRejectedValue(error400)

    // Attempt 1
    await store.syncQueue()
    expect(store.pendingCount).toBe(1)
    expect(store.queue[0].attempts).toBe(1)

    // Attempt 2
    await store.syncQueue()
    expect(store.pendingCount).toBe(1)
    expect(store.queue[0].attempts).toBe(2)

    // Attempt 3 -> discarded
    const res = await store.syncQueue()
    expect(res.failed).toBe(1)
    expect(store.pendingCount).toBe(0)
  })
})
