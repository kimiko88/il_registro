import api from './api'

/**
 * Dashboard service for getting user-specific dashboard data
 */
export default {
    /**
     * Get dashboard stats based on user role
     * @param {string} role - User role
     * @returns {Promise} Dashboard statistics
     */
    async getDashboardStats(role) {
        if (role === 'admin' || role === 'superadmin' || role === 'secretary') {
            const response = await api.get('/admin/dashboard/stats')
            return response.data
        }
        if (role === 'teacher') {
            const response = await api.get('/teachers/dashboard/stats')
            return response.data
        }
        if (role === 'student') {
            const response = await api.get('/students/dashboard/stats')
            return response.data
        }
        if (role === 'parent') {
            const response = await api.get('/parents/dashboard/stats')
            return response.data
        }
        return {}
    }
}
