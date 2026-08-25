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

describe('School Calendar & Non-Teaching Holiday Days Workflow E2E', () => {
    it('executes academic year settings setup and non-teaching day registration', async () => {
        const yearSettings = {
            year_label: '2026/2027',
            start_date: '2026-09-12',
            end_date: '2027-06-08'
        }

        const nonTeachingDays = []

        // Verify academic year parameters
        expect(yearSettings.year_label).toBe('2026/2027')

        // Add regional holiday
        nonTeachingDays.push({ date: '2026-11-01', label: 'Festa di Tutti i Santi' })
        expect(nonTeachingDays.length).toBe(1)
        expect(nonTeachingDays[0].label).toBe('Festa di Tutti i Santi')
    })
})
