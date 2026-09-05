export const localeLoaders = {
  'it-IT': () => Promise.resolve({ default: null }),
  'it': () => Promise.resolve({ default: null }),
  'en-US': () => import('./en-US/index.js'),
  'en': () => import('./en-US/index.js'),
  'de-DE': () => import('./de-DE/index.js'),
  'de': () => import('./de-DE/index.js'),
  'fr-FR': () => import('./fr-FR/index.js'),
  'fr': () => import('./fr-FR/index.js'),
  'es-ES': () => import('./es-ES/index.js'),
  'es': () => import('./es-ES/index.js'),
  'ro-RO': () => import('./ro-RO/index.js'),
  'ro': () => import('./ro-RO/index.js'),
  'sq-AL': () => import('./sq-AL/index.js'),
  'sq': () => import('./sq-AL/index.js'),
  'al': () => import('./sq-AL/index.js'),
  'ru-RU': () => import('./ru-RU/index.js'),
  'ru': () => import('./ru-RU/index.js'),
  'uk-UA': () => import('./uk-UA/index.js'),
  'uk': () => import('./uk-UA/index.js'),
  'ua': () => import('./uk-UA/index.js'),
  'ar-SA': () => import('./ar-SA/index.js'),
  'ar': () => import('./ar-SA/index.js'),
  'zh-CN': () => import('./zh-CN/index.js'),
  'zh': () => import('./zh-CN/index.js'),
  'cn': () => import('./zh-CN/index.js')
}

export const loadedLanguages = ['it-IT', 'it']

export async function loadLocaleMessages(locale, i18nInstance = null) {
  if (!locale) return 'it-IT'
  const norm = locale.trim()
  if (loadedLanguages.includes(norm)) {
    return norm
  }

  const loader = localeLoaders[norm]
  if (!loader) {
    return 'it-IT'
  }

  try {
    const mod = await loader()
    const dict = mod.default || mod
    if (dict && i18nInstance) {
      if (i18nInstance.global && typeof i18nInstance.global.setLocaleMessage === 'function') {
        i18nInstance.global.setLocaleMessage(norm, dict)
      } else if (typeof i18nInstance.setLocaleMessage === 'function') {
        i18nInstance.setLocaleMessage(norm, dict)
      }
    }
    loadedLanguages.push(norm)
    return norm
  } catch (err) {
    console.warn(`[i18n] Failed to dynamically load locale "${norm}":`, err)
    return 'it-IT'
  }
}
