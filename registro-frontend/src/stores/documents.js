import { defineStore } from 'pinia';
import documentService from '../services/documentService';

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
        }
    }
});
