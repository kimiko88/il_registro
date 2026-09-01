import { ref, onMounted, onUnmounted } from 'vue'
import { useThemeStore } from 'src/stores/theme'

export function useSpeechSynthesis() {
  const themeStore = useThemeStore()
  
  const isSupported = ref(typeof window !== 'undefined' && 'speechSynthesis' in window)
  const isPlaying = ref(false)
  const isPaused = ref(false)
  const currentText = ref('')
  const voices = ref([])

  function loadVoices() {
    if (typeof window !== 'undefined' && 'speechSynthesis' in window) {
      voices.value = window.speechSynthesis.getVoices()
    }
  }

  onMounted(() => {
    loadVoices()
    if (typeof window !== 'undefined' && 'speechSynthesis' in window) {
      window.speechSynthesis.onvoiceschanged = loadVoices
    }
  })

  onUnmounted(() => {
    stop()
  })

  function speak(text, options = {}) {
    if (!isSupported.value || !text) return

    stop()

    const cleanText = text.replace(/<[^>]*>?/gm, '').trim()
    if (!cleanText) return

    currentText.value = cleanText

    const utterance = new SpeechSynthesisUtterance(cleanText)
    utterance.rate = options.rate !== undefined ? options.rate : themeStore.ttsRate || 1.0
    utterance.pitch = options.pitch !== undefined ? options.pitch : themeStore.ttsPitch || 1.0
    utterance.lang = options.lang || 'it-IT'

    // Try to find Italian voice if lang is it-IT
    if (voices.value.length > 0) {
      const match = voices.value.find(v => v.lang.startsWith(utterance.lang.substring(0, 2)))
      if (match) {
        utterance.voice = match
      }
    }

    utterance.onstart = () => {
      isPlaying.value = true
      isPaused.value = false
    }

    utterance.onend = () => {
      isPlaying.value = false
      isPaused.value = false
      currentText.value = ''
    }

    utterance.onerror = () => {
      isPlaying.value = false
      isPaused.value = false
      currentText.value = ''
    }

    utterance.onpause = () => {
      isPaused.value = true
    }

    utterance.onresume = () => {
      isPaused.value = false
    }

    window.speechSynthesis.speak(utterance)
  }

  function pause() {
    if (isSupported.value && window.speechSynthesis.speaking) {
      window.speechSynthesis.pause()
      isPaused.value = true
    }
  }

  function resume() {
    if (isSupported.value && window.speechSynthesis.paused) {
      window.speechSynthesis.resume()
      isPaused.value = false
    }
  }

  function stop() {
    if (isSupported.value) {
      window.speechSynthesis.cancel()
      isPlaying.value = false
      isPaused.value = false
      currentText.value = ''
    }
  }

  function toggle(text, options = {}) {
    if (isPlaying.value) {
      if (isPaused.value) {
        resume()
      } else {
        stop()
      }
    } else {
      speak(text, options)
    }
  }

  return {
    isSupported,
    isPlaying,
    isPaused,
    currentText,
    voices,
    speak,
    pause,
    resume,
    stop,
    toggle
  }
}

export default useSpeechSynthesis
