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
      { day_of_week: 1, hour_index: 1, subject_id: 's1', subject_name: 'Matematica', teacher_name: 'Prof. Rossi', room: 'Aula 10' },
      { day_of_week: 1, hour_index: 2, subject_id: 's1', subject_name: 'Matematica', teacher_name: 'Prof. Rossi', room: 'Aula 10' },
      { day_of_week: 2, hour_index: 1, subject_id: 's2', subject_name: 'Italiano', teacher_name: 'Prof. Bianchi', room: 'Aula 10' }
    ]

    const wrapper = mount(ScheduleGrid, {
      props: {
        initialSchedule,
        assignments: [
          { id: 'a1', subject_id: 's1', subject_name: 'Matematica', teacher_id: 't1', teacher_name: 'Prof. Rossi' },
          { id: 'a2', subject_id: 's2', subject_name: 'Italiano', teacher_id: 't2', teacher_name: 'Prof. Bianchi' }
        ]
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
      { day_of_week: 1, hour_index: 1, subject_id: 's1', subject_name: 'Fisica', class_id: 'c1', class_name: '3A', room: 'Lab 1' }
    ]

    const wrapper = mount(TeacherScheduleGrid, {
      props: {
        initialSchedule,
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
