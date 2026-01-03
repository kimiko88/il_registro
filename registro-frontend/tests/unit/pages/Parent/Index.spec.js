import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import ParentIndex from '@/pages/parent/Index.vue'
import { useParentStore } from '@/stores/parent'

describe('Parent/Index.vue', () => {
    let wrapper
    let store

    const mockChildren = [
        { id: 'c1', firstName: 'Luigi', lastName: 'Verdi', class: '3A' },
        { id: 'c2', firstName: 'Anna', lastName: 'Verdi', class: '1B' }
    ]

    beforeEach(() => {
        wrapper = mount(ParentIndex, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            parent: {
                                children: [],
                                loading: false,
                                selectedChildId: null
                            }
                        },
                        stubActions: false
                    })
                ],
                // Stub Quasar components that might be complex
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-btn': true,
                    'q-btn-dropdown': { template: '<div><slot /></div>' },
                    'q-list': { template: '<ul><slot /></ul>' },
                    'q-item': { template: '<li><slot /></li>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-avatar': true,
                    'q-icon': true,
                    'q-timeline': true,
                    'q-timeline-entry': true,
                    'q-spinner': true,
                    'q-page-sticky': true,
                    'q-fab': true,
                    'q-fab-action': true,
                    'router-link': true
                }
            }
        })
        store = useParentStore()
    })

    it('fetches children on mount if empty', () => {
        expect(store.fetchChildren).toHaveBeenCalled()
    })

    it('renders loading spinner', async () => {
        store.loading = true
        await wrapper.vm.$nextTick()
        const spinner = wrapper.findComponent({ name: 'q-spinner' })
        expect(spinner.exists()).toBe(true)
    })

    it('renders empty state when no children', async () => {
        store.loading = false
        store.children = []
        await wrapper.vm.$nextTick()
        expect(wrapper.text()).toContain('Nessun figlio associato')
    })

    it('renders dashboard when child selected', async () => {
        store.children = mockChildren
        store.selectedChildId = 'c1'
        // Computed selectedChild depends on the store logic.
        // If we use createTestingPinia with stubActions: false, the getters should work if state is set.
        // However, selectedChild is likely a getter. 
        // Pinia testing with getters can be tricky if not computed automatically.
        // We can force the getter value by mocking it or if logic is simple it works.
        // store.selectedChild is a getter finding child by selectedChildId.

        // Attempt to verify if getter works or we need to patch it.
        // If it fails, we can patch the store instance.

        // Re-mount with state
        wrapper = mount(ParentIndex, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            parent: {
                                children: mockChildren,
                                selectedChildId: 'c1'
                            }
                        },
                        // We need the real store logic usually for getters
                        stubActions: false
                    })
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'router-link': true,
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-btn-dropdown': { template: '<div><slot /></div>' },
                    'q-page-sticky': true,
                    'q-fab': true,
                    'q-fab-action': true,
                    'q-list': true, // Stub list to avoid rendering children items which might need more stubs
                    'q-timeline': true
                }
            }
        })

        // In many setups, getters are not reactive in testing pinia unless we use real store.
        // But let's see. If dashboard content depends on `selectedChild` ref from `storeToRefs`.
        // The component uses: const { children, selectedChild, ... } = storeToRefs(parentStore)

        // Let's manually mock the getter return if needed.
        // For now, let's assume Pinia testing handles it or we manually set it.

        const pStore = useParentStore()
        // Mock the getter result if needed
        // pStore.selectedChild = mockChildren[0] // This works because store is a reactive object in mock

        // Wait, createTestingPinia mocks actions but preserves getters? No, it usually mocks everything.
        // If strictly mocking, we need to set the value.

        // Let's just set the property directly on the store mock, assuming the component uses it.
        // storeToRefs reads from the store.
        pStore.selectedChild = mockChildren[0]
        await wrapper.vm.$nextTick()

        expect(wrapper.text()).toContain('Media Voti')
        expect(wrapper.text()).toContain('Prossimi Eventi')
    })
})
