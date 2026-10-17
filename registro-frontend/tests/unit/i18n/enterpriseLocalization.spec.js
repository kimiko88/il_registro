import { describe, it, expect } from 'vitest'

const supportedLocales = [
  'it-IT',
  'en-US',
  'de-DE',
  'fr-FR',
  'es-ES',
  'ro-RO',
  'sq-AL',
  'ru-RU',
  'uk-UA',
  'ar-SA',
  'zh-CN'
]

const requiredSubkeys = [
  'header.badgeCompliance',
  'header.badgeAgid',
  'header.title',
  'header.subtitle',
  'header.viewRoleMatrix',
  'header.hideRoleMatrix',
  'matrix.title',
  'matrix.subtitle',
  'matrix.colFeature',
  'matrix.colGovernance',
  'matrix.colOperations',
  'matrix.colEndUsers',
  'matrix.colNorm',
  'governance.authorizedRoles',
  'governance.legalRef',
  'roles.ds',
  'roles.dsga',
  'roles.secretary',
  'roles.teacher',
  'roles.parent',
  'roles.student',
  'roles.dpo',
  'roles.psychologist',
  'tabs.pagopa',
  'tabs.interpelli',
  'tabs.albo',
  'tabs.feq',
  'tabs.maturita',
  'tabs.meals',
  'tabs.inventory',
  'tabs.privacy',
  'tabs.sidi',
  'tabs.psychology',
  'pagopa.cardTitleNotice',
  'pagopa.generateNoticeBtn',
  'interpelli.cardTitlePublish',
  'interpelli.publishBtn',
  'albo.cardTitleAlbo',
  'albo.publishBtn',
  'feq.cardTitleCsc',
  'feq.signBatchBtn',
  'maturita.cardTitleCredits',
  'maturita.calcBtn',
  'meals.cardTitleRollCall',
  'meals.submitRollCallBtn',
  'inventory.cardTitleAsset',
  'inventory.createAssetBtn',
  'privacy.cardTitleTreatment',
  'privacy.saveTreatmentBtn',
  'sidi.cardTitleCert',
  'sidi.syncStudentsBtn',
  'psychology.cardTitleBooking',
  'psychology.bookBtn'
]

function getNestedValue(obj, path) {
  return path.split('.').reduce((acc, part) => (acc && acc[part] !== undefined ? acc[part] : undefined), obj)
}

describe('Enterprise Features i18n Localization Suite', () => {
  supportedLocales.forEach((locale) => {
    describe(`Locale: ${locale}`, () => {
      it(`loads ${locale} dictionary and contains all mandatory enterprise keys`, async () => {
        const mod = await import(`../../../src/i18n/${locale}/index.js`)
        const dict = mod.default
        expect(dict).toBeDefined()
        expect(dict.enterprise).toBeDefined()
        expect(dict.routeTitles?.enterpriseHub).toBeDefined()

        requiredSubkeys.forEach((subkey) => {
          const val = getNestedValue(dict.enterprise, subkey)
          expect(val, `Missing [enterprise.${subkey}] in locale ${locale}`).toBeDefined()
          expect(typeof val).toBe('string')
          expect(val.trim().length).toBeGreaterThan(0)
        })
      })
    })
  })
})
