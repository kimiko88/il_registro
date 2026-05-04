import { defineStore } from 'pinia';
import documentService from '../services/documentService';
import api from '@/services/api';

export const useDocumentsStore = defineStore('documents', {
    state: () => ({
        inbox: [],
        currentDocument: null,
        loading: false,
        error: null,
        pagination: {
            page: 1,
            rowsPerPage: 10,
            rowsNumber: 0
        }
    }),
    actions: {
        async fetchInbox(params = {}) {
            this.loading = true;
            try {
                const response = await documentService.getInbox({
                    page: params.page || this.pagination.page,
                    limit: params.rowsPerPage || this.pagination.rowsPerPage,
                    ...params
                });
                this.inbox = response.data.items;
                this.pagination.rowsNumber = response.data.total;
                this.pagination.page = params.page || this.pagination.page;
                this.pagination.rowsPerPage = params.rowsPerPage || this.pagination.rowsPerPage;
            } catch (err) {
                this.error = err.message;
            } finally {
                this.loading = false;
            }
        },
        async fetchDocument(id) {
            this.loading = true;
            try {
                const response = await documentService.getDocument(id);
                this.currentDocument = response.data;
            } finally {
                this.loading = false;
            }
        },
        async reviewDocument(id, decision, notes) {
            await documentService.reviewDocument(id, decision, notes);
            await this.fetchInbox(); // Refresh
        },
        // Teacher Actions
        async fetchMyDocuments() {
            this.loading = true;
            try {
                // Using the inbox endpoint filtered for my documents
                const response = await documentService.getInbox({ filter: 'mine' });
                this.inbox = response.data.items || [];
            } catch (err) {
                this.error = err.message;
            } finally {
                this.loading = false;
            }
        },
        async fetchTemplates() {
            try {
                const response = await api.get('/documents/templates');
                return response.data || [];
            } catch (err) {
                console.error("Failed to fetch templates:", err);
                return [];
            }
        },
        async createDocument(docData) {
            this.loading = true;
            try {
                const response = await api.post('/documents', docData);
                this.inbox.unshift(response.data);
                return response.data;
            } catch (err) {
                console.error("Failed to create document:", err);
                throw err;
            } finally {
                this.loading = false;
            }
        },

        // Student Actions
        async fetchMyFiles() {
            this.loading = true;
            try {
                const response = await api.get('/documents/my-files');
                this.inbox = response.data || [];
            } catch (err) {
                this.error = err.message;
            } finally {
                this.loading = false;
            }
        },

        async signDocument(docId, pin) {
            this.loading = true;
            try {
                // Real API Call
                const res = await api.post('/signatures/', { document_id: docId, pin: pin });
                return res.data;
            } catch (err) {
                console.error("Signing failed", err);
                throw err;
            } finally {
                this.loading = false;
            }
        }
    }
});
