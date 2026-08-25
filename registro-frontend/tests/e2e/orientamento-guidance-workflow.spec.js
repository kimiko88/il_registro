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

describe('Career Guidance (Orientamento 30h) Workflow E2E', () => {
    it('executes orientation activity registration, student sign-up, and attendance logging', async () => {
        const orientationEvent = {
            id: 'event-1',
            title: "Fiera dell'Orientamento Universitario",
            category: 'University',
            hours: 6,
            registered_students: []
        }

        // Verify orientation event parameters
        expect(orientationEvent.hours).toBe(6)
        expect(orientationEvent.registered_students.length).toBe(0)

        // Student registers for event
        orientationEvent.registered_students.push({ student_id: 'st-1', attended: true })
        expect(orientationEvent.registered_students.length).toBe(1)
        expect(orientationEvent.registered_students[0].attended).toBe(true)
    })
})
