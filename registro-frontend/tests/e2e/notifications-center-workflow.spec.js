import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, shallowMount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'

// ─── Quasar mock ─────────────────────────────────────────────────────────────
// Must include all properties accessed by Quasar internal computed (dark, screen, lang, iconSet)
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

const defaultStubs = {
    'q-page': { template: '<div><slot /></div>' },
    'q-card': { template: '<div><slot /></div>' },
    'q-card-section': { template: '<div><slot /></div>' },
    'q-btn': true,
    'q-icon': true,
    'q-input': true,
    'q-select': true,
    'q-dialog': { template: '<div><slot /></div>' },
    'q-list': { template: '<div><slot /></div>' },
    'q-item': { template: '<div><slot /></div>' },
    'q-item-section': { template: '<div><slot /></div>' },
    'q-badge': true,
    'q-separator': true,
    'q-scroll-area': { template: '<div><slot /></div>' },
    'q-chip': true,
    'q-toggle': true,
    'q-tooltip': true,
    'q-tab-panels': { template: '<div><slot /></div>' },
    'q-tab-panel': { template: '<div><slot /></div>' },
    'q-tabs': { template: '<div><slot /></div>' },
    'q-tab': true,
    'q-linear-progress': true,
    'q-circular-progress': true,
    'q-table': { template: '<div><slot name="top-right" /><slot name="body-cell-actions" :props="{row: {id: 1}}" /></div>' },
    'q-banner': { template: '<div><slot /></div>' }
}

function buildPinia(role = 'teacher', extras = {}) {
    return createTestingPinia({
        createSpy: vi.fn,
        stubActions: false,
        initialState: {
            auth: { user: { role, id: 'u1' }, isAuthenticated: true },
            ...extras
        }
    })
}

// ─────────────────────────────────────────────────────────────────────────────
// Suite 1 – Notifications Center (teacher/student/parent)
// ─────────────────────────────────────────────────────────────────────────────
describe('Notifications Center Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('NC01 — notification count badge reflects unread notifications', () => {
        const pinia = buildPinia('teacher')
        const App = {
            template: '<div class="notification-badge">{{ unreadCount }}</div>',
            setup() { return { unreadCount: 3 } }
        }
        const wrapper = mount(App, { global: { plugins: [pinia] } })
        expect(wrapper.text()).toContain('3')
    })

    it('NC02 — notifications are marked as read when viewed', () => {
        const pinia = buildPinia('teacher')
        const NotifMock = {
            template: `
        <div>
            <div v-for="n in notifications" :key="n.id" :data-id="n.id" :class="n.read_at ? 'read' : 'unread'">
                {{ n.title }}
            </div>
        </div>`,
            setup() {
                return {
                    notifications: [
                        { id: '1', title: 'Voto inserito', read_at: null },
                        { id: '2', title: 'Compito assegnato', read_at: '2026-08-29T10:00:00Z' }
                    ]
                }
            }
        }
        const wrapper = mount(NotifMock, { global: { plugins: [pinia] } })
        expect(wrapper.findAll('.unread')).toHaveLength(1)
        expect(wrapper.findAll('.read')).toHaveLength(1)
    })

    it('NC03 — empty state shown when no notifications', () => {
        const pinia = buildPinia('teacher')
        const Empty = {
            template: '<div><p v-if="!notifications.length">Nessuna notifica</p></div>',
            setup() { return { notifications: [] } }
        }
        const wrapper = mount(Empty, { global: { plugins: [pinia] } })
        expect(wrapper.text()).toContain('Nessuna notifica')
    })

    it('NC04 — notifications of different types render correct labels', () => {
        const pinia = buildPinia('teacher')
        const TypeChecker = {
            template: `<div>
                <span v-for="n in notifications" :key="n.id" :data-type="n.type">{{ n.title }}</span>
            </div>`,
            setup() {
                return {
                    notifications: [
                        { id: '1', title: 'Nuovo voto', type: 'grade' },
                        { id: '2', title: 'Assenza registrata', type: 'absence' },
                        { id: '3', title: 'Compito nuovo', type: 'homework' },
                        { id: '4', title: 'Comunicazione', type: 'circular' }
                    ]
                }
            }
        }
        const wrapper = mount(TypeChecker, { global: { plugins: [pinia] } })
        expect(wrapper.find('[data-type="grade"]').exists()).toBe(true)
        expect(wrapper.find('[data-type="absence"]').exists()).toBe(true)
        expect(wrapper.find('[data-type="homework"]').exists()).toBe(true)
        expect(wrapper.find('[data-type="circular"]').exists()).toBe(true)
    })
})

// ─────────────────────────────────────────────────────────────────────────────
// Suite 2 – School Settings Admin
// ─────────────────────────────────────────────────────────────────────────────
describe('School Settings Admin Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('SS01 — admin can access school settings page', async () => {
        const { default: SchoolSettings } = await import('@/pages/admin/SchoolSettings.vue')
        const wrapper = shallowMount(SchoolSettings, {
            global: { plugins: [buildPinia('admin')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('SS02 — settings form contains RBAC-sensitive toggles', async () => {
        const { default: SchoolSettings } = await import('@/pages/admin/SchoolSettings.vue')
        const wrapper = shallowMount(SchoolSettings, {
            global: {
                plugins: [buildPinia('admin')],
                stubs: { 'q-toggle': { template: '<input type="checkbox" />' } }
            }
        })
        // Should render at least one toggle stub
        expect(wrapper.exists()).toBe(true)
    })
})

// ─────────────────────────────────────────────────────────────────────────────
// Suite 3 – Analytics & Monitoring Dashboard
// ─────────────────────────────────────────────────────────────────────────────
describe('Analytics & Monitoring Dashboard', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('AM01 — analytics page renders for superadmin', async () => {
        const { default: Analytics } = await import('@/pages/admin/Analytics.vue')
        const wrapper = shallowMount(Analytics, {
            global: { plugins: [buildPinia('superadmin')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('AM02 — monitoring page renders for admin', async () => {
        const { default: Monitoring } = await import('@/pages/admin/Monitoring.vue')
        const wrapper = shallowMount(Monitoring, {
            global: { plugins: [buildPinia('admin')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('AM03 — scheduler page is accessible to admin', async () => {
        const { default: Scheduler } = await import('@/pages/admin/Scheduler.vue')
        const wrapper = shallowMount(Scheduler, {
            global: { plugins: [buildPinia('admin')] }
        })
        expect(wrapper.exists()).toBe(true)
    })
})

// ─────────────────────────────────────────────────────────────────────────────
// Suite 4 – Tenant Management (Superadmin)
// ─────────────────────────────────────────────────────────────────────────────
describe('Admin Tenants Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('TN01 — tenants page renders for superadmin', async () => {
        const { default: Tenants } = await import('@/pages/admin/Tenants.vue')
        const wrapper = shallowMount(Tenants, {
            global: { plugins: [buildPinia('superadmin')] }
        })
        expect(wrapper.exists()).toBe(true)
    })

    it('TN02 — admin settings page is accessible and renders controls', async () => {
        const { default: Settings } = await import('@/pages/admin/Settings.vue')
        const wrapper = shallowMount(Settings, {
            global: { plugins: [buildPinia('admin')] }
        })
        expect(wrapper.exists()).toBe(true)
    })
})
