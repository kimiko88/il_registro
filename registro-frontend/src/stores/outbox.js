import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'

// ─── IndexedDB helpers ────────────────────────────────────────────────────────
// Using the native IndexedDB API (no external dependency) wrapped in small
// promise helpers. IndexedDB is:
//  - Asynchronous (non-blocking, unlike localStorage)
//  - Transactional (no corruption on crash/tab close mid-write)
//  - Much larger quota (~50% of free disk space vs ~5–10 MB)

const DB_NAME = 'registro_offline'
const STORE_NAME = 'outbox'
const DB_VERSION = 1

let _db = null

function openDB() {
    if (_db) return Promise.resolve(_db)
    return new Promise((resolve, reject) => {
        const req = indexedDB.open(DB_NAME, DB_VERSION)
        req.onupgradeneeded = (e) => {
            const db = e.target.result
            if (!db.objectStoreNames.contains(STORE_NAME)) {
                db.createObjectStore(STORE_NAME, { keyPath: 'id' })
            }
        }
        req.onsuccess = (e) => {
            _db = e.target.result
            resolve(_db)
        }
        req.onerror = (e) => reject(e.target.error)
    })
}

async function idbGetAll() {
    const db = await openDB()
    return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, 'readonly')
        const req = tx.objectStore(STORE_NAME).getAll()
        req.onsuccess = () => resolve(req.result || [])
        req.onerror = () => reject(req.error)
    })
}

async function idbPut(item) {
    const db = await openDB()
    return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, 'readwrite')
        const req = tx.objectStore(STORE_NAME).put(item)
        req.onsuccess = () => resolve(req.result)
        req.onerror = () => reject(req.error)
    })
}

async function idbDelete(id) {
    const db = await openDB()
    return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, 'readwrite')
        const req = tx.objectStore(STORE_NAME).delete(id)
        req.onsuccess = () => resolve()
        req.onerror = () => reject(req.error)
    })
}

async function idbClear() {
    const db = await openDB()
    return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, 'readwrite')
        const req = tx.objectStore(STORE_NAME).clear()
        req.onsuccess = () => resolve()
        req.onerror = () => reject(req.error)
    })
}

// ─── Pinia Store ──────────────────────────────────────────────────────────────

export const useOutboxStore = defineStore('outbox', () => {
    const queue = ref([])
    const isSyncing = ref(false)
    const lastSyncTime = ref(null)
    const lastSyncResult = ref(null)

    const pendingCount = computed(() => queue.value.length)
    const hasPending = computed(() => queue.value.length > 0)

    /**
     * Load all persisted items from IndexedDB into reactive state.
     * Call once on app boot (e.g. in main.js or App.vue onMounted).
     */
    async function init() {
        try {
            queue.value = await idbGetAll()
            // Sort by timestamp to preserve FIFO order
            queue.value.sort((a, b) => a.timestamp - b.timestamp)
        } catch (e) {
            console.warn('[outbox] Failed to load from IndexedDB:', e)
            queue.value = []
        }
    }

    /**
     * Add a new operation to the outbox.
     * @returns {string} The generated idempotency ID for the item.
     */
    async function enqueue({ url, method = 'post', data = null, params = null, title = '' }) {
        const id = `outbox_${Date.now()}_${Math.random().toString(36).substring(2, 8)}`
        const item = {
            id,
            url,
            method: method.toLowerCase(),
            data,
            params,
            title: title || `${method.toUpperCase()} ${url}`,
            timestamp: Date.now(),
            attempts: 0,
            nextRetryAt: 0
        }
        queue.value.push(item)
        try {
            await idbPut(item)
        } catch (e) {
            console.warn('[outbox] Failed to persist to IndexedDB:', e)
        }
        return id
    }

    async function dequeue(id) {
        const idx = queue.value.findIndex(item => item.id === id)
        if (idx !== -1) queue.value.splice(idx, 1)
        try {
            await idbDelete(id)
        } catch (e) {
            console.warn('[outbox] Failed to delete from IndexedDB:', e)
        }
    }

    async function clear() {
        queue.value = []
        try {
            await idbClear()
        } catch (e) {
            console.warn('[outbox] Failed to clear IndexedDB:', e)
        }
    }

    /**
     * Attempt to replay all queued operations against the API.
     *
     * Improvements over the original:
     *  - Adds `Idempotency-Key` header (= item.id) so the backend safely ignores
     *    duplicates if a previous attempt reached the server but the response was lost.
     *  - Exponential backoff: failed items are skipped until their `nextRetryAt`
     *    timestamp elapses (1s, 2s, 4s … capped at 30s).
     *  - FIFO ordering is preserved: network errors stop iteration immediately.
     */
    async function syncQueue() {
        if (isSyncing.value || queue.value.length === 0) {
            return { synced: 0, remaining: queue.value.length, failed: 0 }
        }

        if (typeof navigator !== 'undefined' && navigator.onLine === false) {
            return { synced: 0, remaining: queue.value.length, failed: 0 }
        }

        isSyncing.value = true
        let synced = 0
        let failed = 0
        const now = Date.now()

        const itemsToProcess = [...queue.value]

        for (const item of itemsToProcess) {
            // Respect exponential backoff — skip items not yet ready to retry.
            if (item.nextRetryAt && now < item.nextRetryAt) {
                continue
            }

            try {
                await api.request({
                    url: item.url,
                    method: item.method,
                    data: item.data,
                    params: item.params,
                    headers: {
                        'x-offline-sync': 'true',
                        // The backend IdempotencyMiddleware uses this to deduplicate
                        // replayed requests whose original response was lost.
                        'Idempotency-Key': item.id
                    }
                })

                await dequeue(item.id)
                synced++
            } catch (err) {
                const isNetworkError =
                    !err.response ||
                    err.code === 'ERR_NETWORK' ||
                    (err.message && /network|fetch|timeout/i.test(err.message))

                if (isNetworkError) {
                    // Transient loss — stop FIFO queue to avoid out-of-order execution.
                    break
                } else {
                    // Permanent/validation error (4xx). Apply exponential backoff.
                    item.attempts = (item.attempts || 0) + 1
                    if (item.attempts >= 5) {
                        // After 5 attempts, discard to avoid head-of-line blocking.
                        await dequeue(item.id)
                        failed++
                        console.warn('[outbox] Discarding permanently failing item:', item.url, err.response?.status)
                    } else {
                        // Schedule next retry with exponential backoff (max 30s).
                        const delay = Math.min(1000 * Math.pow(2, item.attempts - 1), 30_000)
                        item.nextRetryAt = Date.now() + delay
                        try { await idbPut(item) } catch { /* noop */ }
                    }
                }
            }
        }

        lastSyncTime.value = Date.now()
        lastSyncResult.value = { synced, remaining: queue.value.length, failed }
        isSyncing.value = false

        return lastSyncResult.value
    }

    return {
        queue,
        isSyncing,
        lastSyncTime,
        lastSyncResult,
        pendingCount,
        hasPending,
        init,
        enqueue,
        dequeue,
        clear,
        syncQueue
    }
})
