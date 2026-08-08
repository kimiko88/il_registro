import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import { createI18n } from 'vue-i18n'
import messages from '@/i18n'
import Login from '@/pages/Login.vue'

vi.mock('vue-router', () => ({
    useRoute: () => ({ query: {} }),
    useRouter: () => ({ push: vi.fn() })
}))

const i18n = createI18n({
    legacy: false,
    locale: 'it-IT',
    fallbackLocale: 'it-IT',
    messages
})

// Mock useAuth
const mockLogin = vi.fn()
vi.mock('@/composables/useAuth', () => ({
    useAuth: () => ({
        login: mockLogin
    })
}))

describe('Login.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        wrapper = mount(Login, {
            global: {
                plugins: [
                    i18n,
                    Quasar,
                    createTestingPinia({
                        createSpy: vi.fn,
                    })
                ],
                // Stub Quasar components to avoid complex rendering
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-input': {
                        template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
                        props: ['modelValue']
                    },
                    'q-select': {
                        template: '<select :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)"><slot /></select>',
                        props: ['modelValue', 'options']
                    },
                    'q-btn-dropdown': { template: '<div><slot /></div>' },
                    'q-icon': true,
                    'q-checkbox': true,
                    'q-btn': { template: '<button type="submit"></button>' },
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-menu': { template: '<div><slot /></div>' },
                    'q-list': { template: '<div><slot /></div>' }
                }
            }
        })
    })

    it('renders correctly', () => {
        expect(wrapper.text()).toContain('Bentornato')
        expect(wrapper.find('form').exists()).toBe(true)
    })

    it('calls login on form submit', async () => {
        // Fill inputs
        const inputs = wrapper.findAll('input')
        const emailInput = inputs[0] // First input in template
        const passInput = inputs[1]  // Second input

        await emailInput.setValue('test@example.com')
        await passInput.setValue('password123')

        // Submit form
        await wrapper.find('form').trigger('submit')

        expect(mockLogin).toHaveBeenCalledWith('test@example.com', 'password123', false)
    })

    it('shows inline error alert on login error', async () => {
        mockLogin.mockResolvedValue('Invalid credentials')

        const inputs = wrapper.findAll('input')
        await inputs[0].setValue('wrong@example.com')
        await inputs[1].setValue('wrong')

        await wrapper.find('form').trigger('submit')

        // Wait for async login to complete
        await new Promise(resolve => setTimeout(resolve, 0))

        expect(wrapper.text()).toContain('Invalid credentials')
    })
})
