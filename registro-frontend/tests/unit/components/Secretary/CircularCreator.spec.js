import { mount } from '@vue/test-utils'
import { describe, it, expect, vi } from 'vitest'
import CircularCreator from 'src/components/Secretary/CircularCreator.vue'

const mockNotify = vi.fn()
const mockLoading = { show: vi.fn(), hide: vi.fn() }

vi.mock('quasar', async () => {
    const actual = await vi.importActual('quasar')
    return {
        ...actual,
        useQuasar: () => ({
            loading: mockLoading,
            notify: mockNotify
        })
    }
})

describe('CircularCreator', () => {
    it('emits sent event on valid submission', async () => {
        const wrapper = mount(CircularCreator, {
            global: {
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-input': { template: '<input />', props: ['modelValue'] },
                    'q-select': true,
                    'q-file': true,
                    'q-checkbox': { template: '<input type="checkbox" />', props: ['modelValue'], emits: ['update:modelValue'] },
                    'q-editor': true,
                    'q-btn': { template: '<button @click="$emit(\'click\')"></button>', props: ['type'] },
                    'q-icon': true
                }
            }
        })

        // Fill form
        wrapper.vm.form.title = "Test Circular"
        wrapper.vm.form.recipients.teachers = true // Select recipients

        await wrapper.find('form').trigger('submit')

        // Wait for async
        await new Promise(r => setTimeout(r, 1600))

        expect(wrapper.emitted('sent')).toBeTruthy()
        expect(wrapper.emitted('sent')[0][0].title).toBe("Test Circular")
    })
})
