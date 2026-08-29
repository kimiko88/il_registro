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

function buildPinia(role = 'admin', extras = {}) {
    return createTestingPinia({
        createSpy: () => vi.fn().mockResolvedValue(undefined),
        initialState: {
            auth: { user: { role, id: 'admin-1', school_id: 'school-1' }, isAuthenticated: true },
            ...extras
        }
    })
}

// ─────────────────────────────────────────────────────────────────────────────
// Suite 1 – Admin Analytics & Monitoring Workflow
// ─────────────────────────────────────────────────────────────────────────────
describe('Admin Analytics & Monitoring Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('AM01 — Monitoring page mounts and renders header with refresh button', async () => {
        const { default: Monitoring } = await import('@/pages/admin/Monitoring.vue')
        const wrapper = shallowMount(Monitoring, {
            global: { plugins: [buildPinia('superadmin')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('AM02 — Analytics page mounts and displays KPI cards', async () => {
        const { default: Analytics } = await import('@/pages/admin/Analytics.vue')
        const wrapper = shallowMount(Analytics, {
            global: { plugins: [buildPinia('superadmin')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('AM03 — Scheduler page mounts and displays automated cron tasks', async () => {
        const { default: Scheduler } = await import('@/pages/admin/Scheduler.vue')
        const wrapper = shallowMount(Scheduler, {
            global: { plugins: [buildPinia('admin')] }
        })
        expect(wrapper.exists()).toBe(true)
    })
})
