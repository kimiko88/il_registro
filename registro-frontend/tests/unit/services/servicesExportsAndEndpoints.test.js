/**
 * @file servicesExportsAndEndpoints.test.js
 * Verifies all 34 services have dual exports (default + named), proper query param encoding, and timeouts on large operations.
 */
import { describe, it, expect, vi, beforeEach } from 'vitest'
import api from '@/services/api'

import * as adminModule from '@/services/adminService'
import * as apiModule from '@/services/api'
import * as attendanceModule from '@/services/attendanceService'
import * as authModule from '@/services/authService'
import * as colloquiModule from '@/services/colloquiService'
import * as communicationModule from '@/services/communicationService'
import * as competenciesModule from '@/services/competenciesService'
import * as dashboardModule from '@/services/dashboardService'
import * as didacticModule from '@/services/didacticService'
import * as documentModule from '@/services/documentService'
import * as elearningModule from '@/services/elearningService'
import * as extracurricularModule from '@/services/extracurricularService'
import * as gradeModule from '@/services/gradeService'
import * as groupsModule from '@/services/groupsService'
import * as lessonModule from '@/services/lessonService'
import * as monitoringModule from '@/services/monitoringService'
import * as notesModule from '@/services/notesService'
import * as notificationModule from '@/services/notificationService'
import * as pctoModule from '@/services/pctoService'
import * as pdpModule from '@/services/pdpService'
import * as schoolModule from '@/services/schoolService'
import * as schoolSettingsModule from '@/services/schoolSettingsService'
import * as scrutinyModule from '@/services/scrutinyService'
import * as securityModule from '@/services/securityService'
import * as sidiModule from '@/services/sidiService'
import * as studentGoalModule from '@/services/studentGoalService'
import * as substitutionModule from '@/services/substitutionService'
import * as teacherActivityModule from '@/services/teacherActivityService'
import * as tenantsModule from '@/services/tenantsService'
import * as textbookModule from '@/services/textbookService'
import * as tripsModule from '@/services/tripsService'
import * as udaModule from '@/services/udaService'
import * as userModule from '@/services/userService'
import * as verbaliModule from '@/services/verbaliService'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn().mockResolvedValue({ data: {} }),
    post: vi.fn().mockResolvedValue({ data: {} }),
    put: vi.fn().mockResolvedValue({ data: {} }),
    patch: vi.fn().mockResolvedValue({ data: {} }),
    delete: vi.fn().mockResolvedValue({ data: {} }),
  }
}))

describe('Dual Exports on all 34 service modules', () => {
  const serviceModules = [
    { name: 'adminService', mod: adminModule },
    { name: 'attendanceService', mod: attendanceModule },
    { name: 'authService', mod: authModule },
    { name: 'colloquiService', mod: colloquiModule },
    { name: 'communicationService', mod: communicationModule },
    { name: 'competenciesService', mod: competenciesModule },
    { name: 'dashboardService', mod: dashboardModule },
    { name: 'didacticService', mod: didacticModule },
    { name: 'documentService', mod: documentModule },
    { name: 'elearningService', mod: elearningModule },
    { name: 'extracurricularService', mod: extracurricularModule },
    { name: 'gradeService', mod: gradeModule },
    { name: 'groupsService', mod: groupsModule },
    { name: 'lessonService', mod: lessonModule },
    { name: 'monitoringService', mod: monitoringModule },
    { name: 'notesService', mod: notesModule },
    { name: 'notificationService', mod: notificationModule },
    { name: 'pctoService', mod: pctoModule },
    { name: 'pdpService', mod: pdpModule },
    { name: 'schoolService', mod: schoolModule },
    { name: 'schoolSettingsService', mod: schoolSettingsModule },
    { name: 'scrutinyService', mod: scrutinyModule },
    { name: 'securityService', mod: securityModule },
    { name: 'sidiService', mod: sidiModule },
    { name: 'studentGoalService', mod: studentGoalModule },
    { name: 'substitutionService', mod: substitutionModule },
    { name: 'teacherActivityService', mod: teacherActivityModule },
    { name: 'tenantsService', mod: tenantsModule },
    { name: 'textbookService', mod: textbookModule },
    { name: 'tripsService', mod: tripsModule },
    { name: 'udaService', mod: udaModule },
    { name: 'userService', mod: userModule },
    { name: 'verbaliService', mod: verbaliModule },
  ]

  it('verifies all 33 service modules have matching default and named exports', () => {
    serviceModules.forEach(({ name, mod }) => {
      expect(mod.default, `${name} must have a default export`).toBeDefined()
      expect(mod[name], `${name} must have a named export '${name}'`).toBeDefined()
      expect(mod.default).toBe(mod[name])
    })
  })
})

describe('Endpoint Parameter Encoding and Blob Timeouts', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    api.get.mockResolvedValue({ data: 'mock-data' })
    api.post.mockResolvedValue({ data: 'mock-data' })
  })

  it('securityService.downloadCadPackage uses params object and 60s timeout', async () => {
    await securityModule.securityService.downloadCadPackage('2025/2026')
    expect(api.get).toHaveBeenCalledWith('/signatures/cad-preservation/download', {
      params: { academic_year: '2025/2026' },
      responseType: 'blob',
      timeout: 60000
    })
  })

  it('securityService.getRecommendedSubstitutes uses params object', async () => {
    await securityModule.securityService.getRecommendedSubstitutes('c1', '2026-08-23', 2, 'math')
    expect(api.get).toHaveBeenCalledWith('/substitutions/recommend-substitutes', {
      params: { class_id: 'c1', date: '2026-08-23', hour: 2, subject_id: 'math' }
    })
  })

  it('sidiService methods use params object and 60s timeout', async () => {
    await sidiModule.sidiService.downloadStudentsXml('class-1')
    expect(api.get).toHaveBeenCalledWith('/reports/sidi/students', {
      params: { class_id: 'class-1' },
      responseType: 'blob',
      timeout: 60000
    })

    await sidiModule.sidiService.downloadScrutiniXml('class-1', 2)
    expect(api.get).toHaveBeenCalledWith('/reports/sidi/scrutini', {
      params: { class_id: 'class-1', semester: 2 },
      responseType: 'blob',
      timeout: 60000
    })

    await sidiModule.sidiService.downloadAttendanceCsv('class-1')
    expect(api.get).toHaveBeenCalledWith('/reports/sidi/attendance', {
      params: { class_id: 'class-1' },
      responseType: 'blob',
      timeout: 60000
    })
  })

  it('competenciesService.getStudentEvaluations uses params object', async () => {
    await competenciesModule.competenciesService.getStudentEvaluations('s1', 2)
    expect(api.get).toHaveBeenCalledWith('/competencies/student/s1', {
      params: { semester: 2 }
    })
  })

  it('attendanceService.exportAttendance has 60s timeout and blob responseType', async () => {
    await attendanceModule.attendanceService.exportAttendance('c1', '2026-08-23')
    expect(api.get).toHaveBeenCalledWith('/attendance/export', {
      params: { class_id: 'c1', date: '2026-08-23' },
      responseType: 'blob',
      timeout: 60000
    })
  })

  it('scrutinyService.exportPagellaPDF has 60s timeout and blob responseType', async () => {
    await scrutinyModule.scrutinyService.exportPagellaPDF('s1', 'c1', 1)
    expect(api.get).toHaveBeenCalledWith('/scrutiny/export/s1/pdf', {
      params: { class_id: 'c1', semester: 1 },
      responseType: 'blob',
      timeout: 60000
    })
  })

  it('verbaliService.getVerbali returns empty data for invalid id without calling API', async () => {
    const res1 = await verbaliModule.verbaliService.getVerbali(null)
    expect(res1).toEqual({ data: [] })
    expect(api.get).not.toHaveBeenCalled()

    const res2 = await verbaliModule.verbaliService.getVerbali('undefined')
    expect(res2).toEqual({ data: [] })
    expect(api.get).not.toHaveBeenCalled()
  })
})
