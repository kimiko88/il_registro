import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Reports from '@/pages/secretary/Reports.vue'

vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    const mockComponent = {
        template: '<div><slot /></div>'
    }
    return {
        ...actual,
        useQuasar: () => ({
            dark: { isActive: false },
            notify: vi.fn(),
            loading: { show: vi.fn(), hide: vi.fn() }
        }),
        exportFile: vi.fn(),
        QCard: mockComponent,
        QCardSection: mockComponent,
        QIcon: mockComponent,
        QBtn: mockComponent
    }
})

describe('Reports', () => {
    let wrapper

    beforeEach(() => {
        vi.useFakeTimers()
        window.open = vi.fn().mockReturnValue({
            document: {
                write: vi.fn(),
                close: vi.fn()
            },
            print: vi.fn()
        })
        
        // Create mock print-section in DOM
        const printSection = document.createElement('div')
        printSection.id = 'print-section'
        printSection.innerHTML = '<div>Mocked Print Section</div>'
        document.body.appendChild(printSection)

        wrapper = mount(Reports, {
            global: {
                provide: {
                    _q_: {
                        dark: { isActive: false },
                        loading: { show: vi.fn(), hide: vi.fn() },
                        notify: vi.fn(),
                        screen: { lt: { md: false }, gt: { xs: true } },
                        lang: { current: 'it' }
                    }
                },
                mocks: {
                    $q: {
                        dark: { isActive: false },
                        loading: { show: vi.fn(), hide: vi.fn() },
                        notify: vi.fn(),
                        screen: { lt: { md: false }, gt: { xs: true } },
                        lang: { current: 'it' }
                    }
                },
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-toolbar': { template: '<div><slot /></div>' },
                    'q-toolbar-title': { template: '<div><slot /></div>' },
                    'q-table': { template: '<div><slot /></div>' },
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
        const el = document.getElementById('print-section')
        if (el) {
            el.remove()
        }
    })

    it('opens report dialog', () => {
        wrapper.vm.openReport('grades')
        expect(wrapper.vm.currentReport).toBe('grades')
        expect(wrapper.vm.showDialog).toBe(true)
    })

    it('generates report', async () => {
        wrapper.vm.generatePDF()
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
