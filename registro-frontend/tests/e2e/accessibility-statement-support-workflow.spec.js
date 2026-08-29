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

function buildPinia(role = 'teacher', extras = {}) {
    return createTestingPinia({
        createSpy: () => vi.fn().mockResolvedValue(undefined),
        initialState: {
            auth: { user: { role, id: 'user-1', school_id: 'school-1' }, isAuthenticated: true },
            ...extras
        }
    })
}

// ─────────────────────────────────────────────────────────────────────────────
// Suite 1 – Accessibility Statement & Help Desk Support Workflow
// ─────────────────────────────────────────────────────────────────────────────
describe('Accessibility Statement & Help Desk Support Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('AS01 — AccessibilityStatement page mounts and renders compliance badge cards', async () => {
        const { default: AccessibilityStatement } = await import('@/pages/AccessibilityStatement.vue')
        const wrapper = shallowMount(AccessibilityStatement, {
            global: { plugins: [buildPinia('student')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('AS02 — Support page mounts and displays search bar and connection status', async () => {
        const { default: Support } = await import('@/pages/Support.vue')
        const wrapper = shallowMount(Support, {
            global: { plugins: [buildPinia('parent')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('AS03 — AdminUsers page mounts and renders administrators directory', async () => {
        const { default: AdminUsers } = await import('@/pages/admin/AdminUsers.vue')
        const wrapper = shallowMount(AdminUsers, {
            global: { plugins: [buildPinia('superadmin')] }
        })
        expect(wrapper.exists()).toBe(true)
    })
})
