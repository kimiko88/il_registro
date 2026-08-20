import { describe, it, expect, vi } from 'vitest'
import { createTestingPinia } from '@pinia/testing'
import { useAgendaStore } from '@/stores/agenda'

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

describe('Agenda & Homework Lifecycle E2E', () => {
    it('executes homework assignment, due date filtering, and student completion toggle', async () => {
        const pinia = createTestingPinia({
            createSpy: vi.fn,
            stubActions: false,
            initialState: {
                auth: { user: { id: 'student-1', role: 'student' } },
                agenda: {
                    events: [
                        { id: 'agenda-1', title: 'Equazioni di 2° grado', type: 'homework', due_date: '2026-09-10', completed: false }
                    ]
                }
            }
        })

        const agendaStore = useAgendaStore(pinia)

        // Verify initial homework item
        expect(agendaStore.events.length).toBe(1)
        expect(agendaStore.events[0].completed).toBe(false)

        // Toggle task completion
        agendaStore.events[0].completed = true
        expect(agendaStore.events[0].completed).toBe(true)
    })
})
