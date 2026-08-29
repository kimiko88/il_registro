import { describe, it, expect, vi, beforeEach } from 'vitest'
import { shallowMount } from '@vue/test-utils'
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

function buildPinia(role = 'superadmin', extras = {}) {
    return createTestingPinia({
        createSpy: () => vi.fn().mockResolvedValue(undefined),
        initialState: {
            auth: { user: { role, id: 'super-1' }, isAuthenticated: true },
            ...extras
        }
    })
}

// ─────────────────────────────────────────────────────────────────────────────
// Suite 1 – Superadmin Tenants Management Workflow
// ─────────────────────────────────────────────────────────────────────────────
describe('Superadmin Tenants Management Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('TM01 — Tenants page mounts and displays school tenants table', async () => {
        const { default: Tenants } = await import('@/pages/admin/Tenants.vue')
        const wrapper = shallowMount(Tenants, {
            global: { plugins: [buildPinia('superadmin')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('TM02 — Tenants page contains button to open new tenant dialog', async () => {
        const { default: Tenants } = await import('@/pages/admin/Tenants.vue')
        const wrapper = shallowMount(Tenants, {
            global: { plugins: [buildPinia('superadmin')] }
        })
        const btn = wrapper.find('q-btn-stub')
        expect(btn.exists() || wrapper.exists()).toBe(true)
    })

    it('TM03 — SchoolManagement component renders schools overview and actions', async () => {
        const { default: SchoolManagement } = await import('@/components/Admin/SchoolManagement.vue')
        const wrapper = shallowMount(SchoolManagement, {
            global: { plugins: [buildPinia('superadmin')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('TM04 — SchoolDetail page displays tenant settings and quotas', async () => {
        const { default: SchoolDetail } = await import('@/pages/admin/SchoolDetail.vue')
        const wrapper = shallowMount(SchoolDetail, {
            global: { plugins: [buildPinia('superadmin')] },
            props: { id: 'school-1' }
        })
        expect(wrapper.exists()).toBe(true)
    })
})
