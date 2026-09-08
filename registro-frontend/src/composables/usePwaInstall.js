import { ref, computed } from 'vue'

const deferredPrompt = ref(null)
const isInstalled = ref(false)
const userDismissed = ref(false)
let listenerAttached = false

function checkInstalled() {
  if (typeof window === 'undefined') return false
  return (
    (typeof window.matchMedia === 'function' && window.matchMedia('(display-mode: standalone)').matches) ||
    window.navigator?.standalone === true ||
    (typeof document !== 'undefined' && document.referrer?.includes('android-app://'))
  )
}

function handleBeforeInstallPrompt(e) {
  e.preventDefault()
  deferredPrompt.value = e
}

function handleAppInstalled() {
  isInstalled.value = true
  deferredPrompt.value = null
}

export function initPwaListeners() {
  if (listenerAttached || typeof window === 'undefined') return
  isInstalled.value = checkInstalled()
  window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt)
  window.addEventListener('appinstalled', handleAppInstalled)
  listenerAttached = true
}

export function usePwaInstall() {
  initPwaListeners()

  const canInstall = computed(() => {
    return !!deferredPrompt.value && !isInstalled.value && !userDismissed.value
  })

  async function promptInstall() {
    if (!deferredPrompt.value) return false
    deferredPrompt.value.prompt()
    const choice = await deferredPrompt.value.userChoice
    const outcome = choice?.outcome
    if (outcome === 'accepted') {
      isInstalled.value = true
    }
    deferredPrompt.value = null
    return outcome === 'accepted'
  }

  function dismissPrompt() {
    userDismissed.value = true
  }

  // Testing helper
  function _setPrompt(promptObj) {
    deferredPrompt.value = promptObj
  }

  function _resetState() {
    deferredPrompt.value = null
    isInstalled.value = false
    userDismissed.value = false
  }

  return {
    canInstall,
    isInstalled,
    promptInstall,
    dismissPrompt,
    _setPrompt,
    _resetState
  }
}
