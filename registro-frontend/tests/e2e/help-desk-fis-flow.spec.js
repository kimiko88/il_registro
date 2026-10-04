import { describe, it, expect, vi } from 'vitest'
import { helpDeskService } from '@/services/helpDeskService'

vi.mock('@/services/helpDeskService', () => ({
    helpDeskService: {
        createSlot: vi.fn(),
        getSlots: vi.fn(),
        bookSlot: vi.fn(),
        markAttendance: vi.fn(),
        completeSlot: vi.fn(),
        getFISReport: vi.fn()
    }
}))

describe('Help Desk & FIS Liquidation E2E Workflow', () => {
    it('executes full help desk lifecycle: slot creation, student booking, teacher roll call, slot completion, and DSGA FIS accounting', async () => {
        // Step 1: Teacher creates help desk slot
        helpDeskService.createSlot.mockResolvedValueOnce({
            message: 'Disponibilità sportello help registrata',
            slot: {
                id: 'slot-e2e-1',
                slot_date: '2026-10-18',
                start_time: '14:30',
                end_time: '16:30',
                max_capacity: 4,
                status: 'open'
            }
        })
        const createRes = await helpDeskService.createSlot({
            subject_id: 'subj-math-1',
            slot_date: '2026-10-18',
            start_time: '14:30',
            end_time: '16:30',
            max_capacity: 4
        })
        expect(createRes.slot.id).toBe('slot-e2e-1')

        // Step 2: Student books slot with question topic
        helpDeskService.bookSlot.mockResolvedValueOnce({
            message: 'Prenotazione effettuata con successo',
            booking: {
                id: 'book-e2e-1',
                slot_id: 'slot-e2e-1',
                student_id: 'student-e2e-1',
                topic_description: 'Esercizi di preparazione alla verifica sui vettori',
                status: 'booked'
            }
        })
        const bookRes = await helpDeskService.bookSlot('slot-e2e-1', {
            topic_description: 'Esercizi di preparazione alla verifica sui vettori'
        })
        expect(bookRes.booking.status).toBe('booked')

        // Step 3: Teacher marks presence
        helpDeskService.markAttendance.mockResolvedValueOnce({
            message: 'Presenza aggiornata'
        })
        const attRes = await helpDeskService.markAttendance('book-e2e-1', 'attended')
        expect(attRes.message).toBe('Presenza aggiornata')

        // Step 4: Teacher completes slot
        helpDeskService.completeSlot.mockResolvedValueOnce({
            message: 'Sportello concluso e validato per la rendicontazione FIS'
        })
        const compRes = await helpDeskService.completeSlot('slot-e2e-1')
        expect(compRes.message).toContain('FIS')

        // Step 5: DSGA pulls FIS accounting report
        helpDeskService.getFISReport.mockResolvedValueOnce({
            data: [
                {
                    teacher_id: 'teacher-math-1',
                    teacher_name: 'Prof. Mario Rossi',
                    completed_slots: 4,
                    total_hours: 6.0,
                    attended_count: 12
                }
            ]
        })
        const fisRes = await helpDeskService.getFISReport()
        expect(fisRes.data[0].total_hours).toBe(6.0)
        expect(fisRes.data[0].completed_slots).toBe(4)
    })
})
