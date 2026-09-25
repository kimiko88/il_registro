import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import SchedulePreferences from '@/pages/teacher/SchedulePreferences.vue'
import TimetableConstraints from '@/pages/secretary/TimetableConstraints.vue'
import Timetable from '@/pages/secretary/Timetable.vue'
import timetableGenService from '@/services/timetableGenService'
import adminService from '@/services/adminService'
import api from '@/services/api'

vi.mock('@/services/timetableGenService', () => ({
  default: {
    startGeneration: vi.fn(),
    getJobStatus: vi.fn(),
    publishSchedule: vi.fn(),
    getPreferences: vi.fn(),
    savePreferences: vi.fn(),
    getRoomRequirements: vi.fn(),
    saveRoomRequirement: vi.fn(),
    deleteRoomRequirement: vi.fn(),
    getConstraints: vi.fn(),
    saveConstraint: vi.fn()
  }
}))

vi.mock('@/services/adminService', () => ({
  default: {
    getClasses: vi.fn(),
    getTeachers: vi.fn(),
    getTeachersList: vi.fn().mockResolvedValue({ data: [] }),
    getSubjects: vi.fn(),
    getClassSchedule: vi.fn(),
    getTeacherSchedule: vi.fn(),
    getClassSubjects: vi.fn()
  }
}))

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn()
  }
}))

vi.mock('quasar', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useQuasar: () => ({
      dark: { isActive: false },
      loading: { show: vi.fn(), hide: vi.fn() },
      notify: vi.fn(),
      lang: { current: 'it' },
      dialog: vi.fn().mockReturnValue({
        onOk: vi.fn(cb => { if (cb) cb(); return { onCancel: vi.fn(), onDismiss: vi.fn() } })
      })
    })
  }
})

describe('E2E Workflow: Teacher Desiderata, Lab Requirements & Automatic Timetable Generation', () => {
  let piniaTeacher
  let piniaSecretary
  let piniaVicePrincipal

  beforeEach(() => {
    vi.clearAllMocks()

    piniaTeacher = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: { user: { role: 'teacher', id: 'teacher-senior-1', school_id: 'school-1' }, isAuthenticated: true }
      }
    })

    piniaSecretary = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: { user: { role: 'secretary', id: 'sec-1', school_id: 'school-1' }, isAuthenticated: true }
      }
    })

    piniaVicePrincipal = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: { user: { role: 'vice_principal', id: 'vp-1', school_id: 'school-1' }, isAuthenticated: true }
      }
    })

    timetableGenService.getPreferences.mockResolvedValue({ data: [] })
    timetableGenService.getRoomRequirements.mockResolvedValue({ data: [] })
    timetableGenService.getConstraints.mockResolvedValue({ data: [] })

    adminService.getClasses.mockResolvedValue({ data: [{ id: 'class-1', name: '1A' }] })
    adminService.getTeachers.mockResolvedValue({ data: [{ id: 'teacher-senior-1', name: 'Prof Senior' }] })
    adminService.getSubjects.mockResolvedValue({ data: [{ id: 'sub-info', name: 'Informatica' }] })
    adminService.getClassSchedule.mockResolvedValue({ data: [] })
    adminService.getTeacherSchedule.mockResolvedValue({ data: [] })
    adminService.getClassSubjects.mockResolvedValue({ data: [] })
    api.get.mockImplementation((url) => {
      if (url.includes('/subjects')) return Promise.resolve({ data: [{ id: 'sub-info', name: 'Informatica' }] })
      return Promise.resolve({ data: [] })
    })
  })

  it('Step 1: Senior Teacher defines schedule desiderata in 5x8 grid and saves preferences', async () => {
    const wrapper = mount(SchedulePreferences, {
      global: {
        plugins: [piniaTeacher],
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-banner': { template: '<div><slot /></div>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-icon': true,
          'q-chip': { template: '<span><slot /></span>' },
          'q-spinner-dots': true
        }
      }
    })

    await flushPromises()
    expect(wrapper.text()).toContain('Desiderata Orario Scolastico')

    // Teacher marks Monday hour 1 and 2 as preferred
    wrapper.vm.grid[1][1] = 'preferred'
    wrapper.vm.grid[1][2] = 'preferred'

    timetableGenService.savePreferences.mockResolvedValue({ data: { success: true } })
    await wrapper.vm.savePreferences()
    await flushPromises()

    expect(timetableGenService.savePreferences).toHaveBeenCalledWith({
      preferences: expect.arrayContaining([
        expect.objectContaining({ day_of_week: 1, hour_index: 1, preference_type: 'preferred' }),
        expect.objectContaining({ day_of_week: 1, hour_index: 2, preference_type: 'preferred' })
      ])
    })
  })

  it('Step 2: Secretary maps subject to special laboratory room requirements', async () => {
    const wrapper = mount(TimetableConstraints, {
      global: {
        plugins: [piniaSecretary],
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-banner': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-icon': true,
          'q-chip': { template: '<span><slot /></span>' },
          'q-badge': { template: '<span><slot /></span>' },
          'q-select': { template: '<div><slot /></div>' },
          'q-input': { template: '<div><slot /></div>' },
          'q-toggle': { template: '<div><slot /></div>' },
          'q-table': { template: '<div><slot /></div>' },
          'q-tabs': { template: '<div><slot /></div>' },
          'q-tab': { template: '<div><slot /></div>' },
          'q-dialog': { template: '<div><slot /></div>' },
          'q-separator': true,
          'q-spinner-dots': true
        }
      }
    })

    await flushPromises()
    expect(wrapper.text()).toContain('Vincoli e Requisiti Orario Scolastico')

    timetableGenService.saveRoomRequirement.mockResolvedValue({
      data: { id: 'req-info', subject_id: 'sub-info', required_room_type: 'lab_informatica' }
    })

    wrapper.vm.openRoomReqModal()
    wrapper.vm.reqForm = {
      subject_id: 'sub-info',
      required_room_type: 'lab_informatica',
      is_mandatory: true
    }

    await wrapper.vm.saveReq()
    await flushPromises()

    expect(timetableGenService.saveRoomRequirement).toHaveBeenCalledWith({
      subject_id: 'sub-info',
      required_room_type: 'lab_informatica',
      is_mandatory: true
    })
  })

  it('Step 3: Vice Principal triggers automatic timetable generation and publishes result', async () => {
    vi.useFakeTimers()

    const wrapper = mount(Timetable, {
      global: {
        plugins: [piniaVicePrincipal],
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-tabs': { template: '<div><slot /></div>' },
          'q-tab': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-icon': true,
          'q-badge': { template: '<span><slot /></span>' },
          'q-select': { template: '<div><slot /></div>' },
          'q-input': { template: '<div><slot /></div>' },
          'q-table': { template: '<div><slot /></div>' },
          'q-dialog': { template: '<div><slot /></div>' },
          'q-linear-progress': { template: '<div><slot /></div>' },
          'ScheduleGrid': { template: '<div>ScheduleGrid</div>' },
          'TeacherScheduleGrid': { template: '<div>TeacherScheduleGrid</div>' }
        }
      }
    })

    await flushPromises()
    expect(wrapper.vm.canGenerate).toBe(true)

    // 1. Vice Principal opens dialog and triggers generation
    timetableGenService.startGeneration.mockResolvedValue({
      data: { job_id: 'job-xyz-123' }
    })

    wrapper.vm.openGenerateDialog()
    await wrapper.vm.runGeneration()
    await flushPromises()

    expect(timetableGenService.startGeneration).toHaveBeenCalledWith({ time_limit_seconds: 20 })
    expect(wrapper.vm.generationJobId).toBe('job-xyz-123')

    // 2. Mock job status polling completion
    timetableGenService.getJobStatus.mockResolvedValue({
      data: {
        id: 'job-xyz-123',
        status: 'completed',
        result_summary: {
          total_slots: 30,
          assigned_slots: 30,
          coverage_pct: 100.0,
          hard_conflicts: 0
        }
      }
    })

    // Advance timer to trigger polling
    vi.advanceTimersByTime(1600)
    await flushPromises()

    expect(timetableGenService.getJobStatus).toHaveBeenCalledWith('job-xyz-123')
    expect(wrapper.vm.generationResult).not.toBeNull()
    expect(wrapper.vm.generationResult.coverage_pct).toBe(100)

    // 3. Vice Principal publishes schedule
    timetableGenService.publishSchedule.mockResolvedValue({
      data: { message: 'orario pubblicato' }
    })

    await wrapper.vm.publishGeneratedSchedule()
    await flushPromises()

    expect(timetableGenService.publishSchedule).toHaveBeenCalledWith('job-xyz-123')
    vi.useRealTimers()
  })
})
