import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Reports from '@/pages/secretary/Reports.vue'
import { h } from 'vue'
import { QCard } from 'quasar'

// Mock Quasar
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            loading: { show: vi.fn(), hide: vi.fn() },
            notify: vi.fn()
        }),
        QCard: { template: '<div><slot /></div>' },
        QCardSection: { template: '<div><slot /></div>' },
        QIcon: { template: '<div></div>' },
        QBtn: { template: '<div></div>' }
    }
})

describe('Reports', () => {
    let wrapper

    beforeEach(() => {
        vi.useFakeTimers()
        wrapper = mount(Reports, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-select': true,
                    'q-btn': true,
                    // Note: ReportCard is defined inside script setup, so it's a local component.
                    // We might need to rely on stubbing QCard which it renders.
                }
            }
        })
    })

    afterEach(() => {
        vi.useRealTimers()
    })

    it('opens report dialog', () => {
        wrapper.vm.openReport('grades')
        expect(wrapper.vm.currentReport).toBe('grades')
        expect(wrapper.vm.showDialog).toBe(true)
    })

    it('generates report', async () => {
        wrapper.vm.generate('pdf')
        // Mock loading show should be called
        // Since we didn't spy explicitly on the mock returned by useQuasar here easily (unless we exported the mock), 
        // we mainly check state changes after timeout.

        await vi.runAllTimersAsync()

        expect(wrapper.vm.showDialog).toBe(false)
    })

    // Testing computed properties indirectly via rendered output or exposed state if we could.
    // Exposed state 'info' can be manipulated.

    it('ReportCard component renders', () => {
        // Since ReportCard is internal, we can check if QCard with text exists in DOM.
        expect(wrapper.text()).toContain('Centro Reportistica')
        expect(wrapper.text()).toContain('Pagelle / Scrutini')
    })
})
