
import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import Index from '@/pages/secretary/Index.vue'

// Mock Quasar
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            loading: { show: vi.fn(), hide: vi.fn() },
            notify: vi.fn(),
            config: {}
        })
    }
})

// Mock Router
const mockRouter = {
    push: vi.fn()
}
vi.mock('vue-router', () => ({
    useRouter: () => mockRouter
}))

describe('Secretary Dashboard (Index.vue)', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        wrapper = mount(Index, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-avatar': true,
                    'q-list': { template: '<ul><slot /></ul>' },
                    'q-item': { template: '<li><slot /></li>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<span><slot /></span>' },
                    'q-icon': true,
                    'q-chip': true,
                    'q-btn': true // Stub btn to avoid complexities, finding by attributes
                },
                directives: {
                    Ripple: { mounted() { }, updated() { }, unmounted() { } }
                }
            }
        })
    })

    it('renders dashboard title and stats cards', () => {
        expect(wrapper.text()).toContain('Dashboard Segreteria')
        expect(wrapper.text()).toContain('Documenti Pendenti')
        expect(wrapper.text()).toContain('Totale Studenti')
    })

    it('renders pending reviews list', () => {
        expect(wrapper.text()).toContain('Da Revisionare')
        expect(wrapper.text()).toContain('PDP - Giulia Verdi')
    })

    it('navigates to documents on review item click', async () => {
        const item = wrapper.findAll('li').find(w => w.text().includes('PDP - Giulia Verdi'))
        expect(item).toBeDefined()
        if (item) {
            await item.trigger('click')
            expect(mockRouter.push).toHaveBeenCalledWith(expect.stringContaining('/secretary/documents'))
        }
    })
})
