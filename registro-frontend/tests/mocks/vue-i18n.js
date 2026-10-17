export function useI18n() {
  return {
    t: (key, params) => key,
    te: () => true,
    locale: { value: 'it-IT' }
  }
}

export function createI18n() {
  return {
    install: () => {},
    global: {
      t: (key) => key,
      te: () => true,
      locale: { value: 'it-IT' }
    }
  }
}

export default {
  useI18n,
  createI18n
}
