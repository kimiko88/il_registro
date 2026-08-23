import { describe, it, expect } from 'vitest'
import { messages } from '@/i18n'

describe('i18n Dictionary Integrity across all 11 Supported Languages', () => {
  const supportedLocales = ['it-IT', 'en-US', 'de-DE', 'fr-FR', 'es-ES', 'ro-RO', 'sq-AL', 'ru-RU', 'uk-UA', 'ar-SA', 'zh-CN']
  const requiredSections = ['roles', 'layout', 'categories', 'nav', 'notifications', 'help']

  supportedLocales.forEach((locale) => {
    describe(`Locale: ${locale}`, () => {
      it('has message pack defined', () => {
        expect(messages[locale]).toBeDefined()
      })

      requiredSections.forEach((section) => {
        it(`contains required section "${section}"`, () => {
          const dict = messages[locale]
          expect(dict[section]).toBeDefined()
          expect(typeof dict[section]).toBe('object')
        })
      })

      it('contains all essential role translations', () => {
        const roles = messages[locale].roles
        expect(roles.teacher).toBeDefined()
        expect(roles.student).toBeDefined()
        expect(roles.parent).toBeDefined()
        expect(roles.admin).toBeDefined()
        expect(roles.secretary).toBeDefined()
        expect(roles.coordinator).toBeDefined()
      })

      it('contains layout navigation and a11y labels', () => {
        const layout = messages[locale].layout
        expect(layout.skipToContent).toBeDefined()
        expect(layout.mainNav).toBeDefined()
        expect(layout.dsaFontDesc).toBeDefined()
        expect(layout.highContrast).toBeDefined()
      })
    })
  })
})
