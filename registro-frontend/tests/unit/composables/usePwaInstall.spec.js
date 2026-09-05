import { describe, it, expect, vi, beforeEach } from 'vitest'
import { usePwaInstall } from '@/composables/usePwaInstall'

describe('usePwaInstall Composable', () => {
  let pwa

  beforeEach(() => {
    pwa = usePwaInstall()
    pwa._resetState()
  })

  it('initializes with expected default values', () => {
    expect(pwa.canInstall.value).toBe(false)
    expect(pwa.isInstalled.value).toBe(false)
  })

  it('enables canInstall when deferred prompt is captured', () => {
    const mockPrompt = {
      prompt: vi.fn(),
      userChoice: Promise.resolve({ outcome: 'accepted' })
    }

    pwa._setPrompt(mockPrompt)
    expect(pwa.canInstall.value).toBe(true)
  })

  it('prompts install and updates isInstalled on accept', async () => {
    const mockPrompt = {
      prompt: vi.fn(),
      userChoice: Promise.resolve({ outcome: 'accepted' })
    }

    pwa._setPrompt(mockPrompt)
    const result = await pwa.promptInstall()

    expect(mockPrompt.prompt).toHaveBeenCalled()
    expect(result).toBe(true)
    expect(pwa.isInstalled.value).toBe(true)
    expect(pwa.canInstall.value).toBe(false)
  })

  it('allows user to dismiss prompt', () => {
    const mockPrompt = {
      prompt: vi.fn(),
      userChoice: Promise.resolve({ outcome: 'dismissed' })
    }

    pwa._setPrompt(mockPrompt)
    expect(pwa.canInstall.value).toBe(true)

    pwa.dismissPrompt()
    expect(pwa.canInstall.value).toBe(false)
  })
})
