/**
 * @file TeacherGrades.test.js
 * Tests: T28-T36 per Grades.vue
 */
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import Grades from '@/pages/teacher/Grades.vue'

vi.mock('vue-router', () => ({ useRouter: () => ({ push: vi.fn() }), useRoute: () => ({}) }))

function createWrapper({ classes = [], selectedClassId = null } = {}) {
  const wrapper = mount(Grades, {
    global: {
      plugins: [
        createTestingPinia({
          initialState: {
            auth: { user: { id: 'teacher-1', role: 'teacher' }, userRole: 'teacher' },
            classes: { classes, loading: false },
            grades: { grades: [], loading: false },
          },
          stubActions: true,
        })
      ],
    }
  })
  if (selectedClassId !== null) {
    wrapper.vm.selectedClassId = selectedClassId
  }
  return wrapper
}

// T28
describe('T28 — isAssignedClass false when no class selected', () => {
  it('returns false', () => {
    const wrapper = createWrapper({ classes: [{ id: 'c1', name: 'Classe 2A', is_owner: true }] })
    expect(wrapper.vm.isAssignedClass).toBe(false)
  })
})

// T29
describe('T29 — isAssignedClass true when is_owner=true', () => {
  it('returns true', () => {
    const wrapper = createWrapper({
      classes: [{ id: 'c1', name: '2A', is_owner: true }],
      selectedClassId: 'c1'
    })
    expect(wrapper.vm.isAssignedClass).toBe(true)
  })
})

// T30
describe('T30 — isAssignedClass false when is_owner=false (supplenza)', () => {
  it('returns false (supplente)', () => {
    const wrapper = createWrapper({
      classes: [{ id: 'c1', name: '2A', is_owner: false }],
      selectedClassId: 'c1'
    })
    expect(wrapper.vm.isAssignedClass).toBe(false)
  })
})

// T31
describe('T31 — isAssignedClass true when is_owner=undefined (backward compat)', () => {
  it('returns true (field absent, treated as owned)', () => {
    const wrapper = createWrapper({
      classes: [{ id: 'c1', name: '2A' }], // no is_owner field
      selectedClassId: 'c1'
    })
    expect(wrapper.vm.isAssignedClass).toBe(true)
  })
})

// T32
describe('T32 — overlappingTestsCount counts same-day tests', () => {
  it('counts tests on the same day using ISO date split', () => {
    const wrapper = createWrapper({ classes: [{ id: 'c1', name: '2A', is_owner: true }], selectedClassId: 'c1' })
    const today = '2026-08-23'
    wrapper.vm.classTests = [
      { date: `${today}T00:00:00+02:00`, type: 'test' },
      { date: `${today}T08:00:00Z`, type: 'interrogation' },
      { date: '2026-08-24T00:00:00Z', type: 'test' },  // different day — should NOT count
    ]
    wrapper.vm.testForm = { date: today }
    expect(wrapper.vm.overlappingTestsCount).toBe(2)
  })
})

// T33
describe('T33 — overlappingTestsCount is 0 when no tests', () => {
  it('returns 0', () => {
    const wrapper = createWrapper({ classes: [{ id: 'c1', name: '2A', is_owner: true }], selectedClassId: 'c1' })
    wrapper.vm.classTests = []
    wrapper.vm.testForm = { date: '2026-08-23' }
    expect(wrapper.vm.overlappingTestsCount).toBe(0)
  })
})

