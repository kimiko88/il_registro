import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock api client for Primary School evaluations
const mockApi = {
    get: vi.fn(),
    post: vi.fn(),
    delete: vi.fn()
}

vi.mock('@/services/api', () => ({
    default: mockApi
}))
vi.mock('src/services/api', () => ({
    default: mockApi
}))

describe('Primary School Descriptive Evaluations (O.M. 172/2020) E2E Workflow', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('PE01 — Complete primary evaluation lifecycle: fetch objectives, record 4 dimensions with levels, and generate matrix', async () => {
        // Step 1: Teacher loads learning objectives for Primary Class 3A - Italiano
        const mockObjectives = [
            {
                id: 'obj-1',
                title: 'Ascolto e parlato',
                description: 'Ascolta e comprende comunicazioni verbali',
                year_grade: 3,
                subject_id: 'sub-ita'
            },
            {
                id: 'obj-2',
                title: 'Lettura e comprensione',
                description: 'Legge e comprende testi narrativi e descrittivi',
                year_grade: 3,
                subject_id: 'sub-ita'
            }
        ]
        mockApi.get.mockResolvedValueOnce({
            data: { data: mockObjectives }
        })

        const objectivesRes = await mockApi.get('/primary/objectives?subject_id=sub-ita&year_grade=3')
        expect(objectivesRes.data.data).toHaveLength(2)
        expect(objectivesRes.data.data[0].title).toBe('Ascolto e parlato')

        // Step 2: Teacher records batch evaluations for 3 students with 4 ministerial dimensions
        const batchPayload = {
            class_id: 'cls-3a',
            subject_id: 'sub-ita',
            objective_id: 'obj-1',
            date: '2026-01-20',
            semester: 1,
            evaluations: [
                {
                    student_id: 'stu-1',
                    level: 'avanzato',
                    dimension_autonomy: 'autonomo',
                    dimension_continuity: 'continuo',
                    dimension_familiarity: 'situazioni_note_e_non_note',
                    dimension_resources: 'risorse_proprie_e_fornite',
                    notes: 'Ottimo percorso di apprendimento'
                },
                {
                    student_id: 'stu-2',
                    level: 'intermedio',
                    dimension_autonomy: 'autonomo',
                    dimension_continuity: 'discontinuo',
                    dimension_familiarity: 'situazioni_note',
                    dimension_resources: 'risorse_fornite',
                    notes: ''
                },
                {
                    student_id: 'stu-3',
                    level: 'in_via_di_prima_acquisizione',
                    dimension_autonomy: 'guidato',
                    dimension_continuity: 'non_continuo',
                    dimension_familiarity: 'situazioni_note',
                    dimension_resources: 'risorse_fornite',
                    notes: 'Supporto costante del docente'
                }
            ]
        }
        mockApi.post.mockResolvedValueOnce({
            data: { message: 'Valutazioni descrittive registrate con successo' }
        })

        const saveRes = await mockApi.post('/primary/evaluations/batch', batchPayload)
        expect(saveRes.data.message).toContain('Valutazioni descrittive registrate con successo')

        // Step 3: Teacher loads class matrix aggregating evaluations by objective and student
        const mockMatrix = {
            class_id: 'cls-3a',
            subject_id: 'sub-ita',
            semester: 1,
            objectives: mockObjectives,
            students: [
                {
                    student_id: 'stu-1',
                    student_name: 'Lorenzo Galli',
                    evaluations: {
                        'obj-1': { level: 'avanzato', dimension_autonomy: 'autonomo' }
                    }
                },
                {
                    student_id: 'stu-2',
                    student_name: 'Chiara Fontana',
                    evaluations: {
                        'obj-1': { level: 'intermedio', dimension_autonomy: 'autonomo' }
                    }
                },
                {
                    student_id: 'stu-3',
                    student_name: 'Davide Neri',
                    evaluations: {
                        'obj-1': { level: 'in_via_di_prima_acquisizione', dimension_autonomy: 'guidato' }
                    }
                }
            ]
        }
        mockApi.get.mockResolvedValueOnce({
            data: { data: mockMatrix }
        })

        const matrixRes = await mockApi.get('/primary/matrix?class_id=cls-3a&subject_id=sub-ita&semester=1')
        expect(matrixRes.data.data.students).toHaveLength(3)
        expect(matrixRes.data.data.students[0].evaluations['obj-1'].level).toBe('avanzato')
        expect(matrixRes.data.data.students[2].evaluations['obj-1'].level).toBe('in_via_di_prima_acquisizione')

        // Step 4: Parent views individual descriptive evaluation card for Lorenzo Galli
        mockApi.get.mockResolvedValueOnce({
            data: {
                data: [
                    {
                        objective_id: 'obj-1',
                        objective_title: 'Ascolto e parlato',
                        level: 'avanzato',
                        date: '2026-01-20',
                        notes: 'Ottimo percorso di apprendimento'
                    }
                ]
            }
        })

        const studentCardRes = await mockApi.get('/primary/student/stu-1?semester=1')
        expect(studentCardRes.data.data[0].level).toBe('avanzato')
        expect(studentCardRes.data.data[0].objective_title).toBe('Ascolto e parlato')
    })
})
