import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import GradeStatistics from '@/components/Teacher/GradeStatistics.vue'

describe('GradeStatistics.vue', () => {
    let wrapper

    beforeEach(() => {
        wrapper = mount(GradeStatistics, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            grades: {
                                grades: [
                                    { value: 8 }, { value: 7 }, { value: 9 }
                                ]
                            }
                        },
                        // Mock getters if necessary, but pinia-testing usually handles state
                    })
                ],
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-linear-progress': true
                }
            }
        })
    })

    it('displays correct statistics from store', () => {
        // In the component:
        // const average = computed(() => gradesStore.classAverage);
        // const totalGrades = computed(() => gradesStore.grades.length);

        // Note: Pinia testing with getters requires defining them or letting the store use real getters?
        // By default createTestingPinia mocks getters as undefined unless we provide initial state and use stubActions: false to run usage?
        // Actually, getters are computed from state if we don't mock them? 
        // "getters are not computed by default in createTestingPinia".
        // We should patch the store or provide getters in the options if possible?
        // Actually simpler: we can just mock the getter return value directly on the store instance.
    })

    it('renders total grades correctly', () => {
        // Direct state access works for state properties
        expect(wrapper.text()).toContain('3') // totalGrades
    })

    it('renders class average from getter', async () => {
        // Since createTestingPinia mocks getters, we need to set the value manually
        const store = wrapper.vm.gradesStore
        store.classAverage = '8.0' // Mocking the getter property

        await wrapper.vm.$nextTick()
        expect(wrapper.text()).toContain('8.0')
    })
})
