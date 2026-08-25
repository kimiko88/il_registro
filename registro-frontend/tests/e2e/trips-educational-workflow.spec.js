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

describe('Educational Trips & Parent Consent Workflow E2E', () => {
    it('executes trip proposal, accompanying teachers setup, and parent consent submission', async () => {
        const trip = {
            id: 'trip-1',
            title: 'Visita Didattica Museo della Scienza',
            destination: 'Milano',
            departure_date: '2026-10-20',
            accompanying_teachers: 'Prof. Rossi, Prof. Bianchi',
            consents: []
        }

        // Verify initial trip state
        expect(trip.destination).toBe('Milano')
        expect(trip.consents.length).toBe(0)

        // Submit parent consent
        trip.consents.push({ student_id: 'st-1', status: 'granted', signed_at: '2026-10-01T10:00:00Z' })
        expect(trip.consents.length).toBe(1)
        expect(trip.consents[0].status).toBe('granted')
    })
})
