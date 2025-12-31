import { mount } from '@vue/test-utils'
import { describe, it, expect, vi } from 'vitest'
import StudentEnrollmentForm from 'src/components/Secretary/StudentEnrollmentForm.vue'
import { createTestingPinia } from '@pinia/testing'

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

describe('StudentEnrollmentForm', () => {
    it('renders stepper', () => {
        const wrapper = mount(StudentEnrollmentForm, {
            global: {
                plugins: [createTestingPinia()],
                stubs: {
                    'q-stepper': { template: '<div class="q-stepper"><slot /><slot name="navigation" /></div>', methods: { next: vi.fn(), previous: vi.fn() } },
                    'q-step': { template: '<div class="q-step" v-if="name===1"><slot /></div>', props: ['name', 'title'] },
                    'q-stepper-navigation': { template: '<div><slot /></div>' },
                    'q-input': { template: '<input />', props: ['modelValue', 'rules'] },
                    'q-select': true,
                    'q-btn': true,
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-toggle': true,
                    'q-separator': true
                }
            }
        })
        expect(wrapper.find('.q-stepper').exists()).toBe(true)
        expect(wrapper.find('.text-h6').text()).toContain('Nuova Iscrizione Studente')
    })

    // Add more tests for next/prev logic if possible, but stepper is complex to mock fully
})
