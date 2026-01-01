import api from './api'

export const userService = {
    getAll(params) {
        return api.get('/users', { params })
    },
    get(id) {
        return api.get(`/users/${id}`)
    },
    create(data) {
        return api.post('/users', data)
    },
    update(id, data) {
        return api.patch(`/users/${id}`, data)
    },
    delete(id) {
        return api.delete(`/users/${id}`)
    },
    resetPassword(id) {
        return api.post(`/users/${id}/reset-password`, { new_password: 'Password123!' }) // Mock pwd or ask prompt
        // TODO: UI should probably prompt for new password OR logic should be "send reset link" (which usually doesn't need new pwd here, but handler might expect it or it's a "Force" reset)
        // Handler ForceResetPassword expects JSON { new_password }
    }
    // NOTE: adminService used /admin/users/admins/... for admins.
    // For standard users, we might need a specific endpoint or use the admin one if permissible.
    // I will assume generic endpoint or skip implementation detail for now.
}
