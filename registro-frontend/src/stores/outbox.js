import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'

const STORAGE_KEY = 'registro_offline_outbox'
const DB_NAME = 'registro_offline'
const STORE_NAME = 'outbox'
const DB_VERSION = 1

let _db = null

function openDB() {
    if (typeof indexedDB === 'undefined') return Promise.reject(new Error('IndexedDB not supported'))
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

function idbSafePut(item) {
    if (typeof indexedDB === 'undefined') return
    idbPut(item).catch(() => {})
}

function idbSafeDelete(id) {
    if (typeof indexedDB === 'undefined') return
    idbDelete(id).catch(() => {})
}

function idbSafeClear() {
    if (typeof indexedDB === 'undefined') return
    idbClear().catch(() => {})
}

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
        if (typeof localStorage !== 'undefined') {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(queue))
        }
    } catch {
        // Storage quota exceeded or disabled
    }
}

// ─── Pinia Store ──────────────────────────────────────────────────────────────

export const useOutboxStore = defineStore('outbox', () => {
    const queue = ref(loadStoredQueue())
    const isSyncing = ref(false)
    const lastSyncTime = ref(null)
    const lastSyncResult = ref(null)

    const pendingCount = computed(() => queue.value.length)
    const hasPending = computed(() => queue.value.length > 0)

    /**
     * Load persisted items from localStorage and IndexedDB into reactive state.
     */
    async function init() {
        queue.value = loadStoredQueue()
        if (typeof indexedDB !== 'undefined') {
            try {
                const idbItems = await idbGetAll()
                if (idbItems && idbItems.length > 0) {
                    queue.value = idbItems
                    queue.value.sort((a, b) => a.timestamp - b.timestamp)
                    persistQueue(queue.value)
                }
            } catch (e) {
                console.warn('[outbox] Failed to load from IndexedDB:', e)
            }
        }
    }

    /**
     * Add a new operation to the outbox synchronously and persist.
     * @returns {string} The generated idempotency ID for the item.
     */
    function enqueue({ url, method = 'post', data = null, params = null, title = '' }) {
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
        queue.value.push(item)
        persistQueue(queue.value)
        idbSafePut(item)
        return id
    }

    function dequeue(id) {
        const idx = queue.value.findIndex(item => item.id === id)
        if (idx !== -1) {
            queue.value.splice(idx, 1)
            persistQueue(queue.value)
            idbSafeDelete(id)
        }
    }

    function clear() {
        queue.value = []
        try {
            if (typeof localStorage !== 'undefined') {
                localStorage.removeItem(STORAGE_KEY)
            }
        } catch {
            // Storage disabled
        }
        idbSafeClear()
    }

    /**
     * Replay queued operations against the API preserving FIFO ordering.
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

        const itemsToProcess = [...queue.value]

        for (const item of itemsToProcess) {
            // Respect exponential backoff delay before retrying failing item
            if (item.nextRetryAt && Date.now() < item.nextRetryAt) {
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
                        'Idempotency-Key': item.id
                    }
                })

                dequeue(item.id)
                synced++
            } catch (err) {
                const isNetworkError =
                    !err.response ||
                    err.code === 'ERR_NETWORK' ||
                    (err.message && /network|fetch|timeout/i.test(err.message))

                if (isNetworkError) {
                    // Transient connection loss: stop queue iteration to preserve FIFO ordering
                    break
                } else {
                    // Non-network error (e.g. 500 Server Error or 4xx)
                    item.attempts = (item.attempts || 0) + 1
                    if (item.attempts >= 3) {
                        dequeue(item.id)
                        failed++
                        console.warn('[outbox] Discarding permanently failing item:', item.url, err.response?.status)
                    } else {
                        // Apply exponential backoff only for 5xx server errors so 4xx validation errors don't delay queue drain
                        if (err.response?.status >= 500) {
                            const backoffMs = Math.min(1000 * Math.pow(2, item.attempts), 30000)
                            item.nextRetryAt = Date.now() + backoffMs
                        }
                        persistQueue(queue.value)
                        idbSafePut(item)
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
