import { describe, it, expect, vi, beforeEach } from 'vitest'
import visitorService from '@/services/visitorService'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn()
  }
}))

describe('visitorService & Visitor Registry Unit Suite', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  // 1. Visitors
  it('1. listTodayVisitors calls GET /visitors/today with optional date filter', async () => {
    const mockVisitors = [
      { id: 'v-1', name: 'Marco Gialli', badge_number: 'B-01', exit_time: null }
    ]
    api.get.mockResolvedValue({ data: mockVisitors })

    const res = await visitorService.listTodayVisitors('2026-09-24')
    expect(api.get).toHaveBeenCalledWith('/visitors/today', { params: { date: '2026-09-24' } })
    expect(res.data).toEqual(mockVisitors)
  })

  it('2. registerVisitor posts external visitor data', async () => {
    const payload = {
      name: 'Ing. Bianchi',
      document_id: 'CI-12345',
      purpose: 'supplier',
      host_name: 'DSGA',
      badge_number: 'B-05',
      notes: 'Manutenzione server'
    }
    api.post.mockResolvedValue({ data: { id: 'v-new', ...payload } })

    const res = await visitorService.registerVisitor(payload)
    expect(api.post).toHaveBeenCalledWith('/visitors', payload)
    expect(res.data.id).toBe('v-new')
  })

  it('3. recordVisitorExit sends PATCH /visitors/:id/exit', async () => {
    api.patch.mockResolvedValue({ data: { success: true } })

    const res = await visitorService.recordVisitorExit('v-new', 'Badge riconsegnato')
    expect(api.patch).toHaveBeenCalledWith('/visitors/v-new/exit', { notes: 'Badge riconsegnato' })
    expect(res.data.success).toBe(true)
  })

  // 2. Early Exits
  it('4. listTodayEarlyExits calls GET /visitors/early-exits', async () => {
    const mockExits = [
      { id: 'exit-1', student_name: 'Luigi Verdi', delegatee_name: 'Anna Verdi' }
    ]
    api.get.mockResolvedValue({ data: mockExits })

    const res = await visitorService.listTodayEarlyExits('2026-09-24')
    expect(api.get).toHaveBeenCalledWith('/visitors/early-exits', { params: { date: '2026-09-24' } })
    expect(res.data).toEqual(mockExits)
  })

  it('5. recordEarlyExit posts student early exit details', async () => {
    const payload = {
      student_id: 'stud-10',
      delegatee_name: 'Giuseppe Verdi',
      delegate_rel: 'Genitore',
      reason_code: 'visita_medica',
      notes: 'Visita oculistica'
    }
    api.post.mockResolvedValue({ data: { id: 'exit-created', ...payload } })

    const res = await visitorService.recordEarlyExit(payload)
    expect(api.post).toHaveBeenCalledWith('/visitors/early-exits', payload)
    expect(res.data.id).toBe('exit-created')
  })

  it('6. recordStudentReturn sends PATCH /visitors/early-exits/:id/return', async () => {
    api.patch.mockResolvedValue({ data: { success: true } })

    const res = await visitorService.recordStudentReturn('exit-created', 'Rientro regolare')
    expect(api.patch).toHaveBeenCalledWith('/visitors/early-exits/exit-created/return', { notes: 'Rientro regolare' })
    expect(res.data.success).toBe(true)
  })

  // 3. Maintenance Reports
  it('7. listMaintenanceReports calls GET /visitors/maintenance with status filter', async () => {
    const mockReports = [
      { id: 'm-1', location: 'Palestra', priority: 'urgente', status: 'aperto' }
    ]
    api.get.mockResolvedValue({ data: mockReports })

    const res = await visitorService.listMaintenanceReports('aperto')
    expect(api.get).toHaveBeenCalledWith('/visitors/maintenance', { params: { status: 'aperto' } })
    expect(res.data).toEqual(mockReports)
  })

  it('8. createMaintenanceReport posts structural/equipment ticket', async () => {
    const payload = {
      location: 'Aula Magna',
      category: 'informatica',
      description: 'Videoproiettore non si accende',
      priority: 'alta'
    }
    api.post.mockResolvedValue({ data: { id: 'm-new', ...payload } })

    const res = await visitorService.createMaintenanceReport(payload)
    expect(api.post).toHaveBeenCalledWith('/visitors/maintenance', payload)
    expect(res.data.id).toBe('m-new')
  })

  it('9. updateMaintenanceStatus updates status and assigned technician', async () => {
    const patchData = { status: 'in_lavorazione', assigned_to: 'Ditta ElettroService' }
    api.patch.mockResolvedValue({ data: { success: true } })

    const res = await visitorService.updateMaintenanceStatus('m-new', patchData)
    expect(api.patch).toHaveBeenCalledWith('/visitors/maintenance/m-new/status', patchData)
    expect(res.data.success).toBe(true)
  })

  // 4. Overstay calculation helper
  it('10. identifies active visitors and detects overstay (> 3 hours without checkout)', () => {
    const now = new Date('2026-09-24T12:00:00Z').getTime()
    const visitors = [
      { id: 'v1', entry_time: '2026-09-24T11:00:00Z', exit_time: null }, // 1h ago -> normal
      { id: 'v2', entry_time: '2026-09-24T08:00:00Z', exit_time: null }, // 4h ago -> overstay
      { id: 'v3', entry_time: '2026-09-24T08:30:00Z', exit_time: '2026-09-24T09:30:00Z' } // checked out
    ]

    const isOverstay = (entryIso, exitIso) => {
      if (exitIso) return false
      const entryTime = new Date(entryIso).getTime()
      const diffHours = (now - entryTime) / (1000 * 60 * 60)
      return diffHours > 3
    }

    expect(isOverstay(visitors[0].entry_time, visitors[0].exit_time)).toBe(false)
    expect(isOverstay(visitors[1].entry_time, visitors[1].exit_time)).toBe(true)
    expect(isOverstay(visitors[2].entry_time, visitors[2].exit_time)).toBe(false)
  })
})
