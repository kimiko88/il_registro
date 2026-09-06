import { ref } from 'vue'

const needRefresh = ref(false)
const updateRegistration = ref(null)
const userDismissed = ref(false)
let listenerAttached = false

export function initPwaUpdateListeners() {
  if (listenerAttached || typeof window === 'undefined' || !('serviceWorker' in navigator)) {
    return
  }

  try {
    navigator.serviceWorker.ready.then((reg) => {
      updateRegistration.value = reg

      reg.addEventListener('updatefound', () => {
        const installingWorker = reg.installing
        if (!installingWorker) return

        installingWorker.addEventListener('statechange', () => {
          if (installingWorker.state === 'installed' && navigator.serviceWorker.controller) {
            needRefresh.value = true
          }
        })
      })
    }).catch(() => {
      // Service worker registration not available
    })
  } catch {
    // Guard against environments where serviceWorker is disabled
  }

  listenerAttached = true
}

export function usePwaUpdate() {
  initPwaUpdateListeners()

  const updateAndReload = async () => {
    if (updateRegistration.value?.waiting) {
      updateRegistration.value.waiting.postMessage({ type: 'SKIP_WAITING' })
    }
    if (typeof window !== 'undefined') {
      window.location.reload()
    }
  }

  const dismiss = () => {
    userDismissed.value = true
  }

  // Testing helpers
  const _setNeedRefresh = (val) => {
    needRefresh.value = val
    userDismissed.value = false
  }

  const _reset = () => {
    needRefresh.value = false
    userDismissed.value = false
    updateRegistration.value = null
  }

  return {
    needRefresh,
    userDismissed,
    updateAndReload,
    dismiss,
    _setNeedRefresh,
    _reset
  }
}
