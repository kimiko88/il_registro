import { createI18n } from 'vue-i18n'
import { getSavedLocale } from '../utils/locale'
import itIT from './it-IT'
import { loadLocaleMessages, localeLoaders, loadedLanguages } from './loader'

export const messages = {
  'it-IT': itIT,
  'it': itIT
}

export { loadLocaleMessages, localeLoaders, loadedLanguages }

const savedLocale = getSavedLocale(true)

export const i18n = createI18n({
  locale: savedLocale,
  fallbackLocale: 'it-IT',
  legacy: false,
  messages
})

// Asynchronously load saved non-Italian locale at startup if needed
if (savedLocale && savedLocale !== 'it-IT' && savedLocale !== 'it') {
  loadLocaleMessages(savedLocale, i18n).catch(() => {})
}

export default messages
