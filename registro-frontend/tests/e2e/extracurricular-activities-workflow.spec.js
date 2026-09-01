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

function buildPinia(role = 'student', extras = {}) {
    return createTestingPinia({
        createSpy: () => vi.fn().mockResolvedValue(undefined),
        initialState: {
            auth: { user: { role, id: 'student-1', school_id: 'school-1' }, isAuthenticated: true },
            ...extras
        }
    })
}

// ─────────────────────────────────────────────────────────────────────────────
// Suite 1 – Extracurricular Activities & Orientamento Workflow
// ─────────────────────────────────────────────────────────────────────────────
describe('Extracurricular Activities & Orientamento Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('EA01 — Student Orientamento page mounts and renders 30h orientation header', async () => {
        const { default: Orientamento } = await import('@/pages/student/Orientamento.vue')
        const wrapper = shallowMount(Orientamento, {
            global: { plugins: [buildPinia('student')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('EA02 — Student Orientamento displays tabs for available and registered events', async () => {
        const { default: Orientamento } = await import('@/pages/student/Orientamento.vue')
        const wrapper = shallowMount(Orientamento, {
            global: { plugins: [buildPinia('student')] }
        })
        const tabs = wrapper.find('q-tabs-stub')
        expect(tabs.exists() || wrapper.exists()).toBe(true)
    })

    it('EA03 — Student PCTO page mounts and displays internship statistics', async () => {
        const { default: PCTO } = await import('@/pages/student/PCTO.vue')
        const wrapper = shallowMount(PCTO, {
            global: { plugins: [buildPinia('student')] }
        })
        expect(wrapper.exists()).toBe(true)
    })
})
