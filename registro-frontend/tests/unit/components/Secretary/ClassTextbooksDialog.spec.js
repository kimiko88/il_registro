import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import ClassTextbooksDialog from '@/components/Secretary/ClassTextbooksDialog.vue'
import textbookService from '@/services/textbookService'

vi.mock('quasar', async () => {
    const actual = await vi.importActual('quasar')
    return {
        ...actual,
        useQuasar: () => ({
            notify: vi.fn(),
            dialog: vi.fn(() => ({ onOk: vi.fn(cb => cb()) }))
        })
    }
})

vi.mock('@/services/textbookService', () => ({
    default: {
        listByClass: vi.fn(),
        getAll: vi.fn(),
        assignToClass: vi.fn(),
        removeFromClass: vi.fn(),
        create: vi.fn()
    }
}))

describe('ClassTextbooksDialog.vue', () => {
    const targetClass = {
        id: 'cls-1',
        name: '3',
        section: 'B',
        academic_year: '2026/2027'
    }

    const subjectOptions = [
        { label: 'Matematica', value: 'sub-1' },
        { label: 'Italiano', value: 'sub-2' }
    ]

    beforeEach(() => {
        vi.clearAllMocks()
        textbookService.listByClass.mockResolvedValue({
            data: [
                { id: 'ad-1', title: 'Matematica Blu', author: 'Bergamini', subject_name: 'Matematica', is_optional: false }
            ]
        })
        textbookService.getAll.mockResolvedValue({
            data: [
                { id: 'tb-1', title: 'Matematica Blu', author: 'Bergamini' },
                { id: 'tb-2', title: 'I Promessi Sposi', author: 'Manzoni' }
            ]
        })
    })

    const mountComponent = (props = {}) => {
        return mount(ClassTextbooksDialog, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-table': {
                        props: ['rows'],
                        template: `
                            <div>
                                <div v-for="row in rows" :key="row.id">
                                    <span>{{ row.title }}</span>
                                    <slot name="body-cell-actions" :props="{ row }" />
                                </div>
                            </div>
                        `
                    },
                    'q-th': { template: '<th><slot /></th>' },
                    'q-td': { template: '<td><slot /></td>' },
                    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-input': { template: '<div><slot /></div>' },
                    'q-select': { template: '<div><slot /></div>' },
                    'q-checkbox': { template: '<div><slot /></div>' },
                    'q-btn': {
                        props: ['label'],
                        template: '<button class="q-btn">{{ label }}<slot /></button>'
                    },
                    'q-icon': true,
                    'q-avatar': { template: '<div><slot /></div>' },
                    'q-space': true,
                    'q-spinner-dots': true,
                    'q-tooltip': true,
                    'q-item': true,
                    'q-item-section': true
                }
            },
            props: {
                modelValue: true,
                targetClass,
                subjectOptions,
                schoolId: 'school-1',
                ...props
            }
        })
    }

    it('renders dialog header with target class details', async () => {
        const wrapper = mountComponent()
        await wrapper.vm.$nextTick()

        expect(wrapper.text()).toContain('Adozioni Libri - Classe 3B')
        expect(wrapper.text()).toContain('2026/2027')
    })

    it('loads class textbooks and catalog on open', async () => {
        mountComponent()
        await new Promise(resolve => setTimeout(resolve, 10))

        expect(textbookService.listByClass).toHaveBeenCalledWith('cls-1')
        expect(textbookService.getAll).toHaveBeenCalled()
    })

    it('allows removing an adopted textbook and emits updated', async () => {
        textbookService.removeFromClass.mockResolvedValue({})
        const wrapper = mountComponent()
        await new Promise(resolve => setTimeout(resolve, 10))

        await wrapper.vm.removeTextbook({ id: 'ad-1' })
        expect(textbookService.removeFromClass).toHaveBeenCalledWith('ad-1')
        expect(wrapper.emitted('updated')).toBeTruthy()
    })

    it('allows adopting a textbook to class and emits updated', async () => {
        textbookService.assignToClass.mockResolvedValue({})
        const wrapper = mountComponent()
        await new Promise(resolve => setTimeout(resolve, 10))

        wrapper.vm.textbookForm.textbook_id = 'tb-2'
        wrapper.vm.textbookForm.subject_id = 'sub-2'
        wrapper.vm.textbookForm.is_optional = false

        await wrapper.vm.addTextbookToClass()
        expect(textbookService.assignToClass).toHaveBeenCalledWith('cls-1', expect.objectContaining({
            textbook_id: 'tb-2',
            subject_id: 'sub-2'
        }))
        expect(wrapper.emitted('updated')).toBeTruthy()
    })

    it('creates a new textbook into catalog and updates options', async () => {
        textbookService.create.mockResolvedValue({ data: { id: 'tb-new' } })
        const wrapper = mountComponent()

        wrapper.vm.openCreateTextbook()
        expect(wrapper.vm.showCreateTextbookDialog).toBe(true)

        wrapper.vm.newTextbook.title = 'Fisica Moderna'
        wrapper.vm.newTextbook.author = 'Amaldi'
        wrapper.vm.newTextbook.publisher = 'Zanichelli'
        wrapper.vm.newTextbook.isbn = '978-88-08-12345-6'

        await wrapper.vm.createTextbookInCatalog()
        expect(textbookService.create).toHaveBeenCalledWith(expect.objectContaining({
            title: 'Fisica Moderna',
            author: 'Amaldi',
            school_id: 'school-1'
        }))
        expect(wrapper.vm.showCreateTextbookDialog).toBe(false)
        expect(wrapper.vm.textbookForm.textbook_id).toBe('tb-new')
    })
})
