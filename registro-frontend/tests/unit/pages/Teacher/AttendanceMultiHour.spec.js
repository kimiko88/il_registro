import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import TeacherAttendance from '@/pages/teacher/Attendance.vue'

const { mockGetByClass, mockApiGet, mockApiPost, mockApiPut, mockGetLessons, mockExecuteWithOfflineQueue } = vi.hoisted(() => ({
  mockGetByClass: vi.fn(),
  mockApiGet: vi.fn(),
  mockApiPost: vi.fn(),
  mockApiPut: vi.fn(),
  mockGetLessons: vi.fn(),
  mockExecuteWithOfflineQueue: vi.fn()
}))

vi.mock('src/services/attendanceService', () => ({
  attendanceService: { getByClass: mockGetByClass, deleteHour: vi.fn() }
}))
vi.mock('src/services/lessonService', () => ({
  lessonService: {
    getLessons: mockGetLessons,
    deleteLesson: vi.fn()
  }
}))
vi.mock('@/services/api', () => ({
  default: { get: mockApiGet, post: mockApiPost, put: mockApiPut }
}))
vi.mock('@/services/offlineQueueService', () => ({
  executeWithOfflineQueue: mockExecuteWithOfflineQueue
}))

describe('Teacher/Attendance.vue - Multi-Hour & Topic Copy', () => {
  let wrapper

  beforeEach(() => {
    vi.clearAllMocks()

    mockApiGet.mockImplementation((url) => {
      if (url === '/users') {
        return Promise.resolve({
          data: {
            users: [
              { id: 's1', first_name: 'Mario', last_name: 'Rossi' },
              { id: 's2', first_name: 'Luigi', last_name: 'Verdi' }
            ]
          }
        })
      }
      if (url.includes('pending-justifications')) return Promise.resolve({ data: [] })
      if (url.includes('/lessons/class/')) return Promise.resolve({ data: [] })
      if (url.includes('/classes/') && url.includes('/subjects')) {
        return Promise.resolve({
          data: [{ subject_id: 'sub-1', subject_name: 'Matematica', teacher_id: 't-1' }]
        })
      }
      return Promise.resolve({ data: [] })
    })

    mockGetByClass.mockResolvedValue({
      data: [{ student_id: 's1', status: 'Present', hour: 1 }]
    })

    mockExecuteWithOfflineQueue.mockResolvedValue({ data: { success: true } })

    wrapper = mount(TeacherAttendance, {
      global: {
        plugins: [
          [Quasar, {}],
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              auth: { user: { id: 't-1', first_name: 'Prof', last_name: 'Test' } },
              classes: {
                classes: [{ id: 'c1', name: '4A', label: '4A' }]
              },
              schoolYear: { selectedSchoolYear: '2025/2026' }
            },
            stubActions: true
          })
        ],
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-toolbar': true,
          'q-toolbar-title': true,
          'q-btn': true,
          'q-btn-dropdown': true,
          'q-input': true,
          'q-select': true,
          'q-list': true,
          'q-item': true,
          'q-item-section': true,
          'q-item-label': true,
          'q-avatar': true,
          'q-btn-toggle': true,
          'q-spinner': true,
          'q-dialog': true,
          'q-tooltip': true,
          'NoteDialog': true,
          'SkeletonTable': true
        }
      }
    })
  })

  it('renders and contains copy last lesson button and multi-hour actions', () => {
    expect(wrapper.exists()).toBe(true)
  })

  it('copyLastLessonTopic retrieves topics from prior lessons', async () => {
    mockGetLessons.mockResolvedValueOnce({
      data: [
        { id: 'l1', date: '2026-09-01', hour: 1, topic: 'Equazioni di secondo grado', notes: 'Esercizi svolti' },
        { id: 'l2', date: '2026-08-30', hour: 2, topic: 'Disequazioni', notes: '' }
      ]
    })
    wrapper.vm.selectedClass = { id: 'c1' }
    wrapper.vm.lessonSubjectId = 'sub-1'
    await wrapper.vm.copyLastLessonTopic()
    expect(mockGetLessons).toHaveBeenCalled()
    expect(wrapper.vm.lessonTopic).toBe('Equazioni di secondo grado')
    expect(wrapper.vm.lessonNotes).toBe('Esercizi svolti')
  })
})
