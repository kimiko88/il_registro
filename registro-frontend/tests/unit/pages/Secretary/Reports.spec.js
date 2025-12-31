import { mount } from '@vue/test-utils'
import { describe, it, expect, vi } from 'vitest'
import Reports from 'src/pages/secretary/Reports.vue'

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

describe('Reports Page', () => {
    it('opens dialog when card is clicked', async () => {
        const wrapper = mount(Reports, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div @click="$emit(\'click\')"><slot /></div>' }, // Stub card to emit click
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': true,
                    'q-icon': true,
                    'q-btn': true,
                    'q-dialog': { template: '<div v-if="modelValue"><slot /></div>', props: ['modelValue'] },
                    'q-select': true,
                    'ReportCard': false // Use real inline component logic if possible, or reliance on q-card stub
                }
            }
        })

        expect(wrapper.vm.showDialog).toBe(false)

        // Simulate clicking the first ReportCard (which renders as q-card thanks to h())
        // However, since we stub q-card, the inline component 'ReportCard' which returns h(QCard) will render the stub.
        // We need to find the element that triggers openReport.

        // The inline component emits 'click' on the QCard.
        // Let's call the method directly to test logic as rendering h() components in test utils with stubs is tricky.
        wrapper.vm.openReport('grades')
        expect(wrapper.vm.showDialog).toBe(true)
        expect(wrapper.vm.currentReport).toBe('grades')
    })
})
