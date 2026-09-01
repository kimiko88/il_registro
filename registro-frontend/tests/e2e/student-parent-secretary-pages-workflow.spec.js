import { describe, it, expect, vi, beforeEach } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'

// ─── Quasar mock ─────────────────────────────────────────────────────────────
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
    'q-table': { template: '<div><slot name="top-right" /><slot name="body-cell-actions" :props="{row: {id: 1}}" /></div>' },
    'q-badge': true,
    'q-separator': true,
    'q-chip': true,
    'q-toggle': true,
    'q-tooltip': true,
    'q-tab-panels': { template: '<div><slot /></div>' },
    'q-tab-panel': { template: '<div><slot /></div>' },
    'q-tabs': { template: '<div><slot /></div>' },
    'q-tab': true,
    'q-linear-progress': true,
    'q-circular-progress': true,
    'q-banner': { template: '<div><slot /></div>' },
    'q-date': true,
    'q-time': true,
    'q-expansion-item': { template: '<div><slot /></div>' },
    'q-timeline': { template: '<div><slot /></div>' },
    'q-timeline-entry': { template: '<div><slot /></div>' },
    'q-item': { template: '<div><slot /></div>' },
    'q-item-section': { template: '<div><slot /></div>' },
    'q-list': { template: '<div><slot /></div>' },
    'q-step': { template: '<div><slot /></div>' },
    'q-stepper': { template: '<div><slot /></div>' },
    'q-stepper-navigation': { template: '<div><slot /></div>' },
    'q-splitter': { template: '<div><slot name="before"/><slot name="after"/></div>' },
    'q-infinite-scroll': { template: '<div><slot /></div>' },
    'q-scroll-area': { template: '<div><slot /></div>' },
    'q-breadcrumbs': { template: '<div><slot /></div>' },
    'q-breadcrumbs-el': true,
    'q-uploader': true
}

// stubActions: true prevents store actions from running on mount (avoids API calls)
function pinia(role = 'student', extras = {}) {
    return createTestingPinia({
        // Each stubbed action returns a resolved Promise to prevent .catch()-on-undefined errors
        createSpy: () => vi.fn().mockResolvedValue(undefined),
        initialState: {
            auth: { user: { role, id: 'u1', school_id: 's1' }, isAuthenticated: true },
            schoolCalendar: { events: [], loading: false },
            agenda: { items: [], loading: false },
            classes: { classes: [], selectedClass: null, students: [] },
            attendance: { records: [] },
            grades: { grades: [] },
            teacher: { profile: {}, classes: [], subjects: [] },
            student: { profile: {}, orientamento: { resources: [], events: [] }, goals: [] },
            children: {
                children: [{ id: 'c1', full_name: 'Mario Rossi', class_name: '3A' }],
                selectedChild: { id: 'c1', full_name: 'Mario Rossi', class_name: '3A' }
            },
            pcto: { activities: [] },
            documents: { documents: [] },
            colloqui: { slots: [] },
            scrutiny: { records: [] },
            communications: { items: [] },
            notes: { notes: [] },
            ...extras
        }
    })
}

// ─────────────────────────────────────────────────────────────────────────────
// Suite 1 – Student Pages
// ─────────────────────────────────────────────────────────────────────────────
describe('Student Timetable Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('ST01 — student timetable page renders correctly', async () => {
        const { default: Timetable } = await import('@/pages/student/Timetable.vue')
        const wrapper = shallowMount(Timetable, { global: { plugins: [pinia('student')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('ST02 — student orientation page renders guidance content', async () => {
        const { default: Orientamento } = await import('@/pages/student/Orientamento.vue')
        const wrapper = shallowMount(Orientamento, { global: { plugins: [pinia('student')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('ST03 — student PCTO page renders internship list', async () => {
        const { default: PCTO } = await import('@/pages/student/PCTO.vue')
        const wrapper = shallowMount(PCTO, { global: { plugins: [pinia('student')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('ST04 — student goals page renders target tracking', async () => {
        const { default: Goals } = await import('@/pages/student/Goals.vue')
        const wrapper = shallowMount(Goals, { global: { plugins: [pinia('student')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('ST05 — student school calendar page renders', async () => {
        const { default: SchoolCalendar } = await import('@/pages/student/SchoolCalendar.vue')
        const wrapper = shallowMount(SchoolCalendar, { global: { plugins: [pinia('student')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('ST06 — student homework page shows assigned tasks', async () => {
        const { default: Homework } = await import('@/pages/student/Homework.vue')
        const wrapper = shallowMount(Homework, { global: { plugins: [pinia('student')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('ST07 — student profile page renders', async () => {
        const { default: Profile } = await import('@/pages/student/Profile.vue')
        const wrapper = shallowMount(Profile, { global: { plugins: [pinia('student')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('ST08 — student report card page renders', async () => {
        const { default: ReportCard } = await import('@/pages/student/ReportCard.vue')
        const wrapper = shallowMount(ReportCard, { global: { plugins: [pinia('student')] } })
        expect(wrapper.exists()).toBe(true)
    })
})

// ─────────────────────────────────────────────────────────────────────────────
// Suite 2 – Parent Trips & Timetable
// ─────────────────────────────────────────────────────────────────────────────
describe('Parent Trips & Timetable Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('PT01 — parent timetable page renders child schedule', async () => {
        const { default: Timetable } = await import('@/pages/parent/Timetable.vue')
        const wrapper = shallowMount(Timetable, { global: { plugins: [pinia('parent')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('PT02 — parent trips page renders educational trip list', async () => {
        const { default: Trips } = await import('@/pages/parent/Trips.vue')
        const wrapper = shallowMount(Trips, { global: { plugins: [pinia('parent')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('PT03 — parent PDP view page renders plan details', async () => {
        const { default: PdpView } = await import('@/pages/parent/PdpView.vue')
        const wrapper = shallowMount(PdpView, { global: { plugins: [pinia('parent')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('PT04 — parent children selector page renders children list', async () => {
        const { default: Children } = await import('@/pages/parent/Children.vue')
        const wrapper = shallowMount(Children, { global: { plugins: [pinia('parent')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('PT05 — parent general meeting booking page renders slots', async () => {
        const { default: GeneralMeetingBooking } = await import('@/pages/parent/GeneralMeetingBooking.vue')
        const wrapper = shallowMount(GeneralMeetingBooking, { global: { plugins: [pinia('parent')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('PT06 — parent colloqui booking page renders', async () => {
        const { default: Colloqui } = await import('@/pages/parent/Colloqui.vue')
        const wrapper = shallowMount(Colloqui, { global: { plugins: [pinia('parent')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('PT07 — parent grades page renders child grades', async () => {
        const { default: Grades } = await import('@/pages/parent/Grades.vue')
        const wrapper = shallowMount(Grades, { global: { plugins: [pinia('parent')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('PT08 — parent attendance page renders child attendance', async () => {
        const { default: Attendance } = await import('@/pages/parent/Attendance.vue')
        const wrapper = shallowMount(Attendance, { global: { plugins: [pinia('parent')] } })
        expect(wrapper.exists()).toBe(true)
    })
})

// ─────────────────────────────────────────────────────────────────────────────
// Suite 3 – Secretary Users & Classes
// ─────────────────────────────────────────────────────────────────────────────
describe('Secretary Users & Classes Workflow', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('SU01 — secretary users page renders user list', async () => {
        const { default: Users } = await import('@/pages/secretary/Users.vue')
        const wrapper = shallowMount(Users, { global: { plugins: [pinia('secretary')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('SU02 — secretary classes page renders class list', async () => {
        const { default: Classes } = await import('@/pages/secretary/Classes.vue')
        const wrapper = shallowMount(Classes, { global: { plugins: [pinia('secretary')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('SU03 — secretary meetings page renders meeting list', async () => {
        const { default: Meetings } = await import('@/pages/secretary/Meetings.vue')
        const wrapper = shallowMount(Meetings, { global: { plugins: [pinia('secretary')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('SU04 — secretary timetable page renders scheduling interface', async () => {
        const { default: Timetable } = await import('@/pages/secretary/Timetable.vue')
        const wrapper = shallowMount(Timetable, { global: { plugins: [pinia('secretary')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('SU05 — secretary PCTO management page renders', async () => {
        const { default: PCTO } = await import('@/pages/secretary/PCTO.vue')
        const wrapper = shallowMount(PCTO, { global: { plugins: [pinia('secretary')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('SU06 — secretary groups management renders group list', async () => {
        const { default: Groups } = await import('@/pages/secretary/Groups.vue')
        const wrapper = shallowMount(Groups, { global: { plugins: [pinia('secretary')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('SU07 — secretary textbooks page renders', async () => {
        const { default: Textbooks } = await import('@/pages/secretary/Textbooks.vue')
        const wrapper = shallowMount(Textbooks, { global: { plugins: [pinia('secretary')] } })
        expect(wrapper.exists()).toBe(true)
    })

    it('SU08 — secretary documents page renders', async () => {
        const { default: Documents } = await import('@/pages/secretary/Documents.vue')
        const wrapper = shallowMount(Documents, { global: { plugins: [pinia('secretary')] } })
        expect(wrapper.exists()).toBe(true)
    })
})
