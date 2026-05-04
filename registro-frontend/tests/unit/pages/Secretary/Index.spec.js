
import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { Quasar } from 'quasar'
import Index from '@/pages/secretary/Index.vue'

// Mock Router
const mockRouter = {
    push: vi.fn()
}
vi.mock('vue-router', () => ({
    useRouter: () => mockRouter
}))

// Mock Services
const { mockAdminService, mockDocumentService } = vi.hoisted(() => ({
    mockAdminService: {
        getDashboardStats: vi.fn().mockResolvedValue({
            data: {
                pending_documents_count: '5',
                announcements_count: '2',
                total_students: '100',
                total_teachers: '20'
            }
        })
    },
    mockDocumentService: {
        getInbox: vi.fn().mockResolvedValue({
            data: {
                items: [
                    { id: 1, title: 'PDP - Giulia Verdi', author: 'Mario Rossi', created_at: '2025-01-20T10:00:00Z' }
                ]
            }
        })
    }
}))

vi.mock('src/services/adminService', () => ({ default: mockAdminService }))
vi.mock('src/services/documentService', () => ({ default: mockDocumentService }))


describe('Secretary Dashboard (Index.vue)', () => {
    let wrapper

    beforeEach(async () => {
        vi.clearAllMocks()
        wrapper = mount(Index, {
            global: {
                plugins: [Quasar],
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
        // Wait for fetchData
        await new Promise(resolve => setTimeout(resolve, 0))
        await wrapper.vm.$nextTick()
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
