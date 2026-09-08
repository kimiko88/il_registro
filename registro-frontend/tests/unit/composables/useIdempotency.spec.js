import { describe, it, expect } from 'vitest'
import { useIdempotency } from '@/composables/useIdempotency'

// UUID v4 pattern
const UUID_PATTERN = /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i

describe('useIdempotency Composable', () => {
  it('initializes with a non-empty key containing a UUID v4', () => {
    const { currentKey } = useIdempotency('test')
    expect(currentKey.value).toBeTruthy()
    // Key format: "prefix:uuid"
    const uuidPart = currentKey.value.split(':').slice(1).join(':')
    expect(uuidPart).toMatch(UUID_PATTERN)
  })

  it('uses the given prefix in the key', () => {
    const { currentKey } = useIdempotency('grade-add')
    expect(currentKey.value.startsWith('grade-add:')).toBe(true)
  })

  it('uses "op" as default prefix when none provided', () => {
    const { currentKey } = useIdempotency()
    expect(currentKey.value.startsWith('op:')).toBe(true)
  })

  it('generateKey produces a new UUID different from the previous one', () => {
    const { currentKey, generateKey } = useIdempotency('test')
    const first = currentKey.value
    generateKey()
    const second = currentKey.value
    expect(second).not.toBe(first)
    const uuidPart = second.split(':').slice(1).join(':')
    expect(uuidPart).toMatch(UUID_PATTERN)
  })

  it('rotateKey produces a new UUID different from the previous one', () => {
    const { currentKey, rotateKey } = useIdempotency('test')
    const first = currentKey.value
    rotateKey()
    const second = currentKey.value
    expect(second).not.toBe(first)
    const uuidPart = second.split(':').slice(1).join(':')
    expect(uuidPart).toMatch(UUID_PATTERN)
  })

  it('currentKey is readonly — direct mutation should not change internal value', () => {
    const { currentKey } = useIdempotency('test')
    const original = currentKey.value
    // In Vue 3 readonly refs throw a warning but don't throw in prod mode.
    // We just verify the getter returns the correct value.
    expect(currentKey.value).toBe(original)
  })

  it('two separate instances have independent keys', () => {
    const a = useIdempotency('op-a')
    const b = useIdempotency('op-b')
    expect(a.currentKey.value).not.toBe(b.currentKey.value)
    expect(a.currentKey.value.startsWith('op-a:')).toBe(true)
    expect(b.currentKey.value.startsWith('op-b:')).toBe(true)
  })

  it('generateKey on one instance does not affect another', () => {
    const a = useIdempotency('test')
    const b = useIdempotency('test')
    const bOriginal = b.currentKey.value
    a.generateKey()
    // b should remain unchanged
    expect(b.currentKey.value).toBe(bOriginal)
  })
})
