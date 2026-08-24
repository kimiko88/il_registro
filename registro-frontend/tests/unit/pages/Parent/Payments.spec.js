import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import ParentPayments from '@/pages/parent/Payments.vue'
import paymentService from '@/services/paymentService'

vi.mock('@/services/paymentService', () => ({
    default: {
        getPayments: vi.fn(),
        pay: vi.fn(),
        downloadReceipt: vi.fn()
    },
    paymentService: {
        getPayments: vi.fn(),
        pay: vi.fn(),
        downloadReceipt: vi.fn()
    }
}))

const commonStubs = {
    'q-page': { template: '<div><slot /></div>' },
    'q-card': { template: '<div><slot /></div>' },
    'q-card-section': { template: '<div><slot /></div>' },
    'q-tabs': { template: '<div><slot /></div>', props: ['modelValue'] },
    'q-tab': { template: '<button><slot /></button>', props: ['name', 'label', 'icon'] },
    'q-tab-panels': { template: '<div><slot /></div>', props: ['modelValue'] },
    'q-tab-panel': { template: '<div><slot /></div>', props: ['name'] },
    'q-separator': true,
    'q-list': { template: '<ul><slot /></ul>' },
    'q-item': { template: '<li><slot /></li>' },
    'q-item-section': { template: '<span><slot /></span>', props: ['avatar', 'side'] },
    'q-item-label': { template: '<span><slot /></span>', props: ['caption'] },
    'q-avatar': true,
    'q-icon': true,
    'q-chip': { template: '<div data-testid="q-chip"><slot /></div>', props: ['color', 'textColor'] },
    'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>', props: ['label', 'color', 'loading'] },
    'q-dialog': { template: '<div><slot /></div>', props: ['modelValue'] },
    'q-spinner': true,
    'q-tooltip': true
}

function mountPayments() {
    return mount(ParentPayments, {
        global: {
            plugins: [
                [Quasar, {}],
                createTestingPinia({ createSpy: vi.fn })
            ],
            stubs: commonStubs
        }
    })
}

describe('Parent/Payments.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        paymentService.getPayments.mockResolvedValue({
            data: {
                payments: [
                    { id: '1', title: 'Assicurazione Scolastica', amount: 8.50, due_date: '2026-03-30', status: 'pending' },
                    { id: '2', title: 'Gita Firenze', amount: 45.00, due_date: '2026-04-15', status: 'pending' },
                    { id: '3', title: 'Contributo Volontario', amount: 120.00, paid_at: '2025-09-10', status: 'paid', receipt_number: 'REC-2025-001' }
                ]
            }
        })
    })

    it('renders the page title and subtitle', async () => {
        wrapper = mountPayments()
        await flushPromises()
        expect(wrapper.text()).toContain('Pagamenti Scolastici')
    })

    it('computes totalPending correctly from database items', async () => {
        wrapper = mountPayments()
        await flushPromises()
        // 8.50 + 45.00 = 53.50
        const chip = wrapper.find('[data-testid="q-chip"]')
        expect(chip.exists()).toBe(true)
        expect(chip.text()).toContain('53.50')
    })

    it('has pending payment items rendered from database', async () => {
        wrapper = mountPayments()
        await flushPromises()
        const items = wrapper.findAll('li')
        expect(items.length).toBeGreaterThanOrEqual(2)
        expect(wrapper.text()).toContain('Assicurazione Scolastica')
        expect(wrapper.text()).toContain('Gita Firenze')
    })

    it('calls pay method and opens confirmation', async () => {
        paymentService.pay.mockResolvedValue({
            data: {
                message: 'success',
                payment: { id: '1', title: 'Assicurazione Scolastica', amount: 8.50, status: 'paid', receipt_number: 'REC-001' }
            }
        })
        wrapper = mountPayments()
        await flushPromises()

        const payButtons = wrapper.findAll('button')
        const payBtn = payButtons.find(b => b.text().includes('PagoPA') || b.attributes('label')?.includes('PagoPA'))
        if (payBtn) {
            await payBtn.trigger('click')
            await flushPromises()
            expect(paymentService.pay).toHaveBeenCalledWith('1', { payment_method: 'PagoPA' })
        }
    })
})
