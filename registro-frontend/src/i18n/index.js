import { createI18n } from 'vue-i18n'
import { getSavedLocale } from '../utils/locale'
import itIT from './it-IT'
import enUS from './en-US'
import deDE from './de-DE'
import frFR from './fr-FR'
import esES from './es-ES'
import ruRU from './ru-RU'
import ukUA from './uk-UA'
import roRO from './ro-RO'
import sqAL from './sq-AL'
import arSA from './ar-SA'
import zhCN from './zh-CN'

export const messages = {
  'it-IT': itIT,
  'it': itIT,
  'en-US': enUS,
  'en': enUS,
  'de-DE': deDE,
  'de': deDE,
  'fr-FR': frFR,
  'fr': frFR,
  'es-ES': esES,
  'es': esES,
  'ro-RO': roRO,
  'ro': roRO,
  'sq-AL': sqAL,
  'sq': sqAL,
  'al': sqAL,
  'ru-RU': ruRU,
  'ru': ruRU,
  'uk-UA': ukUA,
  'uk': ukUA,
  'ar-SA': arSA,
  'ar': arSA,
  'zh-CN': zhCN,
  'zh': zhCN
}

export const i18n = createI18n({
  locale: getSavedLocale(),
  fallbackLocale: 'it-IT',
  legacy: false,
  messages
})

export default messages

