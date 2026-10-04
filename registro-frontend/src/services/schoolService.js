import api from './api';

export const schoolService = {
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
    getTiers() {
        return api.get('/schools/tiers');
    },
    getTierFeatures(id) {
        return api.get(`/schools/${id}/tier-features`);
    },
    importSchools(file) {
        const formData = new FormData();
        formData.append('file', file);
        return api.post('/schools/import', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
            timeout: 60000
        });
    },
};

export default schoolService;

