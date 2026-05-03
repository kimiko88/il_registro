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
        if (role === 'admin' || role === 'superadmin') {
            const response = await api.get('/admin/dashboard/stats')
            return response.data
        }
        
        // For other roles, this could be expanded later
        // For now, return empty or mock that will be handled in the component
        return {
            stats: [],
            schedule: [],
            announcements: []
        }
    }
}
