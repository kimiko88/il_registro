import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useThemeStore } from '@/stores/theme'

describe('Theme Accessibility Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with default accessibility values', () => {
    const themeStore = useThemeStore()
    expect(themeStore.dsaFont).toBe(false)
    expect(themeStore.highContrast).toBe(false)
  })

  it('should toggle DSA font setting correctly', () => {
    const themeStore = useThemeStore()
    themeStore.toggleDsaFont()
    expect(themeStore.dsaFont).toBe(true)

    themeStore.toggleDsaFont()
    expect(themeStore.dsaFont).toBe(false)
  })

  it('should toggle High Contrast mode correctly', () => {
    const themeStore = useThemeStore()
    themeStore.toggleHighContrast()
    expect(themeStore.highContrast).toBe(true)

    themeStore.toggleHighContrast()
    expect(themeStore.highContrast).toBe(false)
  })
})
