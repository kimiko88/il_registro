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

    it('publishSchedule calls POST /timetable/generate/:jobId/publish', () => {
      timetableGenService.publishSchedule('job-123')
      expect(api.post).toHaveBeenCalledWith('/timetable/generate/job-123/publish')
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
})
