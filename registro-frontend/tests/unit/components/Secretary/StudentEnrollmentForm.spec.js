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
    // Define spies externally so we can assert on them regardless of internal ref overwrites
    const nextSpy = vi.fn()
    const prevSpy = vi.fn()

    beforeEach(() => {
        vi.useFakeTimers()
        nextSpy.mockClear()
        prevSpy.mockClear()

        wrapper = mount(StudentEnrollmentForm, {
            global: {
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-stepper': {
                        // Ensure navigation slot is rendered
                        template: '<div><slot /><slot name="navigation" /></div>',
                        methods: { next: nextSpy, previous: prevSpy }
                    },
                    'q-step': { template: '<div><slot /></div>' },
                    // Ensure class is present for finding
                    'q-stepper-navigation': { template: '<div class="stepper-nav"><slot /></div>' },
                    'q-input': { template: '<div></div>', props: ['modelValue'] },
                    'q-select': { template: '<div></div>' },
                    'q-toggle': { template: '<div></div>' },
                    // Ensure buttons are easily found by tag 'button' since component name might be inferred differently
                    'q-btn': { template: '<button @click="$emit(\'click\')"></button>', name: 'q-btn' },
                    'q-separator': true,
                    'q-icon': true
                }
            }
        })
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
        expect(nextSpy).toHaveBeenCalled()
    })

    it('navigates back', async () => {
        wrapper.vm.step = 2
        await wrapper.vm.$nextTick() // Ensure v-if updates

        // Find buttons inside navigation
        const nav = wrapper.find('.stepper-nav')
        expect(nav.exists()).toBe(true)

        const btns = nav.findAll('button')
        // Button 0 is Next, Button 1 is Back
        expect(btns.length).toBe(2)
        await btns[1].trigger('click')

        expect(prevSpy).toHaveBeenCalled()
    })

    it('completes enrollment with data', async () => {
        wrapper.vm.step = 4
        wrapper.vm.form.firstName = 'New Student'

        // Trigger generic "next" which becomes "complete" at step 4
        await wrapper.vm.nextStep()

        await vi.runAllTimersAsync()

        expect(wrapper.emitted('complete')).toBeTruthy()
        expect(wrapper.emitted('complete')[0][0].firstName).toBe('New Student')
    })
})
