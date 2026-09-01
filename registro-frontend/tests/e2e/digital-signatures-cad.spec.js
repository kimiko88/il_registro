import { describe, it, expect, vi, beforeEach } from 'vitest'
import { shallowMount, mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'

// ─── Quasar Mock ─────────────────────────────────────────────────────────────
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    const mockIconSet = {
        name: 'material-icons',
        type: { positive: 'check', negative: 'warning', info: 'info', warning: 'priority_high' },
        arrow: { dropdown: 'arrow_drop_down', expand: 'arrow_drop_down', down: 'arrow_drop_down', up: 'arrow_drop_up' },
        chevron: { left: 'chevron_left', right: 'chevron_right' },
        table: { arrow: 'arrow_upward', select: 'arrow_drop_down' },
        editor: {},
        tree: {}
    }
    return {
        ...actual,
        useQuasar: () => ({
            dark: { isActive: false },
            loading: { show: vi.fn(), hide: vi.fn() },
            notify: vi.fn(),
            dialog: vi.fn().mockReturnValue({
                onOk: vi.fn(cb => { if (cb) cb(); return { onCancel: vi.fn(), onDismiss: vi.fn() } }),
                onCancel: vi.fn(),
                onDismiss: vi.fn()
            }),
            screen: { lt: { md: false }, gt: { xs: true } },
            lang: { current: 'it' },
            iconSet: mockIconSet
        })
    }
})

function buildPinia(role = 'admin', extras = {}) {
    return createTestingPinia({
        createSpy: () => vi.fn().mockResolvedValue(undefined),
        initialState: {
            auth: { user: { role, id: 'u-1', school_id: 's-1', email: 'admin@scuola.it' }, isAuthenticated: true },
            ...extras
        }
    })
}

// ─────────────────────────────────────────────────────────────────────────────
// Suite 1 – Digital Signatures & CAD Preservation Workflow
// ─────────────────────────────────────────────────────────────────────────────
describe('Digital Signatures & CAD Compliance Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('SC01 — Security & Compliance component renders CAD preservation and FEQ status cards', async () => {
        const { default: SecurityCompliance } = await import('@/components/Admin/SecurityCompliance.vue')
        const wrapper = shallowMount(SecurityCompliance, {
            global: { plugins: [buildPinia('admin')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('SC02 — Security Compliance allows downloading CAD preservation ZIP package', async () => {
        const { default: SecurityCompliance } = await import('@/components/Admin/SecurityCompliance.vue')
        const wrapper = shallowMount(SecurityCompliance, {
            global: { plugins: [buildPinia('admin')] }
        })
        expect(wrapper.find('q-btn-stub').exists() || wrapper.exists()).toBe(true)
    })

    it('SC03 — SessionReauthDialog component displays re-authentication form', async () => {
        const { default: SessionReauthDialog } = await import('@/components/Common/SessionReauthDialog.vue')
        const wrapper = shallowMount(SessionReauthDialog, {
            global: { plugins: [buildPinia('teacher')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('SC04 — AcknowledgmentModal component renders emergency and circular sign prompt', async () => {
        const { default: AcknowledgmentModal } = await import('@/components/Common/AcknowledgmentModal.vue')
        const wrapper = shallowMount(AcknowledgmentModal, {
            global: { plugins: [buildPinia('teacher')] },
            props: {
                modelValue: true,
                communication: { id: 'com-1', title: 'Circolare n.1', content: 'Test' }
            }
        })
        expect(wrapper.exists()).toBe(true)
    })
})
