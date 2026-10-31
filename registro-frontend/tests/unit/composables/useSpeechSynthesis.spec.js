import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { useSpeechSynthesis } from '@/composables/useSpeechSynthesis'

describe('useSpeechSynthesis Composable', () => {
    let pinia
    let mockSpeechSynthesis
    let createdUtterance

    beforeEach(() => {
        pinia = createPinia()
        setActivePinia(pinia)

        mockSpeechSynthesis = {
            speaking: false,
            paused: false,
            getVoices: vi.fn().mockReturnValue([
                { name: 'Italian Voice', lang: 'it-IT' },
                { name: 'English Voice', lang: 'en-US' }
            ]),
            speak: vi.fn((utt) => {
                createdUtterance = utt
                mockSpeechSynthesis.speaking = true
                mockSpeechSynthesis.paused = false
                if (utt.onstart) utt.onstart()
            }),
            cancel: vi.fn(() => {
                mockSpeechSynthesis.speaking = false
                mockSpeechSynthesis.paused = false
            }),
            pause: vi.fn(() => {
                mockSpeechSynthesis.speaking = false
                mockSpeechSynthesis.paused = true
            }),
            resume: vi.fn(() => {
                mockSpeechSynthesis.speaking = true
                mockSpeechSynthesis.paused = false
            }),
            onvoiceschanged: null
        }

        class MockUtterance {
            constructor(text) {
                this.text = text
                this.rate = 1
                this.pitch = 1
                this.lang = 'it-IT'
                this.voice = null
                this.onstart = null
                this.onend = null
                this.onerror = null
                this.onpause = null
                this.onresume = null
            }
        }

        global.window.speechSynthesis = mockSpeechSynthesis
        global.SpeechSynthesisUtterance = MockUtterance
    })

    afterEach(() => {
        vi.clearAllMocks()
    })

    it('detects browser support and initializes state flags', () => {
        const { isSupported, isPlaying, isPaused, currentText } = useSpeechSynthesis()

        expect(isSupported.value).toBe(true)
        expect(isPlaying.value).toBe(false)
        expect(isPaused.value).toBe(false)
        expect(currentText.value).toBe('')
    })

    it('strips HTML tags and speaks clean text with Italian voice match', () => {
        const { speak, isPlaying, currentText } = useSpeechSynthesis()

        speak('<p>Benvenuto nel <strong>Registro Elettronico</strong></p>')

        expect(mockSpeechSynthesis.cancel).toHaveBeenCalled()
        expect(mockSpeechSynthesis.speak).toHaveBeenCalled()
        expect(currentText.value).toBe('Benvenuto nel Registro Elettronico')
        expect(isPlaying.value).toBe(true)
    })

    it('applies ttsRate and ttsPitch from themeStore or custom options', () => {
        const themeStore = useThemeStore()
        themeStore.ttsRate = 1.2
        themeStore.ttsPitch = 0.9

        const { speak } = useSpeechSynthesis()
        speak('Test personalizzazione audio', { rate: 1.5 })

        expect(createdUtterance.rate).toBe(1.5) // explicit option takes precedence
        expect(createdUtterance.pitch).toBe(0.9) // falls back to themeStore
    })

    it('pause, resume, and stop controls manipulate playback state', () => {
        const { speak, pause, resume, stop, isPlaying, isPaused } = useSpeechSynthesis()

        speak('Messaggio per il test dei controlli')
        expect(isPlaying.value).toBe(true)
        expect(mockSpeechSynthesis.speaking).toBe(true)

        pause()
        expect(mockSpeechSynthesis.pause).toHaveBeenCalled()
        expect(isPaused.value).toBe(true)

        resume()
        expect(mockSpeechSynthesis.resume).toHaveBeenCalled()
        expect(isPaused.value).toBe(false)

        stop()
        expect(mockSpeechSynthesis.cancel).toHaveBeenCalled()
        expect(isPlaying.value).toBe(false)
        expect(isPaused.value).toBe(false)
    })

    it('toggle stops when playing and speaks when stopped', () => {
        const { toggle, isPlaying, isPaused } = useSpeechSynthesis()

        // 1. Initially stopped -> calling toggle starts speech
        toggle('Test toggle vocal')
        expect(mockSpeechSynthesis.speak).toHaveBeenCalled()
        expect(isPlaying.value).toBe(true)

        // 2. Currently playing -> calling toggle stops speech
        toggle('Test toggle vocal')
        expect(mockSpeechSynthesis.cancel).toHaveBeenCalled()
        expect(isPlaying.value).toBe(false)
    })

    it('handles utterance onend and clears playback state', () => {
        const { speak, isPlaying, currentText } = useSpeechSynthesis()

        speak('Messaggio breve')
        expect(isPlaying.value).toBe(true)

        // Trigger onend event
        createdUtterance.onend()
        expect(isPlaying.value).toBe(false)
        expect(currentText.value).toBe('')
    })

    it('handles utterance onerror gracefully', () => {
        const { speak, isPlaying } = useSpeechSynthesis()

        speak('Messaggio con errore')
        createdUtterance.onerror(new Error('Audio playback interrupted'))

        expect(isPlaying.value).toBe(false)
    })
})
