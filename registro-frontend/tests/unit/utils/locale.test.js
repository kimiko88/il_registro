import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import {
  SUPPORTED_LOCALES,
  normalizeLocale,
  getBrowserLocale,
  getSavedLocale,
  getQuasarLang,
  applyLocale
} from '@/utils/locale'

describe('locale.js — Multilingual & Accessibility Support', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  afterEach(() => {
    localStorage.clear()
  })

  describe('SUPPORTED_LOCALES', () => {
    it('supports 11 languages including Italian, English, French, German, Spanish, Romanian, Albanian, Russian, Ukrainian, Arabic, Chinese', () => {
      expect(SUPPORTED_LOCALES).toHaveLength(11)
      const values = SUPPORTED_LOCALES.map(l => l.value)
      expect(values).toContain('it-IT')
      expect(values).toContain('en-US')
      expect(values).toContain('de-DE')
      expect(values).toContain('fr-FR')
      expect(values).toContain('es-ES')
      expect(values).toContain('ro-RO')
      expect(values).toContain('sq-AL')
      expect(values).toContain('ru-RU')
      expect(values).toContain('uk-UA')
      expect(values).toContain('ar-SA')
      expect(values).toContain('zh-CN')
    })

    it('marks Arabic as RTL direction', () => {
      const arabic = SUPPORTED_LOCALES.find(l => l.value === 'ar-SA')
      expect(arabic.dir).toBe('rtl')

      const italian = SUPPORTED_LOCALES.find(l => l.value === 'it-IT')
      expect(italian.dir).toBe('ltr')

      const romanian = SUPPORTED_LOCALES.find(l => l.value === 'ro-RO')
      expect(romanian.dir).toBe('ltr')

      const albanian = SUPPORTED_LOCALES.find(l => l.value === 'sq-AL')
      expect(albanian.dir).toBe('ltr')
    })
  })

  describe('normalizeLocale', () => {
    it('normalizes full locale codes', () => {
      expect(normalizeLocale('it-IT')).toBe('it-IT')
      expect(normalizeLocale('en-US')).toBe('en-US')
      expect(normalizeLocale('ro-RO')).toBe('ro-RO')
      expect(normalizeLocale('sq-AL')).toBe('sq-AL')
      expect(normalizeLocale('ar-SA')).toBe('ar-SA')
    })

    it('normalizes short 2-letter codes and uppercase codes', () => {
      expect(normalizeLocale('it')).toBe('it-IT')
      expect(normalizeLocale('EN')).toBe('en-US')
      expect(normalizeLocale('fr')).toBe('fr-FR')
      expect(normalizeLocale('DE')).toBe('de-DE')
      expect(normalizeLocale('es')).toBe('es-ES')
      expect(normalizeLocale('ro')).toBe('ro-RO')
      expect(normalizeLocale('RO')).toBe('ro-RO')
      expect(normalizeLocale('sq')).toBe('sq-AL')
      expect(normalizeLocale('SQ')).toBe('sq-AL')
      expect(normalizeLocale('ru')).toBe('ru-RU')
      expect(normalizeLocale('uk')).toBe('uk-UA')
      expect(normalizeLocale('ar')).toBe('ar-SA')
      expect(normalizeLocale('zh')).toBe('zh-CN')
    })

    it('falls back to it-IT for null, empty or unknown languages', () => {
      expect(normalizeLocale(null)).toBe('it-IT')
      expect(normalizeLocale('')).toBe('it-IT')
      expect(normalizeLocale('unknown-locale')).toBe('it-IT')
    })
  })

  describe('getSavedLocale', () => {
    it('returns it-IT by default when nothing is saved', () => {
      expect(getSavedLocale()).toBe('it-IT')
    })

    it('retrieves and normalizes saved locale from app_language or user_locale', () => {
      localStorage.setItem('app_language', 'en-US')
      expect(getSavedLocale()).toBe('en-US')

      localStorage.removeItem('app_language')
      localStorage.setItem('user_locale', 'fr')
      expect(getSavedLocale()).toBe('fr-FR')
    })

    it('falls back to browser locale when fallbackToBrowser is true and storage is empty', () => {
      const origNav = window.navigator.language
      try {
        Object.defineProperty(window.navigator, 'language', { value: 'de-DE', configurable: true })
        expect(getSavedLocale(true)).toBe('de-DE')
      } finally {
        Object.defineProperty(window.navigator, 'language', { value: origNav, configurable: true })
      }
    })
  })

  describe('getBrowserLocale', () => {
    it('detects and normalizes navigator.language', () => {
      const origNav = window.navigator.language
      try {
        Object.defineProperty(window.navigator, 'language', { value: 'fr-FR', configurable: true })
        expect(getBrowserLocale()).toBe('fr-FR')
        Object.defineProperty(window.navigator, 'language', { value: 'es', configurable: true })
        expect(getBrowserLocale()).toBe('es-ES')
      } finally {
        Object.defineProperty(window.navigator, 'language', { value: origNav, configurable: true })
      }
    })

    it('falls back to it-IT on unknown browser language or error', () => {
      const origNav = window.navigator.language
      try {
        Object.defineProperty(window.navigator, 'language', { value: 'xx-YY', configurable: true })
        expect(getBrowserLocale()).toBe('it-IT')
      } finally {
        Object.defineProperty(window.navigator, 'language', { value: origNav, configurable: true })
      }
    })
  })

  describe('getQuasarLang', () => {
    it('returns correct Quasar language pack', () => {
      const packIt = getQuasarLang('it-IT')
      expect(packIt).toBeDefined()
      expect(packIt.isoName).toBe('it')

      const packEn = getQuasarLang('en')
      expect(packEn).toBeDefined()
      expect(packEn.isoName).toBe('en-US')
    })
  })

  describe('applyLocale', () => {
    it('applies locale, updates i18n, sets DOM attributes and persists to localStorage', () => {
      const mockI18n = {
        global: {
          locale: { value: 'it-IT' }
        }
      }
      const mockQ = {
        lang: {
          set: vi.fn()
        }
      }

      const result = applyLocale('ar-SA', mockI18n, mockQ)
      expect(result).toBe('ar-SA')
      expect(mockI18n.global.locale.value).toBe('ar-SA')
      expect(mockQ.lang.set).toHaveBeenCalled()
      expect(localStorage.getItem('app_language')).toBe('ar-SA')
      expect(localStorage.getItem('user_locale')).toBe('ar-SA')
      expect(document.documentElement.getAttribute('dir')).toBe('rtl')
      expect(document.documentElement.getAttribute('lang')).toBe('ar-SA')
    })

    it('sets LTR for western languages like English or Italian', () => {
      applyLocale('en-US')
      expect(document.documentElement.getAttribute('dir')).toBe('ltr')
      expect(document.documentElement.getAttribute('lang')).toBe('en-US')
    })
  })
})
