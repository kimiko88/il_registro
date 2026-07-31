import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useThemeStore, THEMES } from '@/stores/theme'

describe('Theme Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
    })

    it('initializes with default theme indigo', () => {
        const store = useThemeStore()
        expect(store.currentTheme).toBe('indigo')
        expect(store.activeThemeObj.id).toBe('indigo')
    })

    it('switches theme and persists in localStorage', () => {
        const store = useThemeStore()
        store.setTheme('emerald')

        expect(store.currentTheme).toBe('emerald')
        expect(store.activeThemeObj.name).toBe('Emerald School')
        expect(localStorage.getItem('registrov2_theme')).toBe('emerald')
        expect(document.documentElement.getAttribute('data-theme')).toBe('emerald')
    })

    it('falls back to indigo for invalid theme IDs', () => {
        const store = useThemeStore()
        store.setTheme('non_existent_theme')

        expect(store.currentTheme).toBe('indigo')
        expect(localStorage.getItem('registrov2_theme')).toBe('indigo')
    })

    it('contains themes for students, teachers, and parents', () => {
        expect(THEMES).toHaveLength(5)
        const roles = THEMES.map(t => t.recommendedRole)
        expect(roles).toContain('Studenti')
        expect(roles).toContain('Docenti')
        expect(roles).toContain('Genitori & Docenti')
    })
})
