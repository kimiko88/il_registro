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
    }
}
