import { defineStore } from 'pinia'
import api from '@/services/api'

const STORAGE_KEY = 'registro_offline_outbox'

function loadStoredQueue() {
  try {
    if (typeof localStorage === 'undefined') return []
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw)
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

function persistQueue(queue) {
  try {
    if (typeof localStorage === 'undefined') return
    localStorage.setItem(STORAGE_KEY, JSON.stringify(queue))
  } catch {
    // Storage quota exceeded or disabled
  }
}

export const useOutboxStore = defineStore('outbox', {
  state: () => ({
    queue: loadStoredQueue(),
    isSyncing: false,
    lastSyncTime: null,
    lastSyncResult: null
  }),

  getters: {
    pendingCount: (state) => state.queue.length,
    hasPending: (state) => state.queue.length > 0
  },

  actions: {
    init() {
      this.queue = loadStoredQueue()
    },

    enqueue({ url, method = 'post', data = null, params = null, title = '' }) {
      const id = `outbox_${Date.now()}_${Math.random().toString(36).substring(2, 8)}`
      const item = {
        id,
        url,
        method: method.toLowerCase(),
        data,
        params,
        title: title || `${method.toUpperCase()} ${url}`,
        timestamp: Date.now(),
        attempts: 0
      }
      this.queue.push(item)
      persistQueue(this.queue)
      return id
    },

    dequeue(id) {
      const idx = this.queue.findIndex(item => item.id === id)
      if (idx !== -1) {
        this.queue.splice(idx, 1)
        persistQueue(this.queue)
      }
    },

    clear() {
      this.queue = []
      try {
        if (typeof localStorage !== 'undefined') {
          localStorage.removeItem(STORAGE_KEY)
        }
      } catch {
        // Storage disabled
      }
    },

    async syncQueue() {
      if (this.isSyncing || this.queue.length === 0) {
        return { synced: 0, remaining: this.queue.length, failed: 0 }
      }

      if (typeof navigator !== 'undefined' && navigator.onLine === false) {
        return { synced: 0, remaining: this.queue.length, failed: 0 }
      }

      this.isSyncing = true
      let synced = 0
      let failed = 0

      // Copy queue snapshot to iterate FIFO safely
      const itemsToProcess = [...this.queue]

      for (const item of itemsToProcess) {
        try {
          await api.request({
            url: item.url,
            method: item.method,
            data: item.data,
            params: item.params,
            headers: {
              'x-offline-sync': 'true'
            }
          })

          // Success: remove item from queue
          this.dequeue(item.id)
          synced++
        } catch (err) {
          const isNetworkError = !err.response || err.code === 'ERR_NETWORK' || (err.message && /network|fetch|timeout/i.test(err.message))

          if (isNetworkError) {
            // Transient connection loss: stop queue iteration to preserve FIFO ordering
            break
          } else {
            // Permanent/validation error (e.g. 400 Bad Request, 404, 422)
            item.attempts = (item.attempts || 0) + 1
            if (item.attempts >= 3) {
              // Discard permanently failing item after 3 attempts to avoid head-of-line blocking
              this.dequeue(item.id)
              failed++
            } else {
              persistQueue(this.queue)
            }
          }
        }
      }

      this.lastSyncTime = Date.now()
      this.lastSyncResult = { synced, remaining: this.queue.length, failed }
      this.isSyncing = false

      return this.lastSyncResult
    }
  }
})
