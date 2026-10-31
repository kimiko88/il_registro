import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import Payments from '@/pages/parent/Payments.vue'
import paymentService from '@/services/paymentService'

const { mockNotify, mockDialog } = vi.hoisted(() => ({
    mockNotify: vi.fn(),
    mockDialog: vi.fn()
}))

vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            notify: mockNotify,
            dialog: mockDialog
        })
    }
})

vi.mock('@/services/paymentService', () => ({
    default: {
        getPayments: vi.fn(),
        pay: vi.fn(),
        downloadReceipt: vi.fn()
    }
}))

describe('Parent Payments and PagoPA Workflow E2E', () => {
    let pinia

    const getFreshMockPayments = () => [
        {
            id: 'pay-001',
            title: 'Assicurazione Scolastica e Contributo Volontario 2026/27',
            description: 'Copertura infortuni e RCT alunni',
            amount: 35.00,
            status: 'pending',
            due_date: '2026-11-30T00:00:00Z'
        },
        {
            id: 'pay-002',
            title: 'Visita Didattica Planetario',
            description: 'Quota pullman e biglietto ingresso',
            amount: 15.00,
            status: 'pending',
            due_date: '2026-12-15T00:00:00Z'
        },
        {
            id: 'pay-003',
            title: 'Tassa Iscrizione Anno Precedente',
            amount: 25.00,
            status: 'paid',
            receipt_number: 'REC-2025-00981',
            paid_at: '2025-09-10T14:30:00Z'
        }
    ]

    beforeEach(() => {
        vi.clearAllMocks()
        paymentService.getPayments.mockImplementation(() => Promise.resolve({ data: getFreshMockPayments() }))

        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'parent-1', role: 'parent', name: 'Mario Rossi' }
                },
                children: {
                    children: [{ id: 'child-1', name: 'Lorenzo Rossi' }],
                    selectedChildId: 'child-1'
                }
            }
        })
    })

    const mountComponent = () => {
        return mount(Payments, {
            global: {
                plugins: [pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-chip': { template: '<div class="q-chip"><slot /></div>' },
                    'q-spinner': { template: '<div></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
                    'q-tab': { template: '<div class="q-tab"><slot /></div>' },
                    'q-separator': { template: '<hr />' },
                    'q-tab-panels': { template: '<div class="q-tab-panels"><slot /></div>' },
                    'q-tab-panel': { template: '<div class="q-tab-panel"><slot /></div>' },
                    'q-list': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-avatar': { template: '<div><slot /></div>' },
                    'q-btn': {
                        props: ['label'],
                        template: '<button class="q-btn">{{ label }}<slot /></button>'
                    },
                    'q-icon': true,
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-tooltip': { template: '<span><slot /></span>' }
                }
            }
        })
    }

    it('renders parent payments dashboard with tabs and calculates total pending', async () => {
        const wrapper = mountComponent()
        await vi.waitFor(() => {
            expect(wrapper.vm.loading).toBe(false)
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.find('.q-tabs').exists()).toBe(true)
        // 35.00 + 15.00 = 50.00
        expect(wrapper.vm.totalPending).toBe('50.00')
        expect(wrapper.vm.pendingItems).toHaveLength(2)
        expect(wrapper.vm.historyItems).toHaveLength(1)
    })

    it('executes pay action with PagoPA method, updates status and shows receipt dialog', async () => {
        const wrapper = mountComponent()
        await vi.waitFor(() => {
            expect(wrapper.vm.loading).toBe(false)
            expect(wrapper.vm.pendingItems.length).toBe(2)
        })

        const itemToPay = wrapper.vm.pendingItems[0]
        const mockPaidResponse = {
            id: itemToPay.id,
            title: itemToPay.title,
            amount: itemToPay.amount,
            status: 'paid',
            receipt_number: 'REC-2026-00441',
            paid_at: new Date().toISOString()
        }
        paymentService.pay.mockResolvedValueOnce({
            data: { payment: mockPaidResponse }
        })

        await wrapper.vm.pay(itemToPay)

        expect(paymentService.pay).toHaveBeenCalledWith(itemToPay.id, { payment_method: 'PagoPA' })
        expect(mockNotify).toHaveBeenCalledWith(
            expect.objectContaining({ type: 'positive' })
        )
        expect(wrapper.vm.showReceiptDialog).toBe(true)
        expect(wrapper.vm.completedPayment.id).toBe(itemToPay.id)
        expect(wrapper.vm.completedPayment.status).toBe('paid')
        // Pending list now only has 1 item, total pending is now 15.00
        expect(wrapper.vm.pendingItems).toHaveLength(1)
        expect(wrapper.vm.totalPending).toBe('15.00')
    })

    it('handles payment failure with error notification and keeps item pending', async () => {
        const wrapper = mountComponent()
        await vi.waitFor(() => {
            expect(wrapper.vm.loading).toBe(false)
            expect(wrapper.vm.pendingItems.length).toBe(2)
        })

        paymentService.pay.mockRejectedValueOnce({
            response: { data: { error: 'Circuito PagoPA momentaneamente non disponibile' } }
        })

        const itemToPay = wrapper.vm.pendingItems[0]
        await wrapper.vm.pay(itemToPay)

        expect(mockNotify).toHaveBeenCalledWith(
            expect.objectContaining({
                type: 'negative',
                message: 'Circuito PagoPA momentaneamente non disponibile'
            })
        )
        expect(wrapper.vm.showReceiptDialog).toBe(false)
        expect(wrapper.vm.pendingItems).toHaveLength(2)
    })

    it('triggers receipt download notification for completed payments', async () => {
        const wrapper = mountComponent()
        await vi.waitFor(() => {
            expect(wrapper.vm.loading).toBe(false)
            expect(wrapper.vm.historyItems.length).toBe(1)
        })

        const paidItem = wrapper.vm.historyItems[0]
        wrapper.vm.downloadReceipt(paidItem)

        expect(mockNotify).toHaveBeenCalledWith(
            expect.objectContaining({
                type: 'positive',
                icon: 'download'
            })
        )
    })

    it('opens PagoPA user guide info dialog', () => {
        const wrapper = mountComponent()
        wrapper.vm.openPagoPAGuide()

        expect(mockDialog).toHaveBeenCalledWith(
            expect.objectContaining({
                ok: expect.any(Object)
            })
        )
    })
})
