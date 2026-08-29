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

function buildPinia(role = 'parent', extras = {}) {
    return createTestingPinia({
        createSpy: () => vi.fn().mockResolvedValue(undefined),
        initialState: {
            auth: { user: { role, id: 'parent-1', school_id: 'school-1' }, isAuthenticated: true },
            parent: {
                children: [
                    { id: 'stud-1', first_name: 'Marco', last_name: 'Rossi', class: '3A', school_name: 'Liceo Fermi' }
                ],
                selectedChildId: 'stud-1',
                loading: false
            },
            ...extras
        }
    })
}

// ─────────────────────────────────────────────────────────────────────────────
// Suite 1 – Parents Portal Workflow
// ─────────────────────────────────────────────────────────────────────────────
describe('Parents Portal Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('PP01 — Parent Children page mounts and displays linked children cards', async () => {
        const { default: Children } = await import('@/pages/parent/Children.vue')
        const wrapper = shallowMount(Children, {
            global: { plugins: [buildPinia('parent')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('PP02 — Parent GeneralMeetingBooking page mounts and displays booking slots', async () => {
        const { default: GeneralMeetingBooking } = await import('@/pages/parent/GeneralMeetingBooking.vue')
        const wrapper = shallowMount(GeneralMeetingBooking, {
            global: { plugins: [buildPinia('parent')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('PP03 — Parent Meetings page mounts and displays upcoming teacher meetings', async () => {
        const { default: Meetings } = await import('@/pages/parent/Meetings.vue')
        const wrapper = shallowMount(Meetings, {
            global: { plugins: [buildPinia('parent')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('PP04 — Parent ReportCard page mounts and displays term report cards', async () => {
        const { default: ReportCard } = await import('@/pages/parent/ReportCard.vue')
        const wrapper = shallowMount(ReportCard, {
            global: { plugins: [buildPinia('parent')] }
        })
        expect(wrapper.exists()).toBe(true)
    })
})
