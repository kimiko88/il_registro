import quasarLangIt from 'quasar/lang/it'
import quasarLangEn from 'quasar/lang/en-US'
import quasarLangDe from 'quasar/lang/de'
import quasarLangFr from 'quasar/lang/fr'
import quasarLangEs from 'quasar/lang/es'
import quasarLangRo from 'quasar/lang/ro'
import quasarLangSq from 'quasar/lang/sq'
import quasarLangRu from 'quasar/lang/ru'
import quasarLangUk from 'quasar/lang/uk'
import quasarLangAr from 'quasar/lang/ar'
import quasarLangZh from 'quasar/lang/zh-CN'
import { Quasar } from 'quasar'
import { loadLocaleMessages } from '../i18n/loader'

export const SUPPORTED_LOCALES = [
  { label: 'Italiano', value: 'it-IT', code: 'IT', flag: '🇮🇹', icon: 'flag', dir: 'ltr' },
  { label: 'English', value: 'en-US', code: 'EN', flag: '🇬🇧', icon: 'language', dir: 'ltr' },
  { label: 'Deutsch', value: 'de-DE', code: 'DE', flag: '🇩🇪', icon: 'language', dir: 'ltr' },
  { label: 'Français', value: 'fr-FR', code: 'FR', flag: '🇫🇷', icon: 'language', dir: 'ltr' },
  { label: 'Español', value: 'es-ES', code: 'ES', flag: '🇪🇸', icon: 'language', dir: 'ltr' },
  { label: 'Română', value: 'ro-RO', code: 'RO', flag: '🇷🇴', icon: 'language', dir: 'ltr' },
  { label: 'Shqip', value: 'sq-AL', code: 'SQ', flag: '🇦🇱', icon: 'language', dir: 'ltr' },
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
  'ro-RO': quasarLangRo,
  'ro': quasarLangRo,
  'sq-AL': quasarLangSq,
  'sq': quasarLangSq,
  'al': quasarLangSq,
  'ru-RU': quasarLangRu,
  'ru': quasarLangRu,
  'uk-UA': quasarLangUk,
  'uk': quasarLangUk,
  'ua': quasarLangUk,
  'ar-SA': quasarLangAr,
  'ar': quasarLangAr,
  'zh-CN': quasarLangZh,
  'zh': quasarLangZh,
  'cn': quasarLangZh
}

export function normalizeLocale(lang) {
  if (!lang || typeof lang !== 'string') return 'it-IT'
  const trimmed = lang.trim().toLowerCase()
  if (!trimmed) return 'it-IT'

  const exact = SUPPORTED_LOCALES.find(
    l => l.value.toLowerCase() === trimmed || l.code.toLowerCase() === trimmed
  )
  if (exact) return exact.value

  const byPrefix = SUPPORTED_LOCALES.find(
    l => l.value.toLowerCase().startsWith(trimmed) || trimmed.startsWith(l.value.toLowerCase().slice(0, 2))
  )
  return byPrefix ? byPrefix.value : 'it-IT'
}

export function getBrowserLocale() {
  try {
    if (typeof navigator !== 'undefined' && navigator.language) {
      return normalizeLocale(navigator.language)
    }
  } catch (e) {
    console.warn('Could not detect browser locale:', e)
  }
  return 'it-IT'
}

export function getSavedLocale(fallbackToBrowser = false) {
  try {
    if (typeof localStorage !== 'undefined') {
      const saved = localStorage.getItem('app_language') || localStorage.getItem('user_locale')
      if (saved) return normalizeLocale(saved)
    }
  } catch (e) {
    console.warn('Could not read saved locale from localStorage:', e)
  }

  if (fallbackToBrowser) {
    return getBrowserLocale()
  }

  return 'it-IT'
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
    try {
      if (normalized !== 'it-IT' && normalized !== 'it') {
        loadLocaleMessages(normalized, i18nInstance).catch(() => {})
      }

      if (i18nInstance.global && i18nInstance.global.locale) {
        if (typeof i18nInstance.global.locale === 'object' && 'value' in i18nInstance.global.locale) {
          i18nInstance.global.locale.value = normalized
        } else {
          i18nInstance.global.locale = normalized
        }
      } else if (i18nInstance.locale) {
        if (typeof i18nInstance.locale === 'object' && 'value' in i18nInstance.locale) {
          i18nInstance.locale.value = normalized
        } else {
          i18nInstance.locale = normalized
        }
      }
    } catch (err) {
      console.warn('Failed to set i18n locale:', err)
    }
  }

  // 2. Update Quasar Language Pack
  const quasarPack = QUASAR_LANG_MAP[normalized] || quasarLangIt
  try {
    if ($q && $q.lang && typeof $q.lang.set === 'function') {
      $q.lang.set(quasarPack)
    } else if (typeof Quasar !== 'undefined' && Quasar?.lang && typeof Quasar.lang.set === 'function') {
      Quasar.lang.set(quasarPack)
    }
  } catch (err) {
    console.warn('Failed to set Quasar language pack:', err)
  }

  // 3. Update DOM direction and lang
  if (typeof document !== 'undefined' && document.documentElement) {
    document.documentElement.setAttribute('lang', normalized)
    document.documentElement.setAttribute('dir', isRTL ? 'rtl' : 'ltr')
  }

  // 4. Persist in localStorage
  try {
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('app_language', normalized)
      localStorage.setItem('user_locale', normalized)
    }
  } catch (err) {
    console.warn('Failed to persist locale in localStorage:', err)
  }

  return normalized
}
