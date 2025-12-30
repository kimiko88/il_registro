import api from './api';

export default {
    getInbox(params) {
        return api.get('/documents/inbox', { params });
    },
    getDocument(id) {
        return api.get(`/documents/${id}`);
    },
    reviewDocument(id, decision, notes) {
        return api.post(`/documents/${id}/review`, { decision, notes });
    },
    archiveDocument(id) {
        return api.post(`/documents/${id}/archive`);
    }
};
