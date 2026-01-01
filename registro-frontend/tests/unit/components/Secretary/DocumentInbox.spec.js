import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import DocumentInbox from '@/components/Secretary/DocumentInbox.vue'
import { createTestingPinia } from '@pinia/testing'
import { useDocumentsStore } from '@/stores/documents'

// Mock Quasar
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            notify: vi.fn(),
            dialog: vi.fn().mockImplementation(() => ({
                onOk: (fn) => fn()
            }))
        })
    }
})

describe('DocumentInbox', () => {
    let wrapper
    let store

    beforeEach(() => {
        wrapper = mount(DocumentInbox, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            documents: {
                                inbox: [
                                    { id: 1, title: 'Doc 1', status: 'Pending', favorite: false },
                                    { id: 2, title: 'Doc 2', status: 'Approved', favorite: true }
                                ],
                                loading: false
                            }
                        }
                    })
                ],
                stubs: {
                    // Simplified stub to avoid scoped slot complexity
                    'q-table': {
                        template: '<div><div v-for="row in rows" :key="row.id"><slot name="body-cell-favorite" :row="row" /><slot name="body-cell-status" :value="row.status" /><slot name="body-cell-actions" :row="row" /></div></div>',
                        props: ['rows']
                    },
                    'q-btn-toggle': true,
                    'q-input': true,
                    'q-icon': true,
                    'q-tr': { template: '<div><slot /></div>' },
                    'q-td': { template: '<div><slot /></div>' },
                    'q-space': true,
                    'q-btn': { template: '<button @click="$emit(\'click\')"></button>' },
                    'q-badge': true,
                    'q-tooltip': true
                }
            }
        })
        store = useDocumentsStore()
    })

    it('renders and exposes status color', () => {
        expect(wrapper.vm.getStatusColor('Pending')).toBe('orange')
    })

    it('toggles favorite', () => {
        const row = { favorite: false }
        wrapper.vm.toggleFavorite(row)
        expect(row.favorite).toBe(true)
    })

    it('handles batch approve', () => {
        wrapper.vm.selected = [{ id: 1 }]
        wrapper.vm.batchApprove()
        expect(wrapper.vm.selected.length).toBe(0)
    })

    it('handles batch archive', () => {
        wrapper.vm.selected = [{ id: 1 }]
        wrapper.vm.batchArchive()
        expect(wrapper.vm.selected.length).toBe(0)
    })
})
