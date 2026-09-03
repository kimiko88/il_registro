import { describe, it, expect, beforeEach } from 'vitest'
import { ref } from 'vue'
import { THEMES, getThemeName, getThemeDescription, getThemeRole } from '@/stores/theme'
import itMessages from '@/i18n/it-IT/index.js'
import enMessages from '@/i18n/en-US/index.js'

describe('Sidebar Menu Items & Themes Localization', () => {
  let currentLocale
  let t
  let te

  beforeEach(() => {
    currentLocale = ref('it-IT')

    function getNestedValue(obj, path) {
      if (!obj || !path) return null
      return path.split('.').reduce((prev, curr) => (prev && prev[curr] !== undefined ? prev[curr] : null), obj)
    }

    t = (key) => {
      const messages = currentLocale.value === 'en-US' ? enMessages : itMessages
      return getNestedValue(messages, key) || key
    }

    te = (key) => {
      const messages = currentLocale.value === 'en-US' ? enMessages : itMessages
      return !!getNestedValue(messages, key)
    }
  })

  describe('1. Sidebar Menu Items Localization', () => {
    it('translates Credito Scolastico, Corsi Recupero & PAI, Registro Sostegno & PEI, Ricevimento Generale in it-IT', () => {
      currentLocale.value = 'it-IT'
      expect(t('nav.credits')).toBe('Credito Scolastico')
      expect(t('nav.recovery')).toBe('Corsi Recupero & PAI')
      expect(t('nav.supportRegister')).toBe('Registro Sostegno & PEI')
      expect(t('nav.generalMeetings')).toBe('Ricevimento Generale')
      expect(t('nav.sidi')).toBe('Flussi SIDI')
    })

    it('translates Credito Scolastico, Corsi Recupero & PAI, Registro Sostegno & PEI, Ricevimento Generale in en-US', () => {
      currentLocale.value = 'en-US'
      expect(t('nav.credits')).toBe('School Credits')
      expect(t('nav.recovery')).toBe('Recovery Courses & PAI')
      expect(t('nav.supportRegister')).toBe('Support Register & IEP')
      expect(t('nav.generalMeetings')).toBe('General Parent Meetings')
      expect(t('nav.sidi')).toBe('SIDI Data Flows')
    })
  })

  describe('2. Themes Names, Descriptions & Roles Localization', () => {
    it('localizes all themes in Italian', () => {
      currentLocale.value = 'it-IT'
      for (const theme of THEMES) {
        const name = getThemeName(theme, t, te)
        const desc = getThemeDescription(theme, t, te)
        const role = getThemeRole(theme, t, te)

        expect(name).toBeTruthy()
        expect(desc).toBeTruthy()
        expect(role).toBeTruthy()
      }

      expect(getThemeName(THEMES[0], t, te)).toBe('Modern Indigo')
      expect(getThemeDescription(THEMES[0], t, te)).toBe('Stile classico accademico elegante e bilanciato')
      expect(getThemeRole(THEMES[0], t, te)).toBe('Tutti i Ruoli')
    })

    it('localizes all themes in English', () => {
      currentLocale.value = 'en-US'
      expect(getThemeName(THEMES[0], t, te)).toBe('Modern Indigo')
      expect(getThemeDescription(THEMES[0], t, te)).toBe('Classic academic style, elegant and balanced')
      expect(getThemeRole(THEMES[0], t, te)).toBe('All Roles')

      const arcadeTheme = THEMES.find(th => th.id === 'arcade')
      expect(getThemeName(arcadeTheme, t, te)).toBe('Arcade Gamer 🎮')
      expect(getThemeDescription(arcadeTheme, t, te)).toBe('Tactile 3D gaming layout with Fredoka font and lively animations')
      expect(getThemeRole(arcadeTheme, t, te)).toBe('Students / Gaming')
    })
  })

  describe('3. Dark and Light Mode Descriptions Localization', () => {
    it('localizes dark and light mode descriptions in it-IT and en-US', () => {
      currentLocale.value = 'it-IT'
      expect(t('layout.darkModeDesc')).toContain('OLED Dark')
      expect(t('layout.lightModeDesc')).toContain('Standard Light')

      currentLocale.value = 'en-US'
      expect(t('layout.darkModeDesc')).toBe('OLED Dark • High contrast to rest your eyes')
      expect(t('layout.lightModeDesc')).toBe('Standard Light • Natural brightness and clarity')
    })
  })
})
