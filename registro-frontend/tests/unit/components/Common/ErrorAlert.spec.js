import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import ErrorAlert from '@/components/Common/ErrorAlert.vue'

describe('Common/ErrorAlert.vue', () => {
    it('renders nothing when error is null', () => {
        const wrapper = mount(ErrorAlert, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-banner': true,
                    'q-btn': true
                }
            },
            props: { error: null }
        })

        // q-banner has v-if="error". If stubbed, and v-if is false, it shouldn't render.
        // However, shallowMount/mount w/ stubs might behave differently if not careful.
        // If v-if is on the root node of the component (it is not, root is template?), wait.
        // Template has root q-banner. If v-if is false, wrapper text/html should be empty.
        expect(wrapper.findComponent({ name: 'q-banner' }).exists()).toBe(false)
    })

    it('renders error message when provided', () => {
        const wrapper = mount(ErrorAlert, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-banner': { template: '<div class="banner"><slot /><slot name="action" /></div>' },
                    'q-btn': { template: '<button @click="$emit(\'click\')"></button>' }
                }
            },
            props: { error: 'Something went wrong' }
        })

        expect(wrapper.find('.banner').exists()).toBe(true)
        expect(wrapper.text()).toContain('Something went wrong')
    })

    it('emits dismiss when button clicked', async () => {
        const wrapper = mount(ErrorAlert, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-banner': { template: '<div><slot name="action" /></div>' },
                    'q-btn': { template: '<button class="dismiss-btn" @click="$emit(\'click\')"></button>' }
                }
            },
            props: { error: 'Error' }
        })

        await wrapper.find('.dismiss-btn').trigger('click')
        expect(wrapper.emitted('dismiss')).toBeTruthy()
    })
})
