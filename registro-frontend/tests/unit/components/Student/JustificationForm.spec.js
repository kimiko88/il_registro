import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import JustificationForm from '@/components/Student/JustificationForm.vue'

describe('Student/JustificationForm.vue', () => {
    it('renders correctly', async () => {
        const wrapper = mount(JustificationForm, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-dialog': { template: '<div class="dialog"><slot /></div>', props: ['modelValue'] },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-select': { template: '<input class="q-select" />', props: ['modelValue'] },
                    'q-input': { template: '<input class="q-input" />', props: ['modelValue'] },
                    'q-btn': true
                }
            },
            props: {
                modelValue: true,
                date: '2023-01-01'
            }
        })

        // Force open
        wrapper.vm.isOpen = true
        await wrapper.vm.$nextTick()

        // DEBUG: Fail intentionally to see HTML
        // expect(wrapper.html()).toBe('DEBUG')

        // If correct, these should exist
        expect(wrapper.find('.q-select').exists()).toBe(true)

        // Emit check
        await wrapper.find('form').trigger('submit')
        expect(wrapper.emitted('submit')).toBeTruthy()
    })
})
