import { defineStore } from 'pinia'
import { setCssVar } from 'quasar'

export const THEMES = [
    {
        id: 'indigo',
        name: 'Modern Indigo',
        description: 'Stile classico accademico elegante e bilanciato',
        recommendedRole: 'Tutti i Ruoli',
        badgeColor: 'indigo',
        primary: '#4f46e5',
        gradientStart: '#4f46e5',
        gradientEnd: '#7c3aed',
        icon: 'palette'
    },
    {
        id: 'arcade',
        name: 'Arcade Gamer 🎮',
        description: 'Layout gaming 3D tactile con font Fredoka e pop animato',
        recommendedRole: 'Ragazzi / Gaming',
        badgeColor: 'deep-purple',
        primary: '#8b5cf6',
        gradientStart: '#8b5cf6',
        gradientEnd: '#10b981',
        icon: 'sports_esports'
    },
    {
        id: 'cosmic',
        name: 'Cosmic Explorer 🚀',
        description: 'Vetro galattico, bordi neon e font futuristico Poppins',
        recommendedRole: 'Ragazzi / Space',
        badgeColor: 'cyan',
        primary: '#06b6d4',
        gradientStart: '#06b6d4',
        gradientEnd: '#6366f1',
        icon: 'rocket_launch'
    },
    {
        id: 'bubblepop',
        name: 'Candy Bubble Pop 🍬',
        description: 'Super curve 28px, colori caramella ed effetto rimbalzo',
        recommendedRole: 'Ragazzi / Fun',
        badgeColor: 'pink',
        primary: '#ec4899',
        gradientStart: '#ec4899',
        gradientEnd: '#f97316',
        icon: 'attractions'
    },
    {
        id: 'emerald',
        name: 'Emerald School',
        description: 'Tonalità fresche ed energetiche verde/smeraldo',
        recommendedRole: 'Studenti',
        badgeColor: 'emerald',
        primary: '#059669',
        gradientStart: '#059669',
        gradientEnd: '#0d9488',
        icon: 'eco'
    },
    {
        id: 'sunset',
        name: 'Sunset Velvet',
        description: 'Toni caldi ed avvolgenti viola e rosa tramonto',
        recommendedRole: 'Genitori & Docenti',
        badgeColor: 'purple',
        primary: '#7c3aed',
        gradientStart: '#7c3aed',
        gradientEnd: '#e11d48',
        icon: 'auto_awesome'
    },
    {
        id: 'amber',
        name: 'Amber Horizon',
        description: 'Stile editoriale caldo e rilassante per la lettura',
        recommendedRole: 'Docenti',
        badgeColor: 'amber',
        primary: '#d97706',
        gradientStart: '#d97706',
        gradientEnd: '#ea580c',
        icon: 'wb_sunny'
    },
    {
        id: 'cyber',
        name: 'Cyber Dark',
        description: 'Design scuro ad alto contrasto con accenti neon',
        recommendedRole: 'Tech & Night',
        badgeColor: 'cyan',
        primary: '#06b6d4',
        gradientStart: '#06b6d4',
        gradientEnd: '#6366f1',
        icon: 'dark_mode'
    }
]

export const useThemeStore = defineStore('theme', {
    state: () => ({
        currentTheme: localStorage.getItem('il_registro_theme') || 'indigo',
        dsaFont: localStorage.getItem('il_registro_dsa_font') === 'true',
        fontFamily: localStorage.getItem('il_registro_font_family') || 'default',
        fontSize: localStorage.getItem('il_registro_font_size') || 'normal',
        highContrast: localStorage.getItem('il_registro_high_contrast') === 'true',
        
        // Advanced Accessibility State
        ttsEnabled: localStorage.getItem('il_registro_tts_enabled') === 'true',
        ttsRate: parseFloat(localStorage.getItem('il_registro_tts_rate') || '1.0'),
        ttsPitch: parseFloat(localStorage.getItem('il_registro_tts_pitch') || '1.0'),
        
        readingRuler: localStorage.getItem('il_registro_reading_ruler') === 'true',
        readingRulerHeight: parseInt(localStorage.getItem('il_registro_ruler_height') || '40', 10),
        readingRulerOpacity: parseFloat(localStorage.getItem('il_registro_ruler_opacity') || '0.35'),
        
        lineHeight: localStorage.getItem('il_registro_line_height') || 'normal', // 'normal' | 'relaxed' | 'loose'
        letterSpacing: localStorage.getItem('il_registro_letter_spacing') || 'normal', // 'normal' | 'wide' | 'wider'
        wordSpacing: localStorage.getItem('il_registro_word_spacing') || 'normal', // 'normal' | 'wide' | 'wider'
        
        colorblindMode: localStorage.getItem('il_registro_colorblind_mode') || 'none', // 'none' | 'protanopia' | 'deuteranopia' | 'tritanopia' | 'monochrome'
        highContrastMode: localStorage.getItem('il_registro_contrast_mode') || 'none', // 'none' | 'high_contrast' | 'oled_amber' | 'oled_green' | 'inverted'
        focusHighlight: localStorage.getItem('il_registro_focus_highlight') !== 'false', // default true
        
        keyboardShortcutsHelpOpen: false
    }),
    getters: {
        activeThemeObj: (state) => {
            return THEMES.find(t => t.id === state.currentTheme) || THEMES[0]
        },
        isA11yActive: (state) => {
            return state.dsaFont ||
                state.highContrast ||
                state.highContrastMode !== 'none' ||
                state.colorblindMode !== 'none' ||
                state.readingRuler ||
                state.fontSize !== 'normal' ||
                state.lineHeight !== 'normal' ||
                state.letterSpacing !== 'normal'
        }
    },
    actions: {
        setTheme(themeId) {
            const targetTheme = THEMES.find(t => t.id === themeId) ? themeId : 'indigo'
            this.currentTheme = targetTheme
            localStorage.setItem('il_registro_theme', targetTheme)
            this.applyTheme(targetTheme)
        },
        setFontFamily(family) {
            this.fontFamily = family || 'default'
            localStorage.setItem('il_registro_font_family', this.fontFamily)
            if (this.fontFamily === 'opendyslexic') {
                this.dsaFont = true
                localStorage.setItem('il_registro_dsa_font', 'true')
            } else {
                this.dsaFont = false
                localStorage.setItem('il_registro_dsa_font', 'false')
            }
            this.applyAccessibility()
        },
        setFontSize(size) {
            this.fontSize = size || 'normal'
            localStorage.setItem('il_registro_font_size', this.fontSize)
            this.applyAccessibility()
        },
        applyTheme(themeId = this.currentTheme) {
            const themeObj = THEMES.find(t => t.id === themeId) || THEMES[0]
            if (typeof document !== 'undefined') {
                document.documentElement.setAttribute('data-theme', themeObj.id)
            }
            try {
                setCssVar('primary', themeObj.primary)
            } catch (e) {
                // Ignore if Quasar context is not available during unit tests
            }
            this.applyAccessibility()
        },
        toggleDsaFont(enabled) {
            this.dsaFont = enabled !== undefined ? enabled : !this.dsaFont
            localStorage.setItem('il_registro_dsa_font', this.dsaFont)
            this.fontFamily = this.dsaFont ? 'opendyslexic' : 'default'
            localStorage.setItem('il_registro_font_family', this.fontFamily)
            this.applyAccessibility()
        },
        toggleHighContrast(enabled) {
            this.highContrast = enabled !== undefined ? enabled : !this.highContrast
            localStorage.setItem('il_registro_high_contrast', this.highContrast)
            if (this.highContrast && this.highContrastMode === 'none') {
                this.highContrastMode = 'high_contrast'
                localStorage.setItem('il_registro_contrast_mode', 'high_contrast')
            } else if (!this.highContrast) {
                this.highContrastMode = 'none'
                localStorage.setItem('il_registro_contrast_mode', 'none')
            }
            this.applyAccessibility()
        },
        setHighContrastMode(mode) {
            this.highContrastMode = mode || 'none'
            this.highContrast = mode !== 'none'
            localStorage.setItem('il_registro_contrast_mode', this.highContrastMode)
            localStorage.setItem('il_registro_high_contrast', this.highContrast)
            this.applyAccessibility()
        },
        setColorblindMode(mode) {
            this.colorblindMode = mode || 'none'
            localStorage.setItem('il_registro_colorblind_mode', this.colorblindMode)
            this.applyAccessibility()
        },
        toggleReadingRuler(enabled) {
            this.readingRuler = enabled !== undefined ? enabled : !this.readingRuler
            localStorage.setItem('il_registro_reading_ruler', this.readingRuler)
            this.applyAccessibility()
        },
        setReadingRulerHeight(height) {
            this.readingRulerHeight = height || 40
            localStorage.setItem('il_registro_ruler_height', this.readingRulerHeight)
        },
        setReadingRulerOpacity(opacity) {
            this.readingRulerOpacity = opacity !== undefined ? opacity : 0.35
            localStorage.setItem('il_registro_ruler_opacity', this.readingRulerOpacity)
        },
        setLineHeight(mode) {
            this.lineHeight = mode || 'normal'
            localStorage.setItem('il_registro_line_height', this.lineHeight)
            this.applyAccessibility()
        },
        setLetterSpacing(mode) {
            this.letterSpacing = mode || 'normal'
            localStorage.setItem('il_registro_letter_spacing', this.letterSpacing)
            this.applyAccessibility()
        },
        setWordSpacing(mode) {
            this.wordSpacing = mode || 'normal'
            localStorage.setItem('il_registro_word_spacing', this.wordSpacing)
            this.applyAccessibility()
        },
        toggleTts(enabled) {
            this.ttsEnabled = enabled !== undefined ? enabled : !this.ttsEnabled
            localStorage.setItem('il_registro_tts_enabled', this.ttsEnabled)
        },
        setTtsRate(rate) {
            this.ttsRate = rate || 1.0
            localStorage.setItem('il_registro_tts_rate', this.ttsRate)
        },
        setTtsPitch(pitch) {
            this.ttsPitch = pitch || 1.0
            localStorage.setItem('il_registro_tts_pitch', this.ttsPitch)
        },
        toggleFocusHighlight(enabled) {
            this.focusHighlight = enabled !== undefined ? enabled : !this.focusHighlight
            localStorage.setItem('il_registro_focus_highlight', this.focusHighlight)
            this.applyAccessibility()
        },
        toggleKeyboardShortcutsHelp(open) {
            this.keyboardShortcutsHelpOpen = open !== undefined ? open : !this.keyboardShortcutsHelpOpen
        },
        resetAccessibility() {
            this.dsaFont = false
            this.fontFamily = 'default'
            this.fontSize = 'normal'
            this.highContrast = false
            this.highContrastMode = 'none'
            this.colorblindMode = 'none'
            this.readingRuler = false
            this.lineHeight = 'normal'
            this.letterSpacing = 'normal'
            this.wordSpacing = 'normal'
            this.ttsEnabled = false
            this.focusHighlight = true

            localStorage.removeItem('il_registro_dsa_font')
            localStorage.removeItem('il_registro_font_family')
            localStorage.removeItem('il_registro_font_size')
            localStorage.removeItem('il_registro_high_contrast')
            localStorage.removeItem('il_registro_contrast_mode')
            localStorage.removeItem('il_registro_colorblind_mode')
            localStorage.removeItem('il_registro_reading_ruler')
            localStorage.removeItem('il_registro_line_height')
            localStorage.removeItem('il_registro_letter_spacing')
            localStorage.removeItem('il_registro_word_spacing')
            localStorage.removeItem('il_registro_tts_enabled')
            localStorage.removeItem('il_registro_focus_highlight')

            this.applyAccessibility()
        },
        applyAccessibility() {
            if (typeof document !== 'undefined') {
                const body = document.body
                
                // DSA Font & High Contrast
                body.classList.toggle('dsa-font-active', this.dsaFont || this.fontFamily === 'opendyslexic')
                body.classList.toggle('high-contrast-active', this.highContrast || this.highContrastMode === 'high_contrast')
                
                // Font Family
                body.classList.remove('font-family-lexend', 'font-family-fredoka', 'font-family-roboto')
                if (this.fontFamily === 'lexend') body.classList.add('font-family-lexend')
                if (this.fontFamily === 'fredoka') body.classList.add('font-family-fredoka')
                if (this.fontFamily === 'roboto') body.classList.add('font-family-roboto')

                // Font Size
                body.classList.remove('font-size-large', 'font-size-xlarge')
                if (this.fontSize === 'large') body.classList.add('font-size-large')
                if (this.fontSize === 'xlarge') body.classList.add('font-size-xlarge')

                // Text Spacing (WCAG 1.4.12)
                body.classList.remove('line-height-relaxed', 'line-height-loose')
                if (this.lineHeight === 'relaxed') body.classList.add('line-height-relaxed')
                if (this.lineHeight === 'loose') body.classList.add('line-height-loose')

                body.classList.remove('letter-spacing-wide', 'letter-spacing-wider')
                if (this.letterSpacing === 'wide') body.classList.add('letter-spacing-wide')
                if (this.letterSpacing === 'wider') body.classList.add('letter-spacing-wider')

                body.classList.remove('word-spacing-wide', 'word-spacing-wider')
                if (this.wordSpacing === 'wide') body.classList.add('word-spacing-wide')
                if (this.wordSpacing === 'wider') body.classList.add('word-spacing-wider')

                // Colorblind Modes
                body.classList.remove(
                    'colorblind-protanopia',
                    'colorblind-deuteranopia',
                    'colorblind-tritanopia',
                    'colorblind-monochrome'
                )
                if (this.colorblindMode && this.colorblindMode !== 'none') {
                    body.classList.add(`colorblind-${this.colorblindMode}`)
                }

                // OLED / Advanced Contrast Modes
                body.classList.remove('contrast-oled-amber', 'contrast-oled-green', 'contrast-inverted')
                if (this.highContrastMode === 'oled_amber') body.classList.add('contrast-oled-amber')
                if (this.highContrastMode === 'oled_green') body.classList.add('contrast-oled-green')
                if (this.highContrastMode === 'inverted') body.classList.add('contrast-inverted')

                // Focus Highlight
                body.classList.toggle('focus-visible-high', this.focusHighlight)
            }
        },
        initTheme() {
            this.applyTheme(this.currentTheme)
        }
    }
})
