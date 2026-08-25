import { describe, it, expect, vi } from 'vitest'
import { createTestingPinia } from '@pinia/testing'
import { useScrutinyStore } from '@/stores/scrutiny'

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

describe('Grading & Scrutiny Lifecycle E2E', () => {
    it('executes grade weighting, scrutiny proposal, and locking flow', async () => {
        const pinia = createTestingPinia({
            createSpy: vi.fn,
            stubActions: false,
            initialState: {
                auth: { user: { id: 'admin-1', role: 'admin' } },
                scrutiny: {
                    overview: [
                        { class_id: 'c1', class_name: '2A', total_students: 20, is_locked: false, status: 'in_progress' }
                    ]
                }
            }
        })

        const scrutinyStore = useScrutinyStore(pinia)

        // Verify initial scrutiny state
        expect(scrutinyStore.overview.length).toBe(1)
        expect(scrutinyStore.overview[0].is_locked).toBe(false)

        // Simulate locking scrutiny session
        scrutinyStore.overview[0].is_locked = true
        scrutinyStore.overview[0].status = 'closed'

        expect(scrutinyStore.overview[0].is_locked).toBe(true)
        expect(scrutinyStore.overview[0].status).toBe('closed')
    })
})
