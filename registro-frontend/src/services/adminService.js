import api from './api';

export default {
    getAdmins(params) {
        return api.get('/users/admins', { params });
    },
    createAdmin(data) {
        return api.post('/users/admins', data);
    },
    updateAdmin(id, data) {
        return api.patch(`/users/${id}`, data);
    },
    deleteAdmin(id) {
        return api.delete(`/users/${id}`);
    },
    resetPassword(id) {
        return api.post(`/users/${id}/reset-password`);
    },
};
