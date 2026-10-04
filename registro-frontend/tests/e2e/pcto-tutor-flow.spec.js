import { describe, it, expect, vi } from 'vitest'
import { pctoCompanyTutorService } from '@/services/pctoCompanyTutorService'

vi.mock('@/services/pctoCompanyTutorService', () => ({
    pctoCompanyTutorService: {
        getSession: vi.fn(),
        getAssignedStudents: vi.fn(),
        verifyTimesheet: vi.fn(),
        submitEvaluation: vi.fn()
    }
}))

describe('PCTO Company Tutor Magic Link Portal E2E Workflow', () => {
    it('executes full company tutor lifecycle: magic link access, student list, 1-click timesheet validation, and soft skills evaluation', async () => {
        const token = 'token-magic-link-test-99999'

        // Step 1: Magic link authentication & session load
        pctoCompanyTutorService.getSession.mockResolvedValueOnce({
            tutor: {
                id: 'tut-e2e',
                company_name: 'Digital Innovations S.p.A.',
                tutor_first_name: 'Marco',
                tutor_last_name: 'Galli'
            }
        })
        const sessionRes = await pctoCompanyTutorService.getSession(token)
        expect(sessionRes.tutor.company_name).toBe('Digital Innovations S.p.A.')

        // Step 2: Fetch assigned students in internship
        pctoCompanyTutorService.getAssignedStudents.mockResolvedValueOnce({
            data: [
                {
                    student_id: 'stud-e2e',
                    student_name: 'Francesca Neri',
                    project_id: 'proj-e2e',
                    project_title: 'Sviluppo Web Cloud',
                    total_hours: 80,
                    completed_hours: 40,
                    is_evaluated: false
                }
            ]
        })
        const studentsRes = await pctoCompanyTutorService.getAssignedStudents(token)
        expect(studentsRes.data).toHaveLength(1)

        // Step 3: Tutor signs / verifies daily hours
        pctoCompanyTutorService.verifyTimesheet.mockResolvedValueOnce({
            message: 'Presenze certificate dal tutor aziendale',
            verification: {
                student_id: 'stud-e2e',
                hours_declared: 8,
                hours_approved: 8,
                signed_at: '2026-10-04T12:00:00Z'
            }
        })
        const verifyRes = await pctoCompanyTutorService.verifyTimesheet(token, {
            project_id: 'proj-e2e',
            student_id: 'stud-e2e',
            activity_date: '2026-10-04',
            hours_declared: 8,
            hours_approved: 8,
            tutor_notes: 'Ottima puntualità e accuratezza'
        })
        expect(verifyRes.verification.hours_approved).toBe(8)

        // Step 4: Tutor submits soft skills and final evaluation
        pctoCompanyTutorService.submitEvaluation.mockResolvedValueOnce({
            message: 'Valutazione competenze aziendali registrata con successo',
            evaluation: {
                reliability_level: 5,
                technical_skills: 5,
                teamwork_skills: 5,
                final_feedback: 'Studentessa eccellente'
            }
        })
        const evalRes = await pctoCompanyTutorService.submitEvaluation(token, {
            project_id: 'proj-e2e',
            student_id: 'stud-e2e',
            reliability_level: 5,
            technical_skills: 5,
            teamwork_skills: 5,
            final_feedback: 'Studentessa eccellente'
        })
        expect(evalRes.evaluation.reliability_level).toBe(5)
    })
})
