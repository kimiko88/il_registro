import { describe, it, expect, beforeEach, vi } from 'vitest'
import { usePwaUpdate } from '@/composables/usePwaUpdate'

describe('usePwaUpdate Composable', () => {
  let pwaUpdate

  beforeEach(() => {
    pwaUpdate = usePwaUpdate()
    pwaUpdate._reset()
  })

  it('initializes with needRefresh false and userDismissed false', () => {
    expect(pwaUpdate.needRefresh.value).toBe(false)
    expect(pwaUpdate.userDismissed.value).toBe(false)
  })

  it('updates needRefresh when a new version is detected', () => {
    pwaUpdate._setNeedRefresh(true)
    expect(pwaUpdate.needRefresh.value).toBe(true)
    expect(pwaUpdate.userDismissed.value).toBe(false)
  })

  it('marks userDismissed as true when dismiss is called', () => {
    pwaUpdate._setNeedRefresh(true)
    pwaUpdate.dismiss()
    expect(pwaUpdate.userDismissed.value).toBe(true)
  })

  it('resets state correctly with _reset helper', () => {
    pwaUpdate._setNeedRefresh(true)
    pwaUpdate.dismiss()
    pwaUpdate._reset()
    expect(pwaUpdate.needRefresh.value).toBe(false)
    expect(pwaUpdate.userDismissed.value).toBe(false)
  })
})
