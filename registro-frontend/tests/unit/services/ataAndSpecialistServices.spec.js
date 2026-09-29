import { describe, it, expect, vi, beforeEach } from 'vitest'
import api from '@/services/api'
import substitutionService from '@/services/substitutionService'
import personnelDeskService from '@/services/personnelDeskService'
import pdpService from '@/services/pdpService'
import staffAttendanceService from '@/services/staffAttendanceService'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn()
  }
}))

describe('Specialist and ATA Services Unit Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('substitutionService', () => {
    it('listBySchool calls GET /substitutions with params', () => {
      const params = { date: '2026-10-15' }
      substitutionService.listBySchool(params)
      expect(api.get).toHaveBeenCalledWith('/substitutions', { params })
    })

    it('listMy calls GET /substitutions/my', () => {
      substitutionService.listMy()
      expect(api.get).toHaveBeenCalledWith('/substitutions/my')
    })

    it('create calls POST /substitutions', () => {
      const data = { class_id: 'c1', absent_teacher_id: 't1' }
      substitutionService.create(data)
      expect(api.post).toHaveBeenCalledWith('/substitutions', data)
    })

    it('assign calls PUT /substitutions/:id/assign', () => {
      const data = { substitute_teacher_id: 't2', notes: 'Emergenza' }
      substitutionService.assign('sub-123', data)
      expect(api.put).toHaveBeenCalledWith('/substitutions/sub-123/assign', data)
    })

    it('recommendSubstitutes calls GET /substitutions/recommend-substitutes', () => {
      const params = { class_id: 'c1', subject_id: 's1', date: '2026-10-15', hour: 2 }
      substitutionService.recommendSubstitutes(params)
      expect(api.get).toHaveBeenCalledWith('/substitutions/recommend-substitutes', { params })
    })

    it('signRegister calls POST /substitutions/:id/sign-register', () => {
      substitutionService.signRegister('sub-123', 'Lezione completata')
      expect(api.post).toHaveBeenCalledWith('/substitutions/sub-123/sign-register', { notes: 'Lezione completata' })
    })

    it('todaySummary calls GET /substitutions/today-summary', () => {
      substitutionService.todaySummary('2026-10-15')
      expect(api.get).toHaveBeenCalledWith('/substitutions/today-summary', { params: { date: '2026-10-15' } })
    })
  })

  describe('personnelDeskService', () => {
    it('listRequests calls GET /personnel-desk/requests with status param', () => {
      personnelDeskService.listRequests('pending')
      expect(api.get).toHaveBeenCalledWith('/personnel-desk/requests', { params: { status: 'pending' } })
    })

    it('createRequest calls POST /personnel-desk/requests', () => {
      const payload = { type: 'ferie', days: 3 }
      personnelDeskService.createRequest(payload)
      expect(api.post).toHaveBeenCalledWith('/personnel-desk/requests', payload)
    })

    it('submitRequest calls PATCH /personnel-desk/requests/:id/submit', () => {
      personnelDeskService.submitRequest('req-1')
      expect(api.patch).toHaveBeenCalledWith('/personnel-desk/requests/req-1/submit')
    })

    it('aaReview calls PATCH /personnel-desk/requests/:id/aa-review', () => {
      const reviewData = { outcome: 'regular', notes: 'Documenti conformi' }
      personnelDeskService.aaReview('req-1', reviewData)
      expect(api.patch).toHaveBeenCalledWith('/personnel-desk/requests/req-1/aa-review', reviewData)
    })

    it('dsgaSign calls PATCH /personnel-desk/requests/:id/dsga-sign', () => {
      const vistoData = { outcome: 'favorable', protocol_number: 'PROT-101' }
      personnelDeskService.dsgaSign('req-1', vistoData)
      expect(api.patch).toHaveBeenCalledWith('/personnel-desk/requests/req-1/dsga-sign', vistoData)
    })

    it('dsApprove calls PATCH /personnel-desk/requests/:id/ds-approve', () => {
      const decreeData = { outcome: 'approved', decree_number: 'DEC-2026-44' }
      personnelDeskService.dsApprove('req-1', decreeData)
      expect(api.patch).toHaveBeenCalledWith('/personnel-desk/requests/req-1/ds-approve', decreeData)
    })
  })

  describe('pdpService', () => {
    it('getStudentPlans calls GET /pdp/student/:id', () => {
      pdpService.getStudentPlans('stu-1', '2026/2027')
      expect(api.get).toHaveBeenCalledWith('/pdp/student/stu-1', { params: { year: '2026/2027' } })
    })

    it('createPlan calls POST /pdp', () => {
      const data = { student_id: 'stu-1', diagnosis: 'DSA F81.0' }
      pdpService.createPlan(data)
      expect(api.post).toHaveBeenCalledWith('/pdp', data)
    })

    it('shareWithFamily calls POST /pdp/:id/share', () => {
      pdpService.shareWithFamily('pdp-1', true)
      expect(api.post).toHaveBeenCalledWith('/pdp/pdp-1/share', { share: true })
    })

    it('approveByFamily calls POST /pdp/:id/approve', () => {
      pdpService.approveByFamily('pdp-1', 'Approvo misure dispensative')
      expect(api.post).toHaveBeenCalledWith('/pdp/pdp-1/approve', { comment: 'Approvo misure dispensative' })
    })

    it('getCompensativeMeasures calls GET /pdp/measures/compensative', () => {
      pdpService.getCompensativeMeasures()
      expect(api.get).toHaveBeenCalledWith('/pdp/measures/compensative')
    })

    it('exportPdpPdf calls GET /pdp/:id/pdf with blob responseType', () => {
      pdpService.exportPdpPdf('pdp-1')
      expect(api.get).toHaveBeenCalledWith('/pdp/pdp-1/pdf', expect.objectContaining({
        responseType: 'blob'
      }))
    })
  })

  describe('staffAttendanceService', () => {
    it('getDailySummary calls GET /staff-attendance/summary', async () => {
      api.get.mockResolvedValueOnce({ data: { total: 10, present: 8 } })
      const res = await staffAttendanceService.getDailySummary('2026-10-15')
      expect(api.get).toHaveBeenCalledWith('/staff-attendance/summary', { params: { date: '2026-10-15' } })
      expect(res).toEqual({ total: 10, present: 8 })
    })

    it('setStrikeMode calls POST /staff-attendance/strike-mode', async () => {
      api.post.mockResolvedValueOnce({ data: { success: true } })
      const res = await staffAttendanceService.setStrikeMode({ date: '2026-10-15', is_strike_day: true })
      expect(api.post).toHaveBeenCalledWith('/staff-attendance/strike-mode', { date: '2026-10-15', is_strike_day: true })
      expect(res).toEqual({ success: true })
    })

    it('recordAttendance calls POST /staff-attendance', async () => {
      api.post.mockResolvedValueOnce({ data: { id: 'att-1', status: 'present' } })
      const res = await staffAttendanceService.recordAttendance({ user_id: 'u1', status: 'present' })
      expect(api.post).toHaveBeenCalledWith('/staff-attendance', { user_id: 'u1', status: 'present' })
      expect(res).toEqual({ id: 'att-1', status: 'present' })
    })
  })
})
