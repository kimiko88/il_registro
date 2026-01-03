import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import AttendanceMarking from '@/components/Teacher/AttendanceMarking.vue'

describe('Teacher/AttendanceMarking.vue', () => {
    it('renders correctly and toggles status', async () => {
        const wrapper = mount(AttendanceMarking, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-select': { template: '<div class="select"></div>', props: ['modelValue'] },
                    'q-date': true,
                    'q-table': {
                        template: `
              <div class="table">
                <div v-for="(row, index) in rows" :key="row.id" :class="'row-' + index">
                   <slot name="body-cell-status" :row="row"></slot>
                </div>
              </div>
            `,
                        props: ['rows', 'columns']
                    },
                    'q-td': { template: '<div><slot /></div>' },
                    'q-btn-toggle': {
                        template: '<div class="toggle" @click="$emit(\'update:modelValue\', \'absent\')"></div>',
                        props: ['modelValue']
                    }
                }
            }
        })

        // Initial state
        wrapper.vm.selectedClass = '1A'
        await wrapper.vm.$nextTick()

        // Verify table renders
        expect(wrapper.find('.table').exists()).toBe(true)

        // Find toggle for first student (Mario Rossi, present)
        const firstRowToggle = wrapper.find('.row-0 .toggle')

        // Trigger click to simulate change to 'absent' (hardcoded in stub for simplicity request)
        await firstRowToggle.trigger('click')

        // Verify local state updated
        expect(wrapper.vm.students[0].status).toBe('absent')
    })
})
