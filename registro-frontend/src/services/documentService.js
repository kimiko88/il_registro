import api from './api';

export const documentService = {
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
    },
    getDocumentVersions(id) {
        return api.get(`/documents/${id}/versions`);
    },
    downloadDocument(id) {
        return api.get(`/documents/${id}/download`, {
            responseType: 'blob',
            timeout: 60000
        });
    },
    uploadDocument(formData) {
        return api.post('/documents/upload', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
            timeout: 60000
        });
    }
};

export default documentService;

