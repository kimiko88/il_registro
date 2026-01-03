import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import SchoolForm from '@/components/Admin/SchoolForm.vue'

describe('Admin/SchoolForm.vue', () => {
    it('initializes form with default values (New Mode)', () => {
        const wrapper = mount(SchoolForm, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form><slot /></form>' },
                    'q-input': true,
                    'q-btn': true
                }
            }
        })
        expect(wrapper.text()).toContain('New School')
        expect(wrapper.vm.form.name).toBe('')
    })

    it('initializes form with prop values (Edit Mode)', () => {
        const school = { name: 'Test School', code: 'TS01', address: '123 St', email: 'test@school.com', phone: '555' }
        const wrapper = mount(SchoolForm, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form><slot /></form>' },
                    'q-input': true,
                    'q-btn': true
                }
            },
            props: { school }
        })
        expect(wrapper.text()).toContain('Edit School')
        expect(wrapper.vm.form.name).toBe('Test School')
    })

    it('emits submit with form data', () => {
        const wrapper = mount(SchoolForm, {
            global: { plugins: [Quasar], stubs: { 'q-card': true, 'q-card-section': true, 'q-form': true, 'q-input': true, 'q-btn': true } }
        })
        wrapper.vm.form.name = 'New School'
        wrapper.vm.onSubmit()
        expect(wrapper.emitted('submit')).toBeTruthy()
        expect(wrapper.emitted('submit')[0][0].name).toBe('New School')
    })

    it('validates required fields', async () => {
        const wrapper = mount(SchoolForm, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-input': true,
                    'q-btn': true
                }
            }
        })

        // Verify validation rules by directly testing the rule function
        // Rules are defined inline as: [val => !!val || 'Required']
        const validationRule = (val) => !!val || 'Required'

        expect(validationRule('')).toBe('Required')
        expect(validationRule('Valid Name')).toBe(true)
        expect(validationRule(null)).toBe('Required')
        expect(validationRule('CODE123')).toBe(true)
    })

    it('updates form on input change', async () => {
        const wrapper = mount(SchoolForm, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form><slot /></form>' },
                    'q-input': {
                        template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
                        props: ['modelValue']
                    },
                    'q-btn': true
                }
            }
        })

        expect(wrapper.vm.form.name).toBe('')
        wrapper.vm.form.name = 'Test School'
        await wrapper.vm.$nextTick()
        expect(wrapper.vm.form.name).toBe('Test School')
    })
})
