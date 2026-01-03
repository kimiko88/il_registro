import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar, Notify } from 'quasar'
import Login from '@/pages/Login.vue'

// Mock useAuth
const mockLogin = vi.fn()
vi.mock('@/composables/useAuth', () => ({
    useAuth: () => ({
        login: mockLogin
    })
}))

// Mock Quasar Notify
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        Notify: {
            create: vi.fn()
        }
    }
})

describe('Login.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        wrapper = mount(Login, {
            global: {
                plugins: [
                    [Quasar, {}],
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
                    'q-icon': true,
                    'q-checkbox': true,
                    'q-btn': { template: '<button type="submit"></button>' }
                }
            }
        })
    })

    it('renders correctly', () => {
        expect(wrapper.text()).toContain('Welcome Back')
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

        expect(mockLogin).toHaveBeenCalledWith('test@example.com', 'password123')
    })

    it('shows notification on login error', async () => {
        mockLogin.mockResolvedValue('Invalid credentials')

        const inputs = wrapper.findAll('input')
        await inputs[0].setValue('wrong@example.com')
        await inputs[1].setValue('wrong')

        await wrapper.find('form').trigger('submit')

        // Wait for async login to complete
        await new Promise(resolve => setTimeout(resolve, 0))

        expect(Notify.create).toHaveBeenCalledWith(expect.objectContaining({
            type: 'negative',
            message: 'Invalid credentials'
        }))
    })
})
