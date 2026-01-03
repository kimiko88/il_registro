import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import LoadingSpinner from '@/components/Common/LoadingSpinner.vue'

describe('Common/LoadingSpinner.vue', () => {
    it('renders correctly', () => {
        const wrapper = mount(LoadingSpinner, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-spinner': { template: '<div class="spinner" :color="color" :size="size"></div>', props: ['color', 'size', 'thickness'] }
                }
            }
        })

        const spinner = wrapper.find('.spinner')
        expect(spinner.exists()).toBe(true)
        expect(spinner.attributes('color')).toBe('primary')
        expect(spinner.attributes('size')).toBe('3em')
    })
})
