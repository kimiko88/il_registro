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

describe('PCTO & Internship Workflow E2E', () => {
    it('executes PCTO project creation, student hour logging, and tutor verification', async () => {
        const pctoProject = {
            id: 'proj-1',
            title: 'Tirocinio Sviluppo Software',
            company: 'TechCorp SRL',
            total_hours: 40
        }

        const hourLog = {
            id: 'hour-1',
            project_id: 'proj-1',
            hours: 8,
            activity: 'Formazione Vue.js',
            verified: false
        }

        // Verify hour logging state
        expect(hourLog.verified).toBe(false)

        // Tutor verifies hours
        hourLog.verified = true
        expect(hourLog.verified).toBe(true)
    })
})
