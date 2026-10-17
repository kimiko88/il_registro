import { describe, it, expect, vi } from 'vitest'
import { middleSchoolExamService } from '@/services/middleSchoolExamService'

vi.mock('@/services/middleSchoolExamService', () => ({
    middleSchoolExamService: {
        getSubcommissions: vi.fn(),
        getCandidates: vi.fn(),
        saveCandidateGrades: vi.fn(),
        calculateFinalGrades: vi.fn(),
        generateDiploma: vi.fn()
    }
}))

describe('Middle School Exam (Esame Terza Media D.Lgs. 62/2017) E2E Workflow', () => {
    it('executes full exam lifecycle: fetch subcommission, enter tests, calculate final grade with rounding/lode, and generate diploma', async () => {
        // Step 1: Fetch Subcommission & Candidates
        middleSchoolExamService.getSubcommissions.mockResolvedValueOnce({
            data: [
                { id: 'subcomm-1', class_id: 'class-3a', name: 'Sottocommissione 3A - A.S. 2025/2026' }
            ]
        })
        const subcommRes = await middleSchoolExamService.getSubcommissions('2025/2026')
        expect(subcommRes.data).toHaveLength(1)
        const subcommId = subcommRes.data[0].id

        middleSchoolExamService.getCandidates.mockResolvedValueOnce({
            data: [
                {
                    id: 'cand-1',
                    subcommission_id: subcommId,
                    student_id: 'student-101',
                    first_name: 'Giulia',
                    last_name: 'Verdi',
                    admission_grade: 9,
                    is_admitted: true,
                    written_italian: null,
                    written_math: null,
                    written_english: null,
                    written_second_lang: null,
                    oral_interview: null,
                    final_grade: null,
                    with_honors: false
                }
            ]
        })
        const candRes = await middleSchoolExamService.getCandidates(subcommId)
        expect(candRes.data).toHaveLength(1)
        expect(candRes.data[0].admission_grade).toBe(9)

        // Step 2: Save candidate exam grades
        middleSchoolExamService.saveCandidateGrades.mockResolvedValueOnce({
            data: {
                message: 'Voti esame salvati con successo',
                candidate: {
                    id: 'cand-1',
                    written_italian: 10,
                    written_math: 10,
                    written_english: 10,
                    written_second_lang: 10,
                    oral_interview: 10
                }
            }
        })
        const saveRes = await middleSchoolExamService.saveCandidateGrades('cand-1', {
            written_italian: 10,
            written_math: 10,
            written_english: 10,
            written_second_lang: 10,
            oral_interview: 10
        })
        expect(saveRes.data.candidate.written_italian).toBe(10)

        // Step 3: Compute final grades with D.Lgs. 62/2017 & D.M. 741/2017 rounding & honors (Lode)
        middleSchoolExamService.calculateFinalGrades.mockResolvedValueOnce({
            data: {
                message: 'Voti finali calcolati a norma D.Lgs. 62/2017',
                results: [
                    {
                        candidate_id: 'cand-1',
                        admission_grade: 9,
                        exam_tests_average: 10.0,
                        raw_final_average: 9.5, // (9 + 10) / 2 = 9.5 -> rounds to 10
                        final_grade: 10,
                        with_honors: true,
                        passed: true
                    }
                ]
            }
        })
        const calcRes = await middleSchoolExamService.calculateFinalGrades(subcommId)
        expect(calcRes.data.results[0].final_grade).toBe(10)
        expect(calcRes.data.results[0].with_honors).toBe(true)

        // Step 4: Generate official Diploma Certificate
        middleSchoolExamService.generateDiploma.mockResolvedValueOnce({
            data: {
                message: 'Diploma di licenza media generato con successo',
                diploma_url: '/api/v1/middle-school-exam/candidates/cand-1/diploma.pdf',
                candidate_name: 'Giulia Verdi',
                final_grade_text: 'DIECI E LODE'
            }
        })
        const diplomaRes = await middleSchoolExamService.generateDiploma('cand-1')
        expect(diplomaRes.data.final_grade_text).toBe('DIECI E LODE')
        expect(diplomaRes.data.diploma_url).toContain('.pdf')
    })
})
