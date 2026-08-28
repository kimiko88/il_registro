import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { useA11yAnnouncer } from '@/composables/useA11yAnnouncer'

describe('New Advanced Accessibility Features & Cloud Sync', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('toggles Focus Mode (ADHD / DSA Clean Reading)', () => {
    const themeStore = useThemeStore()
    expect(themeStore.isFocusMode).toBe(false)

    themeStore.toggleFocusMode(true)
    expect(themeStore.isFocusMode).toBe(true)

    themeStore.toggleFocusMode(false)
    expect(themeStore.isFocusMode).toBe(false)
  })

  it('dispatches polite and assertive screen reader announcements', async () => {
    const { announce, politeMessage, assertiveMessage } = useA11yAnnouncer()

    announce('Voto inserito con successo', 'polite')
    await new Promise(resolve => setTimeout(resolve, 60))
    expect(politeMessage.value).toBe('Voto inserito con successo')

    announce('Sessione in scadenza', 'assertive')
    await new Promise(resolve => setTimeout(resolve, 60))
    expect(assertiveMessage.value).toBe('Sessione in scadenza')
  })
})
