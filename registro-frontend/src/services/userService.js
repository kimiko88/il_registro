import api from './api'

export const userService = {
    getAll(params) {
        return api.get('/users', { params })
    },
    getUsers(params) {
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
    forceResetPassword(id, newPassword) {
        return api.post(`/users/${id}/reset-password`, { new_password: newPassword })
    },
    changePassword(userId, currentPassword, newPassword) {
        return api.post(`/users/${userId}/change-password`, {
            current_password: currentPassword,
            new_password: newPassword
        })
    },

    bulkDelete(userIds) {
        return api.post('/users/bulk-delete', { user_ids: userIds })
    },
    bulkImport(file) {
        const formData = new FormData()
        formData.append('file', file)
        return api.post('/users/bulk-import', formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            },
            timeout: 60000
        })
    },
    getGuardians(studentId) {
        return api.get(`/users/${studentId}/guardians`)
    },
    addGuardian(studentId, data) {
        return api.post(`/users/${studentId}/guardians`, data)
    },
    removeGuardian(studentId, guardianId) {
        return api.delete(`/users/${studentId}/guardians/${guardianId}`)
    },
    getStudentFascicolo(studentId) {
        return api.get(`/users/students/${studentId}/fascicolo`)
    },
    getAssignments(userId) {
        return api.get(`/users/${userId}/assignments`)
    },
    addAssignment(userId, data) {
        return api.post(`/users/${userId}/assignments`, data)
    },
    deleteAssignment(userId, assignmentId) {
        return api.delete(`/users/${userId}/assignments/${assignmentId}`)
    },
    setCoordinatedClasses(userId, classIds) {
        return api.put(`/users/${userId}/coordinated-classes`, { class_ids: classIds })
    }
}

export default userService

