import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Textbooks from '@/pages/secretary/Textbooks.vue'
import { textbookService } from '@/services/textbookService'

vi.mock('@/services/textbookService', () => ({
    textbookService: {
        getAll: vi.fn().mockResolvedValue({ data: [] }),
        getSpendingReport: vi.fn(),
        listClassAdoptions: vi.fn(),
        importAIE: vi.fn(),
        exportClassAIE: vi.fn(),
        searchAIECatalog: vi.fn(),
        adoptBook: vi.fn(),
        deleteAdoption: vi.fn()
    },
    default: {
        getAll: vi.fn().mockResolvedValue({ data: [] }),
        getSpendingReport: vi.fn(),
        listClassAdoptions: vi.fn(),
        importAIE: vi.fn(),
        exportClassAIE: vi.fn(),
        searchAIECatalog: vi.fn(),
        adoptBook: vi.fn(),
        deleteAdoption: vi.fn()
    }
}))

vi.mock('@/stores/classes', () => ({
    useClassesStore: () => ({
        classes: [
            { id: 'class-1', name: '1A', section: 'A' },
            { id: 'class-2', name: '2B', section: 'B' }
        ],
        loading: false,
        fetchClasses: vi.fn().mockResolvedValue()
    })
}))

vi.mock('@/services/adminService', () => ({
    default: {
        getClassSubjects: vi.fn().mockResolvedValue({ data: [{ subject_id: 'sub-1', subject_name: 'Matematica' }] }),
        getSubjects: vi.fn().mockResolvedValue({ data: [{ id: 'sub-1', name: 'Matematica' }] })
    }
}))

describe('Secretary Textbooks Page with AIE & Spending Limits', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        textbookService.getSpendingReport.mockResolvedValue({
            data: {
                total_spending: 295.50,
                spending_limit: 300.00,
                tolerance_threshold: 330.00,
                difference: -4.50,
                percentage: 98.5,
                status: 'WITHIN_LIMIT',
                items: [
                    {
                        id: 'adopt-1',
                        book_title: 'Matematica.blu 2.0',
                        price: 28.90,
                        adoption_type: 'nuova_adozione',
                        is_already_owned: false
                    }
                ]
            }
        })
        textbookService.listClassAdoptions.mockResolvedValue({
            data: [
                {
                    id: 'adopt-1',
                    subject_name: 'Matematica',
                    book_title: 'Matematica.blu 2.0',
                    price: 28.90,
                    adoption_type: 'nuova_adozione',
                    is_already_owned: false
                }
            ]
        })
    })

    it('renders the spending limits card and shows within limit status badge', async () => {
        const wrapper = mount(Textbooks, {
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
                    QSelect: { template: '<select><slot /></select>' },
                    QLinearProgress: true,
                    QCircularProgress: { template: '<div class="circular-progress"><slot /></div>' },
                    QBadge: { template: '<span class="badge"><slot /></span>' },
                    QChip: { template: '<span><slot /></span>' },
                    QDialog: { template: '<div><slot /></div>' },
                    QTooltip: true,
                    QIcon: true,
                    QForm: { template: '<form><slot /></form>' },
                    QInput: true,
                    QFile: true,
                    QTabs: true,
                    QTab: true,
                    QTabPanels: true,
                    QTabPanel: { template: '<div><slot /></div>' },
                    QAvatar: true,
                    QSpace: true,
                    QSeparator: true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        await wrapper.vm.$nextTick()
        expect(textbookService.getAll).toHaveBeenCalled()
    })

    it('defines openAdoptionDialog on instance and opens adoption dialog', async () => {
        textbookService.getAll.mockResolvedValue({
            data: [
                { id: 'book-1', title: 'Matematica.blu', author: 'Bergamini', price: 25.0, isbn: '123' }
            ]
        })
        const wrapper = mount(Textbooks, {
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
                    QSelect: { template: '<select><slot /></select>' },
                    QLinearProgress: true,
                    QCircularProgress: { template: '<div class="circular-progress"><slot /></div>' },
                    QBadge: { template: '<span class="badge"><slot /></span>' },
                    QChip: { template: '<span><slot /></span>' },
                    QDialog: { template: '<div><slot /></div>' },
                    QTooltip: true,
                    QIcon: true,
                    QForm: { template: '<form><slot /></form>' },
                    QInput: true,
                    QFile: true,
                    QTabs: true,
                    QTab: true,
                    QTabPanels: true,
                    QTabPanel: { template: '<div><slot /></div>' },
                    QAvatar: true,
                    QSpace: true,
                    QSeparator: true,
                    QCheckbox: true
                }
            }
        })

        await wrapper.vm.$nextTick()
        expect(typeof wrapper.vm.openAdoptionDialog).toBe('function')
        wrapper.vm.selectedClass = 'class-1'
        await wrapper.vm.openAdoptionDialog()
        expect(wrapper.vm.showAdoptionDialog).toBe(true)

        wrapper.vm.adoptionForm.book_id = 'book-1'
        wrapper.vm.adoptionForm.subject_id = 'sub-1'
        await wrapper.vm.submitAdoption()
        expect(textbookService.adoptBook).toHaveBeenCalledWith('class-1', expect.objectContaining({
            book_id: 'book-1',
            subject_id: 'sub-1'
        }))
        expect(wrapper.vm.showAdoptionDialog).toBe(false)
    })
})
