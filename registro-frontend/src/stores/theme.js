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
        highContrast: localStorage.getItem('il_registro_high_contrast') === 'true'
    }),
    getters: {
        activeThemeObj: (state) => {
            return THEMES.find(t => t.id === state.currentTheme) || THEMES[0]
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
            this.applyAccessibility()
        },
        applyAccessibility() {
            if (typeof document !== 'undefined') {
                document.body.classList.toggle('dsa-font-active', this.dsaFont || this.fontFamily === 'opendyslexic')
                document.body.classList.toggle('high-contrast-active', this.highContrast)
                
                // Remove previous font family classes
                document.body.classList.remove('font-family-lexend', 'font-family-fredoka', 'font-family-roboto')
                if (this.fontFamily === 'lexend') document.body.classList.add('font-family-lexend')
                if (this.fontFamily === 'fredoka') document.body.classList.add('font-family-fredoka')
                if (this.fontFamily === 'roboto') document.body.classList.add('font-family-roboto')

                // Remove previous font size classes
                document.body.classList.remove('font-size-large', 'font-size-xlarge')
                if (this.fontSize === 'large') document.body.classList.add('font-size-large')
                if (this.fontSize === 'xlarge') document.body.classList.add('font-size-xlarge')
            }
        },
        initTheme() {
            this.applyTheme(this.currentTheme)
        }
    }
})
