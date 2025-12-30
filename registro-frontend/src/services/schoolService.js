import api from './api';

export default {
    getSchools(params) {
        return api.get('/schools', { params });
    },
    getSchool(id) {
        return api.get(`/schools/${id}`);
    },
    createSchool(data) {
        return api.post('/schools', data);
    },
    updateSchool(id, data) {
        return api.patch(`/schools/${id}`, data);
    },
    deleteSchool(id) {
        return api.delete(`/schools/${id}`);
    },
    importSchools(file) {
        const formData = new FormData();
        formData.append('file', file);
        return api.post('/schools/import', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
        });
    },
};
