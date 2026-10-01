import { describe, it, expect, vi, beforeEach } from 'vitest'
import timetableGenService from '@/services/timetableGenService'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn()
  }
}))

describe('timetableGenService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Timetable Generation Jobs', () => {
    it('startGeneration calls POST /timetable/generate', () => {
      const options = { school_id: 'sch-1', prioritize_seniority: true }
      timetableGenService.startGeneration(options)
      expect(api.post).toHaveBeenCalledWith('/timetable/generate', options)
    })

    it('getJobStatus calls GET /timetable/generate/:jobId', () => {
      timetableGenService.getJobStatus('job-123')
      expect(api.get).toHaveBeenCalledWith('/timetable/generate/job-123')
    })

    it('publishSchedule calls POST /timetable/generate/:jobId/publish without body when empty', () => {
      timetableGenService.publishSchedule('job-123')
      expect(api.post).toHaveBeenCalledWith('/timetable/generate/job-123/publish')
    })

    it('publishSchedule with alternative_id calls POST /timetable/generate/:jobId/publish with body', () => {
      timetableGenService.publishSchedule('job-123', { alternative_id: 2 })
      expect(api.post).toHaveBeenCalledWith('/timetable/generate/job-123/publish', { alternative_id: 2 })
    })

    it('publishSchedule with custom slots calls POST /timetable/generate/:jobId/publish with body', () => {
      const customSlots = [{ class_id: 'c-1', day_of_week: 1, hour_index: 2 }]
      timetableGenService.publishSchedule('job-123', { slots: customSlots })
      expect(api.post).toHaveBeenCalledWith('/timetable/generate/job-123/publish', { slots: customSlots })
    })
  })

  describe('Teacher Preferences', () => {
    it('getPreferences calls GET /timetable/preferences', () => {
      timetableGenService.getPreferences({ teacher_id: 't-1' })
      expect(api.get).toHaveBeenCalledWith('/timetable/preferences', { params: { teacher_id: 't-1' } })
    })

    it('savePreferences calls POST /timetable/preferences', () => {
      const preferences = [{ day_of_week: 1, hour_slot: 1, preference_type: 'preferred' }]
      timetableGenService.savePreferences(preferences)
      expect(api.post).toHaveBeenCalledWith('/timetable/preferences', preferences)
    })
  })

  describe('Room Requirements & Constraints', () => {
    it('getRoomRequirements calls GET /timetable/room-requirements', () => {
      timetableGenService.getRoomRequirements()
      expect(api.get).toHaveBeenCalledWith('/timetable/room-requirements')
    })

    it('saveRoomRequirement calls POST /timetable/room-requirements', () => {
      const req = { subject_id: 'sub-1', required_room_type: 'computer_lab' }
      timetableGenService.saveRoomRequirement(req)
      expect(api.post).toHaveBeenCalledWith('/timetable/room-requirements', req)
    })

    it('deleteRoomRequirement calls DELETE /timetable/room-requirements/:id', () => {
      timetableGenService.deleteRoomRequirement('req-1')
      expect(api.delete).toHaveBeenCalledWith('/timetable/room-requirements/req-1')
    })

    it('getConstraints calls GET /timetable/constraints', () => {
      timetableGenService.getConstraints()
      expect(api.get).toHaveBeenCalledWith('/timetable/constraints')
    })

    it('saveConstraint calls POST /timetable/constraints', () => {
      const c = { constraint_type: 'max_consecutive_hours', value: 3 }
      timetableGenService.saveConstraint(c)
      expect(api.post).toHaveBeenCalledWith('/timetable/constraints', c)
    })

    it('deleteConstraint calls DELETE /timetable/constraints/:id', () => {
      timetableGenService.deleteConstraint('c-1')
      expect(api.delete).toHaveBeenCalledWith('/timetable/constraints/c-1')
    })
  })

  describe('Desiderata Window & Schedule Adjustment', () => {
    it('getDesiderataWindow calls GET /timetable/preferences/window', () => {
      timetableGenService.getDesiderataWindow()
      expect(api.get).toHaveBeenCalledWith('/timetable/preferences/window')
    })

    it('setDesiderataWindow calls POST /timetable/preferences/window', () => {
      timetableGenService.setDesiderataWindow(true)
      expect(api.post).toHaveBeenCalledWith('/timetable/preferences/window', { is_open: true })
    })

    it('savePreferences with params calls POST /timetable/preferences with params', () => {
      const prefs = [{ day_of_week: 2, hour_slot: 1, preference_type: 'unavailable' }]
      timetableGenService.savePreferences(prefs, { academic_year_id: '2024/2025' })
      expect(api.post).toHaveBeenCalledWith('/timetable/preferences', prefs, { params: { academic_year_id: '2024/2025' } })
    })

    it('adjustSchedule calls POST /timetable/generate/:jobId/adjust', () => {
      const adjustment = { slots: [{ day_of_week: 1, hour_index: 3 }] }
      timetableGenService.adjustSchedule('job-99', adjustment)
      expect(api.post).toHaveBeenCalledWith('/timetable/generate/job-99/adjust', adjustment)
    })
  })

  describe('Curriculum Plans & Academic Years', () => {
    it('getAcademicYears calls GET /timetable/academic-years', () => {
      timetableGenService.getAcademicYears()
      expect(api.get).toHaveBeenCalledWith('/timetable/academic-years')
    })

    it('getClassesCurriculumPlans calls GET /timetable/classes-plans with params', () => {
      timetableGenService.getClassesCurriculumPlans({ academic_year: '2024/2025' })
      expect(api.get).toHaveBeenCalledWith('/timetable/classes-plans', { params: { academic_year: '2024/2025' } })
    })

    it('getClassCurriculumPlan calls GET /timetable/classes/:classId/plan', () => {
      timetableGenService.getClassCurriculumPlan('class-1')
      expect(api.get).toHaveBeenCalledWith('/timetable/classes/class-1/plan')
    })

    it('saveClassCurriculumPlan calls PUT /timetable/classes/:classId/plan', () => {
      const plan = { min_hours_per_day: 4, max_hours_per_day: 6, subjects: [] }
      timetableGenService.saveClassCurriculumPlan('class-1', plan)
      expect(api.put).toHaveBeenCalledWith('/timetable/classes/class-1/plan', plan)
    })

    it('inheritClassCurriculumPlan calls POST /timetable/classes/:classId/inherit', () => {
      const data = { source_academic_year: '2023/2024' }
      timetableGenService.inheritClassCurriculumPlan('class-1', data)
      expect(api.post).toHaveBeenCalledWith('/timetable/classes/class-1/inherit', data)
    })

    it('inheritAllClassesCurriculumPlans calls POST /timetable/inherit-all-plans', () => {
      const data = { source_academic_year: '2023/2024' }
      timetableGenService.inheritAllClassesCurriculumPlans(data)
      expect(api.post).toHaveBeenCalledWith('/timetable/inherit-all-plans', data)
    })
  })

  describe('Teacher Quick Preferences (Tabular Representation)', () => {
    it('getTeachersQuickPreferences calls GET /timetable/teachers-quick-preferences with params', () => {
      timetableGenService.getTeachersQuickPreferences({ academic_year_id: '2024/2025' })
      expect(api.get).toHaveBeenCalledWith('/timetable/teachers-quick-preferences', { params: { academic_year_id: '2024/2025' } })
    })

    it('saveTeachersQuickPreferences calls POST /timetable/teachers-quick-preferences', () => {
      const data = { preferences: [{ teacher_id: 't-1', day_off: 1, time_slot_pref: 'early_hours' }] }
      timetableGenService.saveTeachersQuickPreferences(data)
      expect(api.post).toHaveBeenCalledWith('/timetable/teachers-quick-preferences', data)
    })

    it('saveTeacherQuickPreference calls PUT /timetable/teachers-quick-preferences/:teacherId with params', () => {
      const item = { teacher_id: 't-1', day_off: 3, time_slot_pref: 'late_hours' }
      timetableGenService.saveTeacherQuickPreference('t-1', item, { academic_year_id: '2024/2025' })
      expect(api.put).toHaveBeenCalledWith('/timetable/teachers-quick-preferences/t-1', item, { params: { academic_year_id: '2024/2025' } })
    })
  })
})

