
import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import DocumentInbox from '@/components/Secretary/DocumentInbox.vue'

// Mock Store
const mockStore = {
    inbox: [
        { id: 1, title: 'Test Doc 1', status: 'Pending', favorite: false, type: 'PDP' },
        { id: 2, title: 'Test Doc 2', status: 'Approved', favorite: true, type: 'Certificate' }
    ],
    loading: false,
    pagination: {},
    fetchInbox: vi.fn()
}

vi.mock('src/stores/documents', () => ({
    useDocumentsStore: () => mockStore
}))

// Mock Quasar
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            notify: vi.fn(),
            dialog: vi.fn().mockReturnValue({ onOk: (fn) => fn() })
        })
    }
})

describe('DocumentInbox.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        wrapper = mount(DocumentInbox, {
            global: {
                stubs: {
                    'q-table': {
                        template: '<div><slot name="top" /><slot name="top-row" /><div class="rows"><div v-for="row in rows" :key="row.id" class="row-item">{{row.title}}</div></div></div>',
                        props: ['rows', 'columns', 'loading', 'selection', 'filter', 'selected'],
                        emits: ['update:selected']
                    },
                    'q-btn': true,
                    'q-btn-toggle': true,
                    'q-input': true,
                    'q-icon': true,
                    'q-tr': { template: '<tr><slot /></tr>' },
                    'q-td': { template: '<td><slot /></td>' },
                    'q-badge': true,
                    'q-space': true
                }
            }
        })
    })

    it('renders inbox items', () => {
        expect(wrapper.text()).toContain('Test Doc 1')
        expect(wrapper.text()).toContain('Test Doc 2')
    })

    it('shows batch actions when items selected', async () => {
        // Simulate selection by updating 'selected' ref in component
        // Since q-table is stubbed and we sync v-model:selected, we can update wrapper.vm.selected
        wrapper.vm.selected = [mockStore.inbox[0]]
        await wrapper.vm.$nextTick()

        // Check for Export/Approve button stub
        const approveBtn = wrapper.findAll('q-btn-stub').find(w => w.attributes('label') === 'Approva Selezionati')
        expect(approveBtn).toBeDefined()
        expect(approveBtn.exists()).toBe(true)
    })

    it('batch approve triggers dialog', async () => {
        wrapper.vm.selected = [mockStore.inbox[0]]
        await wrapper.vm.$nextTick()

        // Find button? q-btn is true stub.
        // Call method directly
        wrapper.vm.batchApprove()

        // Logic inside batchApprove calls dialog.onOk -> updates status
        // Since dialog auto-confirms in mock, status should change?
        // But store.inbox is a mock object.
        expect(wrapper.vm.selected.length).toBe(0) // Logic clears selection
    })
})
