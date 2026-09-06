import { ref, onMounted, onUnmounted, getCurrentInstance } from 'vue'

const isOnline = ref(typeof navigator !== 'undefined' ? navigator.onLine : true)
const wasOffline = ref(false)

// Callbacks registered by consumers who want to react to online transitions.
const _onlineHandlers = new Set()

let listenersAttached = false
let _healthCheckInterval = null

/**
 * Derive the ping URL from the Vite env variable.
 * Falls back to /api/v1/ping when no VITE_API_URL is set (same-origin dev).
 */
function getPingUrl() {
    try {
        const raw = import.meta.env.VITE_API_URL
        if (!raw) return '/api/v1/ping'
        const base = String(raw).trim().replace(/\/+$/, '')
        return base.endsWith('/api/v1') ? `${base}/ping` : `${base}/api/v1/ping`
    } catch {
        return '/api/v1/ping'
    }
}

/**
 * Actively probe the backend /api/v1/ping endpoint.
 * navigator.onLine is unreliable: it returns true even on captive portals or
 * networks with no actual internet access.
 */
async function checkRealConnectivity() {
    try {
        let res = await fetch(getPingUrl(), {
            method: 'HEAD',
            cache: 'no-store',
            signal: AbortSignal.timeout(3000)
        })
        // Fallback to GET if server or proxy does not support HEAD
        if (res.status === 404 || res.status === 405) {
            res = await fetch(getPingUrl(), {
                method: 'GET',
                cache: 'no-store',
                signal: AbortSignal.timeout(3000)
            })
        }
        if (res.ok) {
            const wasDown = !isOnline.value
            isOnline.value = true
            if (wasDown) {
                // Fire registered callbacks so outbox/offline-sync can react.
                _onlineHandlers.forEach(fn => { try { fn() } catch { /* noop */ } })
            }
        } else {
            isOnline.value = false
            wasOffline.value = true
        }
    } catch {
        isOnline.value = false
        wasOffline.value = true
    }
}

function handleOnline() {
    // Browser says online; confirm with a real probe before declaring online.
    checkRealConnectivity()
}

function handleOffline() {
    isOnline.value = false
    wasOffline.value = true
}

function initNetworkListeners() {
    if (listenersAttached || typeof window === 'undefined') return
    window.addEventListener('online', handleOnline)
    window.addEventListener('offline', handleOffline)
    listenersAttached = true

    // Poll every 15 s to detect captive-portal / silent network loss situations.
    _healthCheckInterval = setInterval(checkRealConnectivity, 15_000)
}

export function useNetworkStatus() {
    if (getCurrentInstance()) {
        onMounted(() => {
            initNetworkListeners()
        })
        onUnmounted(() => {
            // Only tear down when no other component holds a reference.
            // Since we use module-level state we intentionally leave the
            // interval running as long as any component is mounted — the
            // interval is cheap (one HEAD request every 15 s).
        })
    } else {
        initNetworkListeners()
    }

    return {
        isOnline,
        wasOffline,
        /**
         * Register a callback that fires every time connectivity is restored.
         * Returns a cleanup function to deregister.
         */
        onOnline(fn) {
            _onlineHandlers.add(fn)
            return () => _onlineHandlers.delete(fn)
        }
    }
}
