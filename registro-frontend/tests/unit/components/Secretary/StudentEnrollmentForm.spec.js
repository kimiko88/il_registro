import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import StudentEnrollmentForm from '@/components/Secretary/StudentEnrollmentForm.vue'

// Mock Quasar
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            loading: { show: vi.fn(), hide: vi.fn() },
            notify: vi.fn()
        })
    }
})

describe('StudentEnrollmentForm', () => {
    let wrapper

    beforeEach(() => {
        vi.useFakeTimers()
        wrapper = mount(StudentEnrollmentForm, {
            global: {
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-stepper': {
                        template: '<div><slot /></div>',
                        methods: { next: vi.fn(), previous: vi.fn() }
                    },
                    'q-step': { template: '<div><slot /></div>' },
                    'q-stepper-navigation': { template: '<div><slot /></div>' },
                    'q-input': { template: '<div></div>', props: ['modelValue'] },
                    'q-select': { template: '<div></div>' },
                    'q-toggle': { template: '<div></div>' },
                    'q-btn': { template: '<button @click="$emit(\'click\')"></button>' },
                    'q-separator': true,
                    'q-icon': true
                }
            }
        })
        // Mock stepper ref
        wrapper.vm.stepper = { next: vi.fn(), previous: vi.fn() }
    })

    afterEach(() => {
        vi.useRealTimers()
    })

    it('initializes correctly', () => {
        expect(wrapper.vm.step).toBe(1)
        expect(wrapper.vm.form.firstName).toBe('')
    })

    it('advances step', async () => {
        wrapper.vm.nextStep()
        expect(wrapper.vm.stepper.next).toHaveBeenCalled()
    })

    it('completes enrollment', async () => {
        wrapper.vm.step = 4
        wrapper.vm.form.firstName = 'New Student'

        wrapper.vm.nextStep()

        // Should show loading
        // (Mocked Quasar logic assumes inline call to $q.loading.show())

        await vi.runAllTimersAsync()

        expect(wrapper.emitted('complete')).toBeTruthy()
        expect(wrapper.emitted('complete')[0][0].firstName).toBe('New Student')
    })
})
