import api from './api'

/**
 * Admin service for SuperAdmin and Admin operations
 */
export default {
    // ========== Dashboard ==========

    getDashboardStats() {
        return api.get('/admin/dashboard/stats')
    },

    // ========== Schools Management ==========

    getSchools(params) {
        return api.get('/admin/schools', { params })
    },

    getSchool(id) {
        return api.get(`/admin/schools/${id}`)
    },

    createSchool(data) {
        return api.post('/admin/schools', data)
    },

    updateSchool(id, data) {
        return api.put(`/admin/schools/${id}`, data)
    },

    deleteSchool(id) {
        return api.delete(`/admin/schools/${id}`)
    },

    getSchoolClasses(schoolId, academicYear = null) {
        const params = { school_id: schoolId }
        if (academicYear) params.academic_year = academicYear
        return api.get('/classes', { params })
    },

    getClasses(schoolId, academicYear = null) {
        return this.getSchoolClasses(schoolId, academicYear)
    },

    getSchoolUsers(schoolId, role = null) {
        const params = { school_id: schoolId }
        if (role) params.role = role
        return api.get('/users', { params })
    },

    createUser(data) {
        return api.post('/users', data)
    },

    updateUser(id, data) {
        return api.patch(`/users/${id}`, data)
    },

    deleteUser(id) {
        return api.delete(`/users/${id}`)
    },

    // ========== Classes Management ==========

    createClass(data) {
        return api.post('/classes', data)
    },

    updateClass(id, data) {
        return api.put(`/classes/${id}`, data)
    },

    deleteClass(id) {
        return api.delete(`/classes/${id}`)
    },

    getClassSubjects(classId) {
        return api.get(`/classes/${classId}/subjects`)
    },

    assignSubjectToClass(classId, data) {
        return api.post(`/classes/${classId}/subjects`, data)
    },

    removeClassSubject(classId, assignmentId) {
        return api.delete(`/classes/${classId}/subjects/${assignmentId}`)
    },

    getClassSchedule(classId) {
        return api.get(`/classes/${classId}/schedule`)
    },

    saveClassSchedule(classId, data) {
        return api.post(`/classes/${classId}/schedule`, data)
    },

    // ========== Subjects Management ==========

    getSubjects(schoolId) {
        return api.get('/subjects', { params: { school_id: schoolId } })
    },

    createSubject(data) {
        return api.post('/subjects', data)
    },

    updateSubject(id, data) {
        return api.put(`/subjects/${id}`, data)
    },

    deleteSubject(id) {
        return api.delete(`/subjects/${id}`)
    },

    // ========== Teachers Management ==========

    // Get specialized teacher list (with qualification etc)
    getTeachersList(schoolId, subjectId = null) {
        const params = { school_id: schoolId }
        if (subjectId) params.subject_id = subjectId
        return api.get('/teachers', { params })
    },

    getTeacherSubjects(teacherId) {
        return api.get(`/teachers/${teacherId}/subjects`)
    },

    assignSubjectToTeacher(teacherId, subjectId) {
        return api.post(`/teachers/${teacherId}/subjects`, { subject_id: subjectId })
    },

    removeTeacherSubject(teacherId, subjectId) {
        return api.delete(`/teachers/${teacherId}/subjects/${subjectId}`)
    },

    getTeacherSchedule(teacherId) {
        return api.get(`/teachers/${teacherId}/schedule`)
    },

    saveTeacherSchedule(teacherId, data) {
        return api.post(`/teachers/${teacherId}/schedule`, data)
    },

    // ========== Admin Users Management (SuperAdmin only) ==========

    getAdmins(params) {
        return api.get('/admin/users/admins', { params })
    },

    createAdmin(data) {
        return api.post('/admin/users/admins', data)
    },

    updateAdmin(id, data) {
        return api.put(`/admin/users/admins/${id}`, data)
    },

    deleteAdmin(id) {
        return api.delete(`/admin/users/admins/${id}`)
    },

    resetAdminPassword(id) {
        return api.post(`/admin/users/admins/${id}/reset-password`)
    },

    getAdminActivity(id, limit = 50) {
        return api.get(`/admin/users/admins/${id}/activity`, { params: { limit } })
    },

    // ========== Audit Logs ==========

    getAuditLogs(params) {
        return api.get('/admin/audit-logs', { params })
    },

    // ========== System Monitoring (SuperAdmin only) ==========

    getHealthStatus() {
        return api.get('/admin/monitoring/health')
    },

    getAPIMetrics() {
        return api.get('/admin/monitoring/metrics/api')
    },

    getStorageUsage() {
        return api.get('/admin/monitoring/metrics/storage')
    },

    getErrorRate() {
        return api.get('/admin/monitoring/metrics/errors')
    },

    getActiveUsers() {
        return api.get('/admin/monitoring/users/active')
    },

    getUptime() {
        return api.get('/admin/monitoring/uptime')
    },

    // ========== Analytics ==========

    getUserGrowth(params) {
        return api.get('/admin/analytics/users/growth', { params })
    },

    getSchoolDistribution() {
        return api.get('/admin/analytics/schools/distribution')
    },

    getAPIUsage() {
        return api.get('/admin/analytics/api/usage')
    },

    getActiveSchools(limit = 10) {
        return api.get('/admin/analytics/schools/active', { params: { limit } })
    },

    getPerformanceTrends(params) {
        return api.get('/admin/analytics/performance', { params })
    },

    // ========== School Settings ==========

    getSchoolSetting(key) {
        return api.get(`/admin/settings/${key}`)
    },

    updateSchoolSetting(key, value) {
        return api.put(`/admin/settings/${key}`, { value: String(value) })
    }
}
