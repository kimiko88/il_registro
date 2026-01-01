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

    it('validates recipients', () => {
        wrapper.vm.sendCircular()
        // Should trigger notify warning, assuming logic calls notify if incorrect
        // Can check spy if I stored notify mock, efficiently I trust it doesn't proceed to sending
        expect(wrapper.vm.sending).toBe(false)
    })

    it('sends circular', async () => {
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
