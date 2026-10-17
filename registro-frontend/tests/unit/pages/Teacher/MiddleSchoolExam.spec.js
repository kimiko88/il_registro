import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import MiddleSchoolExam from '@/pages/teacher/MiddleSchoolExam.vue'
import { middleSchoolExamService } from '@/services/middleSchoolExamService'

vi.mock('@/services/middleSchoolExamService', () => ({
    middleSchoolExamService: {
        getOrCreateExam: vi.fn(),
        listCandidates: vi.fn(),
        saveAdmission: vi.fn(),
        evaluateCandidate: vi.fn(),
        generateDiploma: vi.fn(),
        updateStatus: vi.fn()
    },
    default: {
        getOrCreateExam: vi.fn(),
        listCandidates: vi.fn(),
        saveAdmission: vi.fn(),
        evaluateCandidate: vi.fn(),
        generateDiploma: vi.fn(),
        updateStatus: vi.fn()
    }
}))

vi.mock('@/stores/classes', () => ({
    useClassesStore: () => ({
        classes: [
            { id: 'class-3a', name: '3A', section: 'A' },
            { id: 'class-3b', name: '3B', section: 'B' }
        ],
        loading: false,
        fetchClasses: vi.fn().mockResolvedValue()
    })
}))

describe('Teacher Middle School Exam Page (D.Lgs. 62/2017)', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        middleSchoolExamService.getOrCreateExam.mockResolvedValue({
            data: {
                id: 'exam-1',
                class_id: 'class-3a',
                academic_year: '2026/2027',
                president_name: 'Prof. Rossi',
                status: 'in_progress'
            }
        })
        middleSchoolExamService.listCandidates.mockResolvedValue({
            data: [
                {
                    id: 'cand-1',
                    student_id: 'stud-1',
                    student_name: 'Giacomo Leopardi',
                    admission_grade: 10,
                    is_admitted: true,
                    grade_italian: 10,
                    grade_math: 10,
                    grade_english: 10,
                    grade_second_lang: 10,
                    grade_interview: 10,
                    exam_mean: 10,
                    final_grade: 10,
                    has_honors: true,
                    outcome: 'licenziato'
                }
            ]
        })
    })

    it('renders the Middle School Exam tabs, candidates table, and evaluation controls', async () => {
        const wrapper = mount(MiddleSchoolExam, {
            global: {
                mocks: {
                    t: (key) => key
                },
                stubs: {
                    QPage: { template: '<div><slot /></div>' },
                    QCard: { template: '<div><slot /></div>' },
                    QCardSection: { template: '<div><slot /></div>' },
                    QTable: { template: '<div class="q-table-stub"><slot name="body-cell-actions" :props="{ row: {} }" /></div>' },
                    QBtn: { template: '<button @click="$emit(\'click\')"><slot /></button>' },
                    QSelect: true,
                    QBadge: { template: '<span class="badge"><slot /></span>' },
                    QChip: { template: '<span><slot /></span>' },
                    QDialog: { template: '<div><slot /></div>' },
                    QInput: true,
                    QCheckbox: true,
                    QIcon: true,
                    QTooltip: true,
                    QSpinner: true,
                    QTabs: true,
                    QTab: true,
                    QTabPanels: true,
                    QTabPanel: { template: '<div><slot /></div>' },
                    QSeparator: true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        await wrapper.vm.$nextTick()
        expect(middleSchoolExamService.getOrCreateExam).toHaveBeenCalled()
    })
})
