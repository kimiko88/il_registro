import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import ConfigurationPanel from '@/components/Common/ConfigurationPanel.vue'

describe('Common/ConfigurationPanel.vue', () => {
    it('renders title and slot content', () => {
        const wrapper = mount(ConfigurationPanel, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-btn': { template: '<button class="save-btn" :disabled="loading"></button>', props: ['loading'] }
                }
            },
            props: {
                title: 'My Config',
                loading: false
            },
            slots: {
                default: '<div class="test-content">Settings</div>'
            }
        })

        expect(wrapper.text()).toContain('My Config')
        expect(wrapper.find('.test-content').exists()).toBe(true)
        expect(wrapper.find('.test-content').text()).toBe('Settings')
    })

    it('emits submit event', async () => {
        const wrapper = mount(ConfigurationPanel, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form class="my-form" @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-btn': true
                }
            },
            props: { title: 'Test' }
        })

        await wrapper.find('.my-form').trigger('submit')
        expect(wrapper.emitted('submit')).toBeTruthy()
    })

    it('passes loading prop to button', () => {
        const wrapper = mount(ConfigurationPanel, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form><slot /></form>' },
                    'q-btn': { template: '<button class="save-btn" :disabled="loading"></button>', props: ['loading'] }
                }
            },
            props: {
                title: 'Test',
                loading: true
            }
        })

        // Check prop on stub
        const btn = wrapper.find('.save-btn')
        // Check raw attribute on element
        expect(btn.element.disabled).toBe(true)
    })
})
