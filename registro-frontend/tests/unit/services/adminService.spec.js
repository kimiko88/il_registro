import { describe, it, expect, vi, beforeEach } from 'vitest'
import adminService from '@/services/adminService'
import api from '@/services/api'

// Mock api
vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        patch: vi.fn(),
        delete: vi.fn()
    }
}))

describe('adminService', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    describe('Dashboard', () => {
        it('getDashboardStats calls api.get', () => {
            adminService.getDashboardStats()
            expect(api.get).toHaveBeenCalledWith('/admin/dashboard/stats')
        })
    })

    describe('Schools', () => {
        it('getSchools calls api.get with params', () => {
            const params = { page: 1 }
            adminService.getSchools(params)
            expect(api.get).toHaveBeenCalledWith('/admin/schools', { params })
        })

        it('getSchool calls api.get with id', () => {
            adminService.getSchool('1')
            expect(api.get).toHaveBeenCalledWith('/admin/schools/1')
        })

        it('createSchool calls api.post', () => {
            const data = { name: 'School' }
            adminService.createSchool(data)
            expect(api.post).toHaveBeenCalledWith('/admin/schools', data)
        })

        it('updateSchool calls api.put', () => {
            const data = { name: 'Updated' }
            adminService.updateSchool('1', data)
            expect(api.put).toHaveBeenCalledWith('/admin/schools/1', data)
        })

        it('deleteSchool calls api.delete', () => {
            adminService.deleteSchool('1')
            expect(api.delete).toHaveBeenCalledWith('/admin/schools/1')
        })

        it('getSchoolClasses calls api.get', () => {
            adminService.getSchoolClasses('school1')
            expect(api.get).toHaveBeenCalledWith('/classes', { params: { school_id: 'school1' } })
        })

        it('getSchoolUsers calls api.get', () => {
            adminService.getSchoolUsers('school1', 'student')
            expect(api.get).toHaveBeenCalledWith('/users', { params: { school_id: 'school1', role: 'student' } })
        })
    })

    describe('Users (Generic)', () => {
        it('createUser calls api.post', () => {
            const data = { name: 'User' }
            adminService.createUser(data)
            expect(api.post).toHaveBeenCalledWith('/users', data)
        })

        it('updateUser calls api.patch', () => {
            const data = { name: 'Updated' }
            adminService.updateUser('1', data)
            expect(api.patch).toHaveBeenCalledWith('/users/1', data)
        })

        it('deleteUser calls api.delete', () => {
            adminService.deleteUser('1')
            expect(api.delete).toHaveBeenCalledWith('/users/1')
        })
    })

    describe('Classes', () => {
        it('createClass calls api.post', () => {
            adminService.createClass({})
            expect(api.post).toHaveBeenCalledWith('/classes', {})
        })

        it('updateClass calls api.put', () => {
            adminService.updateClass('1', {})
            expect(api.put).toHaveBeenCalledWith('/classes/1', {})
        })

        it('deleteClass calls api.delete', () => {
            adminService.deleteClass('1')
            expect(api.delete).toHaveBeenCalledWith('/classes/1')
        })

        it('getClassSubjects calls api.get', () => {
            adminService.getClassSubjects('1')
            expect(api.get).toHaveBeenCalledWith('/classes/1/subjects')
        })

        it('assignSubjectToClass calls api.post', () => {
            adminService.assignSubjectToClass('1', { id: 2 })
            expect(api.post).toHaveBeenCalledWith('/classes/1/subjects', { id: 2 })
        })

        it('removeClassSubject calls api.delete', () => {
            adminService.removeClassSubject('1', '2')
            expect(api.delete).toHaveBeenCalledWith('/classes/1/subjects/2')
        })
    })

    describe('Subjects', () => {
        it('getSubjects calls api.get', () => {
            adminService.getSubjects('school1')
            expect(api.get).toHaveBeenCalledWith('/subjects', { params: { school_id: 'school1' } })
        })
        it('createSubject calls api.post', () => { adminService.createSubject({}); expect(api.post).toHaveBeenCalled() })
        it('updateSubject calls api.put', () => { adminService.updateSubject('1', {}); expect(api.put).toHaveBeenCalled() })
        it('deleteSubject calls api.delete', () => { adminService.deleteSubject('1'); expect(api.delete).toHaveBeenCalled() })
    })

    describe('Teachers', () => {
        it('getTeachersList calls api.get', () => {
            adminService.getTeachersList('school1', 'subject1')
            expect(api.get).toHaveBeenCalledWith('/teachers', { params: { school_id: 'school1', subject_id: 'subject1' } })
        })
        it('getTeacherSubjects', () => { adminService.getTeacherSubjects('1'); expect(api.get).toHaveBeenCalledWith('/teachers/1/subjects') })
        it('assignSubjectToTeacher', () => { adminService.assignSubjectToTeacher('1', 'sub1'); expect(api.post).toHaveBeenCalledWith('/teachers/1/subjects', { subject_id: 'sub1' }) })
        it('removeTeacherSubject', () => { adminService.removeTeacherSubject('1', 'sub1'); expect(api.delete).toHaveBeenCalledWith('/teachers/1/subjects/sub1') })
        it('getTeacherSchedule', () => { adminService.getTeacherSchedule('1'); expect(api.get).toHaveBeenCalledWith('/teachers/1/schedule') })
        it('saveTeacherSchedule', () => { adminService.saveTeacherSchedule('1', { entries: [] }); expect(api.post).toHaveBeenCalledWith('/teachers/1/schedule', { entries: [] }) })
    })

    describe('Admin Users', () => {
        it('getAdmins', () => { adminService.getAdmins({}); expect(api.get).toHaveBeenCalled() })
        it('createAdmin', () => { adminService.createAdmin({}); expect(api.post).toHaveBeenCalled() })
        it('updateAdmin', () => { adminService.updateAdmin('1', {}); expect(api.put).toHaveBeenCalled() })
        it('deleteAdmin', () => { adminService.deleteAdmin('1'); expect(api.delete).toHaveBeenCalled() })
        it('resetAdminPassword', () => { adminService.resetAdminPassword('1'); expect(api.post).toHaveBeenCalled() })
        it('getAdminActivity', () => { adminService.getAdminActivity('1'); expect(api.get).toHaveBeenCalled() })
        it('getAuditLogs', () => { adminService.getAuditLogs({}); expect(api.get).toHaveBeenCalled() })
    })

    describe('Monitoring & Analytics', () => {
        it('getHealthStatus', () => { adminService.getHealthStatus(); expect(api.get).toHaveBeenCalled() })
        it('getAPIMetrics', () => { adminService.getAPIMetrics(); expect(api.get).toHaveBeenCalled() })
        it('getStorageUsage', () => { adminService.getStorageUsage(); expect(api.get).toHaveBeenCalled() })
        it('getErrorRate', () => { adminService.getErrorRate(); expect(api.get).toHaveBeenCalled() })
        it('getActiveUsers', () => { adminService.getActiveUsers(); expect(api.get).toHaveBeenCalled() })
        it('getUptime', () => { adminService.getUptime(); expect(api.get).toHaveBeenCalled() })
        it('getUserGrowth', () => { adminService.getUserGrowth({}); expect(api.get).toHaveBeenCalled() })
        it('getSchoolDistribution', () => { adminService.getSchoolDistribution(); expect(api.get).toHaveBeenCalled() })
        it('getAPIUsage', () => { adminService.getAPIUsage(); expect(api.get).toHaveBeenCalled() })
        it('getActiveSchools', () => { adminService.getActiveSchools(); expect(api.get).toHaveBeenCalled() })
        it('getPerformanceTrends', () => { adminService.getPerformanceTrends({}); expect(api.get).toHaveBeenCalled() })
    })
})
