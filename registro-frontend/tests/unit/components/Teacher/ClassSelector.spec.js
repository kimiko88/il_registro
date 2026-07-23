import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import ClassSelector from '@/components/Teacher/ClassSelector.vue'
import { useClassesStore } from '@/stores/classes'

describe('Teacher/ClassSelector.vue', () => {
    it('fetches classes on mount if empty', () => {
        const wrapper = mount(ClassSelector, {
            global: {
                plugins: [createTestingPinia({ stubActions: false }), Quasar],
                stubs: { 'q-icon': true }
            }
        })
        const store = useClassesStore()
        // Initial state classes is []
        // Expect fetch to be called
        expect(store.fetchAssignedClasses).toHaveBeenCalled()
    })

    it('renders selected class from store', async () => {
        const wrapper = mount(ClassSelector, {
            global: {
                plugins: [createTestingPinia({
                    initialState: {
                        classes: {
                            classes: [
                                { id: 1, name: 'Class 1A', type: 'Math' },
                                { id: 2, name: 'Class 2B', type: 'Italian' }
                            ],
                            selectedClassId: 1
                        }
                    },
                    stubActions: false
                }), Quasar],
                stubs: { 'q-icon': true }
            }
        })

        // We expect the q-select (rendered by Quasar) to show the name
        // But since q-select is complex, we might check the v-model binding indirectly or check internal component instance
        // Or inspect the computed property logic if possible? No, we check interactions.

        // Check computed property via vm (if exposed) or reactivity
        // But we are in setup script.

        // Let's modify the store and see if select updates?
        // Or better, trigger update on select and check store action.
        expect(wrapper.vm.selectedClassModel).toMatchObject({ id: 1, name: 'Class 1A', type: 'Math' })
    })

    it('calls selectClass when selection changes', async () => {
        const wrapper = mount(ClassSelector, {
            global: {
                plugins: [createTestingPinia({
                    initialState: {
                        classes: {
                            classes: [
                                { id: 1, name: 'Class 1A' },
                                { id: 2, name: 'Class 2B' }
                            ],
                            selectedClassId: 1
                        }
                    },
                    stubActions: false
                }), Quasar],
                stubs: { 'q-icon': true }
            }
        })

        const store = useClassesStore()

        // Simulate selection change
        wrapper.vm.selectedClassModel = { id: 2, name: 'Class 2B' }

        expect(store.selectClass).toHaveBeenCalledWith(2)
    })
})
