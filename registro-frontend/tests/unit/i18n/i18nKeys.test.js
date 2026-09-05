import { describe, it, expect } from 'vitest'
import itIT from '@/i18n/it-IT'
import enUS from '@/i18n/en-US'
import deDE from '@/i18n/de-DE'
import frFR from '@/i18n/fr-FR'
import esES from '@/i18n/es-ES'
import roRO from '@/i18n/ro-RO'
import sqAL from '@/i18n/sq-AL'
import ruRU from '@/i18n/ru-RU'
import ukUA from '@/i18n/uk-UA'
import arSA from '@/i18n/ar-SA'
import zhCN from '@/i18n/zh-CN'

const allMessages = {
  'it-IT': itIT,
  'en-US': enUS,
  'de-DE': deDE,
  'fr-FR': frFR,
  'es-ES': esES,
  'ro-RO': roRO,
  'sq-AL': sqAL,
  'ru-RU': ruRU,
  'uk-UA': ukUA,
  'ar-SA': arSA,
  'zh-CN': zhCN
}

describe('i18n Dictionary Integrity across all 11 Supported Languages', () => {
  const supportedLocales = ['it-IT', 'en-US', 'de-DE', 'fr-FR', 'es-ES', 'ro-RO', 'sq-AL', 'ru-RU', 'uk-UA', 'ar-SA', 'zh-CN']
  const requiredSections = ['roles', 'layout', 'categories', 'nav', 'notifications', 'help']

  supportedLocales.forEach((locale) => {
    describe(`Locale: ${locale}`, () => {
      it('has message pack defined', () => {
        expect(allMessages[locale]).toBeDefined()
      })

      requiredSections.forEach((section) => {
        it(`contains required section "${section}"`, () => {
          const dict = allMessages[locale]
          expect(dict[section]).toBeDefined()
          expect(typeof dict[section]).toBe('object')
        })
      })

      it('contains all essential role translations', () => {
        const roles = allMessages[locale].roles
        expect(roles.teacher).toBeDefined()
        expect(roles.student).toBeDefined()
        expect(roles.parent).toBeDefined()
        expect(roles.admin).toBeDefined()
        expect(roles.secretary).toBeDefined()
        expect(roles.coordinator).toBeDefined()
      })

      it('contains layout navigation and a11y labels', () => {
        const layout = allMessages[locale].layout
        expect(layout.skipToContent).toBeDefined()
        expect(layout.mainNav).toBeDefined()
        expect(layout.dsaFontDesc).toBeDefined()
        expect(layout.highContrast).toBeDefined()
      })
    })
  })
})
