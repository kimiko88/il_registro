import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from 'src/stores/auth'
import { useThemeStore } from 'src/stores/theme'
import { useSpeechSynthesis } from 'src/composables/useSpeechSynthesis'

export function useGlobalKeyboardShortcuts() {
  const router = useRouter()
  const authStore = useAuthStore()
  const themeStore = useThemeStore()
  const { speak, stop, isPlaying } = useSpeechSynthesis()

  function handleKeyDown(e) {
    // Ignore if focus is in an input, textarea or contenteditable element, UNLESS modifier key Alt is used
    const activeEl = document.activeElement
    const isInputFocused = activeEl && (
      activeEl.tagName === 'INPUT' ||
      activeEl.tagName === 'TEXTAREA' ||
      activeEl.isContentEditable ||
      activeEl.classList.contains('q-field__native')
    )

    // Help Dialog toggle: ? (Shift + /) when not in an input
    if (e.key === '?' && !isInputFocused) {
      e.preventDefault()
      themeStore.toggleKeyboardShortcutsHelp()
      return
    }

    // Escape handling: stop speech or close help
    if (e.key === 'Escape') {
      if (isPlaying.value) {
        stop()
      }
      if (themeStore.keyboardShortcutsHelpOpen) {
        themeStore.toggleKeyboardShortcutsHelp(false)
      }
      return
    }

    // Alt + modifier shortcuts
    if (e.altKey && !e.ctrlKey && !e.metaKey) {
      const key = e.key.toLowerCase()

      switch (key) {
        case '1': {
          e.preventDefault()
          const role = authStore.userRole || 'teacher'
          router.push(`/${role}/dashboard`).catch(() => {})
          break
        }
        case 'v': {
          e.preventDefault()
          const role = authStore.userRole || 'teacher'
          router.push(`/${role}/grades`).catch(() => {})
          break
        }
        case 'p': {
          e.preventDefault()
          const role = authStore.userRole || 'teacher'
          router.push(`/${role}/attendance`).catch(() => {})
          break
        }
        case 'a': {
          e.preventDefault()
          const role = authStore.userRole || 'teacher'
          if (role === 'teacher' || role === 'student') {
            router.push(`/${role}/agenda`).catch(() => {})
          }
          break
        }
        case 'r': {
          e.preventDefault()
          themeStore.toggleReadingRuler()
          break
        }
        case 't': {
          e.preventDefault()
          if (isPlaying.value) {
            stop()
          } else {
            const selectedText = window.getSelection ? window.getSelection().toString() : ''
            if (selectedText) {
              speak(selectedText)
            } else {
              // Read main content or heading
              const mainH1 = document.querySelector('h1, .text-h4, .text-h5')
              if (mainH1) {
                speak(mainH1.textContent)
              }
            }
          }
          break
        }
        case 's': {
          e.preventDefault()
          const searchInput = document.querySelector('input[type="search"], .q-field input')
          if (searchInput) {
            searchInput.focus()
          }
          break
        }
      }
    }
  }

  onMounted(() => {
    if (typeof window !== 'undefined') {
      window.addEventListener('keydown', handleKeyDown)
    }
  })

  onUnmounted(() => {
    if (typeof window !== 'undefined') {
      window.removeEventListener('keydown', handleKeyDown)
    }
  })

  return {
    handleKeyDown
  }
}

export default useGlobalKeyboardShortcuts
