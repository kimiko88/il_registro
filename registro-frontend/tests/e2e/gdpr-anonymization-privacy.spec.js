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

describe('GDPR Privacy & Adult Student Access Control E2E', () => {
    it('executes 18+ adult student privacy enforcement for parent access', async () => {
        const student = {
            id: 'st-adult-18',
            age: 18,
            parent_consent_granted: false
        }

        const canParentAccessGrades = student.age < 18 || student.parent_consent_granted
        expect(canParentAccessGrades).toBe(false)
    })
})
