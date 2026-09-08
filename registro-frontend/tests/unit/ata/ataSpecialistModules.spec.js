import { describe, it, expect, vi, beforeEach } from 'vitest'
import routes from '@/router/routes'
import substitutionService from '@/services/substitutionService'
import visitorService from '@/services/visitorService'
import staffAttendanceService from '@/services/staffAttendanceService'
import personnelDeskService from '@/services/personnelDeskService'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn().mockResolvedValue({ data: {} }),
    post: vi.fn().mockResolvedValue({ data: {} }),
    put: vi.fn().mockResolvedValue({ data: {} }),
    patch: vi.fn().mockResolvedValue({ data: {} }),
    delete: vi.fn().mockResolvedValue({ data: {} }),
  }
}))

describe('ATA Specialist Modules — Services and Router Suite', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Route Configurations for 4 ATA Modules', () => {
    const mainChildren = routes[0].children

    it('defines emergency-substitutions route with correct roles', () => {
      const r = mainChildren.find(c => c.path === 'ata/emergency-substitutions')
      expect(r).toBeDefined()
      expect(r.meta.roles).toContain('collaboratore_ds')
      expect(r.meta.roles).toContain('dsga')
      expect(r.meta.roles).toContain('admin')
    })

    it('defines visitor-registry route with correct roles', () => {
      const r = mainChildren.find(c => c.path === 'ata/visitor-registry')
      expect(r).toBeDefined()
      expect(r.meta.roles).toContain('collaboratore_scolastico')
      expect(r.meta.roles).toContain('collaboratore_ds')
    })

    it('defines timecard route with correct roles', () => {
      const r = mainChildren.find(c => c.path === 'ata/timecard')
      expect(r).toBeDefined()
      expect(r.meta.roles).toContain('dsga')
      expect(r.meta.roles).toContain('assistente_amministrativo')
      expect(r.meta.roles).toContain('collaboratore_scolastico')
    })

    it('defines personnel-desk route accessible to staff, teachers and management', () => {
      const r = mainChildren.find(c => c.path === 'ata/personnel-desk')
      expect(r).toBeDefined()
      expect(r.meta.roles).toContain('teacher')
      expect(r.meta.roles).toContain('collaboratore_scolastico')
      expect(r.meta.roles).toContain('dsga')
      expect(r.meta.roles).toContain('principal')
    })
  })

  describe('Substitution Service (Emergency Module)', () => {
    it('calls today-summary endpoint with date param', async () => {
      await substitutionService.todaySummary('2026-09-08')
      expect(api.get).toHaveBeenCalledWith('/substitutions/today-summary', {
        params: { date: '2026-09-08' }
      })
    })
  })

  describe('Visitor Service (Visitor Registry Module)', () => {
    it('registers a visitor via POST /visitors', async () => {
      const payload = { name: 'Mario Rossi', purpose: 'parent' }
      await visitorService.registerVisitor(payload)
      expect(api.post).toHaveBeenCalledWith('/visitors', payload)
    })

    it('records visitor exit via PATCH /visitors/:id/exit', async () => {
      await visitorService.recordVisitorExit('v-123', 'uscito')
      expect(api.patch).toHaveBeenCalledWith('/visitors/v-123/exit', { notes: 'uscito' })
    })

    it('records early exit via POST /visitors/early-exits', async () => {
      const payload = { student_id: 's-1', delegatee_name: 'Genitore' }
      await visitorService.recordEarlyExit(payload)
      expect(api.post).toHaveBeenCalledWith('/visitors/early-exits', payload)
    })

    it('creates maintenance report via POST /visitors/maintenance', async () => {
      const payload = { location: 'Aula 1', description: 'Guasto' }
      await visitorService.createMaintenanceReport(payload)
      expect(api.post).toHaveBeenCalledWith('/visitors/maintenance', payload)
    })
  })

  describe('Staff Attendance & Timecard Service', () => {
    it('fetches timecard via GET /staff-attendance/timecard', async () => {
      await staffAttendanceService.getTimecard({ month: '2026-09' })
      expect(api.get).toHaveBeenCalledWith('/staff-attendance/timecard', {
        params: { month: '2026-09' }
      })
    })

    it('creates leave request via POST /staff-attendance/leaves', async () => {
      const payload = { type: 'ferie', start_date: '2026-09-10' }
      await staffAttendanceService.createLeave(payload)
      expect(api.post).toHaveBeenCalledWith('/staff-attendance/leaves', payload)
    })

    it('approves leave request via PATCH /staff-attendance/leaves/:id/approve', async () => {
      await staffAttendanceService.approveLeave('l-1', { notes: 'ok' })
      expect(api.patch).toHaveBeenCalledWith('/staff-attendance/leaves/l-1/approve', { notes: 'ok' })
    })
  })

  describe('Personnel Desk Service (Digital Workflow)', () => {
    it('creates desk request via POST /personnel-desk/requests', async () => {
      const payload = { category: 'ferie', description: 'Ferie' }
      await personnelDeskService.createRequest(payload)
      expect(api.post).toHaveBeenCalledWith('/personnel-desk/requests', payload)
    })

    it('submits request via PATCH /personnel-desk/requests/:id/submit', async () => {
      await personnelDeskService.submitRequest('req-1')
      expect(api.patch).toHaveBeenCalledWith('/personnel-desk/requests/req-1/submit')
    })

    it('reviews as AA via PATCH /personnel-desk/requests/:id/aa-review', async () => {
      const payload = { note: 'Controllato', approve: true }
      await personnelDeskService.aaReview('req-1', payload)
      expect(api.patch).toHaveBeenCalledWith('/personnel-desk/requests/req-1/aa-review', payload)
    })

    it('signs as DSGA via PATCH /personnel-desk/requests/:id/dsga-sign', async () => {
      const payload = { note: 'Visto contabile', approve: true }
      await personnelDeskService.dsgaSign('req-1', payload)
      expect(api.patch).toHaveBeenCalledWith('/personnel-desk/requests/req-1/dsga-sign', payload)
    })

    it('approves as DS via PATCH /personnel-desk/requests/:id/ds-approve', async () => {
      const payload = { decree_num: 'DEC-1', note: 'Approvato', approve: true }
      await personnelDeskService.dsApprove('req-1', payload)
      expect(api.patch).toHaveBeenCalledWith('/personnel-desk/requests/req-1/ds-approve', payload)
    })
  })
})
