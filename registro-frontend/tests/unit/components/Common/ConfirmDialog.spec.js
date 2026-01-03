import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import ConfirmDialog from '@/components/Common/ConfirmDialog.vue'

describe('Common/ConfirmDialog.vue', () => {
    it('renders correctly with default message', () => {
        const wrapper = mount(ConfirmDialog, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-dialog': { template: '<div class="q-dialog-stub"><slot /></div>', props: ['modelValue'] },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-btn': true,
                    'q-avatar': true
                }
            },
            props: {
                modelValue: true
            }
        })

        expect(wrapper.text()).toContain('Are you sure?')
        // Check existence via class we added to stub
        expect(wrapper.find('.q-dialog-stub').exists()).toBe(true)
    })

    it('renders custom message', () => {
        const wrapper = mount(ConfirmDialog, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-btn': true,
                    'q-avatar': true
                }
            },
            props: {
                modelValue: true,
                message: 'Delete this item?'
            }
        })
        expect(wrapper.text()).toContain('Delete this item?')
    })

    it('emits update:modelValue when dialog closes', async () => {
        const wrapper = mount(ConfirmDialog, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-dialog': {
                        template: '<div class="q-dialog-stub" @click="$emit(\'update:modelValue\', false)"><slot /></div>',
                        props: ['modelValue'],
                        emits: ['update:modelValue']
                    },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-btn': true,
                    'q-avatar': true
                }
            },
            props: { modelValue: true }
        })

        // Simulate dialog emitting update
        await wrapper.find('.q-dialog-stub').trigger('click')
        expect(wrapper.emitted('update:modelValue')).toBeTruthy()
        expect(wrapper.emitted('update:modelValue')[0]).toEqual([false])
    })

    it('emits confirm event when confirm button clicked', async () => {
        const wrapper = mount(ConfirmDialog, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-btn': {
                        template: '<button class="confirm-btn" @click="$emit(\'click\')">{{label}}</button>',
                        props: ['label']
                    },
                    'q-avatar': true
                },
                directives: {
                    'close-popup': {}
                }
            },
            props: { modelValue: true }
        })

        const buttons = wrapper.findAll('.confirm-btn')
        const confirmBtn = buttons.find(b => b.text() === 'Confirm')
        await confirmBtn.trigger('click')

        expect(wrapper.emitted('confirm')).toBeTruthy()
    })
})
