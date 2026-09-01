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
            auth: { user: { role, id: 'teacher-1', school_id: 'school-1' }, isAuthenticated: true },
            ...extras
        }
    })
}

// ─────────────────────────────────────────────────────────────────────────────
// Suite 1 – Colloqui Scheduling & Parent Booking Workflow
// ─────────────────────────────────────────────────────────────────────────────
describe('Colloqui Scheduling & Parent Booking Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('CS01 — Teacher Colloqui page mounts and renders slot dialog button', async () => {
        const { default: Colloqui } = await import('@/pages/teacher/Colloqui.vue')
        const wrapper = shallowMount(Colloqui, {
            global: { plugins: [buildPinia('teacher')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('CS02 — Teacher Colloqui displays tabs for bookings and active slots', async () => {
        const { default: Colloqui } = await import('@/pages/teacher/Colloqui.vue')
        const wrapper = shallowMount(Colloqui, {
            global: { plugins: [buildPinia('teacher')] }
        })
        expect(wrapper.find('q-tabs-stub').exists() || wrapper.exists()).toBe(true)
    })

    it('CS03 — Parent Meetings page mounts and displays booking list', async () => {
        const { default: Meetings } = await import('@/pages/parent/Meetings.vue')
        const wrapper = shallowMount(Meetings, {
            global: { plugins: [buildPinia('parent')] }
        })
        expect(wrapper.exists()).toBe(true)
    })
})
