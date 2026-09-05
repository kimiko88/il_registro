import { ref, onMounted } from 'vue'

const isOnline = ref(typeof navigator !== 'undefined' ? navigator.onLine : true)
const wasOffline = ref(false)

let listenersAttached = false

function handleOnline() {
    isOnline.value = true
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
}

export function useNetworkStatus() {
    onMounted(() => {
        initNetworkListeners()
    })

    return {
        isOnline,
        wasOffline
    }
}
