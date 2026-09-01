import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { mount } from '@vue/test-utils'
import { useThemeStore } from '@/stores/theme'
import { useSpeechSynthesis } from '@/composables/useSpeechSynthesis'
import TextToSpeechButton from '@/components/Common/TextToSpeechButton.vue'
import ReadingRuler from '@/components/Common/ReadingRuler.vue'
import KeyboardShortcutsDialog from '@/components/Common/KeyboardShortcutsDialog.vue'
import AccessibilityStatement from '@/pages/AccessibilityStatement.vue'

describe('Suite di Accessibilità (A11y), DSA & AgID Compliance', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
  })

  describe('1. Theme & A11y Store State & Actions', () => {
    it('inizializza le impostazioni di default', () => {
      const themeStore = useThemeStore()
      expect(themeStore.dsaFont).toBe(false)
      expect(themeStore.readingRuler).toBe(false)
      expect(themeStore.highContrast).toBe(false)
      expect(themeStore.colorblindMode).toBe('none')
      expect(themeStore.highContrastMode).toBe('none')
      expect(themeStore.lineHeight).toBe('normal')
      expect(themeStore.letterSpacing).toBe('normal')
      expect(themeStore.wordSpacing).toBe('normal')
      expect(themeStore.focusHighlight).toBe(true)
    })

    it('attiva e disattiva il font OpenDyslexic (DSA)', () => {
      const themeStore = useThemeStore()
      themeStore.toggleDsaFont(true)
      expect(themeStore.dsaFont).toBe(true)
      expect(themeStore.fontFamily).toBe('opendyslexic')
      expect(localStorage.getItem('il_registro_dsa_font')).toBe('true')

      themeStore.toggleDsaFont(false)
      expect(themeStore.dsaFont).toBe(false)
      expect(themeStore.fontFamily).toBe('default')
      expect(localStorage.getItem('il_registro_dsa_font')).toBe('false')
    })

    it('attiva e modifica il righello di lettura (Reading Ruler)', () => {
      const themeStore = useThemeStore()
      themeStore.toggleReadingRuler(true)
      expect(themeStore.readingRuler).toBe(true)
      expect(localStorage.getItem('il_registro_reading_ruler')).toBe('true')

      themeStore.setReadingRulerHeight(48)
      expect(themeStore.readingRulerHeight).toBe(48)
      expect(localStorage.getItem('il_registro_ruler_height')).toBe('48')

      themeStore.setReadingRulerOpacity(0.5)
      expect(themeStore.readingRulerOpacity).toBe(0.5)
      expect(localStorage.getItem('il_registro_ruler_opacity')).toBe('0.5')
    })

    it('imposta la spaziatura testo (WCAG 1.4.12)', () => {
      const themeStore = useThemeStore()
      themeStore.setLineHeight('relaxed')
      expect(themeStore.lineHeight).toBe('relaxed')
      expect(localStorage.getItem('il_registro_line_height')).toBe('relaxed')

      themeStore.setLetterSpacing('wide')
      expect(themeStore.letterSpacing).toBe('wide')
      expect(localStorage.getItem('il_registro_letter_spacing')).toBe('wide')

      themeStore.setWordSpacing('wide')
      expect(themeStore.wordSpacing).toBe('wide')
      expect(localStorage.getItem('il_registro_word_spacing')).toBe('wide')
    })

    it('imposta i filtri daltonismo e contrasti OLED', () => {
      const themeStore = useThemeStore()
      themeStore.setColorblindMode('deuteranopia')
      expect(themeStore.colorblindMode).toBe('deuteranopia')
      expect(localStorage.getItem('il_registro_colorblind_mode')).toBe('deuteranopia')

      themeStore.setHighContrastMode('oled_amber')
      expect(themeStore.highContrastMode).toBe('oled_amber')
      expect(themeStore.highContrast).toBe(true)
      expect(localStorage.getItem('il_registro_contrast_mode')).toBe('oled_amber')
    })

    it('ripristina tutti i valori di default con resetAccessibility()', () => {
      const themeStore = useThemeStore()
      themeStore.toggleDsaFont(true)
      themeStore.toggleReadingRuler(true)
      themeStore.setColorblindMode('protanopia')
      themeStore.setLineHeight('loose')

      expect(themeStore.isA11yActive).toBe(true)

      themeStore.resetAccessibility()
      expect(themeStore.dsaFont).toBe(false)
      expect(themeStore.readingRuler).toBe(false)
      expect(themeStore.colorblindMode).toBe('none')
      expect(themeStore.lineHeight).toBe('normal')
      expect(themeStore.isA11yActive).toBe(false)
    })
  })

  describe('2. Text-to-Speech (TTS) Composable & Button', () => {
    it('fornisce metodi di sintesi vocale se supportato', () => {
      // Mock window.speechSynthesis
      const mockSpeak = vi.fn()
      const mockCancel = vi.fn()
      const mockPause = vi.fn()
      const mockResume = vi.fn()

      global.window.speechSynthesis = {
        speak: mockSpeak,
        cancel: mockCancel,
        pause: mockPause,
        resume: mockResume,
        speaking: false,
        paused: false,
        getVoices: () => [{ name: 'Italian Female', lang: 'it-IT' }]
      }
      global.SpeechSynthesisUtterance = function (text) {
        this.text = text
        this.rate = 1
        this.pitch = 1
        this.lang = 'it-IT'
      }

      const { speak, stop, isSupported } = useSpeechSynthesis()
      expect(isSupported.value).toBe(true)

      speak('Circolare n. 42')
      expect(mockCancel).toHaveBeenCalled()
      expect(mockSpeak).toHaveBeenCalled()

      stop()
      expect(mockCancel).toHaveBeenCalled()
    })

    it('renderizza il componente TextToSpeechButton correttamente', () => {
      const wrapper = mount(TextToSpeechButton, {
        props: {
          text: 'Compito per casa: capitolo 4',
          color: 'primary'
        },
        global: {
          stubs: {
            'q-btn': {
              template: '<button class="q-btn" @click="$emit(\'click\', $event)"><slot /></button>'
            },
            'q-tooltip': true,
            'q-icon': true
          }
        }
      })

      expect(wrapper.find('.tts-button').exists()).toBe(true)
    })
  })

  describe('3. Reading Ruler & Keyboard Shortcuts', () => {
    it('renderizza il righello solo se attivo nello store', async () => {
      const themeStore = useThemeStore()
      themeStore.readingRuler = false

      const wrapper = mount(ReadingRuler)
      expect(wrapper.find('.reading-ruler-container').exists()).toBe(false)

      themeStore.readingRuler = true
      await wrapper.vm.$nextTick()
      expect(wrapper.find('.reading-ruler-container').exists()).toBe(true)
    })

    it('renderizza il dialogo delle scorciatoie con le categorie', () => {
      const themeStore = useThemeStore()
      themeStore.keyboardShortcutsHelpOpen = true

      const wrapper = mount(KeyboardShortcutsDialog, {
        global: {
          mocks: {
            $t: (key) => key === 'a11y.shortcutsTitle' ? 'Scorciatoie da Tastiera & Accessibilità' : key
          },
          stubs: {
            'q-dialog': {
              template: '<div><slot /></div>'
            },
            'q-card': { template: '<div><slot /></div>' },
            'q-card-section': { template: '<div><slot /></div>' },
            'q-card-actions': { template: '<div><slot /></div>' },
            'q-input': { template: '<input />' },
            'q-list': { template: '<div><slot /></div>' },
            'q-item': { template: '<div><slot /></div>' },
            'q-item-section': { template: '<div><slot /></div>' },
            'q-item-label': { template: '<div><slot /></div>' },
            'q-icon': true,
            'q-btn': true,
            'q-separator': true
          }
        }
      })

      expect(wrapper.text()).toContain('Navigazione Globale')
      expect(wrapper.text()).toContain('Accessibilità')
    })
  })

  describe('4. Dichiarazione di Accessibilità AgID & Invio Segnalazioni', () => {
    it('renderizza la dichiarazione con i 4 pilastri AgID e il form di feedback', () => {
      const wrapper = mount(AccessibilityStatement, {
        global: {
          mocks: {
            $t: (key) => key === 'a11y.statementTitle' ? 'Dichiarazione di Accessibilità' : key
          },
          stubs: {
            'q-page': { template: '<div><slot /></div>' },
            'q-card': { template: '<div><slot /></div>' },
            'q-avatar': { template: '<div><slot /></div>' },
            'q-icon': true,
            'q-banner': { template: '<div><slot /><slot name="action" /></div>' },
            'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
            'q-input': true,
            'q-select': true,
            'q-btn': true
          }
        }
      })

      expect(wrapper.text()).toContain('Dichiarazione di Accessibilità')
      expect(wrapper.text()).toContain('WCAG 2.2')
      expect(wrapper.text()).toContain('Difensore Civico per il Digitale')
    })

    it('invia la segnalazione di accessibilità tramite accessibilityService', async () => {
      const { accessibilityService } = await import('@/services/accessibilityService')
      const mockPost = vi.spyOn(accessibilityService, 'submitFeedback').mockResolvedValue({
        id: '123-abc',
        protocol_number: 'A11Y-2026-0827-1001',
        message: 'Segnalazione registrata con successo'
      })

      const res = await accessibilityService.submitFeedback({
        name: 'Maria Bianchi',
        email: 'maria.bianchi@scuola.it',
        barrierType: 'contrast',
        description: 'Contrasto basso sui pulsanti di azione'
      })

      expect(mockPost).toHaveBeenCalledWith({
        name: 'Maria Bianchi',
        email: 'maria.bianchi@scuola.it',
        barrierType: 'contrast',
        description: 'Contrasto basso sui pulsanti di azione'
      })
      expect(res.protocol_number).toBe('A11Y-2026-0827-1001')
    })
  })
})
