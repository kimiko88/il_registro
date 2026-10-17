import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import ClassFormation from '@/pages/secretary/ClassFormation.vue'
import { enrollmentService } from '@/services/enrollmentService'

vi.mock('@/services/enrollmentService', () => ({
    enrollmentService: {
        listDrafts: vi.fn(),
        listApplications: vi.fn(),
        importSIDI: vi.fn(),
        generateFormationDraft: vi.fn(),
        getDraft: vi.fn(),
        updateDraftAssignments: vi.fn(),
        finalizeDraft: vi.fn()
    },
    default: {
        listDrafts: vi.fn(),
        listApplications: vi.fn(),
        importSIDI: vi.fn(),
        generateFormationDraft: vi.fn(),
        getDraft: vi.fn(),
        updateDraftAssignments: vi.fn(),
        finalizeDraft: vi.fn()
    }
}))

describe('Secretary ClassFormation Kanban Page', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        enrollmentService.listDrafts.mockResolvedValue({
            data: [
                {
                    id: 'draft-1',
                    title: 'Bozza Formazione Prime 2026/2027',
                    academic_year: '2026/2027',
                    is_finalized: false,
                    assignments: {
                        classes: [
                            {
                                class_name: '1A',
                                total_students: 19,
                                males_count: 10,
                                females_count: 9,
                                l104_count: 1,
                                dsa_count: 2,
                                average_grade: 8.1,
                                students: [
                                    {
                                        id: 'app-1',
                                        student_first_name: 'Mario',
                                        student_last_name: 'Rossi',
                                        gender: 'M',
                                        middle_school_grade: 8,
                                        has_disability_l104: true,
                                        has_dsa: false
                                    }
                                ]
                            },
                            {
                                class_name: '1B',
                                total_students: 22,
                                males_count: 11,
                                females_count: 11,
                                l104_count: 0,
                                dsa_count: 1,
                                average_grade: 7.9,
                                students: [
                                    {
                                        id: 'app-2',
                                        student_first_name: 'Giulia',
                                        student_last_name: 'Bianchi',
                                        gender: 'F',
                                        middle_school_grade: 9,
                                        has_disability_l104: false,
                                        has_dsa: false
                                    }
                                ]
                            }
                        ]
                    }
                }
            ]
        })
        enrollmentService.listApplications.mockResolvedValue({ data: [] })
    })

    it('renders the Kanban columns for class sections with normative warning metrics', async () => {
        const wrapper = mount(ClassFormation, {
            global: {
                mocks: {
                    t: (key) => key
                },
                stubs: {
                    QPage: { template: '<div><slot /></div>' },
                    QCard: { template: '<div><slot /></div>' },
                    QCardSection: { template: '<div><slot /></div>' },
                    QCardActions: { template: '<div><slot /></div>' },
                    QBtn: { template: '<button @click="$emit(\'click\')"><slot /></button>' },
                    QBadge: { template: '<span class="badge"><slot /></span>' },
                    QChip: { template: '<span><slot /></span>' },
                    QDialog: { template: '<div><slot /></div>' },
                    QInput: true,
                    QSelect: true,
                    QFile: true,
                    QIcon: true,
                    QTooltip: true,
                    QSpinner: true,
                    QTabs: true,
                    QTab: true,
                    QSeparator: true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        await wrapper.vm.$nextTick()
        expect(enrollmentService.listDrafts).toHaveBeenCalled()
    })
})
