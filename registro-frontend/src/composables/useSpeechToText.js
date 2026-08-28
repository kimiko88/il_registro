import { ref, onUnmounted } from 'vue'

export function useSpeechToText(options = {}) {
    const isListening = ref(false)
    const transcript = ref('')
    const isSupported = ref(false)
    const error = ref(null)

    const SpeechRecognition = typeof window !== 'undefined'
        ? (window.SpeechRecognition || window.webkitSpeechRecognition)
        : null

    isSupported.value = !!SpeechRecognition

    let recognition = null

    if (SpeechRecognition) {
        recognition = new SpeechRecognition()
        recognition.continuous = options.continuous ?? true
        recognition.interimResults = options.interimResults ?? true
        recognition.lang = options.lang || 'it-IT'

        recognition.onstart = () => {
            isListening.value = true
            error.value = null
        }

        recognition.onresult = (event) => {
            let currentTranscript = ''
            for (let i = event.resultIndex; i < event.results.length; i++) {
                const result = event.results[i]
                if (result.isFinal) {
                    currentTranscript += result[0].transcript + ' '
                }
            }
            if (currentTranscript.trim()) {
                transcript.value = currentTranscript.trim()
                if (options.onResult) {
                    options.onResult(currentTranscript.trim())
                }
            }
        }

        recognition.onerror = (event) => {
            console.warn('SpeechRecognition error:', event.error)
            error.value = event.error
            isListening.value = false
            if (options.onError) {
                options.onError(event.error)
            }
        }

        recognition.onend = () => {
            isListening.value = false
            if (options.onEnd) {
                options.onEnd()
            }
        }
    }

    function start() {
        if (!recognition) {
            error.value = 'speech_recognition_not_supported'
            return
        }
        try {
            transcript.value = ''
            recognition.start()
        } catch (e) {
            console.warn('Failed to start SpeechRecognition:', e)
        }
    }

    function stop() {
        if (recognition && isListening.value) {
            try {
                recognition.stop()
            } catch (e) {
                console.warn('Failed to stop SpeechRecognition:', e)
            }
        }
    }

    function toggle() {
        if (isListening.value) {
            stop()
        } else {
            start()
        }
    }

    onUnmounted(() => {
        stop()
    })

    return {
        isListening,
        transcript,
        isSupported,
        error,
        start,
        stop,
        toggle
    }
}
