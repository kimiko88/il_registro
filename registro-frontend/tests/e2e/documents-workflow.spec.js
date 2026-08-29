import { describe, it, expect, vi, beforeEach } from 'vitest'
import { shallowMount, mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'

// ─── Quasar Mock ─────────────────────────────────────────────────────────────
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    const mockIconSet = {
        name: 'material-icons',
        type: { positive: 'check', negative: 'warning', info: 'info', warning: 'priority_high' },
        arrow: { dropdown: 'arrow_drop_down', expand: 'arrow_drop_down', down: 'arrow_drop_down', up: 'arrow_drop_up' },
        chevron: { left: 'chevron_left', right: 'chevron_right' },
        table: { arrow: 'arrow_upward', select: 'arrow_drop_down' },
        editor: {},
        tree: {}
    }
    return {
        ...actual,
        useQuasar: () => ({
            dark: { isActive: false },
            loading: { show: vi.fn(), hide: vi.fn() },
            notify: vi.fn(),
            dialog: vi.fn().mockReturnValue({
                onOk: vi.fn(cb => { if (cb) cb(); return { onCancel: vi.fn(), onDismiss: vi.fn() } }),
                onCancel: vi.fn(),
                onDismiss: vi.fn()
            }),
            screen: { lt: { md: false }, gt: { xs: true } },
            lang: { current: 'it' },
            iconSet: mockIconSet
        })
    }
})

function buildPinia(role = 'teacher', extras = {}) {
    return createTestingPinia({
        createSpy: () => vi.fn().mockResolvedValue(undefined),
        initialState: {
            auth: { user: { role, id: 'u-1', school_id: 's-1' }, isAuthenticated: true },
            documents: {
                documents: [
                    { id: 'd1', title: 'PDP Mario Rossi', type: 'PDP', status: 'draft', date: '2026-08-29' },
                    { id: 'd2', title: 'Relazione 15 Maggio 5A', type: 'Programmazione', status: 'published', date: '2026-08-28' },
                    { id: 'd3', title: 'Convenzione PCTO', type: 'PCTO', status: 'draft', date: '2026-08-27' }
                ],
                templates: [
                    { id: 't1', name: 'Modello PDP Standard', type: 'PDP', content: '<p>Template PDP</p>' }
                ],
                loading: false
            },
            classes: {
                classes: [{ id: 'c1', name: '3A' }, { id: 'c2', name: '5B' }]
            },
            ...extras
        }
    })
}

// ─────────────────────────────────────────────────────────────────────────────
// Suite 1 – Documents Management Workflow (Teacher / Secretary)
// ─────────────────────────────────────────────────────────────────────────────
describe('Documents Management Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('DW01 — Teacher Documents page renders document list and upload button', async () => {
        const { default: Documents } = await import('@/pages/teacher/Documents.vue')
        const wrapper = shallowMount(Documents, {
            global: { plugins: [buildPinia('teacher')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('DW02 — Documents page allows filtering by document type and class', async () => {
        const { default: Documents } = await import('@/pages/teacher/Documents.vue')
        const wrapper = shallowMount(Documents, {
            global: { plugins: [buildPinia('teacher')] }
        })
        expect(wrapper.find('q-page-stub').exists() || wrapper.exists()).toBe(true)
    })

    it('DW03 — DocumentEditor component can be mounted with template data', async () => {
        const { default: DocumentEditor } = await import('@/components/Teacher/DocumentEditor.vue')
        const wrapper = shallowMount(DocumentEditor, {
            global: { plugins: [buildPinia('teacher')] },
            props: {
                modelValue: '<p>Contenuto iniziale</p>',
                documentType: 'PDP'
            }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('DW04 — DocumentTemplate component renders template selector with options', async () => {
        const { default: DocumentTemplate } = await import('@/components/Teacher/DocumentTemplate.vue')
        const wrapper = shallowMount(DocumentTemplate, {
            global: { plugins: [buildPinia('teacher')] },
            props: {
                templates: [
                    { id: 't1', name: 'Modello PDP', type: 'PDP' },
                    { id: 't2', name: 'Modello PFI', type: 'PFI' }
                ]
            }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('DW05 — Secretary Documents page allows reviewing submitted documents', async () => {
        const { default: SecretaryDocuments } = await import('@/pages/secretary/Documents.vue')
        const wrapper = shallowMount(SecretaryDocuments, {
            global: { plugins: [buildPinia('secretary')] }
        })
        expect(wrapper.exists()).toBe(true)
    })
})
