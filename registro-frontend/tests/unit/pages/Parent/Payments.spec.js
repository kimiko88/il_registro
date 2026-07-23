import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import ParentPayments from '@/pages/parent/Payments.vue'

const commonStubs = {
    'q-page': { template: '<div><slot /></div>' },
    'q-card': { template: '<div><slot /></div>' },
    'q-card-section': { template: '<div><slot /></div>' },
    'q-banner': { template: '<div data-testid="q-banner"><slot /><slot name="avatar" /></div>' },
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
        wrapper = mountPayments()
    })

    it('renders the page title', () => {
        expect(wrapper.text()).toContain('Pagamenti Scolastici')
    })

    it('shows the simulation mode banner', () => {
        const banner = wrapper.find('[data-testid="q-banner"]')
        expect(banner.exists()).toBe(true)
        expect(banner.text()).toContain('Modalità Dimostrativa')
    })

    it('shows total pending chip', () => {
        const chip = wrapper.find('[data-testid="q-chip"]')
        expect(chip.exists()).toBe(true)
        expect(chip.text()).toContain('Totale da Pagare')
    })

    it('computes totalPending correctly from default items', () => {
        // Default items: 8.50 + 45.00 = 53.50
        const chip = wrapper.find('[data-testid="q-chip"]')
        expect(chip.text()).toContain('53.50')
    })

    it('has pending payment items rendered', () => {
        const items = wrapper.findAll('li')
        // At least 2 pending items + 1 empty state item = 3+
        expect(items.length).toBeGreaterThanOrEqual(2)
    })
})
