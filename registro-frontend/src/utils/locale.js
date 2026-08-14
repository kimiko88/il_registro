import quasarLangIt from 'quasar/lang/it'
import quasarLangEn from 'quasar/lang/en-US'
import quasarLangDe from 'quasar/lang/de'
import quasarLangFr from 'quasar/lang/fr'
import quasarLangEs from 'quasar/lang/es'
import quasarLangRu from 'quasar/lang/ru'
import quasarLangUk from 'quasar/lang/uk'
import quasarLangAr from 'quasar/lang/ar'
import quasarLangZh from 'quasar/lang/zh-CN'
import { Quasar } from 'quasar'

export const SUPPORTED_LOCALES = [
  { label: 'Italiano', value: 'it-IT', code: 'IT', flag: '🇮🇹', icon: 'flag', dir: 'ltr' },
  { label: 'English', value: 'en-US', code: 'EN', flag: '🇬🇧', icon: 'language', dir: 'ltr' },
  { label: 'Deutsch', value: 'de-DE', code: 'DE', flag: '🇩🇪', icon: 'language', dir: 'ltr' },
  { label: 'Français', value: 'fr-FR', code: 'FR', flag: '🇫🇷', icon: 'language', dir: 'ltr' },
  { label: 'Español', value: 'es-ES', code: 'ES', flag: '🇪🇸', icon: 'language', dir: 'ltr' },
  { label: 'Русский', value: 'ru-RU', code: 'RU', flag: '🇷🇺', icon: 'language', dir: 'ltr' },
  { label: 'Українська', value: 'uk-UA', code: 'UK', flag: '🇺🇦', icon: 'language', dir: 'ltr' },
  { label: 'العربية', value: 'ar-SA', code: 'AR', flag: '🇸🇦', icon: 'language', dir: 'rtl' },
  { label: '中文 (简体)', value: 'zh-CN', code: 'ZH', flag: '🇨🇳', icon: 'language', dir: 'ltr' }
]

const QUASAR_LANG_MAP = {
  'it-IT': quasarLangIt,
  'it': quasarLangIt,
  'en-US': quasarLangEn,
  'en': quasarLangEn,
  'de-DE': quasarLangDe,
  'de': quasarLangDe,
  'fr-FR': quasarLangFr,
  'fr': quasarLangFr,
  'es-ES': quasarLangEs,
  'es': quasarLangEs,
  'ru-RU': quasarLangRu,
  'ru': quasarLangRu,
  'uk-UA': quasarLangUk,
  'uk': quasarLangUk,
  'ar-SA': quasarLangAr,
  'ar': quasarLangAr,
  'zh-CN': quasarLangZh,
  'zh': quasarLangZh
}

export function normalizeLocale(lang) {
  if (!lang) return 'it-IT'
  const found = SUPPORTED_LOCALES.find(
    l => l.value.toLowerCase() === lang.toLowerCase() ||
         l.code.toLowerCase() === lang.toLowerCase() ||
         l.value.startsWith(lang)
  )
  return found ? found.value : 'it-IT'
}

export function getSavedLocale() {
  const saved = (typeof localStorage !== 'undefined')
    ? (localStorage.getItem('app_language') || localStorage.getItem('superadmin_language') || localStorage.getItem('user_locale'))
    : null
  return normalizeLocale(saved || 'it-IT')
}

export function getQuasarLang(langCode) {
  const normalized = normalizeLocale(langCode)
  return QUASAR_LANG_MAP[normalized] || quasarLangIt
}

export function applyLocale(langCode, i18nInstance = null, $q = null) {
  const normalized = normalizeLocale(langCode)
  const isRTL = normalized === 'ar-SA' || normalized === 'ar'

  // 1. Update i18n
  if (i18nInstance) {
    if (i18nInstance.global && i18nInstance.global.locale) {
      i18nInstance.global.locale.value = normalized
    } else if (i18nInstance.locale) {
      i18nInstance.locale.value = normalized
    }
  }

  // 2. Update Quasar Language Pack
  const quasarPack = QUASAR_LANG_MAP[normalized] || quasarLangIt
  if ($q && $q.lang && typeof $q.lang.set === 'function') {
    $q.lang.set(quasarPack)
  } else if (Quasar && Quasar.lang && typeof Quasar.lang.set === 'function') {
    Quasar.lang.set(quasarPack)
  }

  // 3. Update DOM direction and lang
  if (typeof document !== 'undefined') {
    document.documentElement.setAttribute('lang', normalized)
    document.documentElement.setAttribute('dir', isRTL ? 'rtl' : 'ltr')
  }

  // 4. Persist in localStorage
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('app_language', normalized)
    localStorage.setItem('superadmin_language', normalized)
    localStorage.setItem('user_locale', normalized)
  }

  return normalized
}
