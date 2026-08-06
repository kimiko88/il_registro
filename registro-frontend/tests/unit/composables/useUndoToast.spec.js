import { describe, it, expect, vi } from 'vitest'
import { useUndoToast } from '@/composables/useUndoToast'

// Mock Quasar Notify
vi.mock('quasar', () => ({
  useQuasar: () => ({
    notify: vi.fn(() => () => {})
  })
}))

describe('useUndoToast Composable', () => {
  it('should initialize notifyWithUndo function', () => {
    const { notifyWithUndo } = useUndoToast()
    expect(typeof notifyWithUndo).toBe('function')
  })

  it('should resolve to false if undo is not clicked within timeout', async () => {
    vi.useFakeTimers()
    const { notifyWithUndo } = useUndoToast(1000)
    const promise = notifyWithUndo('Assenza salvata', vi.fn())

    vi.advanceTimersByTime(1500)
    const result = await promise
    expect(result).toBe(false)
    vi.useRealTimers()
  })
})
