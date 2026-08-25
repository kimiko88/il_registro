import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import AttendanceSummary from '@/components/Teacher/AttendanceSummary.vue'
import { useAttendanceStore } from '@/stores/attendance'

describe('Teacher/AttendanceSummary.vue', () => {
    it('renders statistics from store', () => {
        const wrapper = mount(AttendanceSummary, {
            global: {
                plugins: [createTestingPinia({
                    initialState: {
                        attendance: {
                            // These properties might be getters or state in the real store
                            // If they are getters, pinia-testing usually mocks them if not specified, 
                            // or we can set state and let getters compute if we use createPinia (but here createTestingPinia mocks actions/store)
                            // If they are state:
                            // presentCount: 5, absentCount: 2, lateCount: 1
                            // If getters: mock the getters
                        }
                    },
                    stubActions: false
                }), Quasar],
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' }
                }
            }
        })

        const store = useAttendanceStore()

        // Mock getters or state
        // If the component uses store.presentCount directly, we need to ensure the store mock has it.
        // createTestingPinia creates a mock store. We can set values on it.
        store.presentCount = 10
        store.absentCount = 3
        store.lateCount = 2

        // Trigger re-render (computed/ref update)
        return wrapper.vm.$nextTick().then(() => {
            expect(wrapper.text()).toContain('10')
            expect(wrapper.text()).toMatch(/Present|Presenti/)
            expect(wrapper.text()).toContain('3')
            expect(wrapper.text()).toMatch(/Absent|Assenti/)
            expect(wrapper.text()).toContain('2')
            expect(wrapper.text()).toMatch(/Late|Ritard/)
        })
    })
})
