import { describe, it, expect, beforeEach } from 'vitest'
import { ref } from 'vue'
import itMessages from '@/i18n/it-IT/index.js'
import enMessages from '@/i18n/en-US/index.js'
import deMessages from '@/i18n/de-DE/index.js'
import frMessages from '@/i18n/fr-FR/index.js'
import esMessages from '@/i18n/es-ES/index.js'
import zhMessages from '@/i18n/zh-CN/index.js'

describe('Teacher Settings & Accessibility Localization Suite', () => {
  let currentLocale
  let t

  beforeEach(() => {
    currentLocale = ref('it-IT')

    function getNestedValue(obj, path) {
      if (!obj || !path) return null
      return path.split('.').reduce((prev, curr) => (prev && prev[curr] !== undefined ? prev[curr] : null), obj)
    }

    t = (key) => {
      let messages = itMessages
      if (currentLocale.value === 'en-US') messages = enMessages
      else if (currentLocale.value === 'de-DE') messages = deMessages
      else if (currentLocale.value === 'fr-FR') messages = frMessages
      else if (currentLocale.value === 'es-ES') messages = esMessages
      else if (currentLocale.value === 'zh-CN') messages = zhMessages
      return getNestedValue(messages, key) || key
    }
  })

  describe('1. Settings General Tab (Language & Date Format)', () => {
    it('provides localized labels in it-IT', () => {
      currentLocale.value = 'it-IT'
      expect(t('settingsPage.langAndDateFormat')).toBe('Lingua & Formato Data')
      expect(t('settingsPage.appLanguage')).toBe("Lingua dell'Applicazione")
      expect(t('settingsPage.defaultDateFormat')).toBe('Formato Data Predefinito')
      expect(t('settingsPage.firstDayOfWeek')).toBe('Primo Giorno della Settimana')
      expect(t('settingsPage.monday')).toBe('Lunedì')
      expect(t('settingsPage.sunday')).toBe('Domenica')
      expect(t('settingsPage.timeFormatAndTimezone')).toBe('Fuso Orario & Formato Ora')
    })

    it('provides localized labels in en-US', () => {
      currentLocale.value = 'en-US'
      expect(t('settingsPage.langAndDateFormat')).toBe('Language & Date Format')
      expect(t('settingsPage.appLanguage')).toBe('Application Language')
      expect(t('settingsPage.defaultDateFormat')).toBe('Default Date Format')
      expect(t('settingsPage.firstDayOfWeek')).toBe('First Day of the Week')
      expect(t('settingsPage.monday')).toBe('Monday')
      expect(t('settingsPage.sunday')).toBe('Sunday')
      expect(t('settingsPage.timeFormatAndTimezone')).toBe('Timezone & Time Format')
    })

    it('provides localized labels in zh-CN and de-DE', () => {
      currentLocale.value = 'zh-CN'
      expect(t('settingsPage.langAndDateFormat')).toBe('语言与日期格式')
      expect(t('settingsPage.firstDayOfWeek')).toBe('一周的第一天')
      expect(t('settingsPage.monday')).toBe('星期一')

      currentLocale.value = 'de-DE'
      expect(t('settingsPage.langAndDateFormat')).toBe('Sprache & Datumsformat')
      expect(t('settingsPage.monday')).toBe('Montag')
    })
  })

  describe('2. Security, Notifications, Register Customization & Digital PIN', () => {
    it('localizes security and notifications headers and options in en-US', () => {
      currentLocale.value = 'en-US'
      expect(t('settingsPage.accountSecurityAndPwd')).toBe('Account Security & Change Password')
      expect(t('settingsPage.twoFactorTitle')).toBe('Two-Factor Authentication (2FA)')
      expect(t('settingsPage.teacherNotifPref')).toBe('Teacher Notification Preferences')
      expect(t('settingsPage.registerCustomizationTitle')).toBe('Register & Grade Grid Customization')
      expect(t('settingsPage.digitalSignatureTitle')).toBe('Digital Signature & Quick Lesson PIN')
      expect(t('settingsPage.quickPinTitle')).toBe('Quick Register PIN (4 Digits)')
    })
  })

  describe('3. Accessibility Settings Panel (DSA, Ruler, Spacing, TTS, Colorblind)', () => {
    it('localizes all accessibility tool sections and options in it-IT and en-US', () => {
      currentLocale.value = 'it-IT'
      expect(t('a11y.dsaFontTitle')).toBe('Font OpenDyslexic (DSA)')
      expect(t('a11y.readingRulerTitle')).toBe('Righello di Lettura (Focus Mask)')
      expect(t('a11y.textSpacingTitle')).toBe('Spaziatura Testo & Interlinea')
      expect(t('a11y.ttsNativeTitle')).toBe('Sintesi Vocale Nativa (TTS)')
      expect(t('a11y.colorblindContrastTitle')).toBe('Contrasto Visivo & Filtri per Daltonismo (Color Blindness)')

      currentLocale.value = 'en-US'
      expect(t('a11y.dsaFontTitle')).toBe('OpenDyslexic Font (DSA/Dyslexia)')
      expect(t('a11y.readingRulerTitle')).toBe('Reading Ruler (Focus Mask)')
      expect(t('a11y.textSpacingTitle')).toBe('Text Spacing & Line Height')
      expect(t('a11y.ttsNativeTitle')).toBe('Native Text-to-Speech (TTS)')
      expect(t('a11y.colorblindContrastTitle')).toBe('Visual Contrast & Color Blindness Filters')
      expect(t('a11y.fontOpenDyslexic')).toBe('OpenDyslexic (High Readability DSA)')
    })
  })
})
