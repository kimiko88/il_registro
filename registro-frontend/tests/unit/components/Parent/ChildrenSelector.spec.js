import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import ChildrenSelector from '@/components/Parent/ChildrenSelector.vue'
import { useChildrenStore } from '@/stores/children'

describe('ChildrenSelector.vue', () => {
    let wrapper
    let store

    beforeEach(() => {
        wrapper = mount(ChildrenSelector, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            children: {
                                children: [
                                    { id: '1', firstName: 'Mario', lastName: 'Rossi', className: '3A', avatar: 'img.jpg' },
                                    { id: '2', firstName: 'Luigi', lastName: 'Rossi', className: '1B', avatar: 'img2.jpg' }
                                ],
                                selectedChildId: '1',
                                selectedChild: { id: '1', firstName: 'Mario', lastName: 'Rossi' }
                            }
                        },
                        stubActions: false
                    })
                ],
                stubs: {
                    'q-btn-dropdown': { template: '<div><slot /></div>' },
                    'q-list': { template: '<ul><slot /></ul>' },
                    'q-item': { template: '<li @click="$emit(\'click\')"><slot /></li>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-avatar': { template: '<div><slot /></div>' }
                }
            }
        })
        store = useChildrenStore()
    })

    it('renders correctly', () => {
        expect(wrapper.text()).toContain('Mario Rossi')
        expect(wrapper.text()).toContain('3A')
    })

    it('calls fetchChildren on mount if empty', () => {
        // Re-mount with empty
        wrapper = mount(ChildrenSelector, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            children: {
                                children: []
                            }
                        },
                        stubActions: false
                    })
                ]
            }
        })
        const s = useChildrenStore()
        expect(s.fetchChildren).toHaveBeenCalled()
    })

    it('selects child', async () => {
        const items = wrapper.findAll('li')
        await items[1].trigger('click')
        expect(store.selectChild).toHaveBeenCalledWith('2')
    })
})
