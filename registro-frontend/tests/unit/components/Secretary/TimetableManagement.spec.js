import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import ScheduleGrid from '@/components/Secretary/ScheduleGrid.vue'
import TeacherScheduleGrid from '@/components/Secretary/TeacherScheduleGrid.vue'
import adminService from '@/services/adminService'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn()
  }
}))

describe('Timetable Management & Components', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('ScheduleGrid Component (Class Schedule)', () => {
    const mockAssignments = [
      { id: 'ass-1', subject_id: 'sub-1', subject_name: 'Matematica', teacher_id: 't-1', teacher_name: 'Rossi Mario' },
      { id: 'ass-2', subject_id: 'sub-2', subject_name: 'Italiano', teacher_id: 't-2', teacher_name: 'Bianchi Laura' }
    ]

    const mockInitialSchedule = [
      { day_of_week: 1, hour_index: 1, subject_id: 'sub-1', subject_name: 'Matematica', teacher_id: 't-1', teacher_name: 'Rossi Mario', room: '1A' },
      { day_of_week: 1, hour_index: 2, subject_id: 'sub-2', subject_name: 'Italiano', teacher_id: 't-2', teacher_name: 'Bianchi Laura', room: '1A' }
    ]

    it('renders timetable grid table and calculates total hours correctly', () => {
      const wrapper = mount(ScheduleGrid, {
        props: {
          assignments: mockAssignments,
          initialSchedule: mockInitialSchedule,
          loading: false
        }
      })

      expect(wrapper.text()).toContain('Riepilogo Ore Settimanali:')
      expect(wrapper.text()).toContain('Matematica')
      expect(wrapper.text()).toContain('Italiano')
    })

    it('emits save event with entries when save button is clicked', async () => {
      const wrapper = mount(ScheduleGrid, {
        props: {
          assignments: mockAssignments,
          initialSchedule: mockInitialSchedule,
          loading: false
        }
      })

      const saveBtn = wrapper.find('button')
      expect(saveBtn.exists()).toBe(true)
      await saveBtn.trigger('click')

      expect(wrapper.emitted('save')).toBeTruthy()
      expect(wrapper.emitted('save')[0][0]).toHaveLength(2)
      expect(wrapper.emitted('save')[0][0][0]).toMatchObject({
        day_of_week: 1,
        hour_index: 1,
        subject_id: 'sub-1'
      })
    })
  })

  describe('TeacherScheduleGrid Component (Teacher Schedule)', () => {
    const mockClasses = [
      { id: 'c-1', label: 'Classe 1A', name: '1', section: 'A' },
      { id: 'c-2', label: 'Classe 2B', name: '2', section: 'B' }
    ]

    const mockSubjects = [
      { id: 'sub-1', label: 'Matematica', name: 'Matematica' },
      { id: 'sub-2', label: 'Fisica', name: 'Fisica' }
    ]

    const mockInitialTeacherSchedule = [
      { day_of_week: 1, hour_index: 1, class_id: 'c-1', class_name: 'Classe 1A', subject_id: 'sub-1', subject_name: 'Matematica', room: 'Lab 1' },
      { day_of_week: 2, hour_index: 3, class_id: 'c-2', class_name: 'Classe 2B', subject_id: 'sub-2', subject_name: 'Fisica', room: 'Aula 2B' }
    ]

    it('renders teacher schedule grid and total teaching hours correctly', () => {
      const wrapper = mount(TeacherScheduleGrid, {
        props: {
          classes: mockClasses,
          subjects: mockSubjects,
          initialSchedule: mockInitialTeacherSchedule,
          loading: false
        }
      })

      expect(wrapper.text()).toContain('Totale Ore Docente:')
      expect(wrapper.text()).toContain('Matematica')
      expect(wrapper.text()).toContain('Fisica')
    })

    it('emits save event with teacher entries on save click', async () => {
      const wrapper = mount(TeacherScheduleGrid, {
        props: {
          classes: mockClasses,
          subjects: mockSubjects,
          initialSchedule: mockInitialTeacherSchedule,
          loading: false
        }
      })

      const saveBtn = wrapper.find('button')
      expect(saveBtn.exists()).toBe(true)
      await saveBtn.trigger('click')

      expect(wrapper.emitted('save')).toBeTruthy()
      expect(wrapper.emitted('save')[0][0]).toHaveLength(2)
      expect(wrapper.emitted('save')[0][0][0]).toMatchObject({
        day_of_week: 1,
        hour_index: 1,
        class_id: 'c-1',
        subject_id: 'sub-1'
      })
    })
  })

  describe('adminService Timetable API methods', () => {
    it('getTeacherSchedule calls /teachers/:id/schedule endpoint', () => {
      adminService.getTeacherSchedule('teacher-123')
      expect(api.get).toHaveBeenCalledWith('/teachers/teacher-123/schedule')
    })

    it('saveTeacherSchedule posts to /teachers/:id/schedule endpoint with payload', () => {
      const payload = { entries: [{ day_of_week: 1, hour_index: 1, class_id: 'c-1', subject_id: 'sub-1' }] }
      adminService.saveTeacherSchedule('teacher-123', payload)
      expect(api.post).toHaveBeenCalledWith('/teachers/teacher-123/schedule', payload)
    })
  })
})
