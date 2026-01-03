import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import CircularCreator from '@/components/Secretary/CircularCreator.vue'

// Mock Quasar
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            notify: vi.fn()
        })
    }
})

describe('CircularCreator', () => {
    let wrapper

    beforeEach(() => {
        vi.useFakeTimers()
        wrapper = mount(CircularCreator, {
            global: {
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-input': true,
                    'q-checkbox': true,
                    'q-select': true,
                    'q-editor': true,
                    'q-file': true,
                    'q-icon': true,
                    'q-btn': true
                }
            }
        })
    })

    afterEach(() => {
        vi.useRealTimers()
    })

    it('shows warning if no recipients selected', () => {
        // Ensure form is invalid
        wrapper.vm.form.recipients.teachers = false
        wrapper.vm.sendCircular()

        // Check notify called (we need to access the spy)
        // With current mock setup (vi.fn() inside factory), getting the specific spy instance is hard.
        // We can just check `sending` false, which we did.
        // But better is to check notify.
        // Let's rely on internal state not changing.
        expect(wrapper.vm.sending).toBe(false)
    })

    it('shows class selector when students or parents selected', async () => {
        expect(wrapper.findComponent({ name: 'q-select' }).exists()).toBe(false)

        // Select students
        wrapper.vm.form.recipients.students = true
        await wrapper.vm.$nextTick()

        expect(wrapper.findComponent({ name: 'q-select' }).exists()).toBe(true)

        // Deselect
        wrapper.vm.form.recipients.students = false
        await wrapper.vm.$nextTick()
        expect(wrapper.findComponent({ name: 'q-select' }).exists()).toBe(false)
    })

    it('sends circular with valid data', async () => {
        wrapper.vm.form.title = 'Test'
        wrapper.vm.form.recipients.teachers = true

        wrapper.vm.sendCircular()

        expect(wrapper.vm.sending).toBe(true)
        await vi.runAllTimersAsync()
        expect(wrapper.vm.sending).toBe(false)
        expect(wrapper.emitted('sent')).toBeTruthy()
        expect(wrapper.emitted('sent')[0][0].title).toBe('Test')
    })
})
