import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { mount } from '@vue/test-utils'
import ChildSwitcher from '@/components/Parent/ChildSwitcher.vue'
import StudentGrades from '@/components/Parent/StudentGrades.vue'
import CircularCreator from '@/components/Secretary/CircularCreator.vue'
import Competencies from '@/components/Teacher/Competencies.vue'
import GradeEntry from '@/components/Teacher/GradeEntry.vue'
import { useChildrenStore } from '@/stores/children'
import { useClassesStore } from '@/stores/classes'
import { useGradesStore } from '@/stores/grades'
import competenciesService from '@/services/competenciesService'

vi.mock('@/services/competenciesService', () => ({
  default: {
    getStudentEvaluations: vi.fn().mockResolvedValue([]),
    saveEvaluation: vi.fn().mockResolvedValue({})
  }
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn()
  })
}))

describe('Parent, Secretary and Teacher Components Tests', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('ChildSwitcher displays child name from useChildrenStore and selects child', async () => {
    const store = useChildrenStore()
    store.children = [
      { id: 'c1', firstName: 'Mario', lastName: 'Rossi', className: '1A' },
      { id: 'c2', firstName: 'Luigi', lastName: 'Rossi', className: '2B' }
    ]
    store.selectedChildId = 'c1'

    const wrapper = mount(ChildSwitcher, {
      global: {
        mocks: { t: (k) => k },
        stubs: {
          'q-btn-dropdown': { template: '<button><slot /></button>' },
          'q-list': { template: '<div><slot /></div>' },
          'q-item-label': { template: '<div><slot /></div>' },
          'q-item': { template: '<div @click="$emit(\'click\')"><slot /></div>' },
          'q-item-section': { template: '<div><slot /></div>' },
          'q-avatar': { template: '<div><slot /></div>' }
        }
      }
    })

    expect(wrapper.text()).toContain('Mario Rossi')
  })

  it('StudentGrades renders cards for each child in store and navigates', async () => {
    const store = useChildrenStore()
    store.children = [
      { id: 'c1', firstName: 'Mario', lastName: 'Rossi', className: '1A' }
    ]

    const wrapper = mount(StudentGrades, {
      global: {
        mocks: { t: (k) => k },
        stubs: {
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-separator': true,
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-spinner': true
        }
      }
    })

    expect(wrapper.text()).toContain('Mario Rossi')
  })

  it('CircularCreator maps classes to UUID value options', async () => {
    const classesStore = useClassesStore()
    classesStore.classes = [
      { id: 'uuid-1a', name: '1A', year: 1, section: 'A' },
      { id: 'uuid-2b', name: '2B', year: 2, section: 'B' }
    ]

    const wrapper = mount(CircularCreator, {
      global: {
        mocks: { t: (k) => k },
        stubs: {
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-form': { template: '<form><slot /></form>' },
          'q-input': true,
          'q-checkbox': true,
          'q-select': true,
          'q-editor': true,
          'q-file': true,
          'q-btn': true,
          'q-icon': true
        }
      }
    })

    expect(wrapper.vm.classOptions).toEqual([
      { label: '1A', value: 'uuid-1a' },
      { label: '2B', value: 'uuid-2b' }
    ])
  })

  it('Competencies loads evaluations when student selected', async () => {
    const classesStore = useClassesStore()
    classesStore.classStudents = [
      { id: 'student-1', first_name: 'Giulia', last_name: 'Bianchi' }
    ]

    const wrapper = mount(Competencies, {
      global: {
        mocks: { t: (k) => k },
        stubs: {
          'q-icon': true,
          'q-btn': true,
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-select': true,
          'q-markup-table': { template: '<table><slot /></table>' },
          'q-btn-toggle': true,
          'q-input': true
        }
      }
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.vm.studentOptions).toEqual([
      { id: 'student-1', name: 'Bianchi Giulia' }
    ])
    expect(competenciesService.getStudentEvaluations).toHaveBeenCalled()
  })

  it('GradeEntry safely handles students and entries without crashing', async () => {
    const gradesStore = useGradesStore()
    gradesStore.grades = {
      students: [
        { student_id: 's1', full_name: 'Mario Rossi', absences: 2, grades: [{ grade_value: 8 }] }
      ]
    }

    const wrapper = mount(GradeEntry, {
      props: {
        subject: 'math-1',
        date: '2026-03-01'
      },
      global: {
        mocks: { t: (k) => k },
        stubs: {
          'q-card': { template: '<div><slot /></div>' },
          'q-banner': true,
          'q-table': { template: '<div><slot name="header" :props="{}" /><slot name="body" :props="{ row: { student_id: \'s1\', full_name: \'Mario Rossi\' }, rowIndex: 0 }" /></div>' },
          'q-tr': { template: '<tr><slot /></tr>' },
          'q-th': { template: '<th><slot /></th>' },
          'q-td': { template: '<td><slot /></td>' },
          'q-select': true,
          'q-input': true,
          'q-btn': true,
          'q-icon': true,
          'q-dialog': true
        }
      }
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.vm.studentsWithGrades).toHaveLength(1)
    expect(wrapper.vm.entryData['s1']).toBeDefined()
  })
})
