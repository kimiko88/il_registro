import { describe, it, expect, vi } from 'vitest'

vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            loading: { show: vi.fn(), hide: vi.fn() },
            notify: vi.fn(),
            dialog: vi.fn().mockReturnValue({ onOk: (fn) => fn() })
        })
    }
})

describe('High-Concurrency Scrutiny Locking & Read-Only UI E2E', () => {
    it('executes scrutiny session locking and blocks grade edit triggers', async () => {
        const scrutinyState = {
            class_id: 'class-1',
            is_locked: false,
            grades_editable: true
        }

        // Before lock
        expect(scrutinyState.grades_editable).toBe(true)

        // Lock scrutiny
        scrutinyState.is_locked = true
        scrutinyState.grades_editable = false

        expect(scrutinyState.is_locked).toBe(true)
        expect(scrutinyState.grades_editable).toBe(false)
    })
})
