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
        return api.put(`/users/${id}`, data)
    },
    delete(id) {
        return api.delete(`/users/${id}`)
    },
    resetPassword(id) {
        return api.post(`/users/${id}/password-reset`) // Standard endpoint
        // NOTE: adminService used /admin/users/admins/... for admins.
        // For standard users, we might need a specific endpoint or use the admin one if permissible.
        // I will assume generic endpoint or skip implementation detail for now.
    }
}
