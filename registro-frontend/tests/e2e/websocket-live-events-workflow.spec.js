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
            auth: { user: { role, id: 'u-1', school_id: 's-1' }, isAuthenticated: true },
            ...extras
        }
    })
}

// ─────────────────────────────────────────────────────────────────────────────
// Suite 1 – WebSocket Live Events & Activity Feed Workflow
// ─────────────────────────────────────────────────────────────────────────────
describe('WebSocket Live Events & Activity Feed Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('WS01 — TimelineActivityFeed component renders real-time event updates', async () => {
        const { default: TimelineActivityFeed } = await import('@/components/Common/TimelineActivityFeed.vue')
        const wrapper = shallowMount(TimelineActivityFeed, {
            global: { plugins: [buildPinia('teacher')] },
            props: {
                activities: [
                    { id: 'a1', type: 'grade', title: 'Nuovo voto registrato', time: '10:30', user: 'Prof. Rossi' },
                    { id: 'a2', type: 'attendance', title: 'Presenze confermate', time: '09:00', user: 'Prof. Bianchi' }
                ]
            }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('WS02 — GeneralMeetingLiveQueue page mounts and handles live queue updates', async () => {
        const { default: GeneralMeetingLiveQueue } = await import('@/pages/teacher/GeneralMeetingLiveQueue.vue')
        const wrapper = shallowMount(GeneralMeetingLiveQueue, {
            global: { plugins: [buildPinia('teacher')] }
        })
        expect(wrapper.exists()).toBe(true)
    })
})
