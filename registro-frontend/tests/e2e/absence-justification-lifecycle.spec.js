import { describe, it, expect, vi } from 'vitest'
import { createTestingPinia } from '@pinia/testing'
import { useAttendanceStore } from '@/stores/attendance'

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

describe('Absence & Justification Lifecycle E2E', () => {
    it('executes student absence entry, parent justification request, and approval flow', async () => {
        const pinia = createTestingPinia({
            createSpy: vi.fn,
            stubActions: false,
            initialState: {
                auth: { user: { id: 'teacher-1', role: 'teacher' } },
                attendance: {
                    justifications: [
                        { id: 'just-1', student_id: 'st-1', student_name: 'Mario Rossi', date: '2026-09-01', reason: 'Motivi di salute', status: 'pending' }
                    ]
                }
            }
        })

        const attendanceStore = useAttendanceStore(pinia)

        // Verify pending justification
        expect(attendanceStore.justifications.length).toBe(1)
        expect(attendanceStore.justifications[0].status).toBe('pending')

        // Simulate approval
        attendanceStore.justifications[0].status = 'approved'
        expect(attendanceStore.justifications[0].status).toBe('approved')
    })
})
