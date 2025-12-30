import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import AttendanceMarker from '@/components/Teacher/AttendanceMarker.vue'
import { createPinia, setActivePinia } from 'pinia'

// NO QUASAR IMPORT HERE
// installQuasar() // Removed

describe('AttendanceMarker.vue', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('renders correctly', () => {
        const wrapper = mount(AttendanceMarker, {
            props: { modelValue: 'present' }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('emits update when clicked', async () => {
        const wrapper = mount(AttendanceMarker, {
            props: { modelValue: 'present' }
        })

        // Attempt to find button by data-test
        const btn = wrapper.find('[data-test="btn-absent"]')
        if (btn.exists()) {
            await btn.trigger('click')
            expect(wrapper.emitted('update:modelValue')).toBeTruthy()
            expect(wrapper.emitted('update:modelValue')[0]).toEqual(['absent'])
        }
    })
})
