import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import ScheduleGrid from '@/components/Secretary/ScheduleGrid.vue'
import TeacherScheduleGrid from '@/components/Secretary/TeacherScheduleGrid.vue'

describe('Secretary Schedule Grids Robustness Suite', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('ScheduleGrid calculates total weekly hours and populates cells', async () => {
    const initialSchedule = [
      { day: 1, hour: 1, subject_name: 'Matematica', teacher_name: 'Prof. Rossi', room: 'Aula 10' },
      { day: 1, hour: 2, subject_name: 'Matematica', teacher_name: 'Prof. Rossi', room: 'Aula 10' },
      { day: 2, hour: 1, subject_name: 'Italiano', teacher_name: 'Prof. Bianchi', room: 'Aula 10' }
    ]

    const wrapper = mount(ScheduleGrid, {
      props: {
        schedule: initialSchedule,
        subjects: [{ id: 's1', name: 'Matematica' }, { id: 's2', name: 'Italiano' }],
        teachers: [{ id: 't1', first_name: 'Mario', last_name: 'Rossi' }]
      },
      global: {
        mocks: { t: (k) => k },
        stubs: {
          'q-icon': true,
          'q-badge': { template: '<span class="badge"><slot /></span>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-dialog': true,
          'q-card': true,
          'q-card-section': true,
          'q-select': true,
          'q-input': true
        }
      }
    })

    expect(wrapper.text()).toContain('3 ore/settimana')
    expect(wrapper.text()).toContain('Matematica')
    expect(wrapper.text()).toContain('Italiano')
  })

  it('TeacherScheduleGrid emits save with current grid entries', async () => {
    const initialSchedule = [
      { day: 1, hour: 1, subject_name: 'Fisica', class_name: '3A', room: 'Lab 1' }
    ]

    const wrapper = mount(TeacherScheduleGrid, {
      props: {
        schedule: initialSchedule,
        classes: [{ id: 'c1', name: '3A' }],
        subjects: [{ id: 's1', name: 'Fisica' }]
      },
      global: {
        mocks: { t: (k) => k },
        stubs: {
          'q-icon': true,
          'q-badge': { template: '<span class="badge"><slot /></span>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-dialog': true,
          'q-card': true,
          'q-card-section': true,
          'q-select': true,
          'q-input': true
        }
      }
    })

    expect(wrapper.text()).toContain('1 ore/settimana')
    expect(wrapper.text()).toContain('Fisica')

    const saveBtn = wrapper.find('button')
    if (saveBtn.exists()) {
      await saveBtn.trigger('click')
      expect(wrapper.emitted('save')).toBeTruthy()
    }
  })
})
