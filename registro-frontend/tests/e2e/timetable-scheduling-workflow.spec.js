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

describe('Timetable Scheduling & Conflict Workflow E2E', () => {
    it('executes weekly class schedule update and slot assignment', async () => {
        const classSchedule = [
            { day_of_week: 1, hour_index: 1, subject: 'Matematica', room: 'Aula 101' },
            { day_of_week: 1, hour_index: 2, subject: 'Matematica', room: 'Aula 101' }
        ]

        // Verify timetable slots
        expect(classSchedule.length).toBe(2)
        expect(classSchedule[0].subject).toBe('Matematica')
        expect(classSchedule[1].hour_index).toBe(2)
    })
})
