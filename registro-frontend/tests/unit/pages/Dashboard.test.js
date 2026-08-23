/**
 * @file Dashboard.test.js
 * Tests: T18-T27 per Dashboard.vue
 */
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import Dashboard from '@/pages/Dashboard.vue'

const mockRouter = { push: vi.fn() }
vi.mock('vue-router', () => ({ useRouter: () => mockRouter }))

function createWrapper(roleOverride, classes = []) {
  return mount(Dashboard, {
    global: {
      plugins: [
        createTestingPinia({
          initialState: {
            auth: { user: { first_name: 'Test', role: roleOverride }, userRole: roleOverride },
            classes: { classes, loading: false },
          },
          stubActions: true,
        })
      ],
      stubs: { QPage: { template: '<div><slot /></div>' } }
    }
  })
}

// T18
describe('T18 — actions for teacher role', () => {
  it('returns attendance, grades, lessons, agenda routes', () => {
    const wrapper = createWrapper('teacher')
    const vm = wrapper.vm
    const keys = vm.actions.map(a => a.key)
    expect(keys).toContain('attendance')
    expect(keys).toContain('grades')
    expect(keys).toContain('lessons')
    expect(keys).toContain('agenda')
    expect(keys).not.toContain('users')
  })
})

// T19
describe('T19 — actions for secretary role', () => {
  it('no /admin/* routes', () => {
    const wrapper = createWrapper('secretary')
    const routes = wrapper.vm.actions.map(a => a.route)
    expect(routes.every(r => !r.startsWith('/admin'))).toBe(true)
  })
})

// T20
describe('T20 — actions for student role', () => {
  it('returns student-specific routes', () => {
    const wrapper = createWrapper('student')
    const routes = wrapper.vm.actions.map(a => a.route)
    expect(routes.every(r => r.startsWith('/student'))).toBe(true)
  })
})

// T21
describe('T21 — actions for parent role', () => {
  it('returns parent-specific routes', () => {
    const wrapper = createWrapper('parent')
    const routes = wrapper.vm.actions.map(a => a.route)
    expect(routes.every(r => r.startsWith('/parent'))).toBe(true)
  })
})

// T22
describe('T22 — isDashboardAdmin for principal', () => {
  it('returns true', () => {
    const wrapper = createWrapper('principal')
    expect(wrapper.vm.isDashboardAdmin).toBe(true)
  })
})

// T23
describe('T23 — isDashboardAdmin for vice_principal', () => {
  it('returns true', () => {
    const wrapper = createWrapper('vice_principal')
    expect(wrapper.vm.isDashboardAdmin).toBe(true)
  })
})

// T24
describe('T24 — displaySchedule draft ids are unique', () => {
  it('no duplicate ids when multiple drafts exist for today', () => {
    const wrapper = createWrapper('teacher')
    const today = new Date().toISOString().substring(0, 10)
    // Inject 3 drafts without ids for today
    wrapper.vm.lessonDrafts = [
      { date: today, subject: 'Math', topic: 'Algebra', hour: 1 },
      { date: today, subject: 'Physics', topic: 'Newton', hour: 2 },
      { date: today, subject: 'History', topic: 'Rome', hour: 3 },
    ]
    const ids = wrapper.vm.displaySchedule.map(e => e.id)
    expect(new Set(ids).size).toBe(ids.length) // all unique
  })
})

// T25
describe('T25 — saveLessonDraft without topic', () => {
  it('shows warning and does not write to sessionStorage', () => {
    const wrapper = createWrapper('teacher')
    const before = sessionStorage.getItem('registro_lesson_drafts')
    wrapper.vm.draftForm.topic = ''
    wrapper.vm.saveLessonDraft()
    expect(sessionStorage.getItem('registro_lesson_drafts')).toBe(before)
  })
})

// T26
describe('T26 — loadStoredDrafts reads from sessionStorage', () => {
  it('populates lessonDrafts from sessionStorage', () => {
    const drafts = [{ id: 'draft_1', date: '2026-08-23', subject: 'Test', topic: 'Hello', hour: 1, homework: '' }]
    sessionStorage.setItem('registro_lesson_drafts', JSON.stringify(drafts))
    const wrapper = createWrapper('teacher')
    wrapper.vm.loadStoredDrafts()
    expect(wrapper.vm.lessonDrafts).toHaveLength(1)
    expect(wrapper.vm.lessonDrafts[0].id).toBe('draft_1')
  })
})

// T27
describe('T27 — navigatingAction set during handleActionClick', () => {
  it('sets navigatingAction to action key while navigating', async () => {
    const wrapper = createWrapper('teacher')
    const action = { key: 'grades', route: '/teacher/grades' }
    const pushPromise = Promise.resolve()
    mockRouter.push = vi.fn(() => pushPromise)
    const clickPromise = wrapper.vm.handleActionClick(action)
    expect(wrapper.vm.navigatingAction).toBe('grades')
    await clickPromise
    expect(wrapper.vm.navigatingAction).toBeNull()
  })
})

